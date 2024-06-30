package productcli

import (
	"context"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/config"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/model"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/pkg/logger"
)

const urlapi = "/get_product"

// Client интерфейс клиента для сервиса продуктов.
type Client interface {
	GetProduct(ctx context.Context, skuID int64) (*model.ProductResp, error)
}

type client struct {
	logger logger.Logger
	config *config.ProductServiceSettings
}

// NewProductClient создает новый клиент.
func NewProductClient(l logger.Logger, stngs *config.ProductServiceSettings) Client {
	return &client{
		logger: l,
		config: stngs,
	}
}
