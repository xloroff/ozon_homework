package logger

import (
	"context"

	"go.uber.org/zap"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/pkg/tracer"
)

// contextLoggingKey название ключа контекста в котором будут храниться дополнительные параметры для zap.
type key string

const contextLoggingKey key = "zap_with_ctx"

// AddFieldsToContext добавляет произвольное количество пар ключ-значение в контекст.
func AddFieldsToContext(ctx context.Context, keyValues ...any) context.Context {
	fields := make([]zap.Field, 0, len(keyValues)/2)

	for i := 0; i < len(keyValues)-1; i += 2 {
		key, ok := keyValues[i].(string)
		if !ok {
			continue
		}

		value := keyValues[i+1]
		fields = append(fields, zap.Any(key, value))
	}

	return Append(ctx, fields)
}

// Set записывает поля в контекст по ключу contextLoggingKey.
func Set(ctx context.Context, fields []zap.Field) context.Context {
	return context.WithValue(ctx, contextLoggingKey, fields)
}

// Append добавляет поля к контексту с нужным ключом contextLoggingKey.
func Append(ctx context.Context, fields []zap.Field) context.Context {
	if loggerFields, ok := ctx.Value(contextLoggingKey).([]zap.Field); ok {
		fields = append(fields, loggerFields...)
	}

	return context.WithValue(ctx, contextLoggingKey, fields)
}

// getAllLoggerFields получение полей zap из контекста по ключу contextLoggingKey.
func getAllLoggerFields(ctx context.Context) []zap.Field {
	fields := make([]zap.Field, 0, 2)
	traceID := tracer.GetTraceID(ctx)

	if traceID != "" {
		fields = append(fields, zap.String("trace_id", traceID))
	}

	spanID := tracer.GetSpanID(ctx)
	if spanID != "" {
		fields = append(fields, zap.String("span_id", spanID))
	}

	if loggerFields, ok := ctx.Value(contextLoggingKey).([]zap.Field); ok {
		fields = append(fields, loggerFields...)
	}

	return fields
}
