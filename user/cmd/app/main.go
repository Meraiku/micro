package main

import (
	"context"
	"log"

	"github.com/joho/godotenv"
	"github.com/meraiku/micro/user/internal/app"
)

func main() {
	godotenv.Load(".env")

	ctx := context.TODO()

	app, err := app.New(ctx)
	if err != nil {
		log.Fatalf("failed to create app: %s", err)
	}

	if err := app.Run(ctx); err != nil {
		log.Fatalf("failed to start app: %s", err)
	}
}
