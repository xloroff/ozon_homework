package orderstore

import (
	"context"
	"encoding/json"
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

	// Важно помнить, что открытая транзакция тут больше антипаттерн в outbox и применена она исключительно потому, что это позволительно в рамках задания, а так же чтобы не усложнять логику inbox на стороне сервиса notifier.
	err := dbWrPool.BeginFuncWithTx(ctx, func(tx pgx.Tx) error {
		q := sqlc.New(dbWrPool).WithTx(tx)

		var err error

		orderID, err = q.AddOrder(ctx, sqlc.AddOrderParams{Column1: 1, User: user, Status: model.OrderStatusNew})
		if err != nil {
			span.SetTag("error", true)
			return fmt.Errorf("Ошибка добавлени заказа (строки в таблицу order) - %w", err)
		}

		msg, err := createMessage(ctx, orderID, user, model.OrderStatusNew)
		if err != nil {
			span.SetTag("error", true)
			return fmt.Errorf("Ошибка формирования сообщения в outbox для отправки в брокер - %w", err)
		}

		for _, item := range items {
			err = q.AddOrderItem(ctx, sqlc.AddOrderItemParams{OrderID: orderID, Sku: item.Sku, Count: int32(item.Count)})
			if err != nil {
				span.SetTag("error", true)
				return fmt.Errorf("Ошибка добавлени заказа (строки в таблицу order_item) %d - %w", item.Sku, err)
			}
		}

		err = ms.outboxstore.AddMessage(ctx, tx, msg)
		if err != nil {
			span.SetTag("error", true)
			return fmt.Errorf("Ошибка добавлени заказа в outbox для отправки в брокер - %w", err)
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

func createMessage(_ context.Context, orderID, user int64, status string) (*model.Outbox, error) {
	data := &model.Order{
		ID:     orderID,
		User:   user,
		Status: status,
	}

	payload, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("Ошибка сериализации данных для отправки в брокер - %w", err)
	}

	outbox := &model.Outbox{
		EntityID: fmt.Sprintf("%d", orderID),
		Payload:  string(payload),
	}

	return outbox, nil
}
