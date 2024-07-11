package orderservice

import (
	"context"
	"fmt"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/model"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/pkg/metrics"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/pkg/tracer"
)

// Pay помечает оплату заказа (меняет статус заказа).
func (s *service) Pay(ctx context.Context, orderID int64) error {
	ctx, span := tracer.StartSpanFromContext(ctx, "service.orderservice.pay")
	span.SetTag("component", "orderservice")

	defer span.End()

	order, err := s.orderStore.GetOrder(ctx, orderID)
	if err != nil {
		span.SetTag("error", true)
		s.logger.Debugf(ctx, "OrderService.Pay: Ошибка получения заказа - %v", err)

		return fmt.Errorf("Ошибка получения заказа - %w", err)
	}

	delReserve := orderToReserve(order.Items)

	err = s.stockStore.DelItemFromReserve(ctx, delReserve)
	if err != nil {
		span.SetTag("error", true)
		s.logger.Debugf(ctx, "OrderService.Pay: Ошибка возвращение резервов товарам - %v", err)

		return fmt.Errorf("Ошибка возвращение резервов товарам - %w", err)
	}

	err = s.orderStore.SetStatus(ctx, order.ID, order.User, model.OrderStatusPayed)
	if err != nil {
		span.SetTag("error", true)
		s.logger.Debugf(ctx, "OrderService.Pay: Ошибка смены статуса заказа - %v", err)

		return fmt.Errorf("Ошибка смены статуса заказа - %w", err)
	}

	metrics.UpdateOrderStatusChanged(order.Status, model.OrderStatusPayed)

	return nil
}
