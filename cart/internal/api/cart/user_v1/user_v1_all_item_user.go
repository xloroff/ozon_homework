package user_v1

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"go.uber.org/zap"

	"github.com/gookit/validate"
	"github.com/gorilla/mux"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/internal/initilize"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/internal/model/v1"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/internal/pkg/logger"
)

// GetAllUserItems получение всей корзины пользователя.
func (a *apiv1) GetAllUserItems(settings *initilize.ConfigAPI) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		ctx := r.Context()

		usrID, err := getUserId(ctx, r)
		if err != nil {
			logger.Errorf(ctx, "ApiV1.GetAllUserItems: не удалось обработать входящий запрос - %w", err)
			a.ResponseSenderV1(ctx, w, http.StatusOK, http.StatusBadRequest, nil, err)
			return
		}

		logger.Debugf(ctx, "ApiV1.GetAllUserItems: запрос товаров в корзине пользователя - %v", usrID)

		fullCart, err := a.cartService.GetAllUserItems(ctx, settings, usrID)
		if err != nil {
			a.ResponseSenderV1(ctx, w, http.StatusOK, http.StatusBadRequest, nil, err)
			logger.Errorf(ctx, "ApiV1.GetAllUserItems: ошибка получения корзины пользователя %v - %w", usrID, err)
		}

		ctx = logger.Append(ctx, []zap.Field{zap.Any("cart", fullCart)})
		logger.Debugf(ctx, "ApiV1.GetAllUserItems: получена корзина пользователя - %v", usrID)

		a.ResponseSenderV1(ctx, w, http.StatusOK, http.StatusBadRequest, fullCart, err)
	}
}

// getUserId получает из запроса UserId.
func getUserId(ctx context.Context, r *http.Request) (int64, error) {
	// Читаем переменные.
	vars := mux.Vars(r)

	usrID, err := strconv.ParseInt(vars[v1.UsrID], 10, 64)
	if err != nil {
		logger.Debugf(ctx, "ApiV1.getUserId: ошибка чтения "+v1.UsrID+"- %w", err)
		return 0, fmt.Errorf("Ошибка чтения %s- %w", v1.SkuID, err)
	}

	usr := &v1.UserIdintyfier{
		UserID: usrID,
	}

	// Валидируем входящие данные.
	v := validate.Struct(usr)
	if !v.Validate() {
		err = v.Errors
		logger.Debugf(ctx, "ApiV1.getUserId: ошибка валидации входящих данных - %w", err)
		return 0, fmt.Errorf("Ошибка валидации входящих данных - %w", err)
	}

	return usrID, nil
}