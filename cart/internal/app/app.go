package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/internal/config"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/internal/pkg/logger"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/internal/pkg/server"
)

// Application запуск приложения.
type Application interface {
	Run() error
}

type app struct {
	ctx    context.Context
	cancel context.CancelFunc
}

// NewApp создание приложения.
func NewApp(ctx context.Context) Application {
	ctx, cancel := context.WithCancel(ctx)

	return &app{
		ctx:    ctx,
		cancel: cancel,
	}
}

// Run запуск.
func (a *app) Run() error {
	// Канал для сигналов завершения.
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Стартуем получение конфига
	cnfg, err := config.LoadAPIConfig()
	if err != nil {
		panic(fmt.Errorf("Ошибка получения параметров из конфигурационного файла - %w", err))
	}

	// Стартуем логгер
	l := logger.InitializeLogger(cnfg.LogLevel, cnfg.LogType)

	// Стартуем веб-сервер.
	err = server.NewServer(a.ctx, l, cnfg).Start()
	if err != nil {
		l.Fatalf(a.ctx, "Ошибка запуска сервера - %v", err.Error())
	}

	// Обработка сигналов завершения.
	go func() {
		<-sigChan
		l.Warnf(a.ctx, "Получен сигнал завершения, остановка приложения произведена...")

		defer l.Close()
		defer a.cancel()
	}()

	// Блокировка до завершения контекста.
	<-a.ctx.Done()

	return nil
}
