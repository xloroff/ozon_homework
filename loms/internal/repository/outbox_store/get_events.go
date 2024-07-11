package outboxstore

import (
	"context"
	"fmt"
	"time"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/pkg/metrics"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/repository/outbox_store/sqlc"
)

// GetEvents получение эвентов для отправки в брокер.
func (ks *outboxStorage) GetEvents(ctx context.Context) (msgs []*sqlc.OutboxRow, err error) {
	metrics.UpdateDatabaseRequestsTotal(
		repName,
		"GetEvents",
		"select",
	)
	defer metrics.UpdateDatabaseResponseTime(time.Now().UTC())

	ctx, cancel := context.WithTimeout(ctx, time.Duration(ks.lockTime)*time.Second)
	defer cancel()

	ks.logger.Debugf(ctx, "outboxstore.GetEvents: начинаю получение сообщений  для отправки из outbox")
	defer ks.logger.Debugf(ctx, "outboxstore.GetEvents: закончил получение сообщений  для отправки из outbox")

	dbWrPool := ks.data.GetWriterPool()
	q := sqlc.New(dbWrPool)

	getParams := sqlc.OutboxParams{
		StatusNew:    sqlc.OutboxStatusTypeNew,
		StatusLocked: sqlc.OutboxStatusTypeLocked,
		StatusFailed: sqlc.OutboxStatusTypeFailed,
		LockedTo:     toPGText(fmt.Sprintf("%d", ks.lockTime)),
	}

	msgs, err = q.Outbox(ctx, getParams)
	if err != nil {
		metrics.UpdateDatabaseResponseCode(
			repName,
			"GetEvents",
			"select",
			"error",
		)

		return nil, fmt.Errorf("Ошибка получения сообщений для отправки в брокер - %w", err)
	}

	metrics.UpdateDatabaseResponseCode(
		repName,
		"GetEvents",
		"select",
		"ok",
	)

	return msgs, nil
}
