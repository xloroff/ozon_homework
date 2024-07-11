package orderservice

import (
	"context"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/notifier/internal/pkg/tracer"
)

func (s *service) OrderStatusChanges(ctx context.Context, orderID int64, status string) error {
	ctx, span := tracer.StartSpanFromContext(ctx, "orderservice.order_status_changes")
	span.SetTag("component", "orderservice")
	defer span.End()

	s.logger.Infof(ctx, "Статус заказ: %d изменен на: %s", orderID, status)

	return nil
}
