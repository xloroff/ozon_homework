package cart

import (
	"context"
	"fmt"
	"sort"
	"sync"

	"golang.org/x/time/rate"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/config"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/model"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/pkg/errgroup"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/pkg/tracer"
)

// GetAllUserItems получение корзины пользователя через сервис, вызов клиента для связи с сервисом продуктов и обращение к сервису хранения.
func (s *service) GetAllUserItems(ctx context.Context, userID int64) (*model.FullUserCart, error) {
	ctx, span := tracer.StartSpanFromContext(ctx, "get_all_user_items")
	span.SetTag("component", "cart")

	defer span.End()

	s.logger.Info(ctx, fmt.Sprintf("cartService.GetAllUserItems: начинаю получение корзины userId: %d", userID))
	defer s.logger.Info(ctx, fmt.Sprintf("cartService.GetAllUserItems: закончил получение корзины userId: %d", userID))

	userItems, err := s.cartStore.GetAllUserItems(ctx, userID)
	if err != nil {
		span.SetTag("error", true)
		s.logger.Errorf(ctx, "cartService.GetAllUserItems: ошибка получение содержимого корзины пользователя %d - %v", userID, err)

		return nil, fmt.Errorf("Ошибка получение содержимого корзины пользователя %d - %w", userID, err)
	}

	fullUsrCart, err := s.fullCartReciver(ctx, userItems)
	if err != nil {
		span.SetTag("error", true)
		s.logger.Errorf(ctx, "cartService.GetAllUserItems: ошибка при формировании содержимого корзины и общей стоимости товаров пользователя %d - %v", userID, err)

		return nil, fmt.Errorf("Ошибка при формировании содержимого корзины и общей стоимости товаров пользователя %d- %w", userID, err)
	}

	return fullUsrCart, nil
}

// fullCartReciver идет в сервис продуктов по каждой позиции и формирует результирующую корзину.
func (s *service) fullCartReciver(ctx context.Context, cart *model.Cart) (*model.FullUserCart, error) {
	ctx, span := tracer.StartSpanFromContext(ctx, "service.cart.full_cart_reciver")
	span.SetTag("component", "cart")
	span.SetTag("span.kind", "child")

	defer span.End()

	var result *model.FullUserCart = &model.FullUserCart{}
	result.Items = make([]*model.UserCartItem, 0, len(cart.Items))

	mu := sync.Mutex{}
	erg, ctx := errgroup.NewErrGroup(
		ctx,
		errgroup.WithCancelFirst(),
	)

	limiter := rate.NewLimiter(rate.Limit(config.RPS), 1)

	for skuID, cI := range cart.Items {
		erg.Go(func() error {
			if err := limiter.Wait(ctx); err != nil {
				span.SetTag("error", true)
				return fmt.Errorf("Ошибка ожидания при срабатывании лимита RPS %v - %w", config.RPS, err)
			}

			getProduct, err := s.productCli.GetProduct(ctx, skuID)
			if err != nil {
				span.SetTag("error", true)
				s.logger.Errorf(ctx, "cartService.fullCartReciver: ошибка получения продукта %v - %v", skuID, err)

				return fmt.Errorf("Ошибка получения продукта  %v - %w", skuID, err)
			}

			mu.Lock()
			defer mu.Unlock()

			result.Items = append(result.Items, &model.UserCartItem{
				SkuID: skuID,
				Name:  getProduct.Name,
				Price: getProduct.Price,
				Count: cI.Count,
			})
			result.TotalPrice += uint32(cI.Count) * getProduct.Price

			return nil
		})
	}

	if err := erg.Wait(); err != nil {
		span.SetTag("error", true)
		return nil, fmt.Errorf("Ошибка при формировании результирующей корзины  -  %v", erg.ErrsToString())
	}

	// Нужно сортирнуть результат по ТЗ
	sort.SliceStable(result.Items, func(i, j int) bool {
		return result.Items[i].SkuID < result.Items[j].SkuID
	})

	return result, nil
}
