package orderservice

import (
	"context"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/model"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/pkg/logger"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/repository/order_store"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/repository/stock_store"
)

// Service заведем как сервис под апи.
type Service interface {
	Create(user int64, items model.AllNeedReserve) (int64, error)
	Cancel(orderID int64) error
	Info(orderID int64) (*model.Order, error)
	Pay(orderID int64) error
}

type oService struct {
	ctx        context.Context
	orderStore orderstore.Storage
	stockStore stockstore.Storage
	logger     logger.ILog
}

// NewService создает новый сервис LOMS с хранилищами резервов и хранилищем заказов.
func NewService(ctx context.Context, l logger.ILog, os orderstore.Storage, ss stockstore.Storage) Service {
	return &oService{
		ctx:        ctx,
		orderStore: os,
		stockStore: ss,
		logger:     l,
	}
}
