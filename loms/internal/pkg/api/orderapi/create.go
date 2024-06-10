package orderapi

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/pkg/api/order/v1"
)

// Create создание заказа.
func (a *API) Create(ctx context.Context, req *order.OrderCreateRequest) (*order.OrderCreateResponse, error) {
	orderID, err := a.orderService.Create(req.GetUser(), reqItemstoItems(req.GetItems()))
	if err != nil {
		a.logger.Debugf(ctx, "OrderApi.Create: Ошибка создания заказа - %v", err)

		return nil, status.Errorf(codes.FailedPrecondition, err.Error())
	}

	return &order.OrderCreateResponse{
		OrderId: orderID,
	}, nil
}
