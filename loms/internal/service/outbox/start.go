package outboxservice

import (
	"context"
	"time"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/pkg/tracer"
)

// Start запускает outbox-сервис для отправки в брокер.
func (s *service) Start(ctx context.Context) {
	s.wg.Add(1)

	go func() {
		s.sendOutbox(ctx)
	}()

	s.logger.Warn(s.ctx, "Outbox сервис отправки сообщений в брокер запущен...")
}

func (s *service) sendOutbox(ctx context.Context) {
	ctx, span := tracer.StartSpanFromContext(ctx, "service.outboxservice.sendMessages")
	span.SetTag("component", "outboxservice")
	defer span.End()

	ticker := time.NewTicker(time.Duration(s.periodSend) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-s.sendStatuses:
			s.wg.Done()
		case <-ticker.C:
			err := s.processOutboxEvents(ctx)
			if err != nil {
				s.logger.Errorf(ctx, "Ошибка отправки сообщений в брокер - %v", err)
			}
		}
	}
}
