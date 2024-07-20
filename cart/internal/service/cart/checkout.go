package cart

import (
	"context"
	"fmt"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/model"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/pkg/tracer"
)

// Checkout приобретение товаров через сервис LOMS.
func (s *cService) Checkout(ctx context.Context, userID int64) (*model.OrderCart, error) {
	ctx, span := tracer.StartSpanFromContext(ctx, "service.cart.checkout")
	span.SetTag("component", "cart")
	defer span.End()

	s.logger.Info(ctx, fmt.Sprintf("cartService.Checkout: начинаю создание заказа пользователя: %d", userID))
	defer s.logger.Info(ctx, fmt.Sprintf("cartService.Checkout: закончил создание заказа пользователя: %d", userID))

	cart, err := s.cartStore.GetAllUserItems(ctx, userID)
	if err != nil {
		span.SetTag("error", true)
		s.logger.Errorf(ctx, "cartService.Checkout: ошибка получение содержимого корзины пользователя %d - %v", userID, err)

		return nil, fmt.Errorf("Ошибка получение содержимого корзины пользователя %d - %w", userID, err)
	}

	ord, err := s.lomsCli.AddOrder(ctx, userID, cart)
	if err != nil {
		span.SetTag("error", true)
		s.logger.Errorf(ctx, "cartService.Checkout: Ошибка создания заказа пользователя %v - %v", userID, err)

		return nil, fmt.Errorf("Ошибка создания заказа пользователя %v - %w", userID, err)
	}

	err = s.cartStore.DelCart(ctx, userID)
	if err != nil {
		span.SetTag("error", true)
		s.logger.Errorf(ctx, fmt.Sprintf("Ошибка удаление корзины пользователя %v - %v", userID, err))

		return nil, fmt.Errorf("Ошибка удаление корзины пользователя %v - %w", userID, err)
	}

	return &model.OrderCart{OrderID: ord}, nil
}
