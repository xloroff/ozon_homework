package outboxservice

import (
	"context"
	"encoding/json"
	"fmt"

	"golang.org/x/time/rate"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/config"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/model"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/pkg/errgroup"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/pkg/tracer"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/repository/outbox_store/sqlc"
)

func (s *service) processOutboxEvents(ctx context.Context) error {
	events, err := s.outboxstore.GetEvents(ctx)
	if err != nil {
		return fmt.Errorf("Ошибка получения сообщений из хранилища outbox - %w", err)
	}

	erg, ctx := errgroup.NewErrGroup(ctx)
	limiter := rate.NewLimiter(rate.Limit(config.RPS), 1)

	for _, event := range events {
		erg.Go(func() error {
			if err = limiter.Wait(ctx); err != nil {
				return fmt.Errorf("Ошибка ожидания при срабатывании лимита RPS %v - %w", config.RPS, err)
			}

			err = s.sendEvent(ctx, event)
			if err != nil {
				return err
			}

			return nil
		})
	}

	if err = erg.Wait(); err != nil {
		return fmt.Errorf("Ошибки отправки сообщений в брокер -  %v", erg.ErrsToString())
	}

	return nil
}

func (s *service) sendEvent(ctx context.Context, event *sqlc.OutboxRow) error {
	var metadata model.Metadata

	err := json.Unmarshal(event.Metadata, &metadata)
	if err != nil {
		return fmt.Errorf("OutboxService.sendEventOrder: Ошибка преобразования metadata сообщения перед отправкой в брокер - %w", err)
	}

	ctx, span, err := tracer.StartSpanFromIDs(ctx, metadata.TraceID, metadata.SpanID, "service.outboxService.send_event_order")
	if err != nil {
		s.logger.Warnf(ctx, "OutboxService.sendEvents: Ошибка получениея context_id, span_id для передачи через брокер - %v", err)
		ctx, span = tracer.StartSpanFromContext(ctx, "service.outboxservice.send_event_order")
	}

	span.SetTag("component", "outboxservice")
	defer span.End()

	order := &model.Order{}

	err = json.Unmarshal([]byte(event.Payload.String), order)
	if err != nil {
		span.SetTag("error", true)
		return fmt.Errorf("OutboxService.sendEventOrder: Ошибка преобразования payload сообщения для отправки в брокер - %w", err)
	}

	msgToSend := model.OrderEventPayload{
		ID:       event.ID,
		Time:     event.CreatedAt.Time,
		EntityID: event.EntityID.String,
		Payload:  *order,
	}

	err = s.outboxProducer.Send(ctx, event.EntityID.String, msgToSend)
	if err != nil {
		s.logger.Errorf(ctx, "OutboxService.sendEvents: Ошибка отправки сообщения в брокер - %v", err)
		span.SetTag("error", true)
	}

	err = s.outboxstore.SetStatus(ctx, event)
	if err != nil {
		span.SetTag("error", true)
		return fmt.Errorf("Ошибка сохранения статуса сообщения в хранилище outbox - %w", err)
	}

	return nil
}
