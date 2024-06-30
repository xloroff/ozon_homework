package cartapi

import (
	"bufio"
	"bytes"
	"encoding/json"
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

func TestAddItemTable(t *testing.T) {
	t.Parallel()

	type data struct {
		name        string
		inputJSON   string
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
			name:        "Слой API: Некорректное число товаров",
			userID:      1,
			skuID:       1,
			count:       0,
			errService:  someError,
			wantCodeAns: http.StatusBadRequest,
		},
		{
			name:        "Слой API: Кривой JSON",
			userID:      1,
			skuID:       1,
			count:       1,
			inputJSON:   "bla bla",
			errService:  someError,
			wantCodeAns: http.StatusBadRequest,
		},
		{
			name:        "Слой API: Кривой продукт",
			userID:      1,
			skuID:       1,
			count:       1,
			errService:  model.ErrNotFound,
			wantCodeAns: http.StatusPreconditionFailed,
		},
		{
			name:        "Слой API: Нормальный продукт",
			userID:      1,
			skuID:       1,
			count:       1,
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

			item := &model.AddItem{}

			item.UsrSkuID.SkuID = tt.skuID
			item.AddItemBody.Count = tt.count
			item.UserIdintyfier.UserID = tt.userID

			if !errors.Is(tt.errService, someError) {
				servMock.AddItemMock.When(minimock.AnyContext, item).Then(tt.errService)
			}

			var jsonReq []byte

			var err error

			if tt.inputJSON == "" {
				request := &model.AddItemBody{Count: tt.count}
				jsonReq, err = json.Marshal(request)
				require.NoError(t, err)
			} else {
				jsonReq = []byte(tt.inputJSON)
			}

			var body bytes.Buffer
			bodyWriter := bufio.NewWriter(&body)

			_, err = bodyWriter.Write(jsonReq)
			require.NoError(t, err)

			err = bodyWriter.Flush()
			require.NoError(t, err)

			reader := bufio.NewReader(&body)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/user/%d/cart/%d", tt.userID, tt.skuID), reader)

			vars := map[string]string{
				model.UsrID: fmt.Sprintf("%d", tt.userID),
				model.SkuID: fmt.Sprintf("%d", tt.skuID),
			}

			r = mux.SetURLVars(r, vars)

			api.AddItem(w, r)
			assert.Equal(t, tt.wantCodeAns, w.Result().StatusCode)
		})
	}
}
