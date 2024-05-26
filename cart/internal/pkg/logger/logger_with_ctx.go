package logger

import (
	"context"

	"go.uber.org/zap"
)

// contextLoggingKey название ключа контекста в котором будут хранится дополнительные параметры для zap.
const contextLoggingKey = "zap_with_ctx"

// Set записывает поля в контекст по ключу contextLoggingKey.
func Set(ctx context.Context, fields []zap.Field) context.Context {
	return context.WithValue(ctx, contextLoggingKey, fields)
}

// Append добавляет поля к контексту с нужным ключем contextLoggingKey.
func Append(ctx context.Context, fields []zap.Field) context.Context {
	if loggerFields, ok := ctx.Value(contextLoggingKey).([]zap.Field); ok {
		fields = append(fields, loggerFields...)
	}

	return context.WithValue(ctx, contextLoggingKey, fields)
}

// getAllLoggerFields получение полей zap из контекста по ключу contextLoggingKey.
func getAllLoggerFields(ctx context.Context) []zap.Field {
	if loggerFields, ok := ctx.Value(contextLoggingKey).([]zap.Field); ok {
		return loggerFields
	}
	return nil
}

// getLoggerField получение определенного параметра из контекста по ключу.
func getLoggerField(ctx context.Context, key string) (field zap.Field) {
	if loggerFields, ok := ctx.Value(contextLoggingKey).([]zap.Field); ok {
		for _, field := range loggerFields {
			if field.Key == key {
				return field
			}
		}
	}
	return
}