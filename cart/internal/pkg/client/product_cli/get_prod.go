package productcli

import (
	"context"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/gookit/validate"
	"go.uber.org/zap"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/model"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/pkg/logger"
)

func (c *client) GetProduct(ctx context.Context, skuID int64) (*model.ProductResp, error) {
	c.logger.Debug(ctx, fmt.Sprintf("ClientProductServiceResty.ProductReciver: начало обращения к сервису продуктов, продукт - %d", skuID))
	defer c.logger.Debug(ctx, fmt.Sprintf("ClientProductServiceResty.ProductReciver: конец обращения к сервису продуктов, продукт - %d", skuID))

	connector := resty.New()

	connector.
		// SetRetryCount устанавливает число ретраев.
		SetRetryCount(c.config.ProductServiceRetr).
		// SetRetryWaitTime таймаут между попытками.
		SetRetryWaitTime(200 * time.Millisecond).
		// SetRetryMaxWaitTime сколько ждуем выполнения попытки.
		SetRetryMaxWaitTime(1 * time.Second).
		// SetTimeout устанавливаем общий таймаут.
		SetTimeout(5 * time.Second).
		// SetRetryAfter что делать если попытки закончились.
		SetRetryAfter(func(_ *resty.Client, _ *resty.Response) (time.Duration, error) {
			c.logger.Errorf(ctx, "ClientProductServiceResty.ProductReciver: попытка связи с сервисом продуктов не удалась - %v", connector.Error)

			return 0, fmt.Errorf("Попытка связи с сервисом продуктов не удалась - %w", connector.Error)
		})

	connector.SetBaseURL(c.config.ProductServiceHost)

	req := &model.ProductReq{
		Token: c.config.ProductServiceToken,
		Sku:   skuID,
	}

	result := &model.ProductResp{}

	resp, err := connector.R().
		SetBody(req).
		EnableTrace().
		SetResult(result).
		ForceContentType("application/json").
		Post(urlapi)
	if err != nil {
		return nil, fmt.Errorf("ClientProductServiceResty.ProductReciver:  Ошибка отправки запроса в сервис продуктов - %w", err)
	}

	err = c.responseChecker(ctx, resp)
	if err != nil {
		return nil, err
	}

	// Валидируем ответ
	v := validate.Struct(resp)
	if !v.Validate() {
		err = v.Errors

		return nil, fmt.Errorf("ClientProductServiceResty.responseChecker: Валидация данных в ответе от сервиса продуктов непройдена - %w", err)
	}

	ctx = logger.Set(ctx, []zap.Field{zap.Any("response", result)})
	c.logger.Debugf(ctx, "ClientProductServiceResty.GetProduct: ответ сервиса продуктов")

	return result, nil
}
