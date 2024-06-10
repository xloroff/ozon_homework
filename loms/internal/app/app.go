package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/config"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/pkg/api/orderapi"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/pkg/api/stockapi"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/pkg/grpc_server"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/pkg/http_server"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/pkg/logger"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/repository/order_store"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/repository/stock_store"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/service/order"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/service/stock"
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

	// Стартуем получение конфига.
	cnfg, err := config.LoadAPIConfig()
	if err != nil {
		return fmt.Errorf("Ошибка получения параметров из конфигурационного файла - %w", err)
	}

	// Стартуем логгер.
	l := logger.InitializeLogger(cnfg.LogLevel, cnfg.LogType)

	// Стартуем хранилище заказов.
	ordStrg, err := orderstore.NewOrderStorage(a.ctx, l)
	if err != nil {
		return fmt.Errorf("Ошибка создания хранилища заказов - %w", err)
	}

	// Стартуем хранилище остатков.
	resStg, err := stockstore.NewReserveStorage(a.ctx, l)
	if err != nil {
		return fmt.Errorf("Ошибка создания хранилища остатков - %w", err)
	}

	// Стартуем новый сервис Заказов.
	ordSrvc := orderservice.NewService(a.ctx, l, ordStrg, resStg)

	// Стартуем новый сервис остатков.
	stckSrvc := stockservice.NewService(a.ctx, l, resStg)

	// Стартуем GRPC-сервер.
	servGRPC := grpcserver.NewServer(a.ctx, l, cnfg)

	servGRPC.RegisterAPI([]grpcserver.APIHandler{
		orderapi.NewAPI(a.ctx, l, ordSrvc),
		stockapi.NewAPI(a.ctx, l, stckSrvc),
	})

	go func() {
		err = servGRPC.Start()
		if err != nil {
			l.Errorf(a.ctx, "Ошибка запуска GRPC-сервера - %v", err)
			a.cancel()

			return
		}

		<-sigChan
	}()

	// Создаем HTTP-сервер.
	servHTTP, err := httpserver.NewServer(a.ctx, l, cnfg)
	if err != nil {
		return fmt.Errorf("Ошибка создания экземпляра HTTP-сервера - %w", err)
	}

	// Прикручиваем GRPC приклады.
	err = servHTTP.RegisterAPI()
	if err != nil {
		return fmt.Errorf("Ошибка регистрации прикладов LOMS в Proxy-GRPC - %w", err)
	}

	// Стартуем HTTP-сервер.
	err = servHTTP.Start()
	if err != nil {
		return fmt.Errorf("Ошибка запуска HTTP-сервера - %w", err)
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
