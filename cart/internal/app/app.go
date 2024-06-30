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
	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/pkg/tracer"
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
	cnfg, err := a.initializeConfig()
	if err != nil {
		return err
	}

	l := logger.InitializeLogger(cnfg.LogLevel, cnfg.LogType)
	clsr := closer.NewCloser(l, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	clsr.Add(func() error {
		return l.Close()
	})
	defer clsr.Wait()

	err = a.initializeTracer(cnfg, clsr)
	if err != nil {
		return err
	}

	productCli := productcli.NewProductClient(l, cnfg.ProductServiceSettings)
	memStore := memorystore.NewCartStorage(l)

	lomsCli, err := initializeLomsClient(a.ctx, l, clsr, cnfg.LomsServiceSettings)
	if err != nil {
		return err
	}

	webSrvr := server.NewServer(a.ctx, l, cnfg, productCli, lomsCli, memStore)

	err = startWebServer(webSrvr, clsr)
	if err != nil {
		return err
	}

	return nil
}

func (a *app) initializeConfig() (*config.ApplicationParameters, error) {
	cnfg, err := config.LoadAPIConfig()
	if err != nil {
		return nil, fmt.Errorf("Ошибка получения параметров из конфигурационного файла - %w", err)
	}

	return cnfg, nil
}

func (a *app) initializeTracer(cnfg *config.ApplicationParameters, clsr closer.Closer) error {
	err := tracer.InitTracerProvider(a.ctx, cnfg.JaegerSettings)
	if err != nil {
		return fmt.Errorf("Ошибка создания провайдера трасировки - %w", err)
	}

	clsr.Add(tracer.Close)

	return nil
}

func initializeLomsClient(ctx context.Context, l logger.Logger, clsr closer.Closer, settings *config.LomsServiceSettings) (lomscli.LomsService, error) {
	lomsDialler, err := lomscli.ClientDialler(ctx, l, settings)
	if err != nil {
		clsr.CloseAll()

		return nil, fmt.Errorf("Ошибка создания диаллера для сервиса заказов - %w", err)
	}

	lomsCli := lomscli.NewClient(ctx, l, lomsDialler)

	return lomsCli, nil
}

func startWebServer(webSrvr server.AppServer, clsr closer.Closer) error {
	err := webSrvr.Start()
	if err != nil {
		clsr.CloseAll()

		return fmt.Errorf("Ошибка запуска сервера - %w", err)
	}

	clsr.Add(func() error {
		return webSrvr.Stop()
	})

	return nil
}
