package memorystore

import (
	"context"
	"testing"

	"github.com/go-test/deep"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/model"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/pkg/logger"
)

func TestDelItem(t *testing.T) {
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
		delItemUsr  int64
		DelItemSku  int64
		wantStorage map[int64]*model.Cart
	}

	testData := []data{
		{
			name: "Слой памяти: Удаление итема",
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
			delItemUsr: 1,
			DelItemSku: 100,
			wantStorage: map[int64]*model.Cart{
				1: {
					Items: map[int64]*model.CartItem{
						200: {
							Count: 1,
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
			for _, addItem := range tt.items {
				storage.data[addItem.userID] = &model.Cart{
					Items: map[int64]*model.CartItem{
						addItem.skuID: {
							Count: addItem.count,
						},
					},
				}
			}

			del := model.DelItem{}
			del.UserID = tt.delItemUsr
			del.SkuID = tt.DelItemSku

			if !storage.DelItem(ctx, &del) {
				t.Error("Ошибки удаления корзины быть не должно")
			}

			diff := deep.Equal(storage.data, tt.wantStorage)
			if diff != nil {
				t.Errorf("Хранилища должны совпадать: %+v", diff)
			}
		})
	}
}
