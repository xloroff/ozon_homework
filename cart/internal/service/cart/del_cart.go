package cart

import (
	"context"
	"fmt"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/internal/pkg/logger"
)

// DelCart удалет корзину пользователю через сервис - обращение к сервису хранения.
func (s *cService) DelCart(ctx context.Context, userID int64) error {
	logger.Info(ctx, fmt.Sprintf("cartService.DelCart: начинаю удаление корзины пользователя userId - %d", userID))
	defer logger.Info(ctx, fmt.Sprintf("cartService.DelCart: закончил удаление корзины пользователя userId - %d", userID))

	err := s.cartStore.DelCart(ctx, userID)
	if err != nil {
		logger.Errorf(ctx, "cartService.DelCart: ошибка удаление корзины пользователя %v - %w", userID, err)
		return fmt.Errorf("Ошибка удаление корзины пользователя %v - %w", userID, err)
	}

	return nil
}