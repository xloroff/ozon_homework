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

func TestAddItemTable(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	type data struct {
		userID int64
		skuID  int64
		count  uint16
	}

	type tests struct {
		name        string
		addItm      []data
		wantStorage map[int64]*model.Cart
	}

	testData := []tests{
		{
			name: "Слой памяти: Добавление продукта пользователю",
			addItm: []data{
				{
					userID: 1,
					skuID:  100,
					count:  1,
				},
			},
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
			name: "Слой памяти: увеличение числа продукта пользователю",
			addItm: []data{
				{
					userID: 2,
					skuID:  100,
					count:  1,
				},
				{
					userID: 2,
					skuID:  100,
					count:  2,
				},
			},
			wantStorage: map[int64]*model.Cart{
				2: {
					Items: map[int64]*model.CartItem{
						100: {
							Count: 3,
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

			for _, itm := range tt.addItm {
				item := &model.AddItem{}

				item.UsrSkuID.SkuID = itm.skuID
				item.AddItemBody.Count = itm.count
				item.UserIdintyfier.UserID = itm.userID

				err := storage.AddItem(ctx, item)
				require.NoError(t, err)
			}

			diff := deep.Equal(storage.data, tt.wantStorage)
			if diff != nil {
				t.Errorf("Хранилища должны совпадать: %+v", diff)
			}
		})
	}
}

func TestAddItemConcurrent(t *testing.T) {
	t.Parallel()

	t.Run("Конкурентное добавление продукта пользователю", func(t *testing.T) {
		t.Parallel()

		loggerMock := logger.InitializeLogger("", 1)

		storage := &cartStorage{
			data:   make(map[int64]*model.Cart),
			logger: loggerMock,
		}

		ctx := context.Background()
		userID := int64(1)
		skuID := int64(100)
		itemCount := 1

		var wg sync.WaitGroup

		numGoroutines := 100

		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)

			go func() {
				defer wg.Done()

				it := &model.AddItem{}

				it.UserID = userID
				it.SkuID = skuID
				it.Count = uint16(itemCount)

				err := storage.AddItem(ctx, it)
				require.NoError(t, err)
			}()
		}

		wg.Wait()

		storage.Lock()
		defer storage.Unlock()

		cart, ok := storage.data[userID]
		require.True(t, ok, "Корзина не должна быть пустой")
		require.NotNil(t, cart)

		cartItem, ok := cart.Items[skuID]
		require.True(t, ok, "Должны быть товары в корзине")
		require.NotNil(t, cartItem)
		require.Equal(t, uint16(numGoroutines*itemCount), cartItem.Count, "Суммарное число итемов должно быть равно ЧислоРутин * Количество итемов на старте")
	})
}
