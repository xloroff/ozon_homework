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
		storage.data = map[int64]*model.Cart{
			tt.userID: {
				Items: map[int64]*model.CartItem{
					tt.skuID: {
						Count: tt.count,
					},
				},
			},
		}

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := storage.DelCart(ctx, tt.userID)
			require.NoError(t, err)

			diff := deep.Equal(storage.data, tt.wantStorage)
			if diff != nil {
				t.Errorf("Хранилища должны совпадать: %+v", diff)
			}
		})
	}
}

func TestDelCartConcurrent(t *testing.T) {
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

		storage.data[userID] = &model.Cart{
			Items: model.CartItems{
				100: &model.CartItem{Count: 1},
			},
		}

		var wg sync.WaitGroup

		numGoroutines := 100

		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)

			go func() {
				defer wg.Done()

				err := storage.DelCart(ctx, userID)
				require.NoError(t, err)
			}()
		}

		wg.Wait()

		storage.Lock()
		defer storage.Unlock()

		_, ok := storage.data[userID]
		require.False(t, ok, "Не должно быть корзины")
	})
}
