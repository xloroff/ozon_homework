package server

import (
	"net/http"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/internal/api/cartapi"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/internal/model"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/internal/pkg/server/middleware"
)

// AddHandl добавляем приклады к нашему серверу и обрабатываем на уровне middleware входящие запросы.
func (s *server) AddHandl() error {
	// Включаем логирование для всех прикладов.
	s.router.Use(middleware.Logging(s.ctx, s.logger))

	// Создаем сервис для API.
	api := cartapi.NewAPI(s.logger, s.config.ProductServiceSettings)
	// Cаброутер API.
	// TODO можно валидацию структуры вынести в мидлварю саброутера.
	user := s.router.PathPrefix("/user").Subrouter()

	// Добавление товара в корзину.
	user.HandleFunc("/{"+model.UsrID+"}/cart/{"+model.SkuID+"}", api.AddItem).Methods(http.MethodPost)
	// Получение всей корзины пользователя.
	user.HandleFunc("/{"+model.UsrID+"}/cart/list", api.GetAllUserItems).Methods(http.MethodGet)
	// Удаление товара из корзины.
	user.HandleFunc("/{"+model.UsrID+"}/cart/{"+model.SkuID+"}", api.DelItem).Methods(http.MethodDelete)
	// Полное удаление корзины.
	user.HandleFunc("/{"+model.UsrID+"}/cart", api.DelCart).Methods(http.MethodDelete)

	return nil
}
