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

func TestCreateTable(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	err1 := errors.New("bla bla")
	err2 := errors.New("blo blo")
	err3 := errors.New("ble ble")

	type fields struct {
		stockStorageMock *mock.StockStorageMock
		orderStorageMock *mock.OrderStorageMock
		loggerMock       logger.Logger
	}

	type data struct {
		name            string
		items           *model.AllNeedReserve
		userID          int64
		status          string
		addOrderErr     error
		orderID         int64
		setStatusErr    error
		reserveStockErr error
		wantErr         error
	}

	testData := []data{
		{
			name:   "Ошибка при создании заказа",
			userID: 1,
			items: &model.AllNeedReserve{
				&model.NeedReserve{
					Sku:   1,
					Count: 1,
				},
			},
			addOrderErr: err1,
			wantErr:     err1,
		},
		{
			name:   "Корректное создание заказа",
			userID: 2,
			items: &model.AllNeedReserve{
				&model.NeedReserve{
					Sku:   2,
					Count: 2,
				},
			},
			status:  model.OrderStatusAwaitingPayment,
			orderID: 1,
		},
		{
			name:   "Ошибка резервирования остатков",
			userID: 3,
			items: &model.AllNeedReserve{
				&model.NeedReserve{
					Sku:   3,
					Count: 3,
				},
			},
			status:          model.OrderStatusFailed,
			reserveStockErr: err2,
			wantErr:         err2,
		},
		{
			name:   "Ошибка резервирования остатков и ошибка смены статуса",
			userID: 4,
			items: &model.AllNeedReserve{
				&model.NeedReserve{
					Sku:   4,
					Count: 4,
				},
			},
			status:          model.OrderStatusFailed,
			reserveStockErr: err3,
			setStatusErr:    err2,
			wantErr:         err2,
		},
		{
			name: "Успешное резервирования остатков, но ошибка смены статуса",
			items: &model.AllNeedReserve{
				&model.NeedReserve{
					Sku:   5,
					Count: 5,
				},
			},
			status:       model.OrderStatusAwaitingPayment,
			setStatusErr: err3,
			wantErr:      err3,
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
		fieldsForTableTest.orderStorageMock.AddOrderMock.
			Expect(minimock.AnyContext, tt.userID, resevToOrders(*tt.items)).Return(tt.orderID, tt.addOrderErr)

		fieldsForTableTest.stockStorageMock.AddReserveMock.
			Expect(minimock.AnyContext, *tt.items).Return(tt.reserveStockErr)

		fieldsForTableTest.orderStorageMock.SetStatusMock.
			Expect(minimock.AnyContext, tt.orderID, tt.status).Return(tt.setStatusErr)

		t.Run(tt.name, func(t *testing.T) {
			orderID, err := servO.Create(ctx, tt.userID, *tt.items)
			if tt.wantErr != nil {
				require.NotNil(t, err, "Должна быть ошибка")
				require.ErrorIs(t, err, tt.wantErr, "Не та ошибка")
			} else {
				require.Nil(t, err, "Не должно быть ошибки")
				require.Equal(t, tt.orderID, orderID, "Не совпал номер заказа")
			}
		})
	}
}
