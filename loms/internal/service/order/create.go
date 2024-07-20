package orderservice

import (
	"fmt"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/model"
)

func (s *oService) Create(user int64, items model.AllNeedReserve) (int64, error) {
	orderID, err := s.orderStore.AddOrder(user, resevToOrders(items))
	if err != nil {
		s.logger.Debugf(s.ctx, "OrderService.Create: Ошибка создания заказа - %v", err)

		return 0, fmt.Errorf("Ошибка создания заказа - %w", err)
	}

	reserved := false

	rsrvdErr := s.stockStore.AddReserve(items)
	if rsrvdErr != nil {
		s.logger.Debugf(s.ctx, "OrderService.Create: Ошибка резервирования товара - %v", rsrvdErr)
	} else {
		reserved = true
	}

	var status string
	if reserved {
		status = model.OrderStatusAwaitingPayment
	} else {
		status = model.OrderStatusFailed
	}

	err = s.orderStore.SetStatus(orderID, status)
	if err != nil {
		s.logger.Debugf(s.ctx, "OrderService.Create: Ошибка смены статуса заказа - %v", err)

		return 0, fmt.Errorf("Ошибка смены статуса заказа -  %w", err)
	}

	if !reserved {
		return 0, fmt.Errorf("Ошибка резервирования товара - %w", rsrvdErr)
	}

	return orderID, nil
}
