package memorystore

import (
	"context"
	"fmt"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/model"
)

// AddItem добавляет товар в память.
func (ms *cartStorage) AddItem(ctx context.Context, item *model.AddItem) error {
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

	return nil
}
