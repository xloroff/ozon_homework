package cart

import (
	"context"
	"fmt"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/internal/model"
)

// AddItem добавление итема пользователю через сервис, вызов клиента для связи с сервисом продуктов и обращение к сервису хранения.
func (s *cService) AddItem(ctx context.Context, item *model.AddItem) error {
	s.logger.Info(ctx, fmt.Sprintf("cartService.AddItem: начинаю добавление продукта userId: %d, skuID: %d, count: %v", item.UserID, item.SkuID, item.Count))
	defer s.logger.Info(ctx, fmt.Sprintf("cartService.AddItem: закончил добавление продукта userId: %d, skuID: %d, count: %v", item.UserID, item.SkuID, item.Count))

	_, err := s.productCli.GetProduct(ctx, item.SkuID)
	if err != nil {
		s.logger.Errorf(ctx, "cartService.AddItem: ошибка получения продукта %v - %v", item.SkuID, err)
		return fmt.Errorf("Ошибка получения продукта %v - %w", item.SkuID, err)
	}

	err = s.cartStore.AddItem(ctx, item)
	if err != nil {
		s.logger.Errorf(ctx, "cartService.AddItem: ошибка добавления продукта %v - %v", item.SkuID, err)
		return fmt.Errorf("Ошибка добавления продукта %v - %w", item.SkuID, err)
	}

	return nil
}
