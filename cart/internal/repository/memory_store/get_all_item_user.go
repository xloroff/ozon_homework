package memorystore

import (
	"context"
	"fmt"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/internal/model/v1"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/internal/pkg/logger"
)

// GetAllUserItems получает все данные по пользователю.
func (ms *cartStorage) GetAllUserItems(ctx context.Context, userID int64) (*v1.Cart, error) {
	logger.Info(ctx, fmt.Sprintf("repositoryMemory.GetAllUserItems: начинаю получение корзины пользователя userId - %d", userID))
	defer logger.Info(ctx, fmt.Sprintf("repositoryMemory.GetAllUserItems: закончил получение корзины пользователя userI - %d", userID))

	cart, ok := ms.data[userID]
	if !ok {
		logger.Debugf(ctx, fmt.Sprintf("repositoryMemory.GetAllUserItems: Корзины для пользователя %d не найдено", userID))
		return nil, fmt.Errorf("Корзины для пользователя %d не найдено.", userID)
	}

	return cart, nil
}