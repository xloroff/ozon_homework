package memorystore

import (
	"context"
	"sync"
	"testing"

	"github.com/go-test/deep"
	"github.com/stretchr/testify/require"

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

	for _, tt := range testData {
		storage := cartStorage{
			data:   map[int64]*model.Cart{},
			logger: logger.InitializeLogger("", 1),
		}

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

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

func TestDelItemConcurrent(t *testing.T) {
	t.Parallel()

	t.Run("Конкурентное удаление корзины", func(t *testing.T) {
		t.Parallel()

		loggerMock := logger.InitializeLogger("", 1)
		storage := &cartStorage{
			data:   make(map[int64]*model.Cart),
			logger: loggerMock,
		}
		ctx := context.Background()
		userID := int64(1)
		skuID1 := int64(100)
		skuID2 := int64(101)

		storage.data[userID] = &model.Cart{
			Items: model.CartItems{
				100: &model.CartItem{Count: 1},
				101: &model.CartItem{Count: 1},
			},
		}

		var wg sync.WaitGroup

		numGoroutines := 100

		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)

			go func() {
				defer wg.Done()

				iOne := &model.DelItem{}

				iOne.UserID = userID
				iOne.SkuID = skuID1

				err := storage.DelItem(ctx, iOne)
				require.True(t, err)
			}()

			wg.Add(1)

			go func() {
				defer wg.Done()

				iTwo := &model.DelItem{}

				iTwo.UserID = userID
				iTwo.SkuID = skuID2

				err := storage.DelItem(ctx, iTwo)
				require.True(t, err)
			}()
		}

		wg.Wait()

		storage.Lock()
		defer storage.Unlock()

		_, ok := storage.data[userID]
		require.False(t, ok, "Не должно быть корзины")
	})
}
