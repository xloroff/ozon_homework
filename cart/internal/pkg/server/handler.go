package server

import (
	"net/http"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/api/cartapi"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/model"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/pkg/client/loms_cli"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/pkg/client/product_cli"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/pkg/server/middleware"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/repository/memory_store"
)

// AddHandl добавляем приклады к нашему серверу и обрабатываем на уровне middleware входящие запросы.
func (s *server) AddHandl(productCli productcli.Client, lomsCli lomscli.LomsService, memStore memorystore.Storage) error {
	// Включаем логирование для всех прикладов.
	s.router.Use(middleware.Logging(s.ctx, s.logger))

	// Создаем сервис для API.
	api := cartapi.NewAPI(s.logger, productCli, lomsCli, memStore)
	// Cаброутер API.
	user := s.router.PathPrefix("/user").Subrouter()
	// Саброутер Helthcheck
	s.router.HandleFunc("/healthcheck", Healthcheck).Methods(http.MethodHead, http.MethodGet)

	// Добавление товара в корзину.
	user.HandleFunc("/{"+model.UsrID+"}/cart/{"+model.SkuID+"}", api.AddItem).Methods(http.MethodPost)
	// Получение всей корзины пользователя.
	user.HandleFunc("/{"+model.UsrID+"}/cart/list", api.GetAllUserItems).Methods(http.MethodGet)
	// Удаление товара из корзины.
	user.HandleFunc("/{"+model.UsrID+"}/cart/{"+model.SkuID+"}", api.DelItem).Methods(http.MethodDelete)
	// Полное удаление корзины.
	user.HandleFunc("/{"+model.UsrID+"}/cart", api.DelCart).Methods(http.MethodDelete)
	// Создание заказа.
	user.HandleFunc("/{"+model.UsrID+"}/checkout", api.Checkout).Methods(http.MethodPost)

	return nil
}
