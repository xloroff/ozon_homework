package stockstore

import (
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/model"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/repository/stock_store/sqlc"
)

// GetAvailableForReserve возвращает число единиц товара доступных для резервирования.
func (ms *reserveStorage) GetAvailableForReserve(sku int64) (uint16, error) {
	ms.logger.Debugf(ms.ctx, fmt.Sprintf("stockStore.GetAvailableForReserve: начинаю получение резерва товара - %v", sku))
	defer ms.logger.Debugf(ms.ctx, fmt.Sprintf("stockStore.GetAvailableForReserve: закончил получение резерва товара - %v", sku))

	dbRePool := ms.data.GetReaderPool()
	q := sqlc.New(dbRePool)

	stock, err := q.GetAvailableForReserve(ms.ctx, sku)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, model.ErrReserveNotFound
		}

		return 0, fmt.Errorf("Ошибка получения резерва товара %v - %w", sku, err)
	}

	return uint16(stock.TotalCount - stock.Reserved), nil
}
