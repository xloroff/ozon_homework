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

func TestGetAllUserItemsTable(t *testing.T) {
	t.Parallel()

	type data struct {
		name        string
		userID      int64
		count       uint16
		userCart    *model.FullUserCart
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
			wantCodeAns: http.StatusNotFound,
		},
		{
			name:   "Слой API: Корзина существует",
			userID: 1,
			count:  1,
			userCart: &model.FullUserCart{
				Items: []*model.UserCartItem{
					{
						SkuID: 100,
						Count: 1,
						Name:  "Тестовый товар 1",
						Price: 100,
					},
					{
						SkuID: 200,
						Count: 2,
						Name:  "Тестовый товар 2",
						Price: 200,
					},
					{
						SkuID: 300,
						Count: 3,
						Name:  "Тестовый товар 3",
						Price: 300,
					},
				},
				TotalPrice: 1400,
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
				servMock.GetAllUserItemsMock.Expect(minimock.AnyContext, tt.userID).Return(tt.userCart, tt.errService)
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
			r := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/user/%d/cart/list", tt.userID), reader)

			vars := map[string]string{
				model.UsrID: fmt.Sprintf("%d", tt.userID),
			}

			r = mux.SetURLVars(r, vars)
			api.GetAllUserItems(w, r)
			assert.Equal(t, tt.wantCodeAns, w.Result().StatusCode)

			if w.Result().StatusCode != http.StatusOK {
				return
			}

			bodyAns, err := io.ReadAll(w.Body)
			require.NoError(t, err, "Не получил тело ответа")
			require.Greater(t, len(bodyAns), 0, "Тело ответа слишком маленькое")

			carResult := &model.FullUserCart{}
			err = json.Unmarshal(bodyAns, carResult)
			require.NoError(t, err, "Тело ответа не соответствует требующейся структуре")

			diff := deep.Equal(tt.userCart, carResult)
			if diff != nil {
				t.Errorf("Корзины отличаются: %+v", diff)
			}
		})
	}
}
