package cart

import (
	"context"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/model"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/pkg/client/loms_cli"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/pkg/logger"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/repository/memory_store"
)

// Service заведем как сервис под апи.
type Service interface {
	AddItem(ctx context.Context, item *model.AddItem) error
	GetAllUserItems(ctx context.Context, userID int64) (*model.FullUserCart, error)
	DelItem(ctx context.Context, item *model.DelItem) error
	DelCart(ctx context.Context, userID int64) error
	Checkout(ctx context.Context, userID int64) (*model.OrderCart, error)
}

type service struct {
	productCli ProductClient
	cartStore  memorystore.Storage
	lomsCli    lomscli.LomsService
	logger     logger.Logger
}

// ProductClient определяет интерфейс, который должны реализовывать клиенты.
type ProductClient interface {
	GetProduct(ctx context.Context, skuID int64) (*model.ProductResp, error)
}

// NewService создает новый сервис включая в него клиентов и хранилище.
func NewService(l logger.Logger, product ProductClient, loms lomscli.LomsService, store memorystore.Storage) Service {
	return &service{
		productCli: product,
		lomsCli:    loms,
		cartStore:  store,
		logger:     l,
	}
}
