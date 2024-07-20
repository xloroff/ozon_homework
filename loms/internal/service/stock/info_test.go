package stockservice

import (
	"context"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/model"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/pkg/logger"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/service/stock/mock"
)

func TestInfoTable(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	type fields struct {
		stockStorageMock *mock.StockStorageMock
		loggerMock       logger.Logger
	}

	type data struct {
		name    string
		skuID   int64
		count   uint16
		wantErr error
	}

	testData := []data{
		{
			name:    "Товар отсутствует в остатках",
			skuID:   1,
			wantErr: model.ErrReserveNotFound,
		},
		{
			name:  "Товар есть в остатках",
			skuID: 2958025,
			count: 10,
		},
	}

	ctrl := minimock.NewController(t)

	for _, tt := range testData {
		fieldsForTableTest := fields{
			stockStorageMock: mock.NewStockStorageMock(ctrl),
			loggerMock:       logger.InitializeLogger("", 1),
		}

		servS := NewService(ctx, fieldsForTableTest.loggerMock, fieldsForTableTest.stockStorageMock)

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			fieldsForTableTest.stockStorageMock.GetAvailableForReserveMock.
				When(minimock.AnyContext, tt.skuID).
				Then(tt.count, tt.wantErr)

			count, err := servS.Info(ctx, tt.skuID)
			require.ErrorIs(t, err, tt.wantErr, "Должна быть ошибка", tt.wantErr)
			require.Equal(t, count, tt.count, "Число остатков должно соответствовать.")
		})
	}
}
