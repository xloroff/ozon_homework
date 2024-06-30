package stockapi

import (
	"context"
	"errors"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/model"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/pkg/tracer"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/pkg/api/stock/v1"
)

// Info получение информации о остатке товара.
func (a *API) Info(ctx context.Context, req *stock.StockInfoRequest) (*stock.StockInfoResponse, error) {
	ctx, span := tracer.StartSpanFromContext(ctx, "api.stockapi.info")
	span.SetTag("component", "stockapi")

	defer span.End()

	c, err := a.stockService.Info(ctx, req.GetSku())
	if err != nil {
		span.SetTag("error", true)
		a.logger.Debugf(ctx, "StockApi.Info: Ошибка остатков товара %s - %v", req.GetSku(), err)

		if errors.Is(err, model.ErrReserveNotFound) {
			return nil, status.Errorf(codes.NotFound, fmt.Sprintf("Товар %d отсутствует", req.GetSku()))
		}

		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Неизвестная ошибка - %v", err))
	}

	return &stock.StockInfoResponse{
		Count: uint64(c),
	}, nil
}
