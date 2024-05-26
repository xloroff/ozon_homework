package restycli

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

/* TODO productStorage храним товары которые уже посещали
var productStorage map[int64]*v1.ProductResp = map[int64]*v1.ProductResp{}
*/

// NewStdClient создает новый клиент.
func NewStdClient() Client {
	return &client{}
}