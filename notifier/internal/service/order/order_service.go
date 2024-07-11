package orderservice

import (
	"context"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/notifier/internal/pkg/logger"
)

// Service методы сервиса для работы обработки эвентов по заказам.
type Service interface {
	OrderStatusChanges(ctx context.Context, orderID int64, status string) error
}

type service struct {
	ctx    context.Context
	logger logger.Logger
}

// NewService создает новый сервис для событий заказов.
func NewService(ctx context.Context, l logger.Logger) Service {
	return &service{
		ctx:    ctx,
		logger: l,
	}
}
