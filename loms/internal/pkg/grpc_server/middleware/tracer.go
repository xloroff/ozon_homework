package middleware

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/pkg/tracer"
)

// Tracer создает новый трейс для ручки.
func Tracer(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	md, _ := metadata.FromIncomingContext(ctx)

	var traceID string

	var spanID string

	if traceIDs, exists := md["x-trace-id"]; exists {
		traceID = traceIDs[0]
	}

	if spanIDs, exists := md["x-span-id"]; exists {
		spanID = spanIDs[0]
	}

	if spanID != "" {
		var err error

		spanContext := trace.SpanContextConfig{
			TraceFlags: trace.FlagsSampled,
			Remote:     true,
		}

		spanContext.TraceID, err = trace.TraceIDFromHex(traceID)
		if err != nil {
			return nil, fmt.Errorf("Ошибка получения TraceID из запроса - %w", err)
		}

		spanContext.SpanID, err = trace.SpanIDFromHex(spanID)
		if err != nil {
			return nil, fmt.Errorf("Ошибка получения SpanID из запроса - %w", err)
		}

		ctx = trace.ContextWithSpanContext(ctx, trace.NewSpanContext(spanContext))
	}

	ctx, span := tracer.StartSpanFromContext(ctx, "grpcserver.middleware", trace.WithSpanKind(trace.SpanKindServer))
	span.SetTag("component", "middleware")
	span.SetTag("method", info.FullMethod)

	defer span.End()

	return handler(ctx, req)
}
