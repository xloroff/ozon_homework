package memorystore

import (
	"context"
	"fmt"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/model"
)

// GetAllUserItems получает все данные по пользователю.
func (ms *cartStorage) GetAllUserItems(ctx context.Context, userID int64) (*model.Cart, error) {
	ms.logger.Info(ctx, fmt.Sprintf("repositoryMemory.GetAllUserItems: начинаю получение корзины пользователя userId - %d", userID))
	defer ms.logger.Info(ctx, fmt.Sprintf("repositoryMemory.GetAllUserItems: закончил получение корзины пользователя userI - %d", userID))

	ms.RLock()
	defer ms.RUnlock()

	cart, ok := ms.data[userID]
	if !ok {
		ms.logger.Debugf(ctx, fmt.Sprintf("repositoryMemory.GetAllUserItems: Корзины для пользователя %d не найдено", userID))
		return nil, fmt.Errorf("Корзины для пользователя %d не найдено", userID)
	}

	return cart, nil
}
