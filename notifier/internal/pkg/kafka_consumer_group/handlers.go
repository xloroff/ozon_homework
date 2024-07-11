package kafkaconsumergroup

import (
	"context"

	"github.com/IBM/sarama"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/notifier/internal/pkg/logger"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/notifier/internal/pkg/tracer"
)

// Handler кастомный тип обработчика для имплементации функций.
type Handler func(ctx context.Context, message *sarama.ConsumerMessage) (bool, error)

type consumerGroupHandler struct {
	ctx     context.Context
	logger  logger.Logger
	handler Handler
}

// newConsumerGroupHandler создает обработчик для сообщений kafka.
func newConsumerGroupHandler(ctx context.Context, l logger.Logger, handler Handler) *consumerGroupHandler {
	return &consumerGroupHandler{
		ctx:     ctx,
		logger:  l,
		handler: handler,
	}
}

// Setup Начинаем новую сессию, до ConsumeClaim.
func (h *consumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}

// Cleanup завершает сессию, после того, как все ConsumeClaim завершатся.
func (h *consumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim читаем до тех пор пока сессия не завершилась.
func (h *consumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				return nil
			}

			ctx := session.Context()
			ctx = h.extractTraceAndSpanID(ctx, message.Headers)

			_, err := h.processMessage(ctx, message)
			if err != nil {
				h.logger.Errorf(ctx, "KafkaConsumerGroup.ConsumeClaim: Ошибка обработки сообщения - %v", err)
			}

			session.MarkMessage(message, "")
			session.Commit()

		case <-session.Context().Done():
			return nil
		}
	}
}

func (h *consumerGroupHandler) extractTraceAndSpanID(ctx context.Context, headers []*sarama.RecordHeader) context.Context {
	var traceID string

	var spanID string

	for _, header := range headers {
		switch string(header.Key) {
		case "x-trace-id":
			traceID = string(header.Value)
		case "x-span-id":
			spanID = string(header.Value)
		}
	}

	ctx, _, err := tracer.StartSpanFromIDs(ctx, traceID, spanID, "service.outboxService.send")
	if err != nil {
		h.logger.Debugf(ctx, "KafkaConsumerGroup.ConsumeClaim: Ошибка получениея context_id, span_id переданных через брокер - %v", err)
		ctx, _ = tracer.StartSpanFromContext(ctx, "pkg.kafkaconsumergroup.consumeclaim")

		return ctx
	}

	return ctx
}

func (h *consumerGroupHandler) processMessage(ctx context.Context, message *sarama.ConsumerMessage) (bool, error) {
	ok, err := h.handler(ctx, message)
	return ok, err
}
