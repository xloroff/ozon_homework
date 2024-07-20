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
	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/pkg/kafka_producer"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/pkg/logger"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/pkg/tracer"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/repository/order_store"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/repository/outbox_store"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/repository/stock_store"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/service/order"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/service/outbox"
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
func (a *app) Run() error { // nolint:revive // Тут запуски, никакой сложности для восприятия не несут.
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

	err = a.initializeTracer(cnfg.JaegerSettings, clsr)
	if err != nil {
		return err
	}

	err = a.runMigrations(cnfg, l, clsr)
	if err != nil {
		return err
	}

	pr, err := a.initializeOrderProducer(cnfg.KafkaSettings, l, clsr)
	if err != nil {
		return err
	}

	dbClient, err := a.initializeDBClient(cnfg.BDConSettings, l, clsr)
	if err != nil {
		return err
	}

	outbStg, err := a.initializeOutboxStorage(cnfg.KafkaSettings, l, dbClient, clsr)
	if err != nil {
		return err
	}

	ordStg, err := a.initializeOrderStorage(l, dbClient, outbStg, clsr)
	if err != nil {
		return err
	}

	stoStrg, err := a.initializeStocktorage(l, dbClient, clsr)
	if err != nil {
		return err
	}

	ordSrvc := orderservice.NewService(a.ctx, l, ordStg, stoStrg)
	stckSrvc := stockservice.NewService(a.ctx, l, stoStrg)

	err = a.startOutboxSenderService(cnfg.KafkaSettings, pr, outbStg, l, clsr)
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

func (a *app) initializeTracer(cnfg *config.JaegerSettings, clsr closer.Closer) error {
	err := tracer.InitTracerProvider(a.ctx, cnfg)
	if err != nil {
		return fmt.Errorf("Ошибка создания провайдера трасировки - %w", err)
	}

	clsr.Add(tracer.Close)

	return nil
}

func (a *app) initializeOrderProducer(cnfg *config.KafkaSettings, l logger.Logger, clsr closer.Closer) (kafkaproducer.Producer, error) {
	producer, err := kafkaproducer.NewProducer(cnfg.KafkaAddress, cnfg.OrderTopic, l,
		kafkaproducer.WithIdempotent(),
		kafkaproducer.WithMaxOpenRequests(cnfg.KafkaConnCount),
		kafkaproducer.WithRetryMax(cnfg.KafkaRetryCount),
		kafkaproducer.WithRetryBackoff(cnfg.KafkaBackoffTime),
	)
	if err != nil {
		return nil, fmt.Errorf("Ошибка создания продюсера в сервис Кафка - %w", err)
	}

	clsr.Add(producer.Close)

	return producer, nil
}

func (a *app) runMigrations(cnfg *config.ApplicationParameters, l logger.Logger, clsr closer.Closer) error {
	err := db.MigrationPool(a.ctx, l, cnfg.MigrationFolder, cnfg.BDMaster1ConString)
	if err != nil {
		clsr.CloseAll()
		return fmt.Errorf("Ошибка миграции БД - %w", err)
	}

	return nil
}

func (a *app) initializeDBClient(cnfg *config.BDConSettings, _ logger.Logger, clsr closer.Closer) (db.ClientBD, error) {
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

func (a *app) initializeOutboxStorage(cnfg *config.KafkaSettings, l logger.Logger, dbClient db.ClientBD, clsr closer.Closer) (outbStg outboxstore.Storage, err error) {
	outbStg, err = outboxstore.NewOutboxStorage(a.ctx, l, dbClient, cnfg.BlockTime)
	if err != nil {
		clsr.CloseAll()
		return nil, fmt.Errorf("Ошибка старта хранилища сообщений для отпарвки в брокер - %w", err)
	}

	return outbStg, nil
}

func (a *app) initializeOrderStorage(l logger.Logger, dbClient db.ClientBD, outbStg outboxstore.Storage, clsr closer.Closer) (ordStrg orderstore.Storage, err error) {
	ordStrg, err = orderstore.NewOrderStorage(a.ctx, l, dbClient, outbStg)
	if err != nil {
		clsr.CloseAll()
		return nil, fmt.Errorf("Ошибка старта хранилища заказов - %w", err)
	}

	return ordStrg, nil
}

func (a *app) initializeStocktorage(l logger.Logger, dbClient db.ClientBD, clsr closer.Closer) (stoStrg stockstore.Storage, err error) {
	stoStrg, err = stockstore.NewReserveStorage(a.ctx, l, dbClient)
	if err != nil {
		clsr.CloseAll()
		return nil, fmt.Errorf("Ошибка создания хранилища остатков - %w", err)
	}

	return stoStrg, nil
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

func (a *app) startOutboxSenderService(cnfg *config.KafkaSettings, pr kafkaproducer.Producer, outbStg outboxstore.Storage, l logger.Logger, clsr closer.Closer) error {
	outboxService := outboxservice.NewService(a.ctx, l, cnfg.KafkaSenderPeriod, outbStg, pr)
	outboxService.Start(a.ctx)

	clsr.Add(func() error {
		return outboxService.Stop()
	})

	return nil
}
