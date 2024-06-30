package productcli

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/gookit/validate"
	"go.opentelemetry.io/otel/trace"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/model"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/pkg/logger"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/pkg/metrics"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/pkg/tracer"
)

func (c *client) GetProduct(ctx context.Context, skuID int64) (*model.ProductResp, error) {
	ctx, span := tracer.StartSpanFromContext(ctx, "productcli.get_product", trace.WithSpanKind(trace.SpanKindClient))
	span.SetTag("component", "productcli")
	span.SetTag("method", http.MethodGet)
	span.SetTag("peer.hostname", c.config.ProductServiceHost)
	span.SetTag("http.url", urlapi)

	defer span.End()

	c.logger.Debug(ctx, fmt.Sprintf("ClientProductServiceResty.ProductReciver: начало обращения к сервису продуктов, продукт - %d", skuID))
	defer c.logger.Debug(ctx, fmt.Sprintf("ClientProductServiceResty.ProductReciver: конец обращения к сервису продуктов, продукт - %d", skuID))

	connector := resty.New()

	connector.
		SetPreRequestHook(func(_ *resty.Client, _ *http.Request) error {
			metrics.UpdateExternalRequestsTotal(c.config.ProductServiceHost, urlapi)
			return nil
		}).
		// SetRetryCount устанавливает число ретраев.
		SetRetryCount(c.config.ProductServiceRetr).
		// SetRetryWaitTime таймаут между попытками.
		SetRetryWaitTime(200 * time.Millisecond).
		// SetRetryMaxWaitTime сколько ждуем выполнения попытки.
		SetRetryMaxWaitTime(1 * time.Second).
		// SetTimeout устанавливаем общий таймаут.
		SetTimeout(5 * time.Second).
		// SetRetryAfter что делать если попытки закончились.
		SetRetryAfter(func(_ *resty.Client, resp *resty.Response) (time.Duration, error) {
			c.logger.Errorf(ctx, "ClientProductServiceResty.ProductReciver: попытка связи с сервисом продуктов не удалась - %v", connector.Error)
			metrics.UpdateExternalResponseCode(c.config.ProductServiceHost, urlapi, http.StatusText(resp.StatusCode()))
			metrics.UpdateExternalResponseDuration(resp.Time())

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
		span.SetTag("error", true)
		return nil, fmt.Errorf("ClientProductServiceResty.ProductReciver:  Ошибка отправки запроса в сервис продуктов - %w", err)
	}

	metrics.UpdateExternalResponseCode(c.config.ProductServiceHost, urlapi, http.StatusText(resp.StatusCode()))
	metrics.UpdateExternalResponseDuration(resp.Time())

	err = c.responseChecker(ctx, resp)
	if err != nil {
		return nil, err
	}

	// Валидируем ответ
	v := validate.Struct(resp)
	if !v.Validate() {
		span.SetTag("error", true)

		err = v.Errors

		return nil, fmt.Errorf("ClientProductServiceResty.responseChecker: Валидация данных в ответе от сервиса продуктов непройдена - %w", err)
	}

	ctx = logger.AddFieldsToContext(ctx, "data", result)
	c.logger.Debugf(ctx, "ClientProductServiceResty.GetProduct: ответ сервиса продуктов")

	return result, nil
}
