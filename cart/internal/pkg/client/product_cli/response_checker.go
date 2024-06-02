package productcli

import (
	"context"
	"net/http"

	"github.com/go-resty/resty/v2"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/internal/config"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/internal/model"
)

// responseChecker проверяем ответ от сервиса продуктов и чекаем.
func (c *client) responseChecker(ctx context.Context, httpResp *resty.Response) error {
	switch httpResp.StatusCode() {
	// Ответ подсказывающий, что товар на сервисе существует.
	case http.StatusOK:
		return nil
	// Товара на сервисе товаров нет.
	case http.StatusNotFound:
		c.logger.Info(ctx, "ClientProductServiceResty.responseChecker: получен статус - продукт не найден")

		return model.ErrNotFound
	// Сервис продуктов захлебнулся.
	case http.StatusTooManyRequests:
		c.logger.Warn(ctx, "ClientProductServiceResty.responseChecker: получен статус - слишком много запросов")

		return model.ErrTooManyRequests
	// В ТЗ вынесена эта ошибка - вынес отдельно, если неверно ее понял.
	case config.HTTPCalmStatus:
		c.logger.Warn(ctx, "ClientProductServiceResty.responseChecker: получен статус - слишком много запросов")

		return model.ErrTooManyRequests
	default:
		c.logger.Errorf(ctx, "ClientProductServiceResty.responseChecker: неизвестная ошибка, код статуса - %v", httpResp.StatusCode())

		return model.ErrUnknownError
	}
}
