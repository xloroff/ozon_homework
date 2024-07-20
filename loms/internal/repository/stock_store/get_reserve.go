package stockstore

import (
	"fmt"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/model"
)

// GetAvailableForReserve возвращает число единиц товара доступных для резервирования.
func (ms *reserveStorage) GetAvailableForReserve(sku int64) (uint16, error) {
	ms.logger.Debugf(ms.ctx, fmt.Sprintf("stockStore.GetAvailableForReserve: начинаю получение резерва товара - %v", sku))
	defer ms.logger.Debugf(ms.ctx, fmt.Sprintf("stockStore.GetAvailableForReserve: закончил получение резерва товара - %v", sku))

	ms.RLock()
	defer ms.RUnlock()

	item, ok := ms.data[sku]
	if !ok {
		return 0, model.ErrReserveNotFound
	}

	return item.TotalCount - item.Reserved, nil
}
