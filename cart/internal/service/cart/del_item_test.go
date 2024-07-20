package cart

import (
	"context"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/internal/model"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/internal/pkg/logger"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/internal/service/cart/mock"
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
	productCliMock := mock.NewProductClientMock(ctrl)
	storageMock := mock.NewStorageMock(ctrl)
	l := logger.InitializeLogger("", 1)

	servM := NewService(l, productCliMock, storageMock)

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			item := model.DelItem{}
			item.SkuID = tt.skuID
			item.UserID = tt.userID

			storageMock.DelItemMock.
				When(ctx, &item).
				Then(tt.memdel)

			err := servM.DelItem(ctx, &item)
			require.ErrorIs(t, err, tt.wantErr, "Должна быть ошибка", tt.wantErr)
		})
	}
}
