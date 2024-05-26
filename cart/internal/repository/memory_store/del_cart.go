package memorystore

import (
	"context"
	"fmt"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/internal/pkg/logger"
)

// DelCart удаляет корзину пользователя из памяти из памяти.
func (ms *cartStorage) DelCart(ctx context.Context, userId int64) error {
	logger.Info(ctx, fmt.Sprintf("repositoryMemory.DelCart: начинаю удаление корзины пользователя userId - %d", userId))
	defer logger.Info(ctx, fmt.Sprintf("repositoryMemory.DelCart: закончил удаление корзины пользователя userId - %d", userId))

	_, ok := ms.data[userId]
	if ok {
		delete(ms.data, userId)
	}

	return nil
}