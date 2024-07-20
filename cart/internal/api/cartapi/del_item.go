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

// DelItem удаляет итем из корзины.
func (a *API) DelItem(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	ctx := r.Context()

	item, err := a.getAndCheckUserItem(ctx, r)
	if err != nil {
		a.logger.Errorf(ctx, "Api.DelItem: не удалось обработать входящий запрос - %v", err)

		if errResp := a.responseSenderV1(ctx, w, http.StatusBadRequest, http.StatusBadRequest, nil, err); errResp != nil {
			a.logger.Debugf(ctx, "Api.DelItem: не удалось отправить ответ - %v", errResp)
		}

		return
	}

	ctx = logger.Append(ctx, []zap.Field{zap.Any("request", item)})
	a.logger.Debugf(ctx, "Api.DelItem: запрос удаления товара из корзины пользователя - %v", item.UserID)

	err = a.cartService.DelItem(ctx, item)
	if err != nil {
		a.logger.Errorf(ctx, "Api.DelItem: ошибка удаления товара из корзины пользователя %d - %v", item.UserID, err)

		if errResp := a.responseSenderV1(ctx, w, http.StatusBadRequest, http.StatusBadRequest, nil, err); errResp != nil {
			a.logger.Debugf(ctx, "Api.DelItem: не удалось отправить ответ - %v", errResp)
		}

		return
	}

	if errResp := a.responseSenderV1(ctx, w, http.StatusNoContent, http.StatusBadRequest, nil, err); errResp != nil {
		a.logger.Debugf(ctx, "Api.DelItem: не удалось отправить ответ - %v", errResp)
	}
}

// getAndCheckUserItem извлекает данные из запроса UserID, SkuID.
func (a *API) getAndCheckUserItem(ctx context.Context, r *http.Request) (*model.DelItem, error) {
	// Читаем переменные.
	vars := mux.Vars(r)

	usrID, err := strconv.ParseInt(vars[model.UsrID], base, bitSize)
	if err != nil {
		a.logger.Debugf(ctx, "Api.getAndCheckUserItem: ошибка чтения "+model.UsrID+"- %v", err)

		return nil, fmt.Errorf("Ошибка чтения %s- %w", model.UsrID, err)
	}

	skuID, err := strconv.ParseInt(vars[model.SkuID], base, bitSize)
	if err != nil {
		a.logger.Debugf(ctx, "Api.getAndCheckUserItem: ошибка чтения "+model.SkuID+" - %v", err)

		return nil, fmt.Errorf("Ошибка чтения %s- %w", model.SkuID, err)
	}

	item := &model.DelItem{}
	item.UserIdintyfier.UserID = usrID
	item.UsrSkuID.SkuID = skuID

	// Валидируем входящие данные.
	v := validate.Struct(item)
	if !v.Validate() {
		err = v.Errors
		a.logger.Debugf(ctx, "Api.getAndCheckUserItem: ошибка валидации входящих данных - %v", err)

		return nil, fmt.Errorf("Ошибка валидации входящих данных - %w", err)
	}

	return item, nil
}
