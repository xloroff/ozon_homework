package orderstore

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/pkg/metrics"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/pkg/tracer"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/repository/order_store/sqlc"
)

func (ms *orderStorage) SetStatus(ctx context.Context, orderID, userID int64, status string) error {
	ctx, span := tracer.StartSpanFromContext(ctx, "repository.orderstore.set_status")
	span.SetTag("component", "orderstore")
	span.SetTag("db.type", "sql")
	span.SetTag("db.statement", "update")

	defer span.End()

	metrics.UpdateDatabaseRequestsTotal(
		repName,
		"SetStatus",
		"update",
	)
	defer metrics.UpdateDatabaseResponseTime(time.Now().UTC())

	ms.logger.Debugf(ctx, fmt.Sprintf("orderStore.SetStatus: начинаю смену статуса заказа orderID: %d", orderID))
	defer ms.logger.Debugf(ctx, fmt.Sprintf("orderStore.SetStatus: закончил смену статуса заказа orderID: %d", orderID))

	dbWrPool := ms.data.GetWriterPool()

	err := dbWrPool.BeginFuncWithTx(ctx, func(tx pgx.Tx) error {
		q := sqlc.New(dbWrPool).WithTx(tx)

		var err error

		err = q.SetStatus(ctx, sqlc.SetStatusParams{ID: orderID, Status: sqlc.OrderStatusType(status)})
		if err != nil {
			span.SetTag("error", true)
			metrics.UpdateDatabaseResponseCode(
				repName,
				"SetStatus",
				"update",
				"error",
			)

			return fmt.Errorf("Ошибка обновления статуса заказа %d - %w", orderID, err)
		}

		msg, err := createMessage(ctx, orderID, userID, status)
		if err != nil {
			span.SetTag("error", true)
			return fmt.Errorf("Ошибка формирования сообщения в outbox для отправки в брокер - %w", err)
		}

		err = ms.outboxstore.AddMessage(ctx, tx, msg)
		if err != nil {
			span.SetTag("error", true)
			return fmt.Errorf("Ошибка добавлени измененного статуса заказа в outbox для отправки в брокер - %w", err)
		}

		return nil
	})
	if err != nil {
		span.SetTag("error", true)
		metrics.UpdateDatabaseResponseCode(
			repName,
			"SetStatus",
			"update",
			"error",
		)

		return fmt.Errorf("Ошибка смены статуса заказа - %w", err)
	}

	metrics.UpdateDatabaseResponseCode(
		repName,
		"SetStatus",
		"update",
		"ok",
	)

	return nil
}
