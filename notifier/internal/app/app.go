package app

import (
	"context"
	"fmt"
	"os"
	"syscall"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/notifier/internal/config"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/notifier/internal/consumers/order"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/notifier/internal/pkg/closer"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/notifier/internal/pkg/kafka_consumer_group"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/notifier/internal/pkg/logger"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/notifier/internal/pkg/tracer"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/notifier/internal/service/order"
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

	err = a.initializeTracer(cnfg.JaegerSettings, clsr)
	if err != nil {
		return err
	}

	oS, err := a.initializeServices(cnfg.AppSettings, l, clsr)
	if err != nil {
		return err
	}

	err = a.initializeConsumerGroup(cnfg.KafkaSettings, l, clsr, oS)
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

func (a *app) initializeServices(_ *config.AppSettings, l logger.Logger, _ closer.Closer) (oS orderservice.Service, err error) {
	oS = orderservice.NewService(a.ctx, l)
	return oS, nil
}

func (a *app) initializeConsumerGroup(cnfg *config.KafkaSettings, l logger.Logger, clsr closer.Closer, oS orderservice.Service) (err error) {
	c := orderconsumer.NewConsumer(a.ctx, l, oS)

	cg, err := kafkaconsumergroup.NewConsumerGroup(
		a.ctx,
		l,
		cnfg,
		c,
		kafkaconsumergroup.WithAutoCommit(cnfg.KafkaAutocommitEnabled),
		kafkaconsumergroup.WithOffsetsAutoCommitInterval(cnfg.KafkaAutoCommitTime),
		kafkaconsumergroup.WithHeartbeatInterval(cnfg.KafkaHeartbeatTime),
		kafkaconsumergroup.WithRebalanceTimeout(cnfg.KafkaRebalanceTime),
		kafkaconsumergroup.WithSessionTimeout(cnfg.KafkaSessionTimeout),
	)
	if err != nil {
		return fmt.Errorf("Ошибка создания группы консюмеров брокера сообщений - %w", err)
	}

	go func() {
		err := cg.Start()
		if err != nil {
			clsr.CloseAll()
			return
		}
	}()

	l.Warn(a.ctx, "Сервис запущен...")

	clsr.Add(cg.Stop)

	return nil
}
