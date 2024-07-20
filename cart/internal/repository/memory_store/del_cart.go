package memorystore

import (
	"context"
	"fmt"
)

// DelCart удаляет корзину пользователя из памяти из памяти.
func (ms *cartStorage) DelCart(ctx context.Context, userID int64) error {
	ms.logger.Info(ctx, fmt.Sprintf("repositoryMemory.DelCart: начинаю удаление корзины пользователя userId - %d", userID))
	defer ms.logger.Info(ctx, fmt.Sprintf("repositoryMemory.DelCart: закончил удаление корзины пользователя userId - %d", userID))

	ms.Lock()
	defer ms.Unlock()

	_, ok := ms.data[userID]
	if ok {
		delete(ms.data, userID)
	}

	return nil
}
