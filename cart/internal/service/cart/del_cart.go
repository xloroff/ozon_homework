package cart

import (
	"context"
	"fmt"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/pkg/tracer"
)

// DelCart удалет корзину пользователю через сервис - обращение к сервису хранения.
func (s *cService) DelCart(ctx context.Context, userID int64) error {
	ctx, span := tracer.StartSpanFromContext(ctx, "service.cart.del_cart")
	span.SetTag("component", "cart")
	defer span.End()

	s.logger.Info(ctx, fmt.Sprintf("cartService.DelCart: начинаю удаление корзины пользователя userId - %d", userID))
	defer s.logger.Info(ctx, fmt.Sprintf("cartService.DelCart: закончил удаление корзины пользователя userId - %d", userID))

	err := s.cartStore.DelCart(ctx, userID)
	if err != nil {
		span.SetTag("error", true)
		s.logger.Errorf(ctx, "cartService.DelCart: ошибка удаление корзины пользователя %v - %v", userID, err)

		return fmt.Errorf("Ошибка удаление корзины пользователя %v - %w", userID, err)
	}

	return nil
}
