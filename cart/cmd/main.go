package main

import (
	"context"
	"log"
	_ "net/http/pprof" // nolint:gosec // Это нормально, судя по линту, тут нужно акцептнуть.

	cartapp "gitlab.ozon.dev/xloroff/ozon-hw-go/cart/internal/app"
)

func main() {
	ctx := context.Background()
	// Запуск приложения.
	if err := cartapp.NewApp(ctx).Run(); err != nil {
		log.Panicf("Неудалось запустить приложение Cart - %v", err)
	}
}
