package standartcli

import (
	"context"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/internal/initilize"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/internal/model/v1"
)

// Client интерфейс клиента для сервиса продуктов.
type Client interface {
	GetProduct(ctx context.Context, settings *initilize.ConfigAPI, skuID int64) (*v1.ProductResp, error)
}

type client struct{}

// NewStdClient создает новый клиент.
func NewStdClient() Client {
	return &client{}
}