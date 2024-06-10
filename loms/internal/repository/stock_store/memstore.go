package stockstore

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/model"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/pkg/logger"
)

// Storage имплементирует методы управления хранилищем памяти.
type Storage interface {
	AddReserve(items model.AllNeedReserve) error
	GetAvailableForReserve(sku int64) (uint16, error)
	DelItemFromReserve(items model.AllNeedReserve) error
	CancelReserve(items model.AllNeedReserve) error
}

type reserveStorage struct {
	sync.RWMutex
	ctx    context.Context
	data   model.AllReserveItems
	logger logger.ILog
}

// NewReserveStorage создает хранилище остатков.
func NewReserveStorage(ctx context.Context, l logger.ILog) (Storage, error) {
	memReserve := map[int64]*model.ReserveItem{}

	stockJSON, err := os.ReadFile("stock-data.json")
	if err != nil {
		return nil, fmt.Errorf("Ошибка чтения файла - %w", err)
	}

	stockItems := []model.ReserveItem{}

	err = json.Unmarshal(stockJSON, &stockItems)
	if err != nil {
		return nil, fmt.Errorf("Проблема с обработкой файла остатков - %w", err)
	}

	for _, stockItem := range stockItems {
		st := stockItem

		memReserve[stockItem.Sku] = &st
	}

	return &reserveStorage{
		ctx:    ctx,
		data:   memReserve,
		logger: l,
	}, nil
}
