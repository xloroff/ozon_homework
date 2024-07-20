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

func TestPayTable(t *testing.T) {
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
		name         string
		orderID      int64
		orderStore   *model.Order
		status       string
		getOrderErr  error
		setStatusErr error
		delItemErr   error
		wantErr      error
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
			name:    "Нормальная оплата",
			orderID: 2,
			status:  model.OrderStatusPayed,
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
			status:  model.OrderStatusPayed,
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
			delItemErr: err1,
			wantErr:    err1,
		},
		{
			name:    "Ошибка смены статуса",
			orderID: 4,
			status:  model.OrderStatusPayed,
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
				Expect(minimock.AnyContext, tt.orderID).Return(tt.orderStore, tt.getOrderErr)

			delReserve := orderToReserve(tt.orderStore.Items)

			if len(tt.orderStore.Items) > 0 {
				fieldsForTableTest.stockStorageMock.DelItemFromReserveMock.
					Expect(minimock.AnyContext, delReserve).Return(tt.delItemErr)

				if tt.delItemErr == nil {
					fieldsForTableTest.orderStorageMock.SetStatusMock.
						Expect(minimock.AnyContext, tt.orderID, model.OrderStatusPayed).Return(tt.setStatusErr)
				}
			}

			err := servO.Pay(ctx, tt.orderID)
			if tt.wantErr != nil {
				require.NotNil(t, err, "Должна быть ошибка")
				require.ErrorIs(t, err, tt.wantErr, "Должна быть ошибка", tt.wantErr)
			} else {
				require.Nil(t, err, "Не должно быть ошибки")
			}
		})
	}
}
