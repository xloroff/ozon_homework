package stockservice

import (
	"context"
	"fmt"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/pkg/tracer"
)

// Info возвращает доступное число остатков товара за вычетом зарезервированного.
func (s *service) Info(ctx context.Context, sku int64) (uint16, error) {
	ctx, span := tracer.StartSpanFromContext(ctx, "service.stockservice.info")
	span.SetTag("component", "stockservice")
	defer span.End()

	c, err := s.stockStore.GetAvailableForReserve(ctx, sku)
	if err != nil {
		span.SetTag("error", true)
		s.logger.Debugf(s.ctx, "StockService.Info: Ошибка получения остатков товара - %v", err)

		return 0, fmt.Errorf("StockService.Info: Ошибка получения остатков товара - %w", err)
	}

	return c, nil
}
