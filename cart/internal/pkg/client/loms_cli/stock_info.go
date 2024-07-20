package lomscli

import (
	"context"
	"fmt"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/pb/api/stock/v1"
)

// StockInfo получает остатки товара из сервиса LOMS.
func (c *lomsClient) StockInfo(ctx context.Context, skuID int64) (uint16, error) {
	c.logger.Debugf(ctx, "LomsCli.StockInfo: начал обращение в сервис LOMS, получение остатков товара - %v", skuID)
	defer c.logger.Debugf(ctx, "LomsCli.StockInfo: закончил обращение в сервис LOMS, получение остатков товара - %v", skuID)

	resp, err := c.stock.Info(ctx, &stock.StockInfoRequest{
		Sku: skuID,
	})
	if err != nil {
		c.logger.Errorf(ctx, "LomsCli.StockInfo: Ошибка получения остатков товара %v - %v", skuID, err)

		return 0, fmt.Errorf("LomsCli.StockInfo: Ошибка получения остатков товара %v - %w", skuID, err)
	}

	return uint16(resp.GetCount()), nil
}
