package stockstore

import (
	"fmt"

	"go.uber.org/zap"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/model"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/pkg/logger"
)

// AddReserve добавлеяет резервирование товара.
func (ms *reserveStorage) AddReserve(items model.AllNeedReserve) error {
	ctx := logger.Append(ms.ctx, []zap.Field{zap.Any("items", items)})

	ms.logger.Debugf(ctx, "stockStore.AddReserve: начинаю резервирование товаров")
	defer ms.logger.Debugf(ctx, "stockStore.AddReserve: закончил резервирование товаров")

	ms.Lock()
	defer ms.Unlock()

	for _, item := range items {
		resItm, ok := ms.data[item.Sku]
		if !ok {
			return model.ErrReserveNotFound
		}

		if (resItm.TotalCount - resItm.Reserved) < item.Count {
			return fmt.Errorf("Недостаточное количество товара %d в остатках", item.Sku)
		}
	}

	for _, i := range items {
		ms.data[i.Sku].Reserved += i.Count
	}

	return nil
}
