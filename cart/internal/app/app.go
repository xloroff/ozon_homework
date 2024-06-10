package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/config"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/pkg/client/loms_cli"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/pkg/client/product_cli"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/pkg/logger"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/pkg/server"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/repository/memory_store"
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
		return fmt.Errorf("Ошибка получения параметров из конфигурационного файла - %w", err)
	}

	// Стартуем логгер
	l := logger.InitializeLogger(cnfg.LogLevel, cnfg.LogType)

	// Старт клиентов
	productCli := productcli.NewProductClient(l, cnfg.ProductServiceSettings)
	memStore := memorystore.NewCartStorage(l)

	lomsDialler, err := lomscli.ClientDialler(a.ctx, l, cnfg.LomsServiceSettings)
	if err != nil {
		return fmt.Errorf("Ошибка создания диаллера для сервиса заказов - %w", err)
	}

	lomsCli := lomscli.NewClient(a.ctx, l, lomsDialler)

	// Стартуем веб-сервер.
	err = server.NewServer(a.ctx, l, cnfg, productCli, lomsCli, memStore).Start()
	if err != nil {
		return fmt.Errorf("Ошибка запуска сервера - %w", err)
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
