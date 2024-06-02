package cart

import (
	"context"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/internal/model"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/internal/pkg/logger"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/internal/repository/memory_store"
)

// Service заведем как сервис под апи.
type Service interface {
	AddItem(ctx context.Context, item *model.AddItem) error
	GetAllUserItems(ctx context.Context, userID int64) (*model.FullUserCart, error)
	DelItem(ctx context.Context, item *model.DelItem) error
	DelCart(ctx context.Context, userID int64) error
}

type cService struct {
	productCli ProductClient
	cartStore  memorystore.Storage
	logger     logger.ILog
}

// ProductClient определяет интерфейс, который должны реализовывать клиенты.
type ProductClient interface {
	GetProduct(ctx context.Context, skuID int64) (*model.ProductResp, error)
}

// NewService создает новый сервис включая в него клиентов и хранилище.
func NewService(l logger.ILog, product ProductClient, store memorystore.Storage) Service {
	return &cService{
		productCli: product,
		cartStore:  store,
		logger:     l,
	}
}
