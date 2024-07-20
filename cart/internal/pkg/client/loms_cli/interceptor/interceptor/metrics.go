package interceptor

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/pkg/metrics"
)

// Metrics кастомный интерсептор по сбору метрик.
func (i *Interceptor) Metrics(urlStr, host string) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req any, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		metrics.UpdateExternalRequestsTotal(
			urlStr,
			method,
		)
		defer metrics.UpdateExternalResponseTime(time.Now().UTC())

		err := invoker(ctx, method, req, reply, cc, opts...)
		if err != nil {
			st, _ := status.FromError(err)
			metrics.UpdateExternalResponseCode(
				host,
				method,
				st.Code().String(),
			)

			return err
		}

		metrics.UpdateExternalResponseCode(
			host,
			method,
			codes.OK.String(),
		)

		return nil
	}
}
