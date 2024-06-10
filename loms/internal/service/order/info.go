package orderservice

import (
	"fmt"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/model"
)

// Info возвращает информацию о заказе.
func (s *oService) Info(orderID int64) (*model.Order, error) {
	order, err := s.orderStore.GetOrder(orderID)
	if err != nil {
		s.logger.Debugf(s.ctx, "OrderService.Info: Ошибка получения заказа - %v", err)

		return nil, fmt.Errorf("Ошибка получения заказа - %w", err)
	}

	return order, nil
}
