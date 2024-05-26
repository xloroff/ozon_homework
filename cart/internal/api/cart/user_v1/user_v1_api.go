package user_v1

import (
	"gitlab.ozon.dev/xloroff/ozon-hw-go/internal/initilize"
	productcli_resty "gitlab.ozon.dev/xloroff/ozon-hw-go/internal/pkg/client/resty_cli"
	productcli_stnd "gitlab.ozon.dev/xloroff/ozon-hw-go/internal/pkg/client/standart_cli"
	storage "gitlab.ozon.dev/xloroff/ozon-hw-go/internal/repository/memory_store"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/internal/service/cart"
)

type apiv1 struct {
	cartService cart.Service
}

// NewApiV1 запускает сервис с хранилкой и коммуникацией с внешними сервисами.
func NewApiV1(settings *initilize.ConfigAPI) *apiv1 {
	var cartService cart.Service
	switch settings.ClientVer {
	// Стандартный клиент.
	case 1:
		cartService = cart.NewService(productcli_stnd.NewStdClient(), storage.NewCartStorage())
	// На базе resty.
	case 2:
		cartService = cart.NewService(productcli_resty.NewStdClient(), storage.NewCartStorage())
	default:
		cartService = cart.NewService(productcli_stnd.NewStdClient(), storage.NewCartStorage())
	}

	return &apiv1{
		cartService: cartService,
	}
}