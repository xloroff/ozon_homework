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
	cnfg, err := config.LoadAPIConfig()
	if err != nil {
		return fmt.Errorf("Ошибка получения параметров из конфигурационного файла - %w", err)
	}

	// Стартуем логгер.
	l := logger.InitializeLogger(cnfg.LogLevel, cnfg.LogType)

	clsr := closer.NewCloser(l, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer clsr.Wait()

	clsr.Add(func() error {
		return l.Close()
	})

	// Запускаем миграции перед стартом сервиса.
	err = db.MigrationPool(a.ctx, l, cnfg.MigrationFolder, cnfg.BDMaster1ConString)
	if err != nil {
		clsr.CloseAll()

		return fmt.Errorf("Ошибка миграции БД - %w", err)
	}

	// Стартуем репозитории/связь с БД и хранилища.
	dbClient, err := db.NewClient(a.ctx, cnfg.BDMaster1ConString, cnfg.BDSync1ConString)
	if err != nil {
		clsr.CloseAll()

		return fmt.Errorf("Ошибка создания клиента БД - %w", err)
	}

	clsr.Add(func() error {
		return dbClient.Close()
	})

	ordStrg, err := orderstore.NewOrderStorage(a.ctx, l, dbClient)
	if err != nil {
		clsr.CloseAll()

		return fmt.Errorf("Ошибка создания хранилища заказов - %w", err)
	}

	resStg, err := stockstore.NewReserveStorage(a.ctx, l, dbClient)
	if err != nil {
		clsr.CloseAll()

		return fmt.Errorf("Ошибка создания хранилища остатков - %w", err)
	}

	ordSrvc := orderservice.NewService(a.ctx, l, ordStrg, resStg)

	stckSrvc := stockservice.NewService(a.ctx, l, resStg)

	servGRPC := grpcserver.NewServer(a.ctx, l, cnfg)

	servGRPC.RegisterAPI([]grpcserver.APIHandler{
		orderapi.NewAPI(a.ctx, l, ordSrvc),
		stockapi.NewAPI(a.ctx, l, stckSrvc),
	})

	go func() {
		err = servGRPC.Start()
		if err != nil {
			clsr.CloseAll()

			l.Errorf(a.ctx, "Ошибка запуска GRPC-сервера - %v", err)
			a.cancel()

			return
		}
	}()

	clsr.Add(func() error {
		return servGRPC.Stop()
	})

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
