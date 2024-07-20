package cartapi

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gookit/validate"
	"github.com/gorilla/mux"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/model"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/pkg/logger"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/pkg/metrics"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/pkg/tracer"
)

// GetAllUserItems получение всей корзины пользователя.
func (a *API) GetAllUserItems(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ctx, span := tracer.StartSpanFromContext(ctx, "cartapi.get_all_user_items")
	span.SetTag("component", "cartapi")
	span.SetTag("http.url", getReqURLTemplate(r))
	span.SetTag("http.method", r.Method)

	defer span.End()

	metrics.UpdateRequestsTotalWithURL(r.Method, getReqURLTemplate(r))

	defer r.Body.Close()

	usrID, err := a.getUserID(ctx, r)
	if err != nil {
		span.SetTag("error", true)
		a.logger.Errorf(ctx, "Api.GetAllUserItems: не удалось обработать входящий запрос - %v", err)

		if errResp := a.responseSender(ctx, w, r, http.StatusBadRequest, nil, err); errResp != nil {
			a.logger.Debugf(ctx, "Api.GetAllUserItems: не удалось отправить ответ - %v", errResp)
		}

		return
	}

	a.logger.Debugf(ctx, "Api.GetAllUserItems: запрос товаров в корзине пользователя - %v", usrID)

	fullCart, err := a.cartService.GetAllUserItems(ctx, usrID)
	if err != nil {
		span.SetTag("error", true)
		a.logger.Errorf(ctx, "Api.GetAllUserItems: ошибка получения корзины пользователя %v - %v", usrID, err)

		if errResp := a.responseSender(ctx, w, r, http.StatusNotFound, nil, err); errResp != nil {
			a.logger.Debugf(ctx, "Api.GetAllUserItems: не удалось отправить ответ - %v", errResp)
		}

		return
	}

	ctx = logger.AddFieldsToContext(ctx, "data", fullCart, "user_id", usrID)
	a.logger.Debugf(ctx, "Api.GetAllUserItems: получена корзина пользователя - %v", usrID)

	if errResp := a.responseSender(ctx, w, r, http.StatusOK, fullCart, nil); errResp != nil {
		a.logger.Debugf(ctx, "Api.GetAllUserItems: не удалось отправить ответ - %v", errResp)
	}
}

// getUserID получает из запроса UserId.
func (a *API) getUserID(ctx context.Context, r *http.Request) (int64, error) {
	ctx, span := tracer.StartSpanFromContext(ctx, "api.cartapi.get_user_id")
	span.SetTag("span.kind", "child")
	span.SetTag("component", "cartapi")
	// Читаем переменные.
	vars := mux.Vars(r)

	usrID, err := strconv.ParseInt(vars[model.UsrID], base, bitSize)
	if err != nil {
		span.SetTag("error", true)
		a.logger.Debugf(ctx, "Api.getUserID: ошибка чтения "+model.UsrID+"- %v", err)

		return 0, fmt.Errorf("Ошибка чтения %s- %w", model.SkuID, err)
	}

	usr := &model.UserIdintyfier{
		UserID: usrID,
	}

	// Валидируем входящие данные.
	v := validate.Struct(usr)
	if !v.Validate() {
		span.SetTag("error", true)

		err = v.Errors
		a.logger.Debugf(ctx, "Api.getUserID: ошибка валидации входящих данных - %v", err)

		return 0, fmt.Errorf("Ошибка валидации входящих данных - %w", err)
	}

	return usrID, nil
}
