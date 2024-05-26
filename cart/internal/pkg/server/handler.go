package server

import (
	"net/http"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/internal/api/cart/user_v1"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/internal/model/v1"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/internal/pkg/server/middleware"
)

// AddHandl добавляем приклады к нашему серверу и обрабатываем на уровне middleware входящие запросы.
func (s *server) AddHandl() error {
	// Включаем логирование для всех прикладов.
	s.router.Use(middleware.MiddlewareLogging(s.ctx))

	// Создаем сервис для v1 API.
	apiV1 := user_v1.NewApiV1(s.settings)
	// Cаброутер API версии 1.
	// TODO можно валидацию структуры вынести в мидлварю саброутера.
	userv1 := s.router.PathPrefix("/user").Subrouter()

	// Добавление товара в корзину.
	userv1.HandleFunc("/{"+v1.UsrID+"}/cart/{"+v1.SkuID+"}", apiV1.AddItem(s.settings)).Methods(http.MethodPost)
	// Получение всей корзины пользователя.
	userv1.HandleFunc("/{"+v1.UsrID+"}/cart/list", apiV1.GetAllUserItems(s.settings)).Methods(http.MethodGet)
	// Удаление товара из корзины.
	userv1.HandleFunc("/{"+v1.UsrID+"}/cart/{"+v1.SkuID+"}", apiV1.DelItem).Methods(http.MethodDelete)
	// Полное удаление корзины.
	userv1.HandleFunc("/{"+v1.UsrID+"}/cart", apiV1.DelCart).Methods(http.MethodDelete)

	return nil
}