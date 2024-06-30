package cartapi

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/pkg/client/loms_cli"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/pkg/client/product_cli"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/pkg/logger"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/repository/memory_store"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/service/cart"
)

// API первая версия API cart service.
type API struct {
	cartService cart.Service
	logger      logger.Logger
}

// CustomResponseWriter реализует интерфейс http.ResponseWriter.
type CustomResponseWriter struct {
	responseWriter http.ResponseWriter
	StatusCode     int
}

// NewAPI запускает сервис с хранилкой и коммуникацией с внешними сервисами.
func NewAPI(l logger.Logger, productCli productcli.Client, lomsCli lomscli.LomsService, memStore memorystore.Storage) *API {
	return &API{
		cartService: cart.NewService(l, productCli, lomsCli, memStore),
		logger:      l,
	}
}

func getReqURLTemplate(r *http.Request) string {
	curRoute := mux.CurrentRoute(r)
	if curRoute == nil {
		return r.RequestURI
	}

	pathTemplate, err := curRoute.GetPathTemplate()
	if err != nil {
		return mux.CurrentRoute(r).GetName()
	}

	return pathTemplate
}

// ExtendResponseWriter рвозвращает кастомный http.ResponseWriter.
func ExtendResponseWriter(w http.ResponseWriter) *CustomResponseWriter {
	return &CustomResponseWriter{w, 0}
}

// Write реализует одноименный метод http.ResponseWriter.
func (w *CustomResponseWriter) Write(b []byte) (int, error) {
	part, err := w.responseWriter.Write(b)
	if err != nil {
		return part, fmt.Errorf("Ошибка при записи в http.ResponseWriter - %w", err)
	}

	return part, nil
}

// Header реализует одноименный метод http.ResponseWriter.
func (w *CustomResponseWriter) Header() http.Header {
	return w.responseWriter.Header()
}

// WriteHeader реализует одноименный метод http.ResponseWriter.
func (w *CustomResponseWriter) WriteHeader(statusCode int) {
	w.StatusCode = statusCode
	w.responseWriter.WriteHeader(statusCode)
}
