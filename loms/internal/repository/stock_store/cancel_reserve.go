package stockstore

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/model"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/pkg/logger"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/repository/stock_store/sqlc"
)

// CancelReserve возвращает товару доступность для резервирования.
func (ms *reserveStorage) CancelReserve(items model.AllNeedReserve) error {
	ctx := logger.Append(ms.ctx, []zap.Field{zap.Any("items", items)})

	ms.logger.Debugf(ctx, "stockStore.CancelReserve: начинаю снятие резерва товаров")
	defer ms.logger.Debugf(ctx, "stockStore.CancelReserve: закончил снятие резерва  товаров")

	dbWrPool := ms.data.GetWriterPool()

	err := dbWrPool.BeginFunc(ms.ctx, func(tx pgx.Tx) error {
		q := sqlc.New(dbWrPool).WithTx(tx)

		for _, item := range items {
			if err := ms.cancelItem(ctx, q, item); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("Ошибка отмены резерва товаров - %w", err)
	}

	return nil
}

func (ms *reserveStorage) cancelItem(ctx context.Context, q *sqlc.Queries, item *model.NeedReserve) error {
	stock, err := q.GetAvailableForReserve(ctx, item.Sku)
	if err != nil {
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
	err := q.CancelReserve(ctx, sqlc.CancelReserveParams{
		Sku:      item.Sku,
		Reserved: int32(item.Count),
	})
	if err != nil {
		return fmt.Errorf("Ошибка отмены резерва товара %d - %w", item.Sku, err)
	}

	return nil
}
