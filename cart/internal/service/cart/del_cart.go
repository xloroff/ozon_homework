package cart

import (
	"context"
	"fmt"
)

// DelCart удалет корзину пользователю через сервис - обращение к сервису хранения.
func (s *cService) DelCart(ctx context.Context, userID int64) error {
	s.logger.Info(ctx, fmt.Sprintf("cartService.DelCart: начинаю удаление корзины пользователя userId - %d", userID))
	defer s.logger.Info(ctx, fmt.Sprintf("cartService.DelCart: закончил удаление корзины пользователя userId - %d", userID))

	err := s.cartStore.DelCart(ctx, userID)
	if err != nil {
		s.logger.Errorf(ctx, "cartService.DelCart: ошибка удаление корзины пользователя %v - %v", userID, err)
		return fmt.Errorf("Ошибка удаление корзины пользователя %v - %w", userID, err)
	}

	return nil
}
