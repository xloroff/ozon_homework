package interceptor

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Validate валидация входящих данных.
func (i *Interceptor) Validate() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ any, err error) {
		if v, ok := req.(interface{ Validate() error }); ok {
			if err = v.Validate(); err != nil {
				i.logger.Errorf(ctx, "Ошибка валидации параметров запроса %v - %v", info.FullMethod, err)

				return nil, status.Error(codes.InvalidArgument, err.Error())
			}
		}

		return handler(ctx, req)
	}
}
