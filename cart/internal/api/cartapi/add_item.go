package cartapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gookit/validate"
	"github.com/gorilla/mux"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/model"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/pkg/logger"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/pkg/metrics"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/pkg/tracer"
)

const (
	base    = 10
	bitSize = 64
)

// AddItem добавляет итем в корзину.
func (a *API) AddItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ctx, span := tracer.StartSpanFromContext(ctx, "api.cartapi.add_item")
	span.SetTag("component", "cartapi")

	defer span.End()

	metrics.UpdateRequestsTotalWithURL(r.Method, getReqURLTemplate(r))

	defer r.Body.Close()

	item, err := a.getAndCheckItem(ctx, r)
	if err != nil {
		span.SetTag("error", true)
		a.logger.Errorf(ctx, "Api.AddItem: не удалось обработать входящий запрос - %v", err)

		if errResp := a.responseSender(ctx, w, r, http.StatusBadRequest, nil, err); errResp != nil {
			a.logger.Debugf(ctx, "Api.AddItem: не удалось отправить ответ - %v", errResp)
		}

		return
	}

	ctx = logger.AddFieldsToContext(ctx, "data", item, "user_id", item.UserID)
	a.logger.Debugf(ctx, "Api.AddItem: запрос добавления товара в корзину пользователя - %v", item.UserID)

	err = a.cartService.AddItem(ctx, item)
	if err != nil {
		span.SetTag("error", true)
		a.logger.Errorf(ctx, "Api.AddItem: ошибка добавления товара в корзину - %v", err)

		if errResp := a.responseSender(ctx, w, r, http.StatusPreconditionFailed, nil, err); errResp != nil {
			a.logger.Debugf(ctx, "Api.AddItem: не удалось отправить ответ - %v", errResp)
		}

		return
	}

	if errResp := a.responseSender(ctx, w, r, http.StatusOK, nil, nil); errResp != nil {
		a.logger.Debugf(ctx, "Api.AddItem: не удалось отправить ответ - %v", errResp)
	}
}

// getAndCheckItem извлекает данные из запроса UserID, SkuID, Count.
func (a *API) getAndCheckItem(ctx context.Context, r *http.Request) (*model.AddItem, error) {
	ctx, span := tracer.StartSpanFromContext(ctx, "api.cartapi.get_and_check_item")
	span.SetTag("component", "cartapi")
	span.SetTag("span.kind", "child")
	span.SetTag("http.url", getReqURLTemplate(r))
	span.SetTag("http.method", r.Method)

	defer span.End()

	vars := mux.Vars(r)

	usrID, err := strconv.ParseInt(vars[model.UsrID], base, bitSize)
	if err != nil {
		span.SetTag("error", true)
		a.logger.Debugf(ctx, "Api.getAndCheckItem: ошибка чтения "+model.UsrID+"- %v", err)

		return nil, fmt.Errorf("Ошибка чтения %s- %w", model.UsrID, err)
	}

	skuID, err := strconv.ParseInt(vars[model.SkuID], base, bitSize)
	if err != nil {
		span.SetTag("error", true)
		a.logger.Debugf(ctx, "Api.getAndCheckItem: ошибка чтения "+model.SkuID+" - %v", err)

		return nil, fmt.Errorf("Ошибка чтения %s- %w", model.SkuID, err)
	}

	// Чтение тела реквеста.
	body, err := io.ReadAll(r.Body)
	if err != nil {
		span.SetTag("error", true)
		a.logger.Debugf(ctx, "Api.getAndCheckItem: ошибка чтения тела запроса - %v", err)

		return nil, fmt.Errorf("Ошибка чтения тела запроса - %w", err)
	}

	// Анмаршал и получение запроса пришедшего в приклад.
	var d *model.AddItemBody

	err = json.Unmarshal(body, &d)
	if err != nil {
		span.SetTag("error", true)
		a.logger.Debugf(ctx, "Api.getAndCheckItem: ошибка конвертации входящего json - %v", err)

		return nil, fmt.Errorf("Ошибка конвертации входящего json - %w", err)
	}

	item := &model.AddItem{}
	item.UsrSkuID.SkuID = skuID
	item.UserIdintyfier.UserID = usrID
	item.AddItemBody.Count = d.Count

	// Валидируем входящие данные.
	v := validate.Struct(item)
	if !v.Validate() {
		span.SetTag("error", true)

		err = v.Errors
		a.logger.Debugf(ctx, "Api.getAndCheckItem: ошибка валидации входящих данных - %v", err)

		return nil, fmt.Errorf("Ошибка валидации входящих данных - %w", err)
	}

	return item, nil
}
