package cart

import (
	"context"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/internal/initilize"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/internal/model/v1"
	productcli_resty "gitlab.ozon.dev/xloroff/ozon-hw-go/internal/pkg/client/resty_cli"
	productcli_stnd "gitlab.ozon.dev/xloroff/ozon-hw-go/internal/pkg/client/standart_cli"
	storage "gitlab.ozon.dev/xloroff/ozon-hw-go/internal/repository/memory_store"
)

// Service заведем как сервис под апи.
type Service interface {
	AddItem(ctx context.Context, settings *initilize.ConfigAPI, item *v1.AddItem) error
	GetAllUserItems(ctx context.Context, settings *initilize.ConfigAPI, userID int64) (*v1.FullUserCart, error)
	DelItem(ctx context.Context, item *v1.DelItem) error
	DelCart(ctx context.Context, userID int64) error
}

type cService struct {
	productCli ProductClient
	cartStore  storage.Storage
}

// ProductClient определяет интерфейс, который должны реализовывать клиенты.
type ProductClient interface {
	GetProduct(ctx context.Context, settings *initilize.ConfigAPI, skuID int64) (*v1.ProductResp, error)
}

func NewService(product ProductClient, store storage.Storage) Service {
	// Проверка соответствия product интерфейсу ProductClient.
	switch product.(type) {
	case productcli_resty.Client:
		return &cService{
			productCli: product.(productcli_resty.Client),
			cartStore:  store,
		}
	case productcli_stnd.Client:
		return &cService{
			productCli: product.(productcli_stnd.Client),
			cartStore:  store,
		}
	default:
		return &cService{
			productCli: product.(productcli_stnd.Client),
			cartStore:  store,
		}
	}
}