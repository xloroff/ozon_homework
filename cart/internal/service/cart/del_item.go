package cart

import (
	"context"
	"errors"
	"fmt"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/internal/model/v1"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/internal/pkg/logger"
)

// DelItem удаление итема пользователя через сервис - обращение к сервису хранения.
func (s *cService) DelItem(ctx context.Context, item *v1.DelItem) error {
	logger.Info(ctx, fmt.Sprintf("cartService.DelItem: начинаю удаление продукта userId: %d, skuID: %d", item.UserID, item.SkuID))
	defer logger.Info(ctx, fmt.Sprintf("cartService.DelItem: закончил удаление продукта userId: %d, skuID: %d", item.UserID, item.SkuID))

	if !s.cartStore.DelItem(ctx, item) {
		logger.Errorf(ctx, "cartService.DelItem: ошибка удаления продукта %v - %w", item.SkuID, errors.New("Неизвестнаяя ошибка"))
		return fmt.Errorf("Ошибка удаления продукта  %v - %w", item.SkuID, "Неизвестнаяя ошибка")
	}

	return nil
}