package restycli

import (
	"context"
	"net/http"

	"github.com/go-resty/resty/v2"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/internal/config"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/internal/pkg/logger"
)

// responseChecker проверяем ответ от сервиса продуктов и чекаем.
func responseChecker(ctx context.Context, httpResp *resty.Response) error {
	switch httpResp.StatusCode() {
	// Ответ подсказывающий что товар на сервисе существует.
	case http.StatusOK:
		return nil
	// Товара на сервисе товаров нет.
	case http.StatusNotFound:
		logger.Info(ctx, "ClientProductServiceResty.responseChecker: получен статус - продукт не найден")
		return ErrNotFound

	// Сервис продуктов захлебнулся.
	case http.StatusTooManyRequests:
		logger.Warnf(ctx, "ClientProductServiceResty.responseChecker: получен статус - слишком много запросов")
		return ErrTooManyRequests

	// В ТЗ вынесена эта ошибка - вынес отдельно, если неверно ее понял.
	case config.HTTPCalmStatus:
		logger.Warnf(ctx, "ClientProductServiceResty.responseChecker: получен статус - слишком много запросов")
		return ErrTooManyRequests

	default:
		logger.Errorf(ctx, "ClientProductServiceResty.responseChecker: неизвестная ошибка, код статуса - ", httpResp.StatusCode())
		return ErrUnknownError
	}
}