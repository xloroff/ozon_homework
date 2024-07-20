package orderapi

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/model"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/pkg/api/order/v1"
)

// Cancel отмена заказа и высвобождение резервов.
func (a *API) Cancel(ctx context.Context, req *order.OrderCancelRequest) (*emptypb.Empty, error) {
	err := a.orderService.Cancel(req.GetOrderId())
	if err != nil {
		a.logger.Debugf(ctx, "OrderApi.Cancel: Ошибка отмены заказа - %s - %v", req.GetOrderId(), err)

		if errors.Is(err, model.ErrOrderNotFound) || errors.Is(err, model.ErrOrderCancel) {
			return nil, status.Errorf(codes.NotFound, err.Error())
		}

		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &emptypb.Empty{}, nil
}
