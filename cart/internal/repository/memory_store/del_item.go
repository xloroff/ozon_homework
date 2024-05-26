package memorystore

import (
	"context"
	"fmt"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/internal/model/v1"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/internal/pkg/logger"
)

// DelItem удаляет товар из памяти.
func (ms *cartStorage) DelItem(ctx context.Context, item *v1.DelItem) bool {
	logger.Info(ctx, fmt.Sprintf("repositoryMemory.DelItem: начинаю удаление продукта userId: %d, skuID: %d", item.UserID, item.SkuID))
	defer logger.Info(ctx, fmt.Sprintf("repositoryMemory.DelItem: закончил удаление продукта userId: %d, skuID: %d", item.UserID, item.SkuID))

	cart, ok := ms.data[item.UserID]
	if !ok {
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

	return true
}