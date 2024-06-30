package cart

import (
	"context"
	"fmt"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/model"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/pkg/tracer"
)

// DelItem удаление итема пользователя через сервис - обращение к сервису хранения.
func (s *cService) DelItem(ctx context.Context, item *model.DelItem) error {
	ctx, span := tracer.StartSpanFromContext(ctx, "service.cart.del_item")
	span.SetTag("component", "cart")
	defer span.End()

	s.logger.Info(ctx, fmt.Sprintf("cartService.DelItem: начинаю удаление продукта userId: %d, skuID: %d", item.UserID, item.SkuID))
	defer s.logger.Info(ctx, fmt.Sprintf("cartService.DelItem: закончил удаление продукта userId: %d, skuID: %d", item.UserID, item.SkuID))

	if !s.cartStore.DelItem(ctx, item) {
		span.SetTag("error", true)
		s.logger.Errorf(ctx, "cartService.DelItem: ошибка удаления продукта %v - %v", item.SkuID, model.ErrUnknownError)

		return fmt.Errorf("Ошибка удаления продукта  %v - %w", item.SkuID, model.ErrUnknownError)
	}

	return nil
}
