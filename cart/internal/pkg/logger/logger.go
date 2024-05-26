package logger

import (
	"context"
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger кастомный тип для стилизации ошибок записываемых в лог.
type Logger struct {
	*zap.Logger
}

var logger *Logger

// InitializeLogger создает логгер.
func InitializeLogger(logLevel, folder string) (err error) {
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	fileEncoder := zapcore.NewJSONEncoder(config)
	consoleEncoder := zapcore.NewConsoleEncoder(config)

	// Создаем запись логов в файл.
	writer, err := LogWriter(folder)
	if err != nil {
		return fmt.Errorf("Ошибка создания направления записи логов - %w.", err)
	}
	defaultLogLevel := DebugLevel(logLevel)
	// Стартуем с двумя ядрами одно будет складывать в файл другое в Stdout.
	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, writer, defaultLogLevel),
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), defaultLogLevel),
	)

	logger = &Logger{zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))}

	return nil
}

func Close() error {
	if logger != nil {
		err := logger.Sync()
		if err != nil {
			return fmt.Errorf("Ошибка синхронизации логов - %w.", err)
		}
	}
	return nil
}

// Debug соответствует уровню дебага Debug.
func (l *Logger) Debug(ctx context.Context, m string) {
	l.Logger.Debug(m, getAllLoggerFields(ctx)...)
}

// Debugf соответствует уровню дебага Debug.
func (l *Logger) Debugf(ctx context.Context, m string, args ...any) {
	l.Logger.Debug(fmt.Sprintf(m, args...), getAllLoggerFields(ctx)...)
}

// Info соответствует уровню дебага Info.
func (l *Logger) Info(ctx context.Context, m string) {
	l.Logger.Info(m, getAllLoggerFields(ctx)...)
}

// Infof с доп полями - соответствует уровню дебага Info.
func (l *Logger) Infof(ctx context.Context, m string, args ...any) {
	l.Logger.Info(fmt.Sprintf(m, args...), getAllLoggerFields(ctx)...)
}

// Warn соответствует уровню дебага Warn.
func (l *Logger) Warn(ctx context.Context, m string) {
	l.Logger.Warn(m, getAllLoggerFields(ctx)...)
}

// Warnf с доп полями - соответствует уровню дебага Warn.
func (l *Logger) Warnf(ctx context.Context, m string, args ...any) {
	l.Logger.Warn(fmt.Sprintf(m, args...), getAllLoggerFields(ctx)...)
}

// Error соответствует уровню дебага Error.
func (l *Logger) Error(ctx context.Context, m string) {
	l.Logger.Error(m, getAllLoggerFields(ctx)...)
}

// Errorf с доп полями - соответствует уровню дебага Error.
func (l *Logger) Errorf(ctx context.Context, m string, args ...any) {
	l.Logger.Warn(fmt.Sprintf(m, args...), getAllLoggerFields(ctx)...)
}

// Panic соответствует уровню дебага Panic.
func (l *Logger) Panic(ctx context.Context, m string) {
	l.Logger.Panic(m, getAllLoggerFields(ctx)...)
}

// Panicf с доп полями - соответствует уровню дебага Panic.
func (l *Logger) Panicf(ctx context.Context, m string, args ...any) {
	l.Logger.Panic(fmt.Sprintf(m, args...), getAllLoggerFields(ctx)...)
}

// Fatal соответствует уровню дебага Fatal.
func (l *Logger) Fatal(ctx context.Context, m string) {
	l.Logger.Fatal(m, getAllLoggerFields(ctx)...)
}

// Fatalf с доп полями - соответствует уровню дебага Fatal.
func (l *Logger) Fatalf(ctx context.Context, m string, args ...any) {
	l.Logger.Fatal(fmt.Sprintf(m, args...), getAllLoggerFields(ctx)...)
}

func Info(ctx context.Context, m string) {
	logger.Info(ctx, m)
}

func Infof(ctx context.Context, m string, args ...any) {
	logger.Infof(ctx, m, args...)
}

func Warn(ctx context.Context, m string) {
	logger.Info(ctx, m)
}

func Warnf(ctx context.Context, m string, args ...any) {
	logger.Warnf(ctx, m, args...)
}

func Error(ctx context.Context, m string) {
	logger.Error(ctx, m)
}

func Errorf(ctx context.Context, m string, args ...any) {
	logger.Errorf(ctx, m, args...)
}

func Panic(ctx context.Context, m string) {
	logger.Panic(ctx, m)
}

func Panicf(ctx context.Context, m string, args ...any) {
	logger.Panicf(ctx, m, args...)
}

func Fatal(ctx context.Context, m string) {
	logger.Fatal(ctx, m)
}

func Fatalf(ctx context.Context, m string, args ...any) {
	logger.Fatalf(ctx, m, args...)
}

func Debug(ctx context.Context, m string) {
	logger.Debug(ctx, m)
}

func Debugf(ctx context.Context, m string, args ...any) {
	logger.Debugf(ctx, m, args...)
}