package memorystore

import (
	"context"
	"testing"

	"github.com/go-test/deep"
	"github.com/stretchr/testify/require"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/model"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/pkg/logger"
)

func TestDelCart(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	type data struct {
		name          string
		delCartuserID int64
		userID        int64
		skuID         int64
		count         uint16
		wantStorage   map[int64]*model.Cart
	}

	testData := []data{
		{
			name:          "Слой памяти: Удаление корзины - корзина существует",
			userID:        1,
			delCartuserID: 1,
			skuID:         100,
			count:         1,
			wantStorage:   map[int64]*model.Cart{},
		},
		{
			name:          "Слой памяти: Удаление корзины - корзина не существует",
			userID:        1,
			delCartuserID: 2,
			skuID:         100,
			count:         1,
			wantStorage:   map[int64]*model.Cart{},
		},
	}

	storage := cartStorage{
		data:   map[int64]*model.Cart{},
		logger: logger.InitializeLogger("", 1),
	}

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			storage.data = map[int64]*model.Cart{
				tt.userID: {
					Items: map[int64]*model.CartItem{
						tt.skuID: {
							Count: tt.count,
						},
					},
				},
			}

			err := storage.DelCart(ctx, tt.userID)
			require.NoError(t, err)

			diff := deep.Equal(storage.data, tt.wantStorage)
			if diff != nil {
				t.Errorf("Хранилища должны совпадать: %+v", diff)
			}
		})
	}
}
