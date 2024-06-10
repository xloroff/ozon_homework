package interceptor

import (
	"context"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/pkg/logger"
)

// Interceptor структура интерсептора для gRpc.
type Interceptor struct {
	ctx    context.Context
	logger logger.ILog
}

// NewInterceptor создает новый интерсептор.
func NewInterceptor(ctx context.Context, l logger.ILog) *Interceptor {
	return &Interceptor{
		ctx:    ctx,
		logger: l,
	}
}
