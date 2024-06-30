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

// DelItemFromReserve изменяет остатки и доступное число по окончанию.
func (ms *reserveStorage) DelItemFromReserve(ctx context.Context, items model.AllNeedReserve) error {
	ctx, span := tracer.StartSpanFromContext(ctx, "repository.stockstore.del_item_from_reserve")
	span.SetTag("component", "stockstore")
	span.SetTag("db.type", "sql")
	span.SetTag("db.statement", "update")

	defer span.End()

	metrics.UpdateDatabaseRequestsTotal(
		repName,
		"DelItemFromReserve",
		"update",
	)
	defer metrics.UpdateDatabaseResponseTime(time.Now().UTC())

	ctx = logger.AddFieldsToContext(ctx, "data", items)

	ms.logger.Debugf(ctx, "stockStore.DelItemFromReserve: начинаю помечать товаровы как проданые")
	defer ms.logger.Debugf(ctx, "stockStore.DelItemFromReserve: закончил помечать товаровы как проданые")

	dbWrPool := ms.data.GetWriterPool()

	err := dbWrPool.BeginFuncWithTx(ctx, func(tx pgx.Tx) error {
		q := sqlc.New(dbWrPool).WithTx(tx)

		for _, item := range items {
			if err := ms.delItem(ctx, q, item); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		span.SetTag("error", true)
		metrics.UpdateDatabaseResponseCode(
			repName,
			"DelItemFromReserve",
			"update",
			"error",
		)

		return fmt.Errorf("Ошибка изменения резерва товаров и общего числа товаров - %w", err)
	}

	metrics.UpdateDatabaseResponseCode(
		repName,
		"DelItemFromReserve",
		"update",
		"ok",
	)

	return nil
}

func (ms *reserveStorage) delItem(ctx context.Context, q *sqlc.Queries, item *model.NeedReserve) error {
	ctx, span := tracer.StartSpanFromContext(ctx, "repository.stockstore.del_item")
	span.SetTag("component", "stockstore")
	span.SetTag("span.kind", "child")

	defer span.End()

	stock, err := q.GetAvailableForReserve(ctx, item.Sku)
	if err != nil {
		span.SetTag("error", true)
		return fmt.Errorf("Ошибка получения остатков товара %v - %w", item.Sku, err)
	}

	if err := ms.validateStock(item, stock); err != nil {
		return err
	}

	return ms.delReserve(ctx, q, item)
}

func (ms *reserveStorage) delReserve(ctx context.Context, q *sqlc.Queries, item *model.NeedReserve) error {
	ctx, span := tracer.StartSpanFromContext(ctx, "repository.stockstore.del_reserve")
	span.SetTag("component", "stockstore")
	span.SetTag("span.kind", "child")

	defer span.End()

	err := q.DelItemFromReserve(ctx, sqlc.DelItemFromReserveParams{
		Sku:        item.Sku,
		TotalCount: int32(item.Count),
	})
	if err != nil {
		span.SetTag("error", true)
		return fmt.Errorf("Ошибка отмены резерва товара %d - %w", item.Sku, err)
	}

	return nil
}
