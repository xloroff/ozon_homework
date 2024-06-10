package cart

import (
	"context"
	"fmt"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/model"
)

// Checkout приобретение товаров через сервис LOMS.
func (s *cService) Checkout(ctx context.Context, userID int64) (*model.OrderCart, error) {
	s.logger.Info(ctx, fmt.Sprintf("cartService.Checkout: начинаю создание заказа пользователя: %d", userID))
	defer s.logger.Info(ctx, fmt.Sprintf("cartService.Checkout: закончил создание заказа пользователя: %d", userID))

	cart, err := s.cartStore.GetAllUserItems(ctx, userID)
	if err != nil {
		s.logger.Errorf(ctx, "cartService.Checkout: ошибка получение содержимого корзины пользователя %d - %v", userID, err)
		return nil, fmt.Errorf("Ошибка получение содержимого корзины пользователя %d - %w", userID, err)
	}

	ord, err := s.lomsCli.AddOrder(userID, cart)
	if err != nil {
		s.logger.Errorf(ctx, "cartService.Checkout: Ошибка создания заказа пользователя %v - %v", userID, err)
		return nil, fmt.Errorf("Ошибка создания заказа пользователя %v - %w", userID, err)
	}

	err = s.cartStore.DelCart(ctx, userID)
	if err != nil {
		s.logger.Errorf(ctx, fmt.Sprintf("Ошибка удаление корзины пользователя %v - %v", userID, err))
		return nil, fmt.Errorf("Ошибка удаление корзины пользователя %v - %w", userID, err)
	}

	return &model.OrderCart{OrderID: ord}, nil
}
