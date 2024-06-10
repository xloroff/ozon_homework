package orderservice

import (
	"context"
	"testing"

	"github.com/go-test/deep"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/model"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/pkg/logger"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/service/order/mock"
)

func TestInfoTable(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	type fields struct {
		stockStorageMock *mock.StockStorageMock
		orderStorageMock *mock.OrderStorageMock
		loggerMock       logger.ILog
	}

	type data struct {
		name        string
		orderID     int64
		orderStore  *model.Order
		getOrderErr error
		wantErr     error
	}

	testData := []data{
		{
			name:        "Заказ отсутствует",
			orderID:     1,
			orderStore:  &model.Order{},
			getOrderErr: model.ErrOrderNotFound,
			wantErr:     model.ErrOrderNotFound,
		},
		{
			name:    "Заказ есть",
			orderID: 2,
			orderStore: &model.Order{
				ID:     2,
				User:   12,
				Status: model.OrderStatusNew,
				Items: []*model.OrderItem{
					{
						Sku:   1212,
						Count: 2,
					},
				},
			},
		},
	}

	ctrl := minimock.NewController(t)

	fieldsForTableTest := fields{
		orderStorageMock: mock.NewOrderStorageMock(ctrl),
		stockStorageMock: mock.NewStockStorageMock(ctrl),
		loggerMock:       logger.InitializeLogger("", 1),
	}

	servO := NewService(ctx, fieldsForTableTest.loggerMock, fieldsForTableTest.orderStorageMock, fieldsForTableTest.stockStorageMock)

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			fieldsForTableTest.orderStorageMock.GetOrderMock.
				When(tt.orderID).Then(tt.orderStore, tt.getOrderErr)

			ord, err := servO.Info(tt.orderID)
			if tt.wantErr != nil {
				require.NotNil(t, err, "Должна быть ошибка")
				require.ErrorIs(t, err, tt.wantErr, "Должна быть ошибка", tt.wantErr)
			} else {
				require.Nil(t, err, "Ошибки быть не должно")

				diff := deep.Equal(tt.orderStore, ord)
				if diff != nil {
					t.Errorf("Заказы должны совпасть: %+v", diff)
				}
			}
		})
	}
}
