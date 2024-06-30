package memorystore

import (
	"context"
	"fmt"
	"time"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/model"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/pkg/metrics"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/pkg/tracer"
)

// AddItem добавляет товар в память.
func (ms *cartStorage) AddItem(ctx context.Context, item *model.AddItem) error {
	_, span := tracer.StartSpanFromContext(ctx, "repository.memorystore.add_item")
	span.SetTag("component", "memorystore")
	defer span.End()

	metrics.UpdateDatabaseRequestsTotal(
		storeName,
		"AddItem",
	)
	defer metrics.UpdateDatabaseResponseTime(time.Now().UTC())

	ms.logger.Info(ctx, fmt.Sprintf("repositoryMemory.AddItem: начинаю добавление продукта userId: %d, skuID: %d, count: %v", item.UserID, item.SkuID, item.Count))
	defer ms.logger.Info(ctx, fmt.Sprintf("repositoryMemory.AddItem: закончил добавление продукта userId: %d, skuID: %d, count: %v", item.UserID, item.SkuID, item.Count))

	ms.Lock()
	defer ms.Unlock()

	cart, ok := ms.data[item.UserID]
	if !ok {
		cart = &model.Cart{
			Items: model.CartItems{},
		}
		ms.data[item.UserID] = cart
	}

	cartItem, ok := cart.Items[item.SkuID]
	if !ok {
		cartItem = &model.CartItem{}
		cart.Items[item.SkuID] = cartItem
	}

	cartItem.Count += item.Count

	metrics.UpdateDatabaseResponseCode(
		storeName,
		"AddItem",
		"insert",
		"ok",
	)

	metrics.UpdateInMemoryItemCount(len(ms.data))

	return nil
}
