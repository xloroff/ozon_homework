package stockstore

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/model"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/pkg/logger"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/pkg/metrics"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/pkg/tracer"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/repository/stock_store/sqlc"
)

// AddReserve добавлеяет резервирование товара.
func (ms *reserveStorage) AddReserve(ctx context.Context, items model.AllNeedReserve) error {
	ctx, span := tracer.StartSpanFromContext(ctx, "repository.stockstore.add_reserve")
	span.SetTag("component", "stockstore")
	span.SetTag("db.type", "sql")
	span.SetTag("db.statement", "insert")

	defer span.End()

	metrics.UpdateDatabaseRequestsTotal(
		repName,
		"AddReserve",
		"insert",
	)
	defer metrics.UpdateDatabaseResponseTime(time.Now().UTC())

	ctx = logger.AddFieldsToContext(ctx, "data", items)

	ms.logger.Debugf(ctx, "stockStore.AddReserve: начинаю резервирование товаров")
	defer ms.logger.Debugf(ctx, "stockStore.AddReserve: закончил резервирование товаров")

	dbWrPool := ms.data.GetWriterPool()

	err := dbWrPool.BeginFunc(ctx, func(tx pgx.Tx) error {
		q := sqlc.New(dbWrPool).WithTx(tx)

		for _, item := range items {
			stock, err := q.GetAvailableForReserve(ctx, item.Sku)
			if err != nil {
				span.SetTag("error", true)
				return fmt.Errorf("Ошибка получения остатков товара %v - %w", item.Sku, err)
			}

			if stock.TotalCount-stock.Reserved < int32(item.Count) {
				span.SetTag("error", true)
				return fmt.Errorf("Недостаточное количество товара %d в остатках", item.Sku)
			}

			err = q.AddReserve(ctx, sqlc.AddReserveParams{Sku: item.Sku, Reserved: int32(item.Count)})
			if err != nil {
				span.SetTag("error", true)
				return fmt.Errorf("Ошибка резервирования товара %d - %w", item.Sku, err)
			}
		}

		return nil
	})
	if err != nil {
		span.SetTag("error", true)
		metrics.UpdateDatabaseResponseCode(
			repName,
			"AddReserve",
			"insert",
			"error",
		)

		return fmt.Errorf("Ошибка резервирования товара - %w", err)
	}

	metrics.UpdateDatabaseResponseCode(
		repName,
		"AddReserve",
		"insert",
		"ok",
	)

	return nil
}
