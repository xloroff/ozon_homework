package cart

import (
	"context"
	"errors"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/internal/model"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/internal/pkg/logger"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/internal/service/cart/mock"
)

func TestAddItemTable(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	type fields struct {
		productCliMock *mock.ProductClientMock
		storageMock    *mock.StorageMock
		loggerMock     logger.ILog
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
				f.productCliMock.GetProductMock.Expect(ctx, 1).Return(nil, model.ErrNotFound)
			},
			wantErr: model.ErrNotFound,
		},
		{
			name:   "Слой сервиса: Проблема при добавлении в хранилище",
			userID: 2,
			skuID:  100,
			count:  2,
			prepare: func(f *fields) {
				f.productCliMock.GetProductMock.Expect(ctx, 100).Return(&model.ProductResp{
					Name:  "Тестовый товар",
					Price: 250,
				}, nil)
				f.storageMock.AddItemMock.Return(err1)
			},
			wantErr: err1,
		},
		{
			name:   "Слой сервиса: Продукт добавлен",
			userID: 3,
			skuID:  200,
			count:  1,
			prepare: func(f *fields) {
				f.productCliMock.GetProductMock.Expect(ctx, 200).Return(&model.ProductResp{
					Name:  "Тестовый товар 2",
					Price: 500,
				}, nil)
				f.storageMock.AddItemMock.Return(nil)
			},
			wantErr: nil,
		},
	}

	ctrl := minimock.NewController(t)

	fieldsForTableTest := fields{
		productCliMock: mock.NewProductClientMock(ctrl),
		storageMock:    mock.NewStorageMock(ctrl),
		loggerMock:     logger.InitializeLogger("", 1),
	}

	servM := NewService(fieldsForTableTest.loggerMock, fieldsForTableTest.productCliMock, fieldsForTableTest.storageMock)

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			item := &model.AddItem{}

			item.UsrSkuID.SkuID = tt.skuID
			item.AddItemBody.Count = tt.count
			item.UserIdintyfier.UserID = tt.userID

			tt.prepare(&fieldsForTableTest)

			err := servM.AddItem(ctx, item)
			require.ErrorIs(t, err, tt.wantErr, "Должна быть ошибка", tt.wantErr)
		})
	}
}
