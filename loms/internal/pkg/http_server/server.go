package httpserver

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/config"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/pkg/logger"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/pkg/api/order/v1"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/pkg/api/stock/v1"
)

// Server интерфейс управляющий HTTP-сервером.
type Server interface {
	Start() error
	Stop(dur float64) error
	RegisterAPI() error
}

type server struct {
	ctx        context.Context
	logger     logger.ILog
	mux        *http.ServeMux
	gwmux      *runtime.ServeMux
	conn       *grpc.ClientConn
	httpServer *http.Server
}

// NewServer создает новый экземпляр HTTP сервера.
func NewServer(ctx context.Context, l logger.ILog, cnfg *config.ApplicationParameters) (Server, error) {
	var err error

	s := &server{
		ctx:    ctx,
		logger: l,
		mux:    http.NewServeMux(),
		gwmux:  runtime.NewServeMux(),
	}

	s.conn, err = grpc.NewClient(
		fmt.Sprintf(":%d", cnfg.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		l.Errorf(s.ctx, "Ошибка обращения к GRPC серверу - %v", err)

		return nil, fmt.Errorf("Ошибка обращения к GRPC серверу - %w", err)
	}

	s.httpServer = &http.Server{
		Addr:              fmt.Sprintf(":%d", cnfg.WebPort),
		ReadHeaderTimeout: time.Duration(cnfg.GracefulTimeout * float64(time.Second)),
		ReadTimeout:       time.Duration(cnfg.GracefulTimeout * float64(time.Second)),
		Handler:           s.mux,
	}

	// Публикуем Swagger.
	swggr := http.FileServer(http.Dir("pkg/docs"))
	s.mux.Handle("/docs/", http.StripPrefix("/docs/", swggr))
	s.mux.HandleFunc("/healthcheck", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Прокидываем клиента GRPC.
	s.mux.Handle("/", s.gwmux)

	return s, nil
}

// Start стартует HTTP-сервер.
func (s *server) Start() error {
	s.logger.Warn(s.ctx, "Запуск HTTP сервера...")

	go func() {
		err := s.httpServer.ListenAndServe()
		if err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				s.logger.Fatalf(s.ctx, "Ошибка запуска HTTP сервера - %v", err)
			}
		}
	}()

	s.logger.Warn(s.ctx, "HTTP cервер запущен...")

	return nil
}

// Stop останавливает HTTP-сервер.
func (s *server) Stop(dur float64) error {
	defer s.logger.Warn(s.ctx, "Остановка HTTP-сервера произведена...")

	ctx, cancel := context.WithTimeout(s.ctx, time.Duration(dur*float64(time.Second)))
	defer cancel()

	err := s.httpServer.Shutdown(ctx)
	if err != nil {
		s.logger.Errorf(ctx, "Ошибка остановки HTTP-сервера - %v", err)

		return fmt.Errorf("Ошибка остановки HTTP-сервера - %w", err)
	}

	return nil
}

// RegisterAPI добавляет приклады GRPC-сервера в прокси.
func (s *server) RegisterAPI() error {
	var err error

	err = order.RegisterOrderAPIHandler(s.ctx, s.gwmux, s.conn)
	if err != nil {
		s.logger.Errorf(s.ctx, "Ошибка имплементации сервиса заказов GRPC-proxy - %v", err)

		return fmt.Errorf("Ошибка имплементации сервиса заказов GRPC-proxy -  %w", err)
	}

	err = stock.RegisterStockAPIHandler(s.ctx, s.gwmux, s.conn)
	if err != nil {
		s.logger.Errorf(s.ctx, "Ошибка имплементации сервиса остатков GRPC-proxy - %v", err)

		return fmt.Errorf("Ошибка имплементации сервиса остатков GRPC-proxy -  %w", err)
	}

	return nil
}
