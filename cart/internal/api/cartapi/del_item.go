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

// DelItem удаляет итем из корзины.
func (a *API) DelItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ctx, span := tracer.StartSpanFromContext(ctx, "api.cartapi.del_item")
	span.SetTag("component", "cartapi")
	span.SetTag("http.url", getReqURLTemplate(r))
	span.SetTag("http.method", r.Method)

	defer span.End()

	metrics.UpdateRequestsTotalWithURL(r.Method, getReqURLTemplate(r))

	defer r.Body.Close()

	item, err := a.getAndCheckUserItem(ctx, r)
	if err != nil {
		span.SetTag("error", true)
		a.logger.Errorf(ctx, "Api.DelItem: не удалось обработать входящий запрос - %v", err)

		if errResp := a.responseSender(ctx, w, r, http.StatusBadRequest, nil, err); errResp != nil {
			a.logger.Debugf(ctx, "Api.DelItem: не удалось отправить ответ - %v", errResp)
		}

		return
	}

	ctx = logger.AddFieldsToContext(ctx, "data", item, "user_id", item.UserID)
	a.logger.Debugf(ctx, "Api.DelItem: запрос удаления товара из корзины пользователя - %v", item.UserID)

	err = a.cartService.DelItem(ctx, item)
	if err != nil {
		span.SetTag("error", true)
		a.logger.Errorf(ctx, "Api.DelItem: ошибка удаления товара из корзины пользователя %d - %v", item.UserID, err)

		if errResp := a.responseSender(ctx, w, r, http.StatusPreconditionFailed, nil, err); errResp != nil {
			a.logger.Debugf(ctx, "Api.DelItem: не удалось отправить ответ - %v", errResp)
		}

		return
	}

	if errResp := a.responseSender(ctx, w, r, http.StatusNoContent, nil, nil); errResp != nil {
		a.logger.Debugf(ctx, "Api.DelItem: не удалось отправить ответ - %v", errResp)
	}
}

// getAndCheckUserItem извлекает данные из запроса UserID, SkuID.
func (a *API) getAndCheckUserItem(ctx context.Context, r *http.Request) (*model.DelItem, error) {
	ctx, span := tracer.StartSpanFromContext(ctx, "api.cartapi.get_and_check_user_item")
	span.SetTag("span.kind", "child")
	span.SetTag("component", "cartapi")

	defer span.End()
	// Читаем переменные.
	vars := mux.Vars(r)

	usrID, err := strconv.ParseInt(vars[model.UsrID], base, bitSize)
	if err != nil {
		span.SetTag("error", true)
		a.logger.Debugf(ctx, "Api.getAndCheckUserItem: ошибка чтения "+model.UsrID+"- %v", err)

		return nil, fmt.Errorf("Ошибка чтения %s- %w", model.UsrID, err)
	}

	skuID, err := strconv.ParseInt(vars[model.SkuID], base, bitSize)
	if err != nil {
		span.SetTag("error", true)
		a.logger.Debugf(ctx, "Api.getAndCheckUserItem: ошибка чтения "+model.SkuID+" - %v", err)

		return nil, fmt.Errorf("Ошибка чтения %s- %w", model.SkuID, err)
	}

	item := &model.DelItem{}
	item.UserIdintyfier.UserID = usrID
	item.UsrSkuID.SkuID = skuID

	// Валидируем входящие данные.
	v := validate.Struct(item)
	if !v.Validate() {
		span.SetTag("error", true)

		err = v.Errors
		a.logger.Debugf(ctx, "Api.getAndCheckUserItem: ошибка валидации входящих данных - %v", err)

		return nil, fmt.Errorf("Ошибка валидации входящих данных - %w", err)
	}

	return item, nil
}
