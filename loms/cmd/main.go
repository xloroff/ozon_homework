package main

import (
	"context"
	"log"

	lomsapp "gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/app"
)

func main() {
	ctx := context.Background()
	// Запуск приложения.
	if err := lomsapp.NewApp(ctx).Run(); err != nil {
		log.Panicf("Неудалось запустить приложение LOMS - %v", err)
	}
}
