package stockstore

import (
	"context"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/model"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/pkg/logger"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/pkg/db"
)

const (
	repName = "StockStore"
)

// Storage имплементирует методы управления хранилищем памяти.
type Storage interface {
	AddReserve(ctx context.Context, items model.AllNeedReserve) error
	GetAvailableForReserve(ctx context.Context, sku int64) (uint16, error)
	DelItemFromReserve(ctx context.Context, items model.AllNeedReserve) error
	CancelReserve(ctx context.Context, items model.AllNeedReserve) error
}

type reserveStorage struct {
	ctx    context.Context
	data   db.ClientBD
	logger logger.Logger
}

// NewReserveStorage создает хранилище остатков.
func NewReserveStorage(ctx context.Context, l logger.Logger, bdCli db.ClientBD) (Storage, error) {
	return &reserveStorage{
		ctx:    ctx,
		data:   bdCli,
		logger: l,
	}, nil
}
