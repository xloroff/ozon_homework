package interceptor

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Panic обрабатываем панику.
func (i *Interceptor) Panic() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ any, err error) {
		defer func() {
			if e := recover(); e != nil {
				i.logger.Errorf(ctx, "panic %v, ошибка - %v", info.FullMethod, e)

				err = status.Errorf(codes.Internal, "panic: %v", e)
			}
		}()

		return handler(ctx, req)
	}
}
