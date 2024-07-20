package memorystore

import (
	"context"
	"fmt"
	"time"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/model"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/pkg/metrics"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/pkg/tracer"
)

// DelItem удаляет товар из памяти.
func (ms *cartStorage) DelItem(ctx context.Context, item *model.DelItem) bool {
	_, span := tracer.StartSpanFromContext(ctx, "repository.memorystore.del_item")
	span.SetTag("component", "memorystore")
	defer span.End()

	metrics.UpdateDatabaseRequestsTotal(
		storeName,
		"DelItem",
	)
	defer metrics.UpdateDatabaseResponseTime(time.Now().UTC())

	ms.logger.Info(ctx, fmt.Sprintf("repositoryMemory.DelItem: начинаю удаление продукта userId: %d, skuID: %d", item.UserID, item.SkuID))
	defer ms.logger.Info(ctx, fmt.Sprintf("repositoryMemory.DelItem: закончил удаление продукта userId: %d, skuID: %d", item.UserID, item.SkuID))

	ms.Lock()
	defer ms.Unlock()

	cart, ok := ms.data[item.UserID]
	if !ok {
		metrics.UpdateDatabaseResponseCode(
			storeName,
			"DelItem",
			"delete",
			"not_found",
		)

		return true
	}

	_, ok = cart.Items[item.SkuID]
	if !ok {
		return true
	}

	delete(cart.Items, item.SkuID)

	if len(cart.Items) == 0 {
		delete(ms.data, item.UserID)
	}

	metrics.UpdateDatabaseResponseCode(
		storeName,
		"DelItem",
		"delete",
		"ok",
	)

	metrics.UpdateInMemoryItemCount(len(ms.data))

	return true
}
