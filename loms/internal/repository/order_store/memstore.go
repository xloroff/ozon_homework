package orderstore

import (
	"context"
	"sync"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/model"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/pkg/logger"
)

// Storage имплементирует методы управления хранилищем памяти.
type Storage interface {
	AddOrder(user int64, items model.OrderItems) (int64, error)
	GetOrder(orderID int64) (*model.Order, error)
	SetStatus(orderID int64, status string) error
}

type orderStorage struct {
	sync.RWMutex
	ctx    context.Context
	data   model.AllOrderItems
	logger logger.ILog
}

// NewOrderStorage создает хранилище заказов.
func NewOrderStorage(ctx context.Context, l logger.ILog) (Storage, error) {
	memOrders := map[int64]*model.Order{}

	return &orderStorage{
		ctx:    ctx,
		data:   memOrders,
		logger: l,
	}, nil
}
