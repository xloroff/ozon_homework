package stockapi

import (
	"context"
	"fmt"

	"google.golang.org/grpc"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/pkg/logger"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/service/stock"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/pkg/api/stock/v1"
)

// API структура апи резервов.
type API struct {
	logger logger.ILog
	ctx    context.Context
	stock.UnimplementedStockAPIServer
	stockService stockservice.Service
}

// NewAPI создает новое API резервов.
func NewAPI(ctx context.Context, l logger.ILog, stockService stockservice.Service) *API {
	return &API{
		ctx:          ctx,
		logger:       l,
		stockService: stockService,
	}
}

// RegisterGrpcServer регистрирует на сервере приклады.
func (a *API) RegisterGrpcServer(server *grpc.Server) {
	stock.RegisterStockAPIServer(server, a)
}

// RegisterHTTPHandler регистрирует приклады для HTTP.
func (a *API) RegisterHTTPHandler(mux *runtime.ServeMux, conn *grpc.ClientConn) error {
	err := stock.RegisterStockAPIHandler(a.ctx, mux, conn)
	if err != nil {
		a.logger.Errorf(a.ctx, "StockApiRegisterHttpHandler: Ошибка создания хэндлера - %v", err)

		return fmt.Errorf("Ошибка создания хэндлера - %w", err)
	}

	return nil
}
