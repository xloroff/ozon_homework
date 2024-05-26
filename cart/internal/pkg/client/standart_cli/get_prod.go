package standartcli

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/internal/initilize"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/internal/model/v1"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/internal/pkg/logger"
)

// ProductReciver обращается на сервис продуктов 1 попытка.
func (c *client) ProductReciver(ctx context.Context, settings *initilize.ConfigAPI, skuID int64) (*v1.ProductResp, error) {
	logger.Debug(ctx, fmt.Sprintf("clientProductServiceStandart.ProductReciver: начало обращения к сервису продуктов, продукт - %d", skuID))
	defer logger.Debug(ctx, fmt.Sprintf("clientProductServiceStandart.ProductReciver: конец обращения к сервису продуктов, продукт - %d", skuID))

	/* TODO Проверяем ходили мы к товару или нет (кэш)
	resp, ok := productStorage[skuID]
	if ok {
		logger.Debug(ctx, fmt.Sprintf("clientProductServiceStandart.ProductReciver: товар был получен ранее, обращение в сервис продуктов не требуется - %d", skuID))
		return resp, nil
	}
	*/
	req := &v1.ProductReq{
		Token: settings.ProductServiceToken,
		Sku:   skuID,
	}

	jsonReq, err := json.Marshal(req)
	if err != nil {
		logger.Errorf(ctx, "clientProductServiceStandart.ProductReciver: ошибка формирования запроса в сервис продуктов - %w", err)
		return nil, fmt.Errorf("Ошибка формирования запроса в сервис продуктов - %w", err)
	}

	httpReq, err := http.NewRequest("POST", fmt.Sprintf("%s/get_product", settings.ProductServiceHost), bytes.NewBuffer(jsonReq))
	if err != nil {
		logger.Errorf(ctx, "clientProductServiceStandart.ProductReciver: ошибка создания нового запроса в сервис продуктов - %w", err)
		return nil, fmt.Errorf("Ошибка создания нового запроса в сервис продуктов - %w", err)
	}

	httpReq = httpReq.WithContext(ctx)
	cli := http.DefaultClient
	httpResp, err := cli.Do(httpReq)
	if err != nil {
		logger.Errorf(ctx, "clientProductServiceStandart.ProductReciver: ошибка отправки запроса в сервис продуктов - %w", err)
		return nil, fmt.Errorf("Ошибка отправки запроса в сервис продуктов - %w", err)
	}
	defer httpResp.Body.Close()
	ctx = logger.Set(ctx, []zap.Field{zap.Any("request", req)})
	logger.Debugf(ctx, "clientProductServiceStandart.ProductReciver: запрос в сервис продуктов")

	return responseChecker(ctx, httpResp)
}

// GetProduct получается продукты с сервиса продуктов - делает ретраи.
func (c *client) GetProduct(ctx context.Context, settings *initilize.ConfigAPI, skuID int64) (*v1.ProductResp, error) {
	retrCounter := settings.ProductServiceRetr
	// Начинаем попытки связи с сервисом продуктов
	for i := 1; i <= retrCounter; i++ {
		logger.Infof(ctx, fmt.Sprintf("clientProductServiceStandart.GetProduct: попытка связи %d получение продукта %d - хост %s", i, skuID, settings.ProductServiceHost))

		resp, err := c.ProductReciver(ctx, settings, skuID)
		switch err {
		case nil:
			ctx = logger.Set(ctx, []zap.Field{zap.Any("response", resp)})
			logger.Debugf(ctx, "clientProductServiceStandart.GetProduct: ответ сервиса продуктов")
			return resp, nil
		case ErrTooManyRequests:
			time.Sleep(100 * time.Millisecond)

			if i == retrCounter {
				return nil, err
			}
		case ErrNotFound:
			return nil, ErrNotFound
		default:
			continue
		}
	}

	return nil, ErrUnknownError
}