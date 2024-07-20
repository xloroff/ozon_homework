package cartapi

import "net/http"

// Checkout удаляет корзину.
func (a *API) Checkout(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	ctx := r.Context()

	usrID, err := a.getUserID(ctx, r)
	if err != nil {
		a.logger.Errorf(ctx, "Api.Checkout: не удалось обработать входящий запрос - %v", err)

		if errResp := a.responseSenderV1(ctx, w, http.StatusBadRequest, http.StatusBadRequest, nil, err); errResp != nil {
			a.logger.Debugf(ctx, "Api.Checkout: не удалось отправить ответ - %v", errResp)
		}

		return
	}

	a.logger.Debugf(ctx, "Api.Checkout: запрос на создание заказа - %v", usrID)

	ordr, err := a.cartService.Checkout(ctx, usrID)
	if err != nil {
		a.logger.Errorf(ctx, "Api.Checkout: Ошибка создания заказа пользователя %v - %v", usrID, err)

		if errResp := a.responseSenderV1(ctx, w, http.StatusPreconditionFailed, http.StatusBadRequest, nil, err); errResp != nil {
			a.logger.Debugf(ctx, "Api.Checkout: не удалось отправить ответ - %v", errResp)
		}

		return
	}

	if errResp := a.responseSenderV1(ctx, w, http.StatusOK, http.StatusBadRequest, ordr, err); errResp != nil {
		a.logger.Debugf(ctx, "Api.Checkout: не удалось отправить ответ - %v", errResp)
	}
}
