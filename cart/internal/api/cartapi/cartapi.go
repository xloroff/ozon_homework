package cartapi

import (
	"gitlab.ozon.dev/xloroff/ozon-hw-go/internal/config"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/internal/pkg/client/product_cli"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/internal/pkg/logger"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/internal/repository/memory_store"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/internal/service/cart"
)

// API первая версия API cart service.
type API struct {
	cartService cart.Service
	logger      logger.ILog
}

// NewAPI запускает сервис с хранилкой и коммуникацией с внешними сервисами.
func NewAPI(l logger.ILog, stngs *config.ProductServiceSettings) *API {
	return &API{
		cartService: cart.NewService(l, productcli.NewProductClient(l, stngs), memorystore.NewCartStorage(l)),
		logger:      l,
	}
}
