package orderapi

import (
	"context"
	"fmt"

	"google.golang.org/grpc"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/pkg/logger"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/service/order"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/pkg/api/order/v1"
)

// API структура апи заказов.
type API struct {
	logger logger.ILog
	ctx    context.Context
	order.UnimplementedOrderAPIServer
	orderService orderservice.Service
}

// NewAPI создает новое API заказов.
func NewAPI(ctx context.Context, l logger.ILog, orderService orderservice.Service) *API {
	return &API{
		logger:       l,
		ctx:          ctx,
		orderService: orderService,
	}
}

// RegisterGrpcServer регистрирует на сервере приклады.
func (a *API) RegisterGrpcServer(server *grpc.Server) {
	order.RegisterOrderAPIServer(server, a)
}

// RegisterHTTPHandler регистрирует приклады для HTTP.
func (a *API) RegisterHTTPHandler(mux *runtime.ServeMux, conn *grpc.ClientConn) error {
	err := order.RegisterOrderAPIHandler(a.ctx, mux, conn)
	if err != nil {
		a.logger.Errorf(a.ctx, "OrderApiRegisterHttpHandler: Ошибка создания хэндлера - %v", err)

		return fmt.Errorf("Ошибка создания хэндлера - %w", err)
	}

	return nil
}
