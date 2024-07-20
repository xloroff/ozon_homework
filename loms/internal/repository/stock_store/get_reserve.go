package stockstore

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/model"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/pkg/metrics"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/pkg/tracer"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/repository/stock_store/sqlc"
)

// GetAvailableForReserve возвращает число единиц товара доступных для резервирования.
func (ms *reserveStorage) GetAvailableForReserve(ctx context.Context, sku int64) (uint16, error) {
	ctx, span := tracer.StartSpanFromContext(ctx, "repository.stockstore.get_available_for_reserve")
	span.SetTag("component", "stockstore")
	span.SetTag("db.type", "sql")
	span.SetTag("db.statement", "select")

	defer span.End()

	metrics.UpdateDatabaseRequestsTotal(
		repName,
		"GetAvailableForReserve",
		"select",
	)
	defer metrics.UpdateDatabaseResponseTime(time.Now().UTC())

	ms.logger.Debugf(ctx, fmt.Sprintf("stockStore.GetAvailableForReserve: начинаю получение резерва товара - %v", sku))
	defer ms.logger.Debugf(ctx, fmt.Sprintf("stockStore.GetAvailableForReserve: закончил получение резерва товара - %v", sku))

	dbRePool := ms.data.GetReaderPool()
	q := sqlc.New(dbRePool)

	stock, err := q.GetAvailableForReserve(ctx, sku)
	if err != nil {
		span.SetTag("error", true)
		metrics.UpdateDatabaseResponseCode(
			repName,
			"GetAvailableForReserve",
			"select",
			"error",
		)

		if errors.Is(err, pgx.ErrNoRows) {
			return 0, model.ErrReserveNotFound
		}

		return 0, fmt.Errorf("Ошибка получения резерва товара %v - %w", sku, err)
	}

	metrics.UpdateDatabaseResponseCode(
		repName,
		"GetAvailableForReserve",
		"select",
		"ok",
	)

	return uint16(stock.TotalCount - stock.Reserved), nil
}
