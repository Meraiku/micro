package main

import (
	"context"
	"log"

	"github.com/joho/godotenv"
	"github.com/meraiku/micro/notification/internal/app"
	"github.com/meraiku/micro/pkg/logging"
)

func main() {

	ctx := context.TODO()

	godotenv.Load()

	l := logging.NewLogger(
		logging.WithLevel(logging.LevelDebug),
		logging.WithSource(false),
	)

	ctx = logging.ContextWithLogger(ctx, l)

	app := app.New(ctx)
	if err := app.Run(ctx); err != nil {
		log.Fatalf("failed to run app: %v", err)
	}
}
