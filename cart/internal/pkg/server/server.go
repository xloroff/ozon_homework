package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/internal/config"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/internal/pkg/logger"
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
	logger     logger.ILog
	config     *config.ApplicationParameters
}

// NewServer создает экземпляр http сервера с таймингами работы прикладов.
func NewServer(ctx context.Context, l logger.ILog, cnfg *config.ApplicationParameters) AppServer {
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
	err := s.AddHandl()
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

	// Обработка сигнала завершения сервера.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	// Блокируем выполнение до получения сигнала.
	<-c

	return s.Stop()
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
