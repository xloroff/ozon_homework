package cartapi

import (
	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/pkg/client/loms_cli"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/pkg/client/product_cli"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/pkg/logger"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/repository/memory_store"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/service/cart"
)

// API первая версия API cart service.
type API struct {
	cartService cart.Service
	logger      logger.ILog
}

// NewAPI запускает сервис с хранилкой и коммуникацией с внешними сервисами.
func NewAPI(l logger.ILog, productCli productcli.Client, lomsCli lomscli.LomsService, memStore memorystore.Storage) *API {
	return &API{
		cartService: cart.NewService(l, productCli, lomsCli, memStore),
		logger:      l,
	}
}
