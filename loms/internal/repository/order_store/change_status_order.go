package orderstore

import (
	"fmt"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/model"
)

func (ms *orderStorage) SetStatus(orderID int64, status string) error {
	ms.logger.Debugf(ms.ctx, fmt.Sprintf("orderStore.SetStatus: начинаю смену статуса заказа orderID: %d", orderID))
	defer ms.logger.Debugf(ms.ctx, fmt.Sprintf("orderStore.SetStatus: закончил смену статуса заказа orderID: %d", orderID))

	ms.Lock()
	defer ms.Unlock()

	order, ok := ms.data[orderID]
	if !ok {
		return fmt.Errorf("Заказ %d не найден, ошибка - %w", orderID, model.ErrOrderNotFound)
	}

	order.Status = status

	return nil
}
