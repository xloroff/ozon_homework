package orderconsumer

import (
	"context"
	"fmt"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/notifier/internal/model"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/notifier/internal/pkg/tracer"
)

func (c *consumer) statusChanged(ctx context.Context, p *model.OrderEventPayload) (bool, error) {
	ctx, span := tracer.StartSpanFromContext(ctx, "orderconsumer.status_changed")
	span.SetTag("component", "orderconsumer")
	defer span.End()

	err := c.orderService.OrderStatusChanges(ctx, p.Payload.ID, p.Payload.Status)
	if err != nil {
		span.SetTag("error", true)
		return false, fmt.Errorf("OrderConsumer.statusChanged: Ошибка обработки эвента заказов - %w", err)
	}

	return true, nil
}
