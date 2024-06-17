package stockstore

import (
	"fmt"

	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/model"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/pkg/logger"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/repository/stock_store/sqlc"
)

// AddReserve добавлеяет резервирование товара.
func (ms *reserveStorage) AddReserve(items model.AllNeedReserve) error {
	ctx := logger.Append(ms.ctx, []zap.Field{zap.Any("items", items)})

	ms.logger.Debugf(ctx, "stockStore.AddReserve: начинаю резервирование товаров")
	defer ms.logger.Debugf(ctx, "stockStore.AddReserve: закончил резервирование товаров")

	dbWrPool := ms.data.GetWriterPool()

	err := dbWrPool.BeginFunc(ms.ctx, func(tx pgx.Tx) error {
		q := sqlc.New(dbWrPool).WithTx(tx)

		for _, item := range items {
			stock, err := q.GetAvailableForReserve(ms.ctx, item.Sku)
			if err != nil {
				return fmt.Errorf("Ошибка получения остатков товара %v - %w", item.Sku, err)
			}

			if stock.TotalCount-stock.Reserved < int32(item.Count) {
				return fmt.Errorf("Недостаточное количество товара %d в остатках", item.Sku)
			}

			err = q.AddReserve(ms.ctx, sqlc.AddReserveParams{Sku: item.Sku, Reserved: int32(item.Count)})
			if err != nil {
				return fmt.Errorf("Ошибка резервирования товара %d - %w", item.Sku, err)
			}
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("Ошибка резервирования товара - %w", err)
	}

	return nil
}
