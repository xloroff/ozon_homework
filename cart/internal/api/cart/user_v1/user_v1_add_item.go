package user_v1

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"go.uber.org/zap"

	"github.com/gookit/validate"
	"github.com/gorilla/mux"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/internal/initilize"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/internal/model/v1"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/internal/pkg/logger"
)

// AddItem добавляет итем в корзину.
func (a *apiv1) AddItem(settings *initilize.ConfigAPI) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		ctx := r.Context()

		item, err := getAndCheckItem(ctx, r)
		if err != nil {
			logger.Errorf(ctx, "ApiV1.AddItem: не удалось обработать входящий запрос - %w", err)
			a.ResponseSenderV1(ctx, w, http.StatusOK, http.StatusBadRequest, nil, err)
			return
		}

		ctx = logger.Append(ctx, []zap.Field{zap.Any("request", item)})
		logger.Debugf(ctx, "ApiV1.AddItem: запрос добавления товара в корзину пользователя - %v", item.UserID)

		err = a.cartService.AddItem(ctx, settings, item)
		if err != nil {
			a.ResponseSenderV1(ctx, w, http.StatusOK, http.StatusPreconditionFailed, nil, err)
			logger.Errorf(ctx, "ApiV1.AddItem: ошибка добавления товара в корзину - %w", err)
		}

		a.ResponseSenderV1(ctx, w, http.StatusOK, http.StatusBadRequest, nil, err)
	}
}

// getAndCheckItem извлекает данные из запроса UserID, SkuID, Count.
func getAndCheckItem(ctx context.Context, r *http.Request) (*v1.AddItem, error) {
	// Читаем переменные.
	vars := mux.Vars(r)

	usrID, err := strconv.ParseInt(vars[v1.UsrID], 10, 64)
	if err != nil {
		logger.Debugf(ctx, "ApiV1.getAndCheckItem: ошибка чтения "+v1.UsrID+"- %w", err)
		return nil, fmt.Errorf("Ошибка чтения %s- %w", v1.UsrID, err)
	}

	skuID, err := strconv.ParseInt(vars[v1.SkuID], 10, 64)
	if err != nil {
		logger.Debugf(ctx, "ApiV1.getAndCheckItem: ошибка чтения "+v1.SkuID+" - %w", err)
		return nil, fmt.Errorf("Ошибка чтения %s- %w", v1.SkuID, err)
	}

	// Чтение тела реквеста.
	body, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Debugf(ctx, "ApiV1.getAndCheckItem: ошибка чтения тела запроса - %w", err)
		return nil, fmt.Errorf("Ошибка чтения тела запроса - %w", err)
	}

	// Анмаршал и получение запроса пришедшего в приклад.
	d := &v1.AddItemBody{}
	err = json.Unmarshal(body, d)
	if err != nil {
		logger.Debugf(ctx, "ApiV1.getAndCheckItem: ошибка конвертации входящего json - %w", err)
		return nil, fmt.Errorf("Ошибка конвертации входящего json - %w", err)
	}

	item := &v1.AddItem{}
	item.UsrSkuID.SkuID = skuID
	item.UserIdintyfier.UserID = usrID
	item.AddItemBody.Count = d.Count

	// Валидируем входящие данные.
	v := validate.Struct(item)
	if !v.Validate() {
		err = v.Errors
		logger.Debugf(ctx, "ApiV1.getAndCheckItem: ошибка валидации входящих данных - %w", err)
		return nil, fmt.Errorf("Ошибка валидации входящих данных - %w", err)
	}

	return item, nil
}