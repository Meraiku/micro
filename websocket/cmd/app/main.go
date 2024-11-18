package main

import (
	"context"
	"log"

	"github.com/joho/godotenv"
	"github.com/meraiku/micro/websocket/intrenal/app"
)

func main() {
	ctx := context.TODO()
	godotenv.Load(".env")

	app := app.New(ctx)
	if err := app.Run(ctx); err != nil {
		log.Fatalf("failed to run app: %v", err)
	}
}
