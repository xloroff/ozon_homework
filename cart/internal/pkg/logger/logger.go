package logger

import (
	"context"
	"fmt"
	"io"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ILog интерфейс взаимодействия с логгером.
type ILog interface {
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

// Logger кастомный тип для стилизации ошибок записываемых в лог.
type lLog struct {
	*zap.Logger
}

// InitializeLogger создает логгер.
func InitializeLogger(logLevel string, lType int) ILog {
	configL := zap.NewProductionEncoderConfig()
	configL.EncodeTime = zapcore.ISO8601TimeEncoder
	consoleEncoder := zapcore.NewConsoleEncoder(configL)
	defaultLogLevel := DebugLevel(logLevel)

	var core zapcore.Core

	switch lType {
	case 1:
		// Сливаем логи в /dev/null.
		core = zapcore.NewTee(
			zapcore.NewCore(consoleEncoder, zapcore.AddSync(io.Discard), defaultLogLevel),
		)
	case 2:
		// Создаем запись логов в Stdout.
		core = zapcore.NewTee(
			zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), defaultLogLevel),
		)
	default:
		// Создаем запись логов в Stdout.
		core = zapcore.NewTee(
			zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), defaultLogLevel),
		)
	}

	return &lLog{zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))}
}

// Close останавливает логирование и синхронизирует файлы (записывает накопленные ошибки в файлы и выводы).
func (l *lLog) Close() error {
	if l != nil {
		err := l.Sync()
		if err != nil {
			return fmt.Errorf("Ошибка синхронизации логов - %w", err)
		}
	}

	return nil
}
