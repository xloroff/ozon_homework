package orderservice

import (
	"context"
	"errors"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/model"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/pkg/logger"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/service/order/mock"
)

func TestCancelTable(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	err1 := errors.New("bla bla")
	err2 := errors.New("blo blo")

	type fields struct {
		stockStorageMock *mock.StockStorageMock
		orderStorageMock *mock.OrderStorageMock
		loggerMock       logger.Logger
	}

	type data struct {
		name             string
		orderID          int64
		orderStore       *model.Order
		status           string
		getOrderErr      error
		setStatusErr     error
		cancelReserveErr error
		wantErr          error
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
			name:    "Нормальная отмена",
			orderID: 2,
			status:  model.OrderStatusCancelled,
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
		{
			name:    "Резерв не снимается",
			orderID: 3,
			orderStore: &model.Order{
				ID:     3,
				User:   13,
				Status: model.OrderStatusNew,
				Items: []*model.OrderItem{
					{
						Sku:   1213,
						Count: 3,
					},
				},
			},
			cancelReserveErr: err1,
			wantErr:          err1,
		},
		{
			name:    "Ошибка смены статуса",
			orderID: 4,
			status:  model.OrderStatusCancelled,
			orderStore: &model.Order{
				ID:     4,
				User:   14,
				Status: model.OrderStatusNew,
				Items: []*model.OrderItem{
					{
						Sku:   1214,
						Count: 4,
					},
				},
			},
			setStatusErr: err2,
			wantErr:      err2,
		},
	}

	ctrl := minimock.NewController(t)

	for _, tt := range testData {
		fieldsForTableTest := fields{
			orderStorageMock: mock.NewOrderStorageMock(ctrl),
			stockStorageMock: mock.NewStockStorageMock(ctrl),
			loggerMock:       logger.InitializeLogger("", 1),
		}

		servO := NewService(ctx, fieldsForTableTest.loggerMock, fieldsForTableTest.orderStorageMock, fieldsForTableTest.stockStorageMock)

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			fieldsForTableTest.orderStorageMock.GetOrderMock.
				When(minimock.AnyContext, tt.orderID).Then(tt.orderStore, tt.getOrderErr)

			delReserve := orderToReserve(tt.orderStore.Items)

			if len(tt.orderStore.Items) > 0 {
				fieldsForTableTest.stockStorageMock.CancelReserveMock.
					When(minimock.AnyContext, delReserve).Then(tt.cancelReserveErr)

				if tt.cancelReserveErr == nil {
					fieldsForTableTest.orderStorageMock.SetStatusMock.
						When(minimock.AnyContext, tt.orderID, model.OrderStatusCancelled).Then(tt.setStatusErr)
				}
			}

			err := servO.Cancel(ctx, tt.orderID)
			if tt.wantErr != nil {
				require.NotNil(t, err, "Должна быть ошибка")
				require.ErrorIs(t, err, tt.wantErr, "Должна быть ошибка", tt.wantErr)
			} else {
				require.Nil(t, err, "Не должно быть ошибки")
			}
		})
	}
}
