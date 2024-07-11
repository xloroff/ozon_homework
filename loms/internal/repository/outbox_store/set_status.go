package outboxstore

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/pkg/metrics"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/repository/outbox_store/sqlc"
)

func (ks *outboxStorage) SetStatus(ctx context.Context, msg *sqlc.OutboxRow) error {
	metrics.UpdateDatabaseRequestsTotal(
		repName,
		"SetStatus",
		"update",
	)
	defer metrics.UpdateDatabaseResponseTime(time.Now().UTC())

	ctx, cancel := context.WithTimeout(ctx, time.Duration(ks.lockTime)*time.Second)
	defer cancel()

	ks.logger.Debugf(ctx, "outboxstore.GetEvents: начинаю получение сообщений  для отправки из outbox")
	defer ks.logger.Debugf(ctx, "outboxstore.GetEvents: закончил получение сообщений  для отправки из outbox")

	dbWrPool := ks.data.GetWriterPool()

	err := dbWrPool.BeginFuncWithTx(ctx, func(tx pgx.Tx) error {
		q := sqlc.New(dbWrPool).WithTx(tx)

		var err error

		setParams := sqlc.SetStatusOutboxParams{
			IDMsg:     msg.ID,
			NewStatus: sqlc.OutboxStatusTypeSent,
		}

		err = q.SetStatusOutbox(ctx, setParams)
		if err != nil {
			return fmt.Errorf("Ошибка установки статуса \"%v\" сообщению ID - %v: %w", setParams.NewStatus, msg.ID, err)
		}

		return nil
	})
	if err != nil {
		metrics.UpdateDatabaseResponseCode(
			repName,
			"SetStatus",
			"update",
			"error",
		)

		return fmt.Errorf("Ошибка при исполненнии запроса обновления статуса сообщения - %w", err)
	}

	metrics.UpdateDatabaseResponseCode(
		repName,
		"SetStatus",
		"update",
		"ok",
	)

	return nil
}
