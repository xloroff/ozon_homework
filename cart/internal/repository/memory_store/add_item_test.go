package memorystore

import (
	"context"
	"testing"

	"github.com/go-test/deep"
	"github.com/stretchr/testify/require"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/internal/model"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/internal/pkg/logger"
)

func TestAddItemTable(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	type data struct {
		name        string
		userID      int64
		skuID       int64
		count       uint16
		wantStorage map[int64]*model.Cart
	}

	testData := []data{
		{
			name:   "Слой памяти: Добавление продукта пользователю",
			userID: 1,
			skuID:  100,
			count:  1,
			wantStorage: map[int64]*model.Cart{
				1: {
					Items: map[int64]*model.CartItem{
						100: {
							Count: 1,
						},
					},
				},
			},
		},
		{
			name:   "Слой памяти: увеличение числа продукта пользователю",
			userID: 1,
			skuID:  100,
			count:  1,
			wantStorage: map[int64]*model.Cart{
				1: {
					Items: map[int64]*model.CartItem{
						100: {
							Count: 2,
						},
					},
				},
			},
		},
	}

	storage := cartStorage{
		data:   map[int64]*model.Cart{},
		logger: logger.InitializeLogger("", 1),
	}

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			item := &model.AddItem{}

			item.UsrSkuID.SkuID = tt.skuID
			item.AddItemBody.Count = tt.count
			item.UserIdintyfier.UserID = tt.userID

			err := storage.AddItem(ctx, item)
			require.NoError(t, err)

			diff := deep.Equal(storage.data, tt.wantStorage)
			if diff != nil {
				t.Errorf("Хранилища должны совпадать: %+v", diff)
			}
		})
	}
}
