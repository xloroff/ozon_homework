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

// DelItemFromReserve изменяет остатки и доступное число по окончанию.
func (ms *reserveStorage) DelItemFromReserve(items model.AllNeedReserve) error {
	ctx := logger.Append(ms.ctx, []zap.Field{zap.Any("items", items)})

	ms.logger.Debugf(ctx, "stockStore.DelItemFromReserve: начинаю помечать товаровы как проданые")
	defer ms.logger.Debugf(ctx, "stockStore.DelItemFromReserve: закончил помечать товаровы как проданые")

	dbWrPool := ms.data.GetWriterPool()

	err := dbWrPool.BeginFuncWithTx(ms.ctx, func(tx pgx.Tx) error {
		q := sqlc.New(dbWrPool).WithTx(tx)

		for _, item := range items {
			if err := ms.delItem(ctx, q, item); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("Ошибка изменения резерва товаров и общего числа товаров - %w", err)
	}

	return nil
}

func (ms *reserveStorage) delItem(ctx context.Context, q *sqlc.Queries, item *model.NeedReserve) error {
	stock, err := q.GetAvailableForReserve(ctx, item.Sku)
	if err != nil {
		return fmt.Errorf("Ошибка получения остатков товара %v - %w", item.Sku, err)
	}

	if err := ms.validateStock(item, stock); err != nil {
		return err
	}

	return ms.delReserve(ctx, q, item)
}

func (ms *reserveStorage) delReserve(ctx context.Context, q *sqlc.Queries, item *model.NeedReserve) error {
	err := q.DelItemFromReserve(ctx, sqlc.DelItemFromReserveParams{
		Sku:        item.Sku,
		TotalCount: int32(item.Count),
	})
	if err != nil {
		return fmt.Errorf("Ошибка отмены резерва товара %d - %w", item.Sku, err)
	}

	return nil
}
