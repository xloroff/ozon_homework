package server

import (
	"net/http"
	"net/http/pprof"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/api/cartapi"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/model"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/pkg/client/loms_cli"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/pkg/client/product_cli"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/pkg/server/middleware"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/repository/memory_store"
)

// AddHandl добавляем приклады к нашему серверу и обрабатываем на уровне middleware входящие запросы.
func (s *server) AddHandl(productCli productcli.Client, lomsCli lomscli.LomsService, memStore memorystore.Storage) error {
	// Создаем сервис для API.
	api := cartapi.NewAPI(s.logger, productCli, lomsCli, memStore)
	// Cаброутер API.
	user := s.router.PathPrefix("/user").Subrouter()
	// Метрики и логи собираем только непосредственно по запросам к сервису.
	user.Use(middleware.Metrics)
	user.Use(middleware.Logging(s.ctx, s.logger))
	// Саброутер Helthcheck
	s.router.HandleFunc("/healthcheck", Healthcheck).Methods(http.MethodHead, http.MethodGet)
	// Сабрроутер метрик.
	s.router.Path("/metrics").Handler(promhttp.Handler())

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

	// Профайлинг
	s.router.Handle("/debug/pprof/", http.HandlerFunc(pprof.Index))
	s.router.Handle("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
	s.router.Handle("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
	s.router.Handle("/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
	s.router.Handle("/debug/pprof/trace", http.HandlerFunc(pprof.Trace))
	s.router.Handle("/debug/pprof/{cmd}", http.HandlerFunc(pprof.Index))

	return nil
}
