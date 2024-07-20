package interceptor

import (
	"context"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/pkg/logger"
)

// Interceptor структура интерсептора для gRpc.
type Interceptor struct {
	ctx    context.Context
	logger logger.Logger
}

// NewInterceptor создает новый интерсептор.
func NewInterceptor(ctx context.Context, l logger.Logger) *Interceptor {
	return &Interceptor{
		ctx:    ctx,
		logger: l,
	}
}
