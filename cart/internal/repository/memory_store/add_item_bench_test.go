package memorystore

import (
	"context"
	"testing"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/internal/model"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/internal/pkg/logger"
)

func BenchmarkAddItemOneToOne(b *testing.B) {
	ctx := context.Background()
	l := logger.InitializeLogger("", 1)
	storage := NewCartStorage(l)

	b.StopTimer()

	b.Run("Бенчмарк: Добавляем один итем одному пользователю", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			item := &model.AddItem{}

			item.UsrSkuID.SkuID = 100
			item.AddItemBody.Count = 1
			item.UserIdintyfier.UserID = 1

			b.StartTimer()

			err := storage.AddItem(ctx, item)
			if err != nil {
				b.Fail()
			}
		}
	})
}

func BenchmarkAddItemOneToMany(b *testing.B) {
	ctx := context.Background()
	l := logger.InitializeLogger("", 1)
	storage := NewCartStorage(l)

	b.StopTimer()

	b.Run("Бенчмарк: Добавляем один итем многим пользователям", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			item := &model.AddItem{}

			item.UsrSkuID.SkuID = 100
			item.AddItemBody.Count = 1
			item.UserIdintyfier.UserID = int64(i)

			b.StartTimer()

			err := storage.AddItem(ctx, item)
			if err != nil {
				b.Fail()
			}
		}
	})
}

func BenchmarkAddItemManyToOne(b *testing.B) {
	ctx := context.Background()
	l := logger.InitializeLogger("", 1)
	storage := NewCartStorage(l)

	b.StopTimer()

	b.Run("Добавляем разные итемы одному пользователям", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			item := &model.AddItem{}

			item.UsrSkuID.SkuID = int64(i)
			item.AddItemBody.Count = 1
			item.UserIdintyfier.UserID = 1

			b.StartTimer()

			err := storage.AddItem(ctx, item)
			if err != nil {
				b.Fail()
			}
		}
	})
}

func BenchmarkAddItemManyToMany(b *testing.B) {
	ctx := context.Background()
	l := logger.InitializeLogger("", 1)
	storage := NewCartStorage(l)

	b.StopTimer()

	b.Run("Добавляем разные итемы многим пользоваелям", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			item := &model.AddItem{}

			item.UsrSkuID.SkuID = int64(i)
			item.AddItemBody.Count = 1
			item.UserIdintyfier.UserID = int64(i)

			b.StartTimer()

			err := storage.AddItem(ctx, item)
			if err != nil {
				b.Fail()
			}
		}
	})
}

func BenchmarkAddItemManyToManyWithCount(b *testing.B) {
	ctx := context.Background()
	l := logger.InitializeLogger("", 1)
	storage := NewCartStorage(l)

	b.StopTimer()

	b.Run("Добавляем разные итемы многим пользоваелям с разным количеством", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			item := &model.AddItem{}

			item.UsrSkuID.SkuID = int64(i)
			item.AddItemBody.Count = uint16(i)
			item.UserIdintyfier.UserID = int64(i)

			b.StartTimer()

			err := storage.AddItem(ctx, item)
			if err != nil {
				b.Fail()
			}
		}
	})
}
