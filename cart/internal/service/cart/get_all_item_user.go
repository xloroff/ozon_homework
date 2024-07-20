package cart

import (
	"context"
	"fmt"
	"sort"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/model"
)

// GetAllUserItems получение корзины пользователя через сервис, вызов клиента для связи с сервисом продуктов и обращение к сервису хранения.
func (s *cService) GetAllUserItems(ctx context.Context, userID int64) (*model.FullUserCart, error) {
	s.logger.Info(ctx, fmt.Sprintf("cartService.GetAllUserItems: начинаю получение корзины userId: %d", userID))
	defer s.logger.Info(ctx, fmt.Sprintf("cartService.GetAllUserItems: закончил получение корзины userId: %d", userID))

	userItems, err := s.cartStore.GetAllUserItems(ctx, userID)
	if err != nil {
		s.logger.Errorf(ctx, "cartService.GetAllUserItems: ошибка получение содержимого корзины пользователя %d - %v", userID, err)
		return nil, fmt.Errorf("Ошибка получение содержимого корзины пользователя %d - %w", userID, err)
	}

	fullUsrCart, err := s.fullCartReciver(ctx, userItems)
	if err != nil {
		s.logger.Errorf(ctx, "cartService.GetAllUserItems: ошибка при формировании содержимого корзины и общей стоимости товаров пользователя %d - %v", userID, err)
		return nil, fmt.Errorf("Ошибка при формировании содержимого корзины и общей стоимости товаров пользователя %d- %w", userID, err)
	}

	return fullUsrCart, nil
}

// fullCartReciver идет в сервис продуктов по каждой позиции и формирует результирующую корзину.
func (s *cService) fullCartReciver(ctx context.Context, cart *model.Cart) (*model.FullUserCart, error) {
	var result *model.FullUserCart = &model.FullUserCart{}
	result.Items = make([]*model.UserCartItem, 0, len(cart.Items))

	for skuID, cI := range cart.Items {
		getProduct, err := s.productCli.GetProduct(ctx, skuID)
		if err != nil {
			s.logger.Errorf(ctx, "cartService.fullCartReciver: ошибка получения продукта %v - %v", skuID, err)
			return nil, fmt.Errorf("Ошибка получения продукта  %v - %w", skuID, err)
		}

		result.Items = append(result.Items, &model.UserCartItem{
			SkuID: skuID,
			Name:  getProduct.Name,
			Price: getProduct.Price,
			Count: cI.Count,
		})
		result.TotalPrice += uint32(cI.Count) * getProduct.Price
	}

	// Нужно сортирнуть результат по ТЗ
	sort.SliceStable(result.Items, func(i, j int) bool {
		return result.Items[i].SkuID < result.Items[j].SkuID
	})

	return result, nil
}
