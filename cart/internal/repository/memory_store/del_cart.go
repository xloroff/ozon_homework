package memorystore

import (
	"context"
	"fmt"
	"time"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/pkg/metrics"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/pkg/tracer"
)

// DelCart удаляет корзину пользователя из памяти из памяти.
func (ms *cartStorage) DelCart(ctx context.Context, userID int64) error {
	_, span := tracer.StartSpanFromContext(ctx, "repository.memorystore.del_cart")
	span.SetTag("component", "memorystore")
	defer span.End()

	metrics.UpdateDatabaseRequestsTotal(
		storeName,
		"DelCart",
	)
	defer metrics.UpdateDatabaseResponseTime(time.Now().UTC())

	ms.logger.Info(ctx, fmt.Sprintf("repositoryMemory.DelCart: начинаю удаление корзины пользователя userId - %d", userID))
	defer ms.logger.Info(ctx, fmt.Sprintf("repositoryMemory.DelCart: закончил удаление корзины пользователя userId - %d", userID))

	ms.Lock()
	defer ms.Unlock()

	_, ok := ms.data[userID]
	if ok {
		delete(ms.data, userID)

		metrics.UpdateDatabaseResponseCode(
			storeName,
			"DelCart",
			"delete",
			"ok",
		)
	} else {
		metrics.UpdateDatabaseResponseCode(
			storeName,
			"DelCart",
			"delete",
			"not_found",
		)
	}

	metrics.UpdateInMemoryItemCount(len(ms.data))

	return nil
}
