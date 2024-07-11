package tracer

import (
	"context"
	"fmt"
	"net/url"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdk_trace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	otel_trace "go.opentelemetry.io/otel/trace"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/notifier/internal/config"
)

// Span кастомный тип Спанов для трасировки.
type Span interface {
	otel_trace.Span
	SetTag(key string, value any)
}

// Обертка над trace.Span, чтобы удовлетворять интерфейсу Span.
type spanWrapper struct {
	otel_trace.Span
}

var tracer otel_trace.Tracer = NewTracer()

// InitTracerProvider создает провайдера трассировок.
func InitTracerProvider(ctx context.Context, cfg *config.JaegerSettings) error {
	u, err := url.Parse("http://" + cfg.JaegerHost + ":" + cfg.JaegerPort)
	if err != nil {
		return fmt.Errorf("Ошибка форматирования url Jaeger сервиса - %w", err)
	}

	exporter, err := otlptracehttp.New(ctx, otlptracehttp.WithEndpointURL(u.String()))
	if err != nil {
		return fmt.Errorf("Ошибка создания экспортера Jaeger - %w", err)
	}

	tracerProvider := sdk_trace.NewTracerProvider(
		sdk_trace.WithBatcher(exporter),
		sdk_trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(config.AppName),
		)),
	)

	otel.SetTracerProvider(tracerProvider)

	return nil
}

// NewTracer создает и возвращает Tracer.
func NewTracer() otel_trace.Tracer {
	return otel.Tracer("default")
}

// StartSpanFromContext создает новый span из контекста и возвращает его.
func StartSpanFromContext(ctx context.Context, name string, opts ...otel_trace.SpanStartOption) (context.Context, Span) {
	ctx, span := tracer.Start(ctx, name, opts...)
	return ctx, &spanWrapper{span}
}

// GetTraceID возвращает trace ID из контекста.
func GetTraceID(ctx context.Context) string {
	spanCtx := otel_trace.SpanContextFromContext(ctx)
	if spanCtx.HasTraceID() {
		return spanCtx.TraceID().String()
	}

	return ""
}

// GetSpanID возвращает span ID из контекста.
func GetSpanID(ctx context.Context) string {
	spanCtx := otel_trace.SpanContextFromContext(ctx)
	if spanCtx.HasSpanID() {
		return spanCtx.SpanID().String()
	}

	return ""
}

// SetTag добавляет тэг в спан.
func (s *spanWrapper) SetTag(key string, value any) {
	var attr attribute.KeyValue
	switch v := value.(type) {
	case string:
		attr = attribute.String(key, v)
	case int:
		attr = attribute.Int(key, v)
	case int64:
		attr = attribute.Int64(key, v)
	case float64:
		attr = attribute.Float64(key, v)
	case bool:
		attr = attribute.Bool(key, v)
	default:
		attr = attribute.String(key, fmt.Sprintf("%v", v))
	}
	s.SetAttributes(attr)
}

// StartSpanFromIDs замещает trace_id и span_id взятыми извне в контексту и возвращает его.
func StartSpanFromIDs(ctx context.Context, traceIDStr, spanIDStr, name string, opts ...otel_trace.SpanStartOption) (context.Context, Span, error) {
	if traceIDStr != "" {
		spanContext := otel_trace.SpanContextConfig{
			TraceFlags: otel_trace.FlagsSampled,
			Remote:     true,
		}

		traceID, err := otel_trace.TraceIDFromHex(traceIDStr)
		if err != nil {
			return nil, nil, fmt.Errorf("Ошибка парсинга TraceID - %w", err)
		}

		spanID, err := otel_trace.SpanIDFromHex(spanIDStr)
		if err != nil {
			return nil, nil, fmt.Errorf("Ошибка парсинга SpanID - %w", err)
		}

		spanContext.TraceID = traceID
		spanContext.SpanID = spanID

		ctx = otel_trace.ContextWithSpanContext(ctx, otel_trace.NewSpanContext(spanContext))
	}

	ctx, span := StartSpanFromContext(ctx, name, opts...)

	return ctx, span, nil
}

// Close закрывает трейсер.
func Close() error {
	var tracerProvider *sdk_trace.TracerProvider
	tracerProvider, _ = otel.GetTracerProvider().(*sdk_trace.TracerProvider)

	err := tracerProvider.ForceFlush(context.Background())
	if err != nil {
		return fmt.Errorf("Ошибка экспорта остатков span - %w", err)
	}

	err = tracerProvider.Shutdown(context.Background())
	if err != nil {
		return fmt.Errorf("Ошибка закрытия провайдера трейсов - %w", err)
	}

	return nil
}
