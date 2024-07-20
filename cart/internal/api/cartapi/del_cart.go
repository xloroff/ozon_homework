package cartapi

import (
	"net/http"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/pkg/metrics"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/pkg/tracer"
)

// DelCart удаляет корзину.
func (a *API) DelCart(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ctx, span := tracer.StartSpanFromContext(ctx, "api.cartapi.del_cart")
	span.SetTag("component", "cartapi")
	span.SetTag("http.url", getReqURLTemplate(r))
	span.SetTag("http.method", r.Method)

	defer span.End()

	metrics.UpdateRequestsTotalWithURL(r.Method, getReqURLTemplate(r))

	defer r.Body.Close()

	usrID, err := a.getUserID(ctx, r)
	if err != nil {
		span.SetTag("error", true)
		a.logger.Errorf(ctx, "Api.DelCart: не удалось обработать входящий запрос - %v", err)

		if errResp := a.responseSender(ctx, w, r, http.StatusBadRequest, nil, err); errResp != nil {
			a.logger.Debugf(ctx, "Api.DelCart: не удалось отправить ответ - %v", errResp)
		}

		return
	}

	a.logger.Debugf(ctx, "Api.DelCart: запрос удаления корзины пользователя - %v", usrID)

	err = a.cartService.DelCart(ctx, usrID)
	if err != nil {
		span.SetTag("error", true)
		a.logger.Errorf(ctx, "Api.DelCart: ошибка удаления корзины пользователя %v - %v", usrID, err)

		if errResp := a.responseSender(ctx, w, r, http.StatusPreconditionFailed, nil, err); errResp != nil {
			a.logger.Debugf(ctx, "Api.DelCart: не удалось отправить ответ - %v", errResp)
		}

		return
	}

	if errResp := a.responseSender(ctx, w, r, http.StatusNoContent, nil, nil); errResp != nil {
		a.logger.Debugf(ctx, "Api.DelCart: не удалось отправить ответ - %v", errResp)
	}
}
