package cartapi

import (
	"net/http"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/pkg/metrics"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/pkg/tracer"
)

// Checkout удаляет корзину но создает заказ.
func (a *API) Checkout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ctx, span := tracer.StartSpanFromContext(ctx, "api.cartapi.checkout")
	span.SetTag("component", "cartapi")
	span.SetTag("http.url", getReqURLTemplate(r))
	span.SetTag("http.method", r.Method)

	defer span.End()

	metrics.UpdateRequestsTotalWithURL(r.Method, getReqURLTemplate(r))

	defer r.Body.Close()

	usrID, err := a.getUserID(ctx, r)
	if err != nil {
		span.SetTag("error", true)

		a.logger.Errorf(ctx, "Api.Checkout: не удалось обработать входящий запрос - %v", err)

		if errResp := a.responseSender(ctx, w, r, http.StatusBadRequest, nil, err); errResp != nil {
			a.logger.Debugf(ctx, "Api.Checkout: не удалось отправить ответ - %v", errResp)
		}

		return
	}

	a.logger.Debugf(ctx, "Api.Checkout: запрос на создание заказа - %v", usrID)

	ordr, err := a.cartService.Checkout(ctx, usrID)
	if err != nil {
		span.SetTag("error", true)
		a.logger.Errorf(ctx, "Api.Checkout: Ошибка создания заказа пользователя %v - %v", usrID, err)

		if errResp := a.responseSender(ctx, w, r, http.StatusPreconditionFailed, nil, err); errResp != nil {
			a.logger.Debugf(ctx, "Api.Checkout: не удалось отправить ответ - %v", errResp)
		}

		return
	}

	if errResp := a.responseSender(ctx, w, r, http.StatusOK, ordr, nil); errResp != nil {
		a.logger.Debugf(ctx, "Api.Checkout: не удалось отправить ответ - %v", errResp)
	}
}
