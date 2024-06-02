package memorystore

import (
	"context"
	"testing"

	"github.com/go-test/deep"
	"github.com/stretchr/testify/require"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/internal/model"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/internal/pkg/logger"
)

func TestGetAllUserItems(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	type itm struct {
		userID int64
		skuID  int64
		count  uint16
	}

	type data struct {
		name        string
		items       []itm
		getItemUsr  int64
		wantStorage *model.Cart
	}

	testData := []data{
		{
			name: "Слой памяти: Получение всех итемов пользователя",
			items: []itm{
				{
					userID: 1,
					skuID:  100,
					count:  1,
				},
				{
					userID: 1,
					skuID:  200,
					count:  1,
				},
			},
			getItemUsr: 1,
			wantStorage: &model.Cart{
				Items: map[int64]*model.CartItem{
					200: {
						Count: 1,
					},
					100: {
						Count: 1,
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
			for _, addItem := range tt.items {
				if _, ok := storage.data[addItem.userID]; !ok {
					storage.data[addItem.userID] = &model.Cart{Items: map[int64]*model.CartItem{}}
				}

				storage.data[addItem.userID].Items[addItem.skuID] = &model.CartItem{
					Count: addItem.count,
				}
			}

			usrStorage, err := storage.GetAllUserItems(ctx, tt.getItemUsr)
			require.NoError(t, err)

			diff := deep.Equal(usrStorage, tt.wantStorage)
			if diff != nil {
				t.Errorf("Хранилища должны совпадать: %+v", diff)
			}
		})
	}
}
