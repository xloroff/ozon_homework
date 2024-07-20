package orderstore

import (
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/model"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/repository/order_store/sqlc"
)

func (ms *orderStorage) GetOrder(orderID int64) (*model.Order, error) {
	ms.logger.Debugf(ms.ctx, fmt.Sprintf("orderStore.GetOrder: начинаю получение заказа orderID: %d", orderID))
	defer ms.logger.Debugf(ms.ctx, fmt.Sprintf("orderStore.GetOrder: закончил получение заказа orderID: %d", orderID))

	dbRePool := ms.data.GetReaderPool()
	q := sqlc.New(dbRePool)

	order, err := q.GetOrder(ms.ctx, orderID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrOrderNotFound
		}

		return nil, fmt.Errorf("Ошибка получения заказа %d - %w", orderID, err)
	}

	items, err := q.GetOrderItemsByOrderIDs(ms.ctx, orderID)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("Ошибка получения товаров в заказе %d - %w", orderID, err)
		}
	}

	return &model.Order{
		ID:     order.ID,
		User:   order.User,
		Status: string(order.Status),
		Items:  toOrderItems(items),
	}, nil
}

func toOrderItems(items []sqlc.GetOrderItemsByOrderIDsRow) model.OrderItems {
	res := make(model.OrderItems, 0, len(items))

	for _, item := range items {
		res = append(res, &model.OrderItem{
			ID:    item.ID,
			Sku:   item.Sku,
			Count: uint16(item.Count),
		})
	}

	return res
}
