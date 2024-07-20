package middleware

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/pkg/metrics"
)

// Metrics обновляет метрики сервера.
func Metrics(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	metrics.UpdateRequestsTotal(info.FullMethod)

	defer metrics.UpdateResponseTime(time.Now().UTC())

	resp, err = handler(ctx, req)
	if err != nil {
		st, _ := status.FromError(err)
		metrics.UpdateResponseCode(info.FullMethod, st.Code().String())
	} else {
		metrics.UpdateResponseCode(info.FullMethod, codes.OK.String())
	}

	return resp, err
}
