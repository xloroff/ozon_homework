package cartapi

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/model"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/pkg/logger"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/service/cart/mock"
)

func TestDelItemTable(t *testing.T) {
	t.Parallel()

	type data struct {
		name        string
		userID      int64
		skuID       int64
		count       uint16
		errService  error
		wantCodeAns int
	}

	someError := errors.New("Сервисы не нужны")

	testData := []data{
		{
			name:        "Слой API: Некорректный UserID",
			userID:      0,
			skuID:       100,
			count:       1,
			errService:  someError,
			wantCodeAns: http.StatusBadRequest,
		},
		{
			name:        "Слой API: Некорректный SkuID",
			userID:      1,
			skuID:       0,
			count:       1,
			errService:  someError,
			wantCodeAns: http.StatusBadRequest,
		},
		{
			name:        "Слой API: Нормальный продукт",
			userID:      1,
			skuID:       1,
			count:       1,
			errService:  nil,
			wantCodeAns: http.StatusNoContent,
		},
		{
			name:        "Слой API: Ошибка удаления продукта",
			userID:      1,
			skuID:       1,
			count:       1,
			errService:  errors.New("Ошибка удаления"),
			wantCodeAns: http.StatusBadRequest,
		},
	}

	ctrl := minimock.NewController(t)

	for _, tt := range testData {
		servMock := mock.NewServiceMock(ctrl)
		l := logger.InitializeLogger("", 1)

		api := API{
			cartService: servMock,
			logger:      l,
		}

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			item := &model.DelItem{}

			item.UsrSkuID.SkuID = tt.skuID
			item.UserIdintyfier.UserID = tt.userID

			if !errors.Is(tt.errService, someError) {
				servMock.DelItemMock.Expect(minimock.AnyContext, item).Return(tt.errService)
			}

			var err error

			var body bytes.Buffer
			bodyWriter := bufio.NewWriter(&body)

			_, err = bodyWriter.Write([]byte{})
			require.NoError(t, err)

			err = bodyWriter.Flush()
			require.NoError(t, err)

			reader := bufio.NewReader(&body)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/user/%d/cart/%d", tt.userID, tt.skuID), reader)

			vars := map[string]string{
				model.UsrID: fmt.Sprintf("%d", tt.userID),
				model.SkuID: fmt.Sprintf("%d", tt.skuID),
			}

			r = mux.SetURLVars(r, vars)

			api.DelItem(w, r)
			assert.Equal(t, tt.wantCodeAns, w.Result().StatusCode)
		})
	}
}
