package orderstore

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/model"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/pkg/metrics"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/pkg/tracer"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/repository/order_store/sqlc"
)

func (ms *orderStorage) GetOrder(ctx context.Context, orderID int64) (*model.Order, error) {
	ctx, span := tracer.StartSpanFromContext(ctx, "repository.orderstore.get_order")
	span.SetTag("component", "orderstore")
	span.SetTag("db.type", "sql")
	span.SetTag("db.statement", "select")

	defer span.End()

	metrics.UpdateDatabaseRequestsTotal(
		repName,
		"GetOrder",
		"select",
	)
	defer metrics.UpdateDatabaseResponseTime(time.Now().UTC())

	ms.logger.Debugf(ctx, fmt.Sprintf("orderStore.GetOrder: начинаю получение заказа orderID: %d", orderID))
	defer ms.logger.Debugf(ctx, fmt.Sprintf("orderStore.GetOrder: закончил получение заказа orderID: %d", orderID))

	dbRePool := ms.data.GetReaderPool()
	q := sqlc.New(dbRePool)

	order, err := q.GetOrder(ctx, orderID)
	if err != nil {
		span.SetTag("error", true)
		metrics.UpdateDatabaseResponseCode(
			repName,
			"GetOrder",
			"select",
			"error",
		)

		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrOrderNotFound
		}

		return nil, fmt.Errorf("Ошибка получения заказа %d - %w", orderID, err)
	}

	metrics.UpdateDatabaseResponseCode(
		repName,
		"GetOrder",
		"select",
		"ok",
	)

	metrics.UpdateDatabaseRequestsTotal(
		repName,
		"GetOrderItemsByOrderIDs",
		"select",
	)

	items, err := q.GetOrderItemsByOrderIDs(ctx, orderID)
	if err != nil {
		span.SetTag("error", true)
		metrics.UpdateDatabaseResponseCode(
			repName,
			"GetOrderItemsByOrderIDs",
			"select",
			"error",
		)

		if !errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("Ошибка получения товаров в заказе %d - %w", orderID, err)
		}
	}

	metrics.UpdateDatabaseResponseCode(
		repName,
		"GetOrderItemsByOrderIDs",
		"select",
		"ok",
	)

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
