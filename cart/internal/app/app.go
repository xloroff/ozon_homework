package app

import (
	"context"
	"fmt"
	"os"
	"syscall"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/config"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/pkg/client/loms_cli"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/pkg/client/product_cli"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/pkg/closer"
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
	// Стартуем получение конфига
	cnfg, err := config.LoadAPIConfig()
	if err != nil {
		return fmt.Errorf("Ошибка получения параметров из конфигурационного файла - %w", err)
	}

	// Стартуем логгер
	l := logger.InitializeLogger(cnfg.LogLevel, cnfg.LogType)

	clsr := closer.NewCloser(l, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer clsr.Wait()

	clsr.Add(func() error {
		return l.Close()
	})

	// Старт клиентов
	productCli := productcli.NewProductClient(l, cnfg.ProductServiceSettings)
	memStore := memorystore.NewCartStorage(l)

	lomsDialler, err := lomscli.ClientDialler(a.ctx, l, cnfg.LomsServiceSettings)
	if err != nil {
		clsr.CloseAll()

		return fmt.Errorf("Ошибка создания диаллера для сервиса заказов - %w", err)
	}

	lomsCli := lomscli.NewClient(a.ctx, l, lomsDialler)

	// Стартуем веб-сервер.
	webSrvr := server.NewServer(a.ctx, l, cnfg, productCli, lomsCli, memStore)

	err = webSrvr.Start()
	if err != nil {
		clsr.CloseAll()

		return fmt.Errorf("Ошибка запуска сервера - %w", err)
	}

	clsr.Add(func() error {
		return webSrvr.Stop()
	})

	return nil
}
