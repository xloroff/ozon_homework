package stockservice

import (
	"context"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/pkg/logger"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/repository/stock_store"
)

// Service заведем как сервис под апи.
type Service interface {
	Info(sku int64) (uint16, error)
}

type sService struct {
	ctx        context.Context
	stockStore stockstore.Storage
	logger     logger.ILog
}

// NewService создает новый сервис LOMS с хранилищами резервов и хранилищем заказов.
func NewService(ctx context.Context, l logger.ILog, ss stockstore.Storage) Service {
	return &sService{
		ctx:        ctx,
		stockStore: ss,
		logger:     l,
	}
}
