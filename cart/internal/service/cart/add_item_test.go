package cart

import (
	"context"
	"errors"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/model"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/pkg/logger"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/service/cart/mock"
)

func TestAddItemTable(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	type fields struct {
		productCliMock *mock.ProductClientMock
		storageMock    *mock.StorageMock
		lomsCliMock    *mock.LomsClientMock
		loggerMock     logger.Logger
	}

	type data struct {
		name    string
		userID  int64
		skuID   int64
		count   uint16
		prepare func(f *fields)
		wantErr error
	}

	err1 := errors.New("bla bla")

	testData := []data{
		{
			name:   "Слой сервиса: Несуществующий продукт 1",
			userID: 1,
			skuID:  1,
			count:  1,
			prepare: func(f *fields) {
				f.productCliMock.GetProductMock.Expect(minimock.AnyContext, 1).Return(nil, model.ErrNotFound)
			},
			wantErr: model.ErrNotFound,
		},
		{
			name:   "Слой сервиса: Проблема при добавлении в хранилище",
			userID: 2,
			skuID:  100,
			count:  2,
			prepare: func(f *fields) {
				f.productCliMock.GetProductMock.Expect(minimock.AnyContext, 100).Return(&model.ProductResp{
					Name:  "Тестовый товар",
					Price: 250,
				}, nil)
				f.storageMock.AddItemMock.Return(err1)
				f.lomsCliMock.StockInfoMock.Return(500, nil)
			},
			wantErr: err1,
		},
		{
			name:   "Слой сервиса: Продукт добавлен",
			userID: 3,
			skuID:  200,
			count:  1,
			prepare: func(f *fields) {
				f.productCliMock.GetProductMock.Expect(minimock.AnyContext, 200).Return(&model.ProductResp{
					Name:  "Тестовый товар 2",
					Price: 500,
				}, nil)
				f.storageMock.AddItemMock.Return(nil)
				f.lomsCliMock.StockInfoMock.Return(500, nil)
			},
			wantErr: nil,
		},
		{
			name:   "Слой сервиса: Продукт не добавлен (недостаточно остатков)",
			userID: 4,
			skuID:  300,
			count:  100,
			prepare: func(f *fields) {
				f.productCliMock.GetProductMock.Expect(minimock.AnyContext, 300).Return(&model.ProductResp{
					Name:  "Тестовый товар 3",
					Price: 500,
				}, nil)
				f.lomsCliMock.StockInfoMock.Return(10, nil)
			},
			wantErr: model.ErrDontHaveReserveCount,
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

			item := &model.AddItem{}

			item.UsrSkuID.SkuID = tt.skuID
			item.AddItemBody.Count = tt.count
			item.UserIdintyfier.UserID = tt.userID

			tt.prepare(f)

			err := s.AddItem(ctx, item)
			require.ErrorIs(t, err, tt.wantErr, "Должна быть ошибка", tt.wantErr)
		})
	}
}
