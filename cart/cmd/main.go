package main

import (
	"context"
	"log"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/internal/app"
)

func main() {
	ctx := context.Background()
	// Запуск приложения.
	if err := app.NewApp(ctx).Run(); err != nil {
		log.Fatalf("Неудалось запустить приложение: %v", err)
	}
}
