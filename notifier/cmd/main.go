package main

import (
	"context"
	"log"

	notifierapp "gitlab.ozon.dev/xloroff/ozon-hw-go/notifier/internal/app"
)

func main() {
	ctx := context.Background()
	// Запуск приложения.
	if err := notifierapp.NewApp(ctx).Run(); err != nil {
		log.Panicf("Неудалось запустить приложение NOTIFIER - %v", err)
	}
}
