package orderstore

import (
	"fmt"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/model"
)

func (ms *orderStorage) GetOrder(orderID int64) (*model.Order, error) {
	ms.logger.Debugf(ms.ctx, fmt.Sprintf("orderStore.GetOrder: начинаю получение заказа orderID: %d", orderID))
	defer ms.logger.Debugf(ms.ctx, fmt.Sprintf("orderStore.GetOrder: закончил получение заказа orderID: %d", orderID))

	ms.RLock()
	defer ms.RUnlock()

	order, ok := ms.data[orderID]
	if !ok {
		return nil, fmt.Errorf("Заказ %d не найден, ошибка - %w", orderID, model.ErrOrderNotFound)
	}

	return order, nil
}
