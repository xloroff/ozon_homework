package cartapi

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-test/deep"
	"github.com/gojuno/minimock/v3"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/model"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/pkg/logger"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/service/cart/mock"
)

func TestCheckoutTable(t *testing.T) {
	t.Parallel()

	type data struct {
		name        string
		userID      int64
		count       uint16
		userOrder   *model.OrderCart
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
			name:        "Слой API: Корзина не существует",
			userID:      1,
			count:       1,
			errService:  model.ErrNotFound,
			wantCodeAns: http.StatusBadRequest,
		},
		{
			name:   "Слой API: Корзина существует",
			userID: 1,
			count:  1,
			userOrder: &model.OrderCart{
				OrderID: 1,
			},
			errService:  nil,
			wantCodeAns: http.StatusOK,
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
				servMock.CheckoutMock.Expect(minimock.AnyContext, tt.userID).Return(tt.userOrder, tt.errService)
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
			r := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/user/%d/checkout", tt.userID), reader)

			vars := map[string]string{
				model.UsrID: fmt.Sprintf("%d", tt.userID),
			}

			r = mux.SetURLVars(r, vars)
			api.Checkout(w, r)
			assert.Equal(t, tt.wantCodeAns, w.Result().StatusCode)

			if w.Result().StatusCode != http.StatusOK {
				return
			}

			bodyAns, err := io.ReadAll(w.Body)
			require.NoError(t, err, "Не получил тело ответа")
			require.Greater(t, len(bodyAns), 0, "Тело ответа слишком маленькое")

			orderResult := &model.OrderCart{}
			err = json.Unmarshal(bodyAns, orderResult)
			require.NoError(t, err, "Тело ответа не соответствует требующейся структуре")

			diff := deep.Equal(tt.userOrder, orderResult)
			if diff != nil {
				t.Errorf("Корзины отличаются: %+v", diff)
			}
		})
	}
}
