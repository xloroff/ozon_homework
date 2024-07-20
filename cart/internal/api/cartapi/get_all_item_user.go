package cartapi

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gookit/validate"
	"github.com/gorilla/mux"
	"go.uber.org/zap"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/model"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/pkg/logger"
)

// GetAllUserItems получение всей корзины пользователя.
func (a *API) GetAllUserItems(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	ctx := r.Context()

	usrID, err := a.getUserID(ctx, r)
	if err != nil {
		a.logger.Errorf(ctx, "Api.GetAllUserItems: не удалось обработать входящий запрос - %v", err)

		if errResp := a.responseSenderV1(ctx, w, http.StatusOK, http.StatusBadRequest, nil, err); errResp != nil {
			a.logger.Debugf(ctx, "Api.GetAllUserItems: не удалось отправить ответ - %v", errResp)
		}

		return
	}

	a.logger.Debugf(ctx, "Api.GetAllUserItems: запрос товаров в корзине пользователя - %v", usrID)

	fullCart, err := a.cartService.GetAllUserItems(ctx, usrID)
	if err != nil {
		a.logger.Errorf(ctx, "Api.GetAllUserItems: ошибка получения корзины пользователя %v - %v", usrID, err)

		if errResp := a.responseSenderV1(ctx, w, http.StatusOK, http.StatusNotFound, nil, err); errResp != nil {
			a.logger.Debugf(ctx, "Api.GetAllUserItems: не удалось отправить ответ - %v", errResp)
		}

		return
	}

	ctx = logger.Append(ctx, []zap.Field{zap.Any("cart", fullCart)})
	a.logger.Debugf(ctx, "Api.GetAllUserItems: получена корзина пользователя - %v", usrID)

	if errResp := a.responseSenderV1(ctx, w, http.StatusOK, http.StatusBadRequest, fullCart, err); errResp != nil {
		a.logger.Debugf(ctx, "Api.GetAllUserItems: не удалось отправить ответ - %v", errResp)
	}
}

// getUserID получает из запроса UserId.
func (a *API) getUserID(ctx context.Context, r *http.Request) (int64, error) {
	// Читаем переменные.
	vars := mux.Vars(r)

	usrID, err := strconv.ParseInt(vars[model.UsrID], base, bitSize)
	if err != nil {
		a.logger.Debugf(ctx, "Api.getUserID: ошибка чтения "+model.UsrID+"- %v", err)

		return 0, fmt.Errorf("Ошибка чтения %s- %w", model.SkuID, err)
	}

	usr := &model.UserIdintyfier{
		UserID: usrID,
	}

	// Валидируем входящие данные.
	v := validate.Struct(usr)
	if !v.Validate() {
		err = v.Errors
		a.logger.Debugf(ctx, "Api.getUserID: ошибка валидации входящих данных - %v", err)

		return 0, fmt.Errorf("Ошибка валидации входящих данных - %w", err)
	}

	return usrID, nil
}
