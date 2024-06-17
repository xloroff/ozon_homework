package orderstore

import (
	"fmt"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/repository/order_store/sqlc"
)

func (ms *orderStorage) SetStatus(orderID int64, status string) error {
	ms.logger.Debugf(ms.ctx, fmt.Sprintf("orderStore.SetStatus: начинаю смену статуса заказа orderID: %d", orderID))
	defer ms.logger.Debugf(ms.ctx, fmt.Sprintf("orderStore.SetStatus: закончил смену статуса заказа orderID: %d", orderID))

	dbWrPool := ms.data.GetWriterPool()
	q := sqlc.New(dbWrPool)

	err := q.SetStatus(ms.ctx, sqlc.SetStatusParams{ID: orderID, Status: sqlc.OrderStatusType(status)})
	if err != nil {
		return fmt.Errorf("Ошибка обновления статуса заказа %d - %w", orderID, err)
	}

	return nil
}
