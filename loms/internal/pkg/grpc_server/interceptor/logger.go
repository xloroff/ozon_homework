package interceptor

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// Logger кастомный интерсептор по сбору логов запросов.
func (i *Interceptor) Logger() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		raw, _ := protojson.Marshal((req).(proto.Message))

		i.logger.Debugf(ctx, "Тип реквеста: %v, состав - %v", info.FullMethod, string(raw))

		return handler(ctx, req)
	}
}
