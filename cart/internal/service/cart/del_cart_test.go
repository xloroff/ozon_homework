package cart

import (
	"context"
	"errors"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/pkg/logger"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/service/cart/mock"
)

func TestDelCartTable(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	err1 := errors.New("bla bla")

	type data struct {
		name    string
		userID  int64
		wantErr error
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
			wantErr: err1,
		},
		{
			name:   "Слой сервиса: Корзина существуют",
			userID: 2,
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
			fieldsForTableTest.storageMock.DelCartMock.
				When(ctx, tt.userID).
				Then(tt.wantErr)

			err := servM.DelCart(ctx, tt.userID)
			require.ErrorIs(t, err, tt.wantErr, "Должна быть ошибка", tt.wantErr)
		})
	}
}
