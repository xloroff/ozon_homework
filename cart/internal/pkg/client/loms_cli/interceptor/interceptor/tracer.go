package interceptor

import (
	"context"

	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/pkg/tracer"
)

// Tracer кастомный интерсептор по добавлению метаданных трасировки.
func (i *Interceptor) Tracer() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req any, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		ctx, span := tracer.StartSpanFromContext(ctx, "lomscli.interceptor", trace.WithSpanKind(trace.SpanKindClient))
		span.SetTag("component", "interceptor")
		span.SetTag("method", method)
		span.SetTag("peer.address", cc.Target())

		defer span.End()

		ctx = metadata.AppendToOutgoingContext(ctx, "x-trace-id", tracer.GetTraceID(ctx))
		ctx = metadata.AppendToOutgoingContext(ctx, "x-span-id", tracer.GetSpanID(ctx))

		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
