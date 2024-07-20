package orderstore

import (
	"fmt"

	"go.uber.org/zap"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/model"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/pkg/logger"
)

func (ms *orderStorage) AddOrder(user int64, items model.OrderItems) (int64, error) {
	ctx := logger.Append(ms.ctx, []zap.Field{zap.Any("items", items)})

	ms.logger.Debugf(ctx, fmt.Sprintf("orderStore.AddOrder: начинаю добавление заказа userId: %d", user))
	defer ms.logger.Debugf(ctx, fmt.Sprintf("orderStore.AddOrder: закончил добавление заказа userId: %d", user))

	ms.Lock()
	defer ms.Unlock()

	order := &model.Order{
		ID:     ms.generateID(),
		User:   user,
		Status: model.OrderStatusNew,
		Items:  items,
	}

	ms.data[order.ID] = order

	return order.ID, nil
}

func (ms *orderStorage) generateID() int64 {
	var maxID int64
	for orderID := range ms.data {
		if orderID > maxID {
			maxID = orderID
		}
	}

	maxID++

	return maxID
}
