package cart

import (
	"context"
	"errors"
	"sync"
	"testing"

	"github.com/go-test/deep"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/model"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/pkg/logger"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/service/cart/mock"
)

func TestGetAllUserItemsTable(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	err1 := errors.New("bla bla")
	err2 := errors.New("blo blo")

	type data struct {
		name     string
		userID   int64
		cart     *model.Cart
		stor     *model.FullUserCart
		prodResp map[int64]model.ProductResp
		wantErr  error
	}

	type fields struct {
		productCliMock *mock.ProductClientMock
		storageMock    *mock.StorageMock
		lomsCliMock    *mock.LomsClientMock
		loggerMock     logger.Logger
	}

	testData := []data{
		{
			name:     "Слой сервиса: Корзина не существует",
			cart:     &model.Cart{},
			stor:     nil,
			prodResp: nil,
			userID:   1,
			wantErr:  err1,
		},
		{
			name:     "Слой сервиса: Корзина существует, нет нужного товара",
			cart:     &model.Cart{},
			stor:     nil,
			prodResp: nil,
			userID:   2,
			wantErr:  err2,
		},
		{
			name:   "Слой сервиса: Корзина существует",
			userID: 3,
			cart: &model.Cart{Items: map[int64]*model.CartItem{
				100: {Count: 1},
				200: {Count: 2},
				300: {Count: 3},
			}},
			stor: &model.FullUserCart{
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
			prodResp: map[int64]model.ProductResp{
				100: {
					Name:  "Тестовый товар 1",
					Price: 100,
				},
				200: {
					Name:  "Тестовый товар 2",
					Price: 200,
				},
				300: {
					Name:  "Тестовый товар 3",
					Price: 300,
				},
			},
			wantErr: nil,
		},
	}

	ctrl := minimock.NewController(t)

	newService := func() (*fields, Service) {
		fieldsForTableTest := &fields{
			productCliMock: mock.NewProductClientMock(ctrl),
			storageMock:    mock.NewStorageMock(ctrl),
			loggerMock:     logger.InitializeLogger("", 1),
			lomsCliMock:    mock.NewLomsClientMock(ctrl),
		}
		servM := NewService(fieldsForTableTest.loggerMock, fieldsForTableTest.productCliMock, fieldsForTableTest.lomsCliMock, fieldsForTableTest.storageMock)

		return fieldsForTableTest, servM
	}

	for _, tt := range testData {
		f, s := newService()

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			f.storageMock.GetAllUserItemsMock.
				When(minimock.AnyContext, tt.userID).
				Then(tt.cart, tt.wantErr)

			var wg sync.WaitGroup

			var mu sync.Mutex

			for k := range tt.cart.Items {
				wg.Add(1)

				k := k
				tt := tt

				go func() {
					defer wg.Done()

					resp := tt.prodResp[k]

					mu.Lock()
					f.productCliMock.GetProductMock.
						When(minimock.AnyContext, k).
						Then(&resp, tt.wantErr)
					mu.Unlock()
				}()
			}

			wg.Wait()

			cart, err := s.GetAllUserItems(ctx, tt.userID)
			require.ErrorIs(t, err, tt.wantErr, "Должна быть ошибка", tt.wantErr)

			diff := deep.Equal(cart, tt.stor)
			if diff != nil {
				t.Errorf("Корзины должны совпадать: %+v", diff)
			}
		})
	}
}
