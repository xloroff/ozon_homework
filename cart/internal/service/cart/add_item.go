package cart

import (
	"context"
	"fmt"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/internal/initilize"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/internal/model/v1"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/internal/pkg/logger"
)

// AddItem добавление итема пользователю через сервис, вызов клиента для связи с сервисом продуктов и обращение к сервису хранения.
func (s *cService) AddItem(ctx context.Context, settings *initilize.ConfigAPI, item *v1.AddItem) error {
	logger.Info(ctx, fmt.Sprintf("cartService.AddItem: начинаю добавление продукта userId: %d, skuID: %d, count: %v", item.UserID, item.SkuID, item.Count))
	defer logger.Info(ctx, fmt.Sprintf("cartService.AddItem: закончил добавление продукта userId: %d, skuID: %d, count: %v", item.UserID, item.SkuID, item.Count))

	_, err := s.productCli.GetProduct(ctx, settings, item.SkuID)
	if err != nil {
		logger.Errorf(ctx, "cartService.AddItem: ошибка получения продукта %v - %w", item.SkuID, err)
		return fmt.Errorf("Ошибка получения продукта %v - %w", item.SkuID, err)
	}

	err = s.cartStore.AddItem(ctx, item)
	if err != nil {
		logger.Errorf(ctx, "cartService.AddItem: ошибка добавления продукта %v - %w", item.SkuID, err)
		return fmt.Errorf("Ошибка добавления продукта %v - %w", item.SkuID, err)
	}

	return nil
}