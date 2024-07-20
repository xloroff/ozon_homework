package lomscli

import (
	"context"
	"fmt"
	"net/url"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/config"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/model"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/pb/api/order/v1"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/pb/api/stock/v1"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/pkg/client/loms_cli/interceptor/interceptor"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/pkg/logger"
)

// LomsService интерфейс для сервиса Loms.
type LomsService interface {
	AddOrder(ctx context.Context, userID int64, cart *model.Cart) (int64, error)
	StockInfo(ctx context.Context, skuID int64) (uint16, error)
}

// lomsClient структура клиента для сервиса loms.
type lomsClient struct {
	ctx    context.Context
	logger logger.Logger
	order  order.OrderAPIClient
	stock  stock.StockAPIClient
}

// NewClient создает нового клиента для сервиса LOMS.
func NewClient(ctx context.Context, l logger.Logger, conn *grpc.ClientConn) LomsService {
	return &lomsClient{
		ctx:    ctx,
		logger: l,
		order:  order.NewOrderAPIClient(conn),
		stock:  stock.NewStockAPIClient(conn),
	}
}

// ClientDialler объявляет соединение для клиента сервиса LOMS.
func ClientDialler(ctx context.Context, l logger.Logger, stngs *config.LomsServiceSettings) (*grpc.ClientConn, error) {
	u, err := url.Parse(stngs.ProductServiceHost + ":" + fmt.Sprintf("%d", stngs.ProductServicePort))
	if err != nil {
		l.Error(ctx, fmt.Sprintf("Ошибка форматирования url LOMS сервиса - %v", err))

		return nil, fmt.Errorf("Ошибка форматирования url LOMS сервиса - %w", err)
	}

	i := interceptor.NewInterceptor(ctx, l)

	conn, err := grpc.NewClient(
		u.String(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			i.Tracer(),
			i.Metrics(u.String(), stngs.ProductServiceHost),
			i.Logger(),
		),
	)
	if err != nil {
		l.Error(ctx, fmt.Sprintf("Ошибка регистрации коннекта с сервисом LOMS - %v", err))

		return nil, fmt.Errorf("Ошибка регистрации коннекта с сервисом LOMS - %w", err)
	}

	return conn, nil
}
