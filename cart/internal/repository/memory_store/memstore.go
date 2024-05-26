package memorystore

import (
	"context"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/internal/model/v1"
)

// Storage имплементирует методы управления хранилищем памяти.
type Storage interface {
	AddItem(ctx context.Context, item *v1.AddItem) error
	GetAllUserItems(ctx context.Context, userID int64) (*v1.Cart, error)
	DelItem(ctx context.Context, item *v1.DelItem) bool
	DelCart(ctx context.Context, userID int64) error
}

type cartStorage struct {
	data map[int64]*v1.Cart
}

// NewCartStorage создаем хранилище.
func NewCartStorage() Storage {
	var memStorage map[int64]*v1.Cart = map[int64]*v1.Cart{}

	return &cartStorage{
		data: memStorage,
	}
}