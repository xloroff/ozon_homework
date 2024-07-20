package orderconsumer

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/IBM/sarama"
	"go.opentelemetry.io/otel/trace"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/notifier/internal/model"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/notifier/internal/pkg/logger"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/notifier/internal/pkg/tracer"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/notifier/internal/service/order"
)

// Consumer определяет методы обработки сообщений из Kafka дл заказов.
type Consumer interface {
	OrderHandle(ctx context.Context, msg *sarama.ConsumerMessage) (bool, error)
}

type consumer struct {
	ctx          context.Context
	logger       logger.Logger
	orderService orderservice.Service
}

// NewConsumer создает новый Consumer для заказов.
func NewConsumer(ctx context.Context, l logger.Logger, os orderservice.Service) Consumer {
	return &consumer{
		ctx:          ctx,
		logger:       l,
		orderService: os,
	}
}

// OrderHandle обрабатывает сообщения из брокера для заказов.
func (c *consumer) OrderHandle(ctx context.Context, msg *sarama.ConsumerMessage) (bool, error) {
	ctx, span := tracer.StartSpanFromContext(ctx, "orderconsumer.handle", trace.WithSpanKind(trace.SpanKindConsumer))
	span.SetTag("component", "orderconsumer")

	defer span.End()

	c.logger.Debugf(ctx, "OrderConsumer.Handle: Начинаю обработку нового эвента в заказах. Offset: %d. Partition: %d", msg.Offset, msg.Partition)
	defer c.logger.Debugf(ctx, "OrderConsumer.Handle: Закончил обработку нового эвента в заказах. Offset: %d. Partition: %d", msg.Offset, msg.Partition)

	payload := model.OrderEventPayload{}

	err := json.Unmarshal(msg.Value, &payload)
	if err != nil {
		span.SetTag("error", true)
		c.logger.Errorf(ctx, "OutboxService.sendEventOrder: Ошибка преобразования полученного из брокера сообщения - %v", err)

		return false, fmt.Errorf("OutboxService.sendEventOrder: Ошибка преобразования полученного из брокера сообщения - %w", err)
	}

	return c.statusChanged(ctx, &payload)
}
