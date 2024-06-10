package orderapi

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/model"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/pkg/api/order/v1"
)

// Info получение информации о заказе и его статусе.
func (a *API) Info(ctx context.Context, req *order.OrderInfoRequest) (*order.OrderInfoResponse, error) {
	ordrResult, err := a.orderService.Info(req.GetOrderId())
	if err != nil {
		a.logger.Debugf(ctx, "OrderApi.Info: Ошибка получения заказа - %s - %v", req.GetOrderId(), err)

		if errors.Is(err, model.ErrOrderNotFound) {
			return nil, status.Errorf(codes.NotFound, err.Error())
		}

		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &order.OrderInfoResponse{
		Order: orderToResponse(req.GetOrderId(), ordrResult),
	}, nil
}
