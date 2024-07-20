package cart

import (
	"context"
	"errors"
	"testing"

	"github.com/go-test/deep"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/model"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/pkg/logger"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/service/cart/mock"
)

func TestCheckoutTable(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	err1 := errors.New("bla bla")
	err2 := errors.New("blo blo")

	type data struct {
		name    string
		userID  int64
		cart    *model.Cart
		order   *model.OrderCart
		wantErr error
		lomsErr error
	}

	type fields struct {
		productCliMock *mock.ProductClientMock
		storageMock    *mock.StorageMock
		lomsCliMock    *mock.LomsClientMock
		loggerMock     logger.ILog
	}

	testData := []data{
		{
			name:    "Слой сервиса: Корзина не существует",
			cart:    &model.Cart{},
			order:   nil,
			userID:  1,
			wantErr: err1,
		},
		{
			name:    "Слой сервиса: Корзина существует, нет нужного товара",
			cart:    &model.Cart{},
			order:   nil,
			userID:  2,
			wantErr: err2,
		},
		{
			name:   "Слой сервиса: Корзина существует",
			userID: 3,
			cart: &model.Cart{Items: map[int64]*model.CartItem{
				100: {Count: 1},
				200: {Count: 2},
				300: {Count: 3},
			}},
			order: &model.OrderCart{OrderID: 1},
		},
		{
			name:   "Слой сервиса: Корзина существует (не удалось зарегать заказ в LOMS)",
			userID: 4,
			cart: &model.Cart{Items: map[int64]*model.CartItem{
				100: {Count: 1},
				200: {Count: 2},
				300: {Count: 3},
			}},
			order:   &model.OrderCart{OrderID: 1},
			lomsErr: err1,
		},
	}

	ctrl := minimock.NewController(t)

	fieldsForTableTest := fields{
		productCliMock: mock.NewProductClientMock(ctrl),
		storageMock:    mock.NewStorageMock(ctrl),
		loggerMock:     logger.InitializeLogger("", 1),
		lomsCliMock:    mock.NewLomsClientMock(ctrl),
	}

	servM := NewService(fieldsForTableTest.loggerMock, fieldsForTableTest.productCliMock, fieldsForTableTest.lomsCliMock, fieldsForTableTest.storageMock)

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			fieldsForTableTest.storageMock.GetAllUserItemsMock.
				When(ctx, tt.userID).
				Then(tt.cart, tt.wantErr)

			if tt.order != nil {
				fieldsForTableTest.lomsCliMock.AddOrderMock.
					When(tt.userID, tt.cart).
					Then(1, tt.lomsErr)

				tt.wantErr = tt.lomsErr
			}

			if tt.lomsErr == nil && tt.wantErr == nil {
				fieldsForTableTest.storageMock.DelCartMock.
					When(ctx, tt.userID).
					Then(nil)
			}

			order, err := servM.Checkout(ctx, tt.userID)
			require.ErrorIs(t, err, tt.wantErr, "Должна быть ошибка", tt.wantErr)

			if tt.lomsErr == nil {
				diff := deep.Equal(order, tt.order)
				if diff != nil {
					t.Errorf("Корзины должны совпадать: %+v", diff)
				}
			}
		})
	}
}
