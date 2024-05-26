package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/internal/initilize"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/internal/pkg/logger"
)

// AppServer интерфейс для сервера.
type AppServer interface {
	Start() error
	Stop() error
}

type server struct {
	ctx        context.Context
	router     *mux.Router
	httpServer *http.Server
	settings   *initilize.ConfigAPI
}

// NewServer создает экземпляр http сервера с таймингами работы прикладов.
func NewServer(ctx context.Context, stgs *initilize.ConfigAPI) AppServer {
	logger.Warn(ctx, "Запуск сервера...")

	// Создаем новый роутер.
	r := mux.NewRouter()

	s := &server{
		ctx:      ctx,
		router:   r,
		settings: stgs,
	}

	s.httpServer = &http.Server{
		Addr:         ":" + fmt.Sprint(stgs.Port),
		WriteTimeout: time.Duration(stgs.GracefulTimeout * float64(time.Second)),
		ReadTimeout:  time.Duration(stgs.GracefulTimeout * float64(time.Second)),
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}

	// Добавляем точки API и мидлвари.
	err := s.AddHandl()
	if err != nil {
		s.Stop()
	}

	return s
}

// Start запуск сервера.
func (s *server) Start() error {
	// Запускаем сервер в рутине.
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil {
			logger.Warnf(s.ctx, "Проблема с сервером - ", err)
		}
	}()

	logger.Warn(s.ctx, "Сервер запущен...")

	// Обработка сигнала завершения сервера.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	// Блокируем выполнение до получения сигнала.
	<-c

	return s.Stop()
}

func (s *server) Stop() error {
	defer func() {
		logger.Warn(s.ctx, "Остановка веб-сервера произведена...")
	}()

	// Таймаут ожидания соединений перед завершением сервера.
	wait := time.Duration(s.settings.GracefulTimeout) * time.Second

	// Создаем контекст с таймаутом для завершения сервера.
	ctx, cancel := context.WithTimeout(s.ctx, wait)
	defer cancel()

	// Завершаем сервер.
	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("Ошибка завершения веб-сервера: %v", err)
	}

	return nil
}
