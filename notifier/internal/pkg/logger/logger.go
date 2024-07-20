package logger

import (
	"context"
	"fmt"
	"io"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/notifier/internal/config"
)

// Logger интерфейс взаимодействия с логгером.
type Logger interface {
	Info(ctx context.Context, m string)
	Infof(ctx context.Context, m string, args ...any)
	Warn(ctx context.Context, m string)
	Warnf(ctx context.Context, m string, args ...any)
	Error(ctx context.Context, m string)
	Errorf(ctx context.Context, m string, args ...any)
	Panic(ctx context.Context, m string)
	Panicf(ctx context.Context, m string, args ...any)
	Fatal(ctx context.Context, m string)
	Fatalf(ctx context.Context, m string, args ...any)
	Debug(ctx context.Context, m string)
	Debugf(ctx context.Context, m string, args ...any)
	Close() error
}

// logger кастомный тип для стилизации ошибок записываемых в лог.
type logger struct {
	*zap.Logger
}

// InitializeLogger создает логгер.
func InitializeLogger(logLevel string, lType int) Logger {
	configL := zap.NewProductionEncoderConfig()
	configL.EncodeTime = zapcore.ISO8601TimeEncoder
	defaultLogLevel := DebugLevel(logLevel)
	configL.TimeKey = "timestamp"
	jsonEncoder := zapcore.NewJSONEncoder(configL)

	var core zapcore.Core

	switch lType {
	case 1:
		// Сливаем логи в /dev/null.
		core = zapcore.NewTee(
			zapcore.NewCore(jsonEncoder, zapcore.AddSync(io.Discard), defaultLogLevel),
		)
	case 2:
		// Создаем запись логов в Stdout.
		core = zapcore.NewTee(
			zapcore.NewCore(jsonEncoder, zapcore.AddSync(os.Stdout), defaultLogLevel),
		)
	default:
		// Создаем запись логов в Stdout.
		core = zapcore.NewTee(
			zapcore.NewCore(jsonEncoder, zapcore.AddSync(os.Stdout), defaultLogLevel),
		)
	}

	l := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	l = l.With(zap.String("app", config.AppName))

	return &logger{l}
}

// Close останавливает логирование и синхронизирует файлы (записывает накопленные ошибки в файлы и выводы).
func (l *logger) Close() error {
	if l != nil {
		err := l.Sync()
		if err != nil {
			return fmt.Errorf("Ошибка синхронизации логов - %w", err)
		}
	}

	return nil
}
