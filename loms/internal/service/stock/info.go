package stockservice

import "fmt"

// Info возвращает доступное число остатков товара за вычетом зарезервированного.
func (s *sService) Info(sku int64) (uint16, error) {
	c, err := s.stockStore.GetAvailableForReserve(sku)
	if err != nil {
		s.logger.Debugf(s.ctx, "StockService.Info: Ошибка получения остатков товара - %v", err)

		return 0, fmt.Errorf("StockService.Info: Ошибка получения остатков товара - %w", err)
	}

	return c, nil
}
