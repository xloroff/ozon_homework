package app

import (
	"context"
	"fmt"
	"os"
	"syscall"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/config"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/pkg/api/orderapi"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/pkg/api/stockapi"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/pkg/closer"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/pkg/grpc_server"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/pkg/http_server"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/pkg/logger"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/pkg/tracer"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/repository/order_store"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/repository/stock_store"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/service/order"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/service/stock"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/pkg/db"
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

	err = a.runMigrations(cnfg, l, clsr)
	if err != nil {
		return err
	}

	dbClient, err := a.initializeDBClients(cnfg, l, clsr)
	if err != nil {
		return err
	}

	ordSrvc, stckSrvc, err := a.initializeServices(cnfg, l, dbClient, clsr)
	if err != nil {
		return err
	}

	err = a.startGRPCServer(cnfg, l, ordSrvc, stckSrvc, clsr)
	if err != nil {
		return err
	}

	err = a.startHTTPServer(cnfg, l, clsr)
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

func (a *app) runMigrations(cnfg *config.ApplicationParameters, l logger.Logger, clsr closer.Closer) error {
	err := db.MigrationPool(a.ctx, l, cnfg.MigrationFolder, cnfg.BDMaster1ConString)
	if err != nil {
		clsr.CloseAll()
		return fmt.Errorf("Ошибка миграции БД - %w", err)
	}

	return nil
}

func (a *app) initializeDBClients(cnfg *config.ApplicationParameters, _ logger.Logger, clsr closer.Closer) (db.ClientBD, error) {
	dbClient, err := db.NewClient(a.ctx, cnfg.BDMaster1ConString, cnfg.BDSync1ConString)
	if err != nil {
		clsr.CloseAll()
		return nil, fmt.Errorf("Ошибка создания клиента БД - %w", err)
	}

	clsr.Add(func() error {
		return dbClient.Close()
	})

	return dbClient, nil
}

func (a *app) initializeServices(_ *config.ApplicationParameters, l logger.Logger, dbClient db.ClientBD, clsr closer.Closer) (ordSrvc orderservice.Service, stckSrvc stockservice.Service, err error) {
	ordStrg, err := orderstore.NewOrderStorage(a.ctx, l, dbClient)
	if err != nil {
		clsr.CloseAll()
		return nil, nil, fmt.Errorf("Ошибка создания хранилища заказов - %w", err)
	}

	resStg, err := stockstore.NewReserveStorage(a.ctx, l, dbClient)
	if err != nil {
		clsr.CloseAll()
		return nil, nil, fmt.Errorf("Ошибка создания хранилища остатков - %w", err)
	}

	ordSrvc = orderservice.NewService(a.ctx, l, ordStrg, resStg)
	stckSrvc = stockservice.NewService(a.ctx, l, resStg)

	return ordSrvc, stckSrvc, nil
}

func (a *app) startGRPCServer(cnfg *config.ApplicationParameters, l logger.Logger, ordSrvc orderservice.Service, stckSrvc stockservice.Service, clsr closer.Closer) error {
	servGRPC := grpcserver.NewServer(a.ctx, l, cnfg)

	servGRPC.RegisterAPI([]grpcserver.APIHandler{
		orderapi.NewAPI(a.ctx, l, ordSrvc),
		stockapi.NewAPI(a.ctx, l, stckSrvc),
	})

	go func() {
		err := servGRPC.Start()
		if err != nil {
			clsr.CloseAll()
			l.Errorf(a.ctx, "Ошибка запуска GRPC-сервера - %v", err)
			a.cancel()
		}
	}()

	clsr.Add(func() error {
		return servGRPC.Stop()
	})

	return nil
}

func (a *app) startHTTPServer(cnfg *config.ApplicationParameters, l logger.Logger, clsr closer.Closer) error {
	servHTTP, err := httpserver.NewServer(a.ctx, l, cnfg)
	if err != nil {
		clsr.CloseAll()
		return fmt.Errorf("Ошибка создания экземпляра HTTP-сервера - %w", err)
	}

	err = servHTTP.RegisterAPI()
	if err != nil {
		clsr.CloseAll()
		return fmt.Errorf("Ошибка регистрации прикладов LOMS в Proxy-GRPC - %w", err)
	}

	err = servHTTP.Start()
	if err != nil {
		clsr.CloseAll()
		return fmt.Errorf("Ошибка запуска HTTP-сервера - %w", err)
	}

	clsr.Add(func() error {
		return servHTTP.Stop(cnfg.GracefulTimeout)
	})

	return nil
}
