package stockstore

import (
	"context"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/model"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/pkg/logger"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/pkg/db"
)

// Storage имплементирует методы управления хранилищем памяти.
type Storage interface {
	AddReserve(items model.AllNeedReserve) error
	GetAvailableForReserve(sku int64) (uint16, error)
	DelItemFromReserve(items model.AllNeedReserve) error
	CancelReserve(items model.AllNeedReserve) error
}

type reserveStorage struct {
	ctx    context.Context
	data   db.ClientBD
	logger logger.ILog
}

// NewReserveStorage создает хранилище остатков.
func NewReserveStorage(ctx context.Context, l logger.ILog, bdCli db.ClientBD) (Storage, error) {
	return &reserveStorage{
		ctx:    ctx,
		data:   bdCli,
		logger: l,
	}, nil
}
