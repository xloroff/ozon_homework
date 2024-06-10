package cartapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"go.uber.org/zap"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/model"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/pkg/logger"
)

func (a *API) responseSenderV1(ctx context.Context, w http.ResponseWriter, normalStatus, badStatus int, data any, errResult error) error {
	w.Header().Add("Content-Type", "application/json")

	var ans any

	if errResult != nil {
		w.WriteHeader(badStatus)

		ans = model.APIResponse{
			Error: errResult.Error(),
		}
	} else {
		w.WriteHeader(normalStatus)

		ans = data
	}

	if data != nil {
		err := json.NewEncoder(w).Encode(&ans)
		if err != nil {
			return fmt.Errorf("Api.ResponseSender: ошибка формирования ответа - %w", err)
		}
	}

	ctx = logger.Set(ctx, []zap.Field{zap.Any("response", ans)})
	a.logger.Debugf(ctx, "Api.ResponseSender: отправлен ответ")

	return nil
}
