package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/config"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/pkg/client/loms_cli"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/pkg/client/product_cli"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/pkg/logger"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/repository/memory_store"
)

const // Максимальное время ожидания нового соединения.
maxTimeWait = 20

// AppServer интерфейс для сервера.
type AppServer interface {
	Start() error
	Stop() error
}

type server struct {
	ctx        context.Context
	router     *mux.Router
	httpServer *http.Server
	logger     logger.Logger
	config     *config.ApplicationParameters
}

// NewServer создает экземпляр http сервера с таймингами работы прикладов.
func NewServer(ctx context.Context, l logger.Logger, cnfg *config.ApplicationParameters, productCli productcli.Client, lomsCli lomscli.LomsService, memStore memorystore.Storage) AppServer {
	l.Warn(ctx, "Запуск сервера...")

	// Создаем новый роутер.
	r := mux.NewRouter()

	s := &server{
		ctx:    ctx,
		router: r,
		logger: l,
		config: cnfg,
	}

	s.httpServer = &http.Server{
		Addr:         ":" + fmt.Sprint(cnfg.Port),
		WriteTimeout: time.Duration(cnfg.GracefulTimeout * float64(time.Second)),
		ReadTimeout:  time.Duration(cnfg.GracefulTimeout * float64(time.Second)),
		IdleTimeout:  time.Second * time.Duration(maxTimeWait),
		Handler:      r,
	}

	// Добавляем точки API и мидлвари.
	err := s.AddHandl(productCli, lomsCli, memStore)
	if err != nil {
		err = s.Stop()
		if err != nil {
			l.Debugf(ctx, "server.NewServer", err)
		}
	}

	return s
}

// Start запуск сервера.
func (s *server) Start() error {
	// Запускаем сервер в рутине.
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil {
			s.logger.Warnf(s.ctx, "Проблема с сервером - ", err)
		}
	}()

	s.logger.Warn(s.ctx, "Сервер запущен...")

	return nil
}

func (s *server) Stop() error {
	defer func() {
		s.logger.Warn(s.ctx, "Остановка веб-сервера произведена...")
	}()

	// Таймаут ожидания соединений перед завершением сервера.
	wait := time.Duration(s.config.GracefulTimeout) * time.Second

	// Создаем контекст с таймаутом для завершения сервера.
	ctx, cancel := context.WithTimeout(s.ctx, wait)
	defer cancel()

	// Завершаем сервер.
	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("Ошибка завершения веб-сервера: %w", err)
	}

	return nil
}
