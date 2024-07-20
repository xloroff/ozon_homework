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

// CancelReserve возвращает товару доступность для резервирования.
func (ms *reserveStorage) CancelReserve(ctx context.Context, items model.AllNeedReserve) error {
	ctx, span := tracer.StartSpanFromContext(ctx, "repository.stockstore.cancel_reserve")
	span.SetTag("component", "stockstore")
	span.SetTag("db.type", "sql")
	span.SetTag("db.statement", "insert")

	defer span.End()

	metrics.UpdateDatabaseRequestsTotal(
		repName,
		"CancelReserve",
		"update",
	)
	defer metrics.UpdateDatabaseResponseTime(time.Now().UTC())

	ctx = logger.AddFieldsToContext(ctx, "data", items, "user_id")

	ms.logger.Debugf(ctx, "stockStore.CancelReserve: начинаю снятие резерва товаров")
	defer ms.logger.Debugf(ctx, "stockStore.CancelReserve: закончил снятие резерва  товаров")

	dbWrPool := ms.data.GetWriterPool()

	err := dbWrPool.BeginFunc(ctx, func(tx pgx.Tx) error {
		q := sqlc.New(dbWrPool).WithTx(tx)

		for _, item := range items {
			if err := ms.cancelItem(ctx, q, item); err != nil {
				span.SetTag("error", true)
				return err
			}
		}

		return nil
	})
	if err != nil {
		span.SetTag("error", true)
		metrics.UpdateDatabaseResponseCode(
			repName,
			"CancelReserve",
			"update",
			"error",
		)

		return fmt.Errorf("Ошибка отмены резерва товаров - %w", err)
	}

	metrics.UpdateDatabaseResponseCode(
		repName,
		"CancelReserve",
		"update",
		"ok",
	)

	return nil
}

func (ms *reserveStorage) cancelItem(ctx context.Context, q *sqlc.Queries, item *model.NeedReserve) error {
	ctx, span := tracer.StartSpanFromContext(ctx, "repository.stockstore.cancel_item")
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

	return ms.cReserve(ctx, q, item)
}

func (ms *reserveStorage) validateStock(item *model.NeedReserve, stock sqlc.Stock) error {
	if stock.TotalCount < int32(item.Count) {
		return fmt.Errorf("Количество зарезервированного товара %d которое вы пытаетесь вернуть больше общего количества %d", item.Sku, stock.TotalCount)
	}

	if stock.Reserved < int32(item.Count) {
		return fmt.Errorf("Количество зарезервированного товара %d неподходит для возврата", item.Sku)
	}

	return nil
}

func (ms *reserveStorage) cReserve(ctx context.Context, q *sqlc.Queries, item *model.NeedReserve) error {
	ctx, span := tracer.StartSpanFromContext(ctx, "repository.stockstore.c_reserve")
	span.SetTag("component", "stockstore")
	span.SetTag("span.kind", "child")

	defer span.End()

	err := q.CancelReserve(ctx, sqlc.CancelReserveParams{
		Sku:      item.Sku,
		Reserved: int32(item.Count),
	})
	if err != nil {
		span.SetTag("error", true)
		return fmt.Errorf("Ошибка отмены резерва товара %d - %w", item.Sku, err)
	}

	return nil
}
