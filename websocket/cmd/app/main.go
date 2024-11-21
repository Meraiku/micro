package main

import (
	"context"
	"os"

	"github.com/joho/godotenv"
	"github.com/meraiku/micro/pkg/logging"
	"github.com/meraiku/micro/websocket/intrenal/app"
)

func main() {
	ctx := context.TODO()
	godotenv.Load(".env")

	var enableLogstash bool

	logAddr := os.Getenv("LOGSTASH_ADDR")
	if logAddr != "" {
		enableLogstash = true
	}

	l := logging.NewLogger(
		logging.WithLevel(logging.LevelDebug),
		logging.WithSource(false),
		logging.WithLogstash(enableLogstash, logAddr),
	)

	app, err := app.New(ctx)
	if err != nil {
		l.Error(
			"failed to create app",
			logging.Err(err),
		)
		return
	}

	if err := app.Run(ctx); err != nil {
		l.Error(
			"failed to run app",
			logging.Err(err),
		)
	}
}
