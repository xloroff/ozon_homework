package logger

import (
	"context"
	"fmt"
)

// Info соответствует уровню дебага Info, принимает только одно поле.
func (l *lLog) Info(ctx context.Context, m string) {
	l.Logger.Info(m, getAllLoggerFields(ctx)...)
}

// Infof соответствует уровню дебага Info, принимает дополнительные аргументы к текстовому полю.
func (l *lLog) Infof(ctx context.Context, m string, args ...any) {
	l.Logger.Info(fmt.Sprintf(m, args...), getAllLoggerFields(ctx)...)
}

// Warn соответствует уровню дебага Warn, принимает только одно поле.
func (l *lLog) Warn(ctx context.Context, m string) {
	l.Logger.Warn(m, getAllLoggerFields(ctx)...)
}

// Warnf соответствует уровню дебага Warn, принимает дополнительные аргументы к текстовому полю.
func (l *lLog) Warnf(ctx context.Context, m string, args ...any) {
	l.Logger.Warn(fmt.Sprintf(m, args...), getAllLoggerFields(ctx)...)
}

// Error соответствует уровню дебага Error, принимает только одно поле.
func (l *lLog) Error(ctx context.Context, m string) {
	l.Logger.Error(m, getAllLoggerFields(ctx)...)
}

// Errorf соответствует уровню дебага Error, принимает дополнительные аргументы к текстовому полю.
func (l *lLog) Errorf(ctx context.Context, m string, args ...any) {
	l.Logger.Warn(fmt.Sprintf(m, args...), getAllLoggerFields(ctx)...)
}

// Panic соответствует уровню дебага Panic, принимает только одно поле и инициирует panic.
func (l *lLog) Panic(ctx context.Context, m string) {
	l.Logger.Panic(m, getAllLoggerFields(ctx)...)
}

// Panicf соответствует уровню дебага Panic, принимает дополнительные аргументы к текстовому полю и инициирует panic.
func (l *lLog) Panicf(ctx context.Context, m string, args ...any) {
	l.Logger.Panic(fmt.Sprintf(m, args...), getAllLoggerFields(ctx)...)
}

// Fatal соответствует уровню дебага Fatal, принимает только одно поле.
func (l *lLog) Fatal(ctx context.Context, m string) {
	l.Logger.Fatal(m, getAllLoggerFields(ctx)...)
}

// Fatalf соответствует уровню дебага Fatal, принимает дополнительные аргументы к текстовому полю.
func (l *lLog) Fatalf(ctx context.Context, m string, args ...any) {
	l.Logger.Fatal(fmt.Sprintf(m, args...), getAllLoggerFields(ctx)...)
}

// Debug соответствует уровню дебага Debug, принимает только одно поле.
func (l *lLog) Debug(ctx context.Context, m string) {
	l.Logger.Debug(m, getAllLoggerFields(ctx)...)
}

// Debugf соответствует уровню дебага Debug, принимает дополнительные аргументы к текстовому полю.
func (l *lLog) Debugf(ctx context.Context, m string, args ...any) {
	l.Logger.Debug(fmt.Sprintf(m, args...), getAllLoggerFields(ctx)...)
}