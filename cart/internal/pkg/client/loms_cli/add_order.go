package lomscli

import (
	"context"
	"fmt"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/model"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/pb/api/order/v1"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/pkg/logger"
)

// AddOrder создает заказ на сервисе заказов.
func (c *lomsClient) AddOrder(ctx context.Context, userID int64, cart *model.Cart) (int64, error) {
	ctx = logger.AddFieldsToContext(ctx, "data", cart, "user_id", userID)
	c.logger.Debugf(ctx, "LomsCli.AddOrder: начал обращение в сервис LOMS, создание заказа пользователь - %v", userID)
	defer c.logger.Debugf(ctx, "LomsCli.AddOrder: закончил обращение в сервис LOMS, создание заказа пользователь - %v", userID)

	resp, err := c.order.Create(ctx, cartToOrder(userID, cart))
	if err != nil {
		return 0, fmt.Errorf("Ошибка создания заказа - %w", err)
	}

	return resp.GetOrderId(), nil
}

func cartToOrder(user int64, cart *model.Cart) *order.OrderCreateRequest {
	itms := make([]*order.OrderCreateRequest_Item, 0, len(cart.Items))

	i := &order.OrderCreateRequest{
		User:  user,
		Items: itms,
	}

	for sku, c := range cart.Items {
		i.Items = append(i.Items, convItem(sku, c))
	}

	return i
}

func convItem(sku int64, c *model.CartItem) *order.OrderCreateRequest_Item {
	return &order.OrderCreateRequest_Item{
		Sku:   sku,
		Count: uint64(c.Count),
	}
}
