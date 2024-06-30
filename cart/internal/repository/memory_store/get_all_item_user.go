package memorystore

import (
	"context"
	"fmt"
	"time"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/model"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/pkg/metrics"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/pkg/tracer"
)

// GetAllUserItems получает все данные по пользователю.
func (ms *cartStorage) GetAllUserItems(ctx context.Context, userID int64) (*model.Cart, error) {
	_, span := tracer.StartSpanFromContext(ctx, "repository.memorystore.get_all_user_items")
	span.SetTag("component", "memorystore")
	defer span.End()

	metrics.UpdateDatabaseRequestsTotal(
		storeName,
		"GetAllUserItems",
	)
	defer metrics.UpdateDatabaseResponseTime(time.Now().UTC())

	ms.logger.Info(ctx, fmt.Sprintf("repositoryMemory.GetAllUserItems: начинаю получение корзины пользователя userId - %d", userID))
	defer ms.logger.Info(ctx, fmt.Sprintf("repositoryMemory.GetAllUserItems: закончил получение корзины пользователя userI - %d", userID))

	ms.RLock()
	defer ms.RUnlock()

	cart, ok := ms.data[userID]
	if !ok {
		metrics.UpdateDatabaseResponseCode(
			storeName,
			"GetAllUserItems",
			"select",
			"not_found",
		)

		ms.logger.Debugf(ctx, fmt.Sprintf("repositoryMemory.GetAllUserItems: Корзины для пользователя %d не найдено", userID))

		return nil, fmt.Errorf("Корзины для пользователя %d не найдено", userID)
	}

	metrics.UpdateDatabaseResponseCode(
		storeName,
		"GetAllUserItems",
		"select",
		"ok",
	)

	return cart, nil
}
