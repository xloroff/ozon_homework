package orderservice

import (
	"context"
	"fmt"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/model"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/pkg/metrics"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/pkg/tracer"
)

// Cancel отменяет заказ, снимает резервы и меняет статус.
func (s *service) Cancel(ctx context.Context, orderID int64) error {
	ctx, span := tracer.StartSpanFromContext(ctx, "service.orderservice.cancel")
	span.SetTag("component", "orderservice")

	defer span.End()

	orderStorage, err := s.orderStore.GetOrder(ctx, orderID)
	if err != nil {
		span.SetTag("error", true)
		s.logger.Debugf(ctx, "OrderService.Cancel: Ошибка получения заказа - %v", err)

		return fmt.Errorf("Ошибка получения заказа - %w", err)
	}

	switch orderStorage.Status {
	case model.OrderStatusPayed:
		span.SetTag("error", true)
		s.logger.Debugf(ctx, "OrderService.Cancel: Ошибка отмены заказа статус заказа %s - %v", orderStorage.Status, model.ErrOrderCancel)

		return model.ErrOrderCancel
	case model.OrderStatusCancelled:
		span.SetTag("error", true)
		s.logger.Debugf(ctx, "OrderService.Cancel: Ошибка отмены заказа статус заказа %s - %v", orderStorage.Status, model.ErrOrderCancel)

		return model.ErrOrderCancel
	}

	order := orderToReserve(orderStorage.Items)

	err = s.stockStore.CancelReserve(ctx, order)
	if err != nil {
		span.SetTag("error", true)
		s.logger.Debugf(ctx, "OrderService.Cancel: Ошибка возвращение резервов товарам - %v", err)

		return fmt.Errorf("Ошибка возвращение резервов товарам - %w", err)
	}

	err = s.orderStore.SetStatus(ctx, orderID, orderStorage.User, model.OrderStatusCancelled)
	if err != nil {
		span.SetTag("error", true)
		s.logger.Debugf(ctx, "OrderService.Cancel: Ошибка смены статуса заказа - %v", err)

		return fmt.Errorf("Ошибка смены статуса заказа - %w", err)
	}

	metrics.UpdateOrderStatusChanged(orderStorage.Status, model.OrderStatusCancelled)

	return nil
}
