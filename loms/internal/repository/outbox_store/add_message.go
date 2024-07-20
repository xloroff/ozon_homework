package outboxstore

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/model"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/pkg/metrics"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/pkg/tracer"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/repository/outbox_store/sqlc"
)

// AddMessage добавляет сообщение для отправки в очередь по паттерну outbox.
func (ks *outboxStorage) AddMessage(ctx context.Context, tx pgx.Tx, message *model.Outbox) error {
	ctx, span := tracer.StartSpanFromContext(ctx, "repository.outboxstore.add_message")
	span.SetTag("component", "outboxstore")
	span.SetTag("db.type", "sql")
	span.SetTag("db.statement", "insert")

	defer span.End()

	metrics.UpdateDatabaseRequestsTotal(
		repName,
		"AddMessage",
		"insert",
	)
	defer metrics.UpdateDatabaseResponseTime(time.Now().UTC())

	traceID := tracer.GetTraceID(ctx)
	spanID := tracer.GetSpanID(ctx)

	metadata := model.Metadata{TraceID: traceID, SpanID: spanID}

	rawMetadata, err := json.Marshal(metadata)
	if err != nil {
		return fmt.Errorf("Ошибка формирования метаданных о сообщении для сохранения в БД - %w", err)
	}

	dbWrPool := ks.data.GetWriterPool()
	q := sqlc.New(dbWrPool).WithTx(tx)

	_, err = q.AddOutbox(ctx,
		sqlc.AddOutboxParams{
			EntityID: message.EntityID,
			Payload:  message.Payload,
			Metadata: rawMetadata,
		})
	if err != nil {
		span.SetTag("error", true)
		metrics.UpdateDatabaseResponseCode(
			repName,
			"AddMessage",
			"insert",
			"error",
		)

		return fmt.Errorf("Ошибка добавления сообщения в БД - %w", err)
	}

	metrics.UpdateDatabaseResponseCode(
		repName,
		"AddMessage",
		"insert",
		"ok",
	)

	return nil
}
