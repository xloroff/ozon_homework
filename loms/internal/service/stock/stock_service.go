package stockservice

import (
	"context"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/pkg/logger"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/repository/stock_store"
)

// Service заведем как сервис под апи.
type Service interface {
	Info(ctx context.Context, sku int64) (uint16, error)
}

type sService struct {
	ctx        context.Context
	stockStore stockstore.Storage
	logger     logger.Logger
}

// NewService создает новый сервис LOMS с хранилищами резервов и хранилищем заказов.
func NewService(ctx context.Context, l logger.Logger, ss stockstore.Storage) Service {
	return &sService{
		ctx:        ctx,
		stockStore: ss,
		logger:     l,
	}
}
