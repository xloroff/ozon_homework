package orderapi

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/model"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/pkg/tracer"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/pkg/api/order/v1"
)

// Pay помечает заказ как оплаченный.
func (a *API) Pay(ctx context.Context, req *order.OrderPayRequest) (*emptypb.Empty, error) {
	ctx, span := tracer.StartSpanFromContext(ctx, "api.orderapi.pay")
	span.SetTag("component", "orderapi")

	defer span.End()

	err := a.orderService.Pay(ctx, req.GetOrderId())
	if err != nil {
		span.SetTag("error", true)
		a.logger.Debugf(ctx, "OrderApi.Pay: Ошибка смены статуса оплаты заказа - %s - %v", req.GetOrderId(), err)

		if errors.Is(err, model.ErrOrderNotFound) {
			return nil, status.Errorf(codes.NotFound, err.Error())
		}

		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &emptypb.Empty{}, nil
}
