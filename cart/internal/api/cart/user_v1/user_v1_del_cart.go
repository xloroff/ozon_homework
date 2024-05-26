package user_v1

import (
	"net/http"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/internal/pkg/logger"
)

// DelCart удаляет корзину.
func (a *apiv1) DelCart(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	ctx := r.Context()

	usrID, err := getUserId(ctx, r)
	if err != nil {
		logger.Errorf(ctx, "ApiV1.DelCart: не удалось обработать входящий запрос - %w", err)
		a.ResponseSenderV1(ctx, w, http.StatusBadRequest, http.StatusBadRequest, nil, err)
		return
	}

	logger.Debugf(ctx, "ApiV1.DelCart: запрос удаления корзины пользователя - %v", usrID)

	err = a.cartService.DelCart(ctx, usrID)
	if err != nil {
		a.ResponseSenderV1(ctx, w, http.StatusBadRequest, http.StatusBadRequest, nil, err)
		logger.Errorf(ctx, "ApiV1.DelCart: ошибка удаления корзины пользователя %v - %w", usrID, err)
	}

	a.ResponseSenderV1(ctx, w, http.StatusNoContent, http.StatusBadRequest, nil, err)
}