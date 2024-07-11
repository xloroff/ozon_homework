package orderservice

import (
	"context"
	"fmt"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/model"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/pkg/metrics"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/pkg/tracer"
)

func (s *service) Create(ctx context.Context, user int64, items model.AllNeedReserve) (int64, error) {
	ctx, span := tracer.StartSpanFromContext(ctx, "service.orderservice.create")
	span.SetTag("component", "orderservice")

	defer span.End()

	orderID, err := s.orderStore.AddOrder(ctx, user, resevToOrders(items))
	if err != nil {
		span.SetTag("error", true)
		metrics.UpdateOrdersCreatedError()
		s.logger.Debugf(ctx, "OrderService.Create: Ошибка создания заказа - %v", err)

		return 0, fmt.Errorf("Ошибка создания заказа - %w", err)
	}

	metrics.UpdateOrdersCreated()

	reserved := false

	rsrvdErr := s.stockStore.AddReserve(ctx, items)
	if rsrvdErr != nil {
		span.SetTag("error", true)
		s.logger.Debugf(ctx, "OrderService.Create: Ошибка резервирования товара - %v", rsrvdErr)
	} else {
		reserved = true
	}

	var status string
	if reserved {
		status = model.OrderStatusAwaitingPayment
	} else {
		status = model.OrderStatusFailed
	}

	err = s.orderStore.SetStatus(ctx, orderID, user, status)
	if err != nil {
		span.SetTag("error", true)
		s.logger.Debugf(ctx, "OrderService.Create: Ошибка смены статуса заказа - %v", err)

		return 0, fmt.Errorf("Ошибка смены статуса заказа -  %w", err)
	}

	metrics.UpdateOrderStatusChanged(model.OrderStatusNew, status)

	if !reserved {
		span.SetTag("error", true)
		return 0, fmt.Errorf("Ошибка резервирования товара - %w", rsrvdErr)
	}

	return orderID, nil
}
