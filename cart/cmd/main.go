package main

import (
	"context"
	"log"

	cartapp "gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/app"
)

func main() {
	ctx := context.Background()
	// Запуск приложения.
	if err := cartapp.NewApp(ctx).Run(); err != nil {
		log.Panicf("Неудалось запустить приложение Cart - %v", err)
	}
}
