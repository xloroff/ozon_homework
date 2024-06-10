package memorystore

import (
	"context"
	"sync"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/model"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/pkg/logger"
)

// Storage имплементирует методы управления хранилищем памяти.
type Storage interface {
	AddItem(ctx context.Context, item *model.AddItem) error
	GetAllUserItems(ctx context.Context, userID int64) (*model.Cart, error)
	DelItem(ctx context.Context, item *model.DelItem) bool
	DelCart(ctx context.Context, userID int64) error
}

type cartStorage struct {
	sync.RWMutex
	data   map[int64]*model.Cart
	logger logger.ILog
}

// NewCartStorage создаем хранилище.
func NewCartStorage(l logger.ILog) Storage {
	var memStorage map[int64]*model.Cart = map[int64]*model.Cart{}

	return &cartStorage{
		data:   memStorage,
		logger: l,
	}
}
