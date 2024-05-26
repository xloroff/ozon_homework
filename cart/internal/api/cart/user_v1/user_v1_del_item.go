package user_v1

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gookit/validate"
	"github.com/gorilla/mux"
	"go.uber.org/zap"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/internal/model/v1"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/internal/pkg/logger"
)

func (a *apiv1) DelItem(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	ctx := r.Context()

	item, err := getAndCheckUserItem(ctx, r)
	if err != nil {
		logger.Errorf(ctx, "ApiV1.DelItem: не удалось обработать входящий запрос - %w", err)
		a.ResponseSenderV1(ctx, w, http.StatusBadRequest, http.StatusBadRequest, nil, err)
		return
	}

	ctx = logger.Append(ctx, []zap.Field{zap.Any("request", item)})
	logger.Debugf(ctx, "ApiV1.DelItem: запрос удаления товара из корзины пользователя - %v", item.UserID)

	err = a.cartService.DelItem(ctx, item)
	if err != nil {
		a.ResponseSenderV1(ctx, w, http.StatusBadRequest, http.StatusBadRequest, nil, err)
		logger.Errorf(ctx, "ApiV1.DelItem: ошибка удаления товара из корзины пользователя %d - %w", item.UserID, err)
	}

	a.ResponseSenderV1(ctx, w, http.StatusNoContent, http.StatusBadRequest, nil, err)
}

// getAndCheckUserItem извлекает данные из запроса UserID, SkuID.
func getAndCheckUserItem(ctx context.Context, r *http.Request) (*v1.DelItem, error) {
	// Читаем переменные.
	vars := mux.Vars(r)

	usrID, err := strconv.ParseInt(vars[v1.UsrID], 10, 64)
	if err != nil {
		logger.Debugf(ctx, "ApiV1.getAndCheckUserItem: ошибка чтения "+v1.UsrID+"- %w", err)
		return nil, fmt.Errorf("Ошибка чтения %s- %w", v1.UsrID, err)
	}

	skuID, err := strconv.ParseInt(vars[v1.SkuID], 10, 64)
	if err != nil {
		logger.Debugf(ctx, "ApiV1.getAndCheckUserItem: ошибка чтения "+v1.SkuID+" - %w", err)
		return nil, fmt.Errorf("Ошибка чтения %s- %w", v1.SkuID, err)
	}

	item := &v1.DelItem{}
	item.UserIdintyfier.UserID = usrID
	item.UsrSkuID.SkuID = skuID

	// Валидируем входящие данные.
	v := validate.Struct(item)
	if !v.Validate() {
		err = v.Errors
		logger.Debugf(ctx, "ApiV1.getAndCheckUserItem: ошибка валидации входящих данных - %w", err)
		return nil, fmt.Errorf("Ошибка валидации входящих данных - %w", err)
	}

	return item, nil
}