package orderstore

import (
	"context"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/model"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/pkg/logger"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/pkg/db"
)

// Storage имплементирует методы управления хранилищем памяти.
type Storage interface {
	AddOrder(user int64, items model.OrderItems) (int64, error)
	GetOrder(orderID int64) (*model.Order, error)
	SetStatus(orderID int64, status string) error
}

type orderStorage struct {
	ctx    context.Context
	data   db.ClientBD
	logger logger.ILog
}

// NewOrderStorage создает хранилище заказов.
func NewOrderStorage(ctx context.Context, l logger.ILog, bdCli db.ClientBD) (Storage, error) {
	return &orderStorage{
		ctx:    ctx,
		data:   bdCli,
		logger: l,
	}, nil
}
