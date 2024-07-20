package lomscli

import (
	"fmt"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/pb/api/stock/v1"
)

// StockInfo получает остатки товара из сервиса LOMS.
func (c *lomsClient) StockInfo(skuID int64) (uint16, error) {
	c.logger.Debugf(c.ctx, "LomsCli.StockInfo: начал обращение в сервис LOMS, получение остатков товара - %v", skuID)
	defer c.logger.Debugf(c.ctx, "LomsCli.StockInfo: закончил обращение в сервис LOMS, получение остатков товара - %v", skuID)

	resp, err := c.stock.Info(c.ctx, &stock.StockInfoRequest{
		Sku: skuID,
	})
	if err != nil {
		c.logger.Errorf(c.ctx, "LomsCli.StockInfo: Ошибка получения остатков товара %v - %v", skuID, err)

		return 0, fmt.Errorf("LomsCli.StockInfo: Ошибка получения остатков товара %v - %w", skuID, err)
	}

	return uint16(resp.GetCount()), nil
}
