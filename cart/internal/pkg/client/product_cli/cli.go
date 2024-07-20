package productcli

import (
	"context"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/internal/config"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/internal/model"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/internal/pkg/logger"
)

const urlapi = "/get_product"

// Client интерфейс клиента для сервиса продуктов.
type Client interface {
	GetProduct(ctx context.Context, skuID int64) (*model.ProductResp, error)
}

type client struct {
	logger logger.ILog
	config *config.ProductServiceSettings
}

// NewProductClient создает новый клиент.
func NewProductClient(l logger.ILog, stngs *config.ProductServiceSettings) Client {
	return &client{
		logger: l,
		config: stngs,
	}
}
