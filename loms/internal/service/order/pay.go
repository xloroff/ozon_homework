package orderservice

import (
	"fmt"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/model"
)

// Pay помечает оплату заказа (меняет статус заказа).
func (s *oService) Pay(orderID int64) error {
	order, err := s.orderStore.GetOrder(orderID)
	if err != nil {
		s.logger.Debugf(s.ctx, "OrderService.Pay: Ошибка получения заказа - %v", err)

		return fmt.Errorf("Ошибка получения заказа - %w", err)
	}

	delReserve := orderToReserve(order.Items)

	err = s.stockStore.DelItemFromReserve(delReserve)
	if err != nil {
		s.logger.Debugf(s.ctx, "OrderService.Pay: Ошибка возвращение резервов товарам - %v", err)

		return fmt.Errorf("Ошибка возвращение резервов товарам - %w", err)
	}

	err = s.orderStore.SetStatus(order.ID, model.OrderStatusPayed)
	if err != nil {
		s.logger.Debugf(s.ctx, "OrderService.Pay: Ошибка смены статуса заказа - %v", err)

		return fmt.Errorf("Ошибка смены статуса заказа - %w", err)
	}

	return nil
}
