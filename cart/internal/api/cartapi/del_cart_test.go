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

func TestDelCartTable(t *testing.T) {
	t.Parallel()

	type data struct {
		name        string
		userID      int64
		count       uint16
		errService  error
		wantCodeAns int
	}

	someError := errors.New("Сервисы не нужны")

	testData := []data{
		{
			name:        "Слой API: Некорректный UserID",
			userID:      0,
			count:       1,
			errService:  someError,
			wantCodeAns: http.StatusBadRequest,
		},
		{
			name:        "Слой API: Корректный UserID",
			userID:      1,
			count:       1,
			errService:  nil,
			wantCodeAns: http.StatusNoContent,
		},
		{
			name:        "Слой API: Ошибка удаление корзины",
			userID:      1,
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

			if !errors.Is(tt.errService, someError) {
				servMock.DelCartMock.Expect(minimock.AnyContext, tt.userID).Return(tt.errService)
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
			r := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/user/%d/cart", tt.userID), reader)

			vars := map[string]string{
				model.UsrID: fmt.Sprintf("%d", tt.userID),
			}

			r = mux.SetURLVars(r, vars)

			api.DelCart(w, r)
			assert.Equal(t, tt.wantCodeAns, w.Result().StatusCode)
		})
	}
}
