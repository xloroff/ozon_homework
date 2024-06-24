package cart

import (
	"context"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/model"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/pkg/logger"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/service/cart/mock"
)

func TestDelItemTable(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	type data struct {
		name    string
		userID  int64
		skuID   int64
		wantErr error
		memdel  bool
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
			userID:  1,
			skuID:   100,
			wantErr: model.ErrUnknownError,
			memdel:  false,
		},
		{
			name:    "Слой сервиса: Корзина cуществует в ней есть нужный товар",
			userID:  2,
			skuID:   200,
			wantErr: nil,
			memdel:  true,
		},
		{
			name:    "Слой сервиса: Корзина cуществует в ней отсутствует нужный товар",
			userID:  3,
			skuID:   300,
			wantErr: model.ErrUnknownError,
			memdel:  false,
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

			item := model.DelItem{}
			item.SkuID = tt.skuID
			item.UserID = tt.userID

			f.storageMock.DelItemMock.
				When(ctx, &item).
				Then(tt.memdel)

			err := s.DelItem(ctx, &item)
			require.ErrorIs(t, err, tt.wantErr, "Должна быть ошибка", tt.wantErr)
		})
	}
}
