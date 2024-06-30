package orderstore

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/model"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/pkg/logger"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/pkg/metrics"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/pkg/tracer"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/repository/order_store/sqlc"
)

// AddOrder добавляет заказ в хранилище.
func (ms *orderStorage) AddOrder(ctx context.Context, user int64, items model.OrderItems) (int64, error) {
	ctx, span := tracer.StartSpanFromContext(ctx, "repository.orderstore.add_order")
	span.SetTag("component", "orderstore")
	span.SetTag("db.type", "sql")
	span.SetTag("db.statement", "insert")

	defer span.End()

	metrics.UpdateDatabaseRequestsTotal(
		repName,
		"AddOrder",
		"insert",
	)
	defer metrics.UpdateDatabaseResponseTime(time.Now().UTC())

	ctx = logger.AddFieldsToContext(ctx, "data", items, "user_id", user)

	ms.logger.Debugf(ctx, fmt.Sprintf("orderStore.AddOrder: начинаю добавление заказа userId: %d", user))
	defer ms.logger.Debugf(ctx, fmt.Sprintf("orderStore.AddOrder: закончил добавление заказа userId: %d", user))

	dbWrPool := ms.data.GetWriterPool()

	var orderID int64

	err := dbWrPool.BeginFuncWithTx(ctx, func(tx pgx.Tx) error {
		q := sqlc.New(dbWrPool).WithTx(tx)

		var err error

		orderID, err = q.AddOrder(ctx, sqlc.AddOrderParams{Column1: 1, User: user, Status: model.OrderStatusNew})
		if err != nil {
			span.SetTag("error", true)
			return fmt.Errorf("Ошибка добавлени заказа (строки в таблицу order) - %w", err)
		}

		for _, item := range items {
			err = q.AddOrderItem(ctx, sqlc.AddOrderItemParams{OrderID: orderID, Sku: item.Sku, Count: int32(item.Count)})
			if err != nil {
				span.SetTag("error", true)
				return fmt.Errorf("Ошибка добавлени заказа (строки в таблицу order_item) %d - %w", item.Sku, err)
			}
		}

		return nil
	})
	if err != nil {
		span.SetTag("error", true)
		metrics.UpdateDatabaseResponseCode(
			repName,
			"AddOrder",
			"insert",
			"error",
		)

		return 0, fmt.Errorf("Ошибка создания заказа - %w", err)
	}

	metrics.UpdateDatabaseResponseCode(
		repName,
		"AddOrder",
		"insert",
		"ok",
	)

	return orderID, nil
}
