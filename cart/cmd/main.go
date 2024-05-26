package main

import (
	"context"
	"log"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/internal/app"
)

func main() {
	ctx := context.Background()
	cart := app.NewApp(ctx)
	if err := cart.Run(); err != nil {
		log.Fatalf("Неудалось запустить приложение: %v", err)
	}
}