package cartapi

import (
	"net/http"
)

// DelCart удаляет корзину.
func (a *API) DelCart(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	ctx := r.Context()

	usrID, err := a.getUserID(ctx, r)
	if err != nil {
		a.logger.Errorf(ctx, "Api.DelCart: не удалось обработать входящий запрос - %v", err)

		if errResp := a.responseSenderV1(ctx, w, http.StatusBadRequest, http.StatusBadRequest, nil, err); errResp != nil {
			a.logger.Debugf(ctx, "Api.DelCart: не удалось отправить ответ - %v", errResp)
		}

		return
	}

	a.logger.Debugf(ctx, "Api.DelCart: запрос удаления корзины пользователя - %v", usrID)

	err = a.cartService.DelCart(ctx, usrID)
	if err != nil {
		a.logger.Errorf(ctx, "Api.DelCart: ошибка удаления корзины пользователя %v - %v", usrID, err)

		if errResp := a.responseSenderV1(ctx, w, http.StatusBadRequest, http.StatusBadRequest, nil, err); errResp != nil {
			a.logger.Debugf(ctx, "Api.DelCart: не удалось отправить ответ - %v", errResp)
		}

		return
	}

	if errResp := a.responseSenderV1(ctx, w, http.StatusNoContent, http.StatusBadRequest, nil, err); errResp != nil {
		a.logger.Debugf(ctx, "Api.DelCart: не удалось отправить ответ - %v", errResp)
	}
}
