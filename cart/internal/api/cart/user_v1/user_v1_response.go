package user_v1

import (
	"context"
	"encoding/json"
	"net/http"

	"go.uber.org/zap"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/internal/model/v1"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/internal/pkg/logger"
)

func (a *apiv1) ResponseSenderV1(ctx context.Context, w http.ResponseWriter, normalStatus, badStatus int, data any, errResult error) {
	w.Header().Add("Content-Type", "application/json")
	var ans any

	if errResult != nil {
		w.WriteHeader(badStatus)
		ans = v1.ApiResponse{
			Error: errResult.Error(),
		}
	} else {
		w.WriteHeader(normalStatus)
		ans = data
	}

	if data != nil {
		err := json.NewEncoder(w).Encode(&ans)
		if err != nil {
			logger.Debugf(ctx, "ApiV1.ResponseSender: ошибка формирования ответа - %w", err)
		}
	}

	ctx = logger.Set(ctx, []zap.Field{zap.Any("response", ans)})
	logger.Debugf(ctx, "ApiV1.ResponseSender: отправлен ответ")
}