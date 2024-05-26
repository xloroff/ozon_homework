package memorystore

import (
	"context"
	"fmt"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/internal/model/v1"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/internal/pkg/logger"
)

// AddItem добавляет товар в память.
func (ms *cartStorage) AddItem(ctx context.Context, item *v1.AddItem) error {
	logger.Info(ctx, fmt.Sprintf("repositoryMemory.AddItem: начинаю добавление продукта userId: %d, skuID: %d, count: %v", item.UserID, item.SkuID, item.Count))
	defer logger.Info(ctx, fmt.Sprintf("repositoryMemory.AddItem: закончил добавление продукта userId: %d, skuID: %d, count: %v", item.UserID, item.SkuID, item.Count))

	cart, ok := ms.data[item.UserID]
	if !ok {
		cart = &v1.Cart{
			Items: v1.CartItems{},
		}
		ms.data[item.UserID] = cart
	}

	cartItem, ok := cart.Items[item.SkuID]
	if !ok {
		cartItem = &v1.CartItem{}
		cart.Items[item.SkuID] = cartItem
	}

	cartItem.Count += item.Count

	return nil
}