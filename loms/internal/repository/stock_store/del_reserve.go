package stockstore

import (
	"fmt"

	"go.uber.org/zap"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/model"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/pkg/logger"
)

// DelItemFromReserve изменяет остатки и доступное число по окончанию.
func (ms *reserveStorage) DelItemFromReserve(items model.AllNeedReserve) error {
	ctx := logger.Append(ms.ctx, []zap.Field{zap.Any("items", items)})

	ms.logger.Debugf(ctx, "stockStore.DelItemFromReserve: начинаю помечать товаровы как проданые")
	defer ms.logger.Debugf(ctx, "stockStore.DelItemFromReserve: закончил помечать товаровы как проданые")

	ms.Lock()
	defer ms.Unlock()

	for _, item := range items {
		itmReserved, ok := ms.data[item.Sku]
		if !ok {
			return fmt.Errorf("Товар %d не найден", item.Sku)
		}

		if itmReserved.TotalCount < item.Count {
			return fmt.Errorf("Недостаточное количество товара %d в остатках", item.Sku)
		}

		if itmReserved.Reserved < item.Count {
			return fmt.Errorf("Количество зарезервированного товара %d неподходит для оформления", item.Sku)
		}
	}

	for _, item := range items {
		ms.data[item.Sku].TotalCount -= item.Count
		ms.data[item.Sku].Reserved -= item.Count
	}

	return nil
}
