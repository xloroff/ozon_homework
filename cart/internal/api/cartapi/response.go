package cartapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/model"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/pkg/logger"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/pkg/metrics"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/pkg/tracer"
)

func (a *API) responseSender(ctx context.Context, w http.ResponseWriter, r *http.Request, code int, data any, errRespons error) error {
	ctx, span := tracer.StartSpanFromContext(ctx, "api.cartapi.response_sender")
	span.SetTag("component", "cartapi")

	defer span.End()

	ctx = logger.AddFieldsToContext(ctx, "event", "ss")
	ew := ExtendResponseWriter(w)
	ew.Header().Add("Content-Type", "application/json")
	span.SetTag("http.status_code", http.StatusText(ew.StatusCode))
	span.SetTag("http.url", getReqURLTemplate(r))
	span.SetTag("http.method", r.Method)

	var ans any

	if data != nil {
		ans = data

		err := json.NewEncoder(ew).Encode(&ans)
		if err != nil {
			ew.WriteHeader(http.StatusInternalServerError)

			return fmt.Errorf("Api.ResponseSender: ошибка формирования ответа - %w", err)
		}

		ctx = logger.AddFieldsToContext(ctx, "data", ans)
	}

	if errRespons != nil {
		ans = model.APIResponse{
			Error: errRespons.Error(),
		}

		ctx = logger.AddFieldsToContext(ctx, "error.object", errRespons)
	}

	ew.WriteHeader(code)
	a.logger.Debugf(ctx, "Api.ResponseSender: отправлен ответ")

	metrics.UpdateResponseCode(r.Method, getReqURLTemplate(r), http.StatusText(ew.StatusCode))

	return nil
}
