package interceptor

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// Logger кастомный интерсептор по сбору логов запросов.
func (i *Interceptor) Logger() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req any, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		raw, _ := protojson.Marshal((req).(proto.Message))
		i.logger.Debugf(ctx, "Тип реквеста в сервис LOMS: %v, состав - %v", method, string(raw))

		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
