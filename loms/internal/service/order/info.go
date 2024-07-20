package orderservice

import (
	"context"
	"fmt"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/model"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/pkg/tracer"
)

// Info возвращает информацию о заказе.
func (s *oService) Info(ctx context.Context, orderID int64) (*model.Order, error) {
	ctx, span := tracer.StartSpanFromContext(ctx, "service.orderservice.info")
	span.SetTag("component", "orderservice")

	defer span.End()

	order, err := s.orderStore.GetOrder(ctx, orderID)
	if err != nil {
		span.SetTag("error", true)
		s.logger.Debugf(ctx, "OrderService.Info: Ошибка получения заказа - %v", err)

		return nil, fmt.Errorf("Ошибка получения заказа - %w", err)
	}

	return order, nil
}
