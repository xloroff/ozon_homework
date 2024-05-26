package standartcli

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gookit/validate"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/internal/config"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/internal/model/v1"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/internal/pkg/logger"
)

// responseChecker проверяем ответ от сервиса продуктов и чекаем.
func responseChecker(ctx context.Context, httpResp *http.Response) (*v1.ProductResp, error) {
	switch httpResp.StatusCode {
	// Ответ подсказывающий, что товар на сервисе существует.
	case http.StatusOK:
		// Читаем тело ответа.
		result, err := io.ReadAll(httpResp.Body)
		if err != nil {
			logger.Errorf(ctx, "clientProductServiceStandart.responseChecker: ошибка чтения тела ответа от сервиса продуктов - %w", err)
			return nil, fmt.Errorf("Ошибка чтения тела ответа от сервиса продуктов - %w", err)
		}

		// Анмаршалим ответ - считаем что получена информация успешно.
		otv := &v1.ProductResp{}
		err = json.Unmarshal(result, otv)
		if err != nil {
			logger.Errorf(ctx, "clientProductServiceStandart.responseChecker: ошибка обработки ответа от сервиса продуктов - %w", err)
			return nil, fmt.Errorf("Ошибка обработки ответа от сервиса продуктов - %w", err)
		}

		// Валидируем ответ.
		v := validate.Struct(otv)
		if !v.Validate() {
			err = v.Errors
			logger.Errorf(ctx, "clientProductServiceStandart.responseChecker: валидация данных в ответе от сервиса продуктов непройдена - %w", err)
			return nil, fmt.Errorf("Валидация данных в ответе от сервиса продуктов непройдена - %w", err)
		}

		return otv, nil
	// Товара на сервисе товаров нет.
	case http.StatusNotFound:
		logger.Info(ctx, "clientProductServiceStandart.responseChecker: получен статус - продукт не найден")
		return nil, ErrNotFound

	// Сервис продуктов захлебнулся.
	case http.StatusTooManyRequests:
		logger.Warnf(ctx, "clientProductServiceStandart.responseChecker: получен статус - слишком много запросов")
		return nil, ErrTooManyRequests
	// В ТЗ вынесена эта ошибка - вынес отдельно, если неверно ее понял.
	case config.HTTPCalmStatus:
		logger.Warnf(ctx, "clientProductServiceStandart.responseChecker: получен статус - слишком много запросов")
		return nil, ErrTooManyRequests
	default:
		logger.Errorf(ctx, "clientProductServiceStandart.responseChecker: неизвестная ошибка, код статуса - ", httpResp.StatusCode)
		return nil, ErrUnknownError
	}
}