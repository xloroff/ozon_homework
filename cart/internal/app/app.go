package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/internal/config"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/internal/initilize"
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
	// Грузим базовые настройки приложения.
	configAPI, err := initilize.LoadApiConfig(config.ConfigDirPath, config.ConfigType, config.AppConfigName)
	if err != nil {
		return fmt.Errorf("Загрузка конфигурационного файла приложения завершилась ошибкой - %w", err.Error())
	}

	// Стартуем кастомные логи.
	if err = logger.InitializeLogger(configAPI.LogLevel, configAPI.LogFolder); err != nil {
		return fmt.Errorf("Инициализации функции логирования завершилась ошибкой - %w", err.Error())
	}

	// Канал для сигналов завершения.
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	serv := server.NewServer(a.ctx, &configAPI)
	err = serv.Start()
	if err != nil {
		logger.Fatalf(a.ctx, "Ошибка запуска сервера - %w", err.Error())
	}

	// Обработка сигналов завершения.
	go func() {
		<-sigChan
		defer logger.Close()
		defer a.cancel()

		logger.Warnf(a.ctx, "Получен сигнал завершения, остановка приложения произведена...")
	}()

	// Блокировка до завершения контекста.
	<-a.ctx.Done()

	return nil
}