package orderservice

import (
	"fmt"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/model"
)

// Cancel отменяет заказ, снимает резервы и меняет статус.
func (s *oService) Cancel(orderID int64) error {
	orderStorage, err := s.orderStore.GetOrder(orderID)
	if err != nil {
		s.logger.Debugf(s.ctx, "OrderService.Cancel: Ошибка получения заказа - %v", err)

		return fmt.Errorf("Ошибка получения заказа - %w", err)
	}

	switch orderStorage.Status {
	case model.OrderStatusPayed:
		s.logger.Debugf(s.ctx, "OrderService.Cancel: Ошибка отмены заказа статус заказа %s - %v", orderStorage.Status, model.ErrOrderCancel)

		return model.ErrOrderCancel
	case model.OrderStatusCancelled:
		s.logger.Debugf(s.ctx, "OrderService.Cancel: Ошибка отмены заказа статус заказа %s - %v", orderStorage.Status, model.ErrOrderCancel)

		return model.ErrOrderCancel
	}

	order := orderToReserve(orderStorage.Items)

	err = s.stockStore.CancelReserve(order)
	if err != nil {
		s.logger.Debugf(s.ctx, "OrderService.Cancel: Ошибка возвращение резервов товарам - %v", err)

		return fmt.Errorf("Ошибка возвращение резервов товарам - %w", err)
	}

	err = s.orderStore.SetStatus(orderID, model.OrderStatusCancelled)
	if err != nil {
		s.logger.Debugf(s.ctx, "OrderService.Cancel: Ошибка смены статуса заказа - %v", err)

		return fmt.Errorf("Ошибка смены статуса заказа - %w", err)
	}

	return nil
}
