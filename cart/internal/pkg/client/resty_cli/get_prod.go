package restycli

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/gookit/validate"
	"go.uber.org/zap"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/internal/initilize"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/internal/model/v1"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/internal/pkg/logger"
)

func (c *client) GetProduct(ctx context.Context, settings *initilize.ConfigAPI, skuID int64) (*v1.ProductResp, error) {
	logger.Debug(ctx, fmt.Sprintf("ClientProductServiceResty.ProductReciver: начало обращения к сервису продуктов, продукт - %d", skuID))
	defer logger.Debug(ctx, fmt.Sprintf("ClientProductServiceResty.ProductReciver: конец обращения к сервису продуктов, продукт - %d", skuID))

	connector := resty.New()

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	connector.SetTransport(transport)

	connector.
		// SetRetryCount устанавливает число ретраев.
		SetRetryCount(settings.ProductServiceRetr).
		// SetRetryWaitTime таймаут между попытками.
		SetRetryWaitTime(200 * time.Millisecond).
		// SetRetryMaxWaitTime сколько ждуем выполнения попытки.
		SetRetryMaxWaitTime(1 * time.Second).
		// SetTimeout устанавливаем общий таймаут.
		SetTimeout(5 * time.Second).
		// SetRetryAfter что делать если попытки закончились.
		SetRetryAfter(func(_ *resty.Client, _ *resty.Response) (time.Duration, error) {
			logger.Errorf(ctx, "ClientProductServiceResty.ProductReciver: попытка связи с сервисом продуктов не удалась - %w", connector.Error)
			return 0, fmt.Errorf("Попытка связи с сервисом продуктов не удалась - %w", connector.Error)
		})

	connector.SetBaseURL(settings.ProductServiceHost)
	urlapi := "/get_product"

	req := &v1.ProductReq{
		Token: settings.ProductServiceToken,
		Sku:   skuID,
	}

	result := &v1.ProductResp{}

	resp, err := connector.R().
		SetBody(req).
		EnableTrace().
		SetResult(result).
		ForceContentType("application/json").
		Post(urlapi)

	if err != nil {
		logger.Errorf(ctx, "ClientProductServiceResty.ProductReciver: ошибка отправки запроса в сервис продуктов - %w", err)
		return nil, fmt.Errorf("Ошибка отправки запроса в сервис продуктов - %w", err)
	}

	err = responseChecker(ctx, resp)
	if err != nil {
		return nil, err
	}

	// Валидируем ответ
	v := validate.Struct(resp)
	if !v.Validate() {
		err = v.Errors
		logger.Errorf(ctx, "ClientProductServiceResty.responseChecker: валидация данных в ответе от сервиса продуктов непройдена - %w", err)
		return nil, fmt.Errorf("Валидация данных в ответе от сервиса продуктов непройдена - %w", err)
	}

	ctx = logger.Set(ctx, []zap.Field{zap.Any("response", result)})
	logger.Debugf(ctx, "ClientProductServiceResty.GetProduct: ответ сервиса продуктов")

	return result, nil
}