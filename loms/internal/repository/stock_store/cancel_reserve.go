package stockstore

import (
	"fmt"

	"go.uber.org/zap"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/model"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/pkg/logger"
)

// CancelReserve возвращает товару доступность для резервирования.
func (ms *reserveStorage) CancelReserve(items model.AllNeedReserve) error {
	ctx := logger.Append(ms.ctx, []zap.Field{zap.Any("items", items)})

	ms.logger.Debugf(ctx, "stockStore.CancelReserve: начинаю снятие резерва товаров")
	defer ms.logger.Debugf(ctx, "stockStore.CancelReserve: закончил снятие резерва  товаров")

	ms.Lock()
	defer ms.Unlock()

	for _, item := range items {
		itmReserved, ok := ms.data[item.Sku]
		if !ok {
			return fmt.Errorf("Товар %d не найден", item.Sku)
		}

		if itmReserved.Reserved < item.Count {
			return fmt.Errorf("Количество зарезервированного товара %d неподходит для возврата", item.Sku)
		}
	}

	for _, item := range items {
		ms.data[item.Sku].Reserved -= item.Count
	}

	return nil
}
