package orderstore

import (
	"fmt"

	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/model"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/pkg/logger"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/repository/order_store/sqlc"
)

// AddOrder добавляет заказ в хранилище.
func (ms *orderStorage) AddOrder(user int64, items model.OrderItems) (int64, error) {
	ctx := logger.Append(ms.ctx, []zap.Field{zap.Any("items", items)})

	ms.logger.Debugf(ctx, fmt.Sprintf("orderStore.AddOrder: начинаю добавление заказа userId: %d", user))
	defer ms.logger.Debugf(ctx, fmt.Sprintf("orderStore.AddOrder: закончил добавление заказа userId: %d", user))

	dbWrPool := ms.data.GetWriterPool()

	var orderID int64

	err := dbWrPool.BeginFuncWithTx(ms.ctx, func(tx pgx.Tx) error {
		q := sqlc.New(dbWrPool).WithTx(tx)

		var err error

		orderID, err = q.AddOrder(ctx, sqlc.AddOrderParams{Column1: 1, User: user, Status: model.OrderStatusNew})
		if err != nil {
			return fmt.Errorf("Ошибка добавлени заказа (строки в таблицу order) - %w", err)
		}

		for _, item := range items {
			err = q.AddOrderItem(ctx, sqlc.AddOrderItemParams{OrderID: orderID, Sku: item.Sku, Count: int32(item.Count)})
			if err != nil {
				return fmt.Errorf("Ошибка добавлени заказа (строки в таблицу order_item) %d - %w", item.Sku, err)
			}
		}

		return nil
	})
	if err != nil {
		return 0, fmt.Errorf("Ошибка создания заказа - %w", err)
	}

	return orderID, nil
}
