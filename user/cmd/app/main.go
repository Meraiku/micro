package main

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/meraiku/micro/user/internal/app"
)

func main() {
	godotenv.Load(".env")

	app, err := app.New()
	if err != nil {
		log.Fatalf("failed to create app: %s", err)
	}

	fmt.Println("Starting app...")
	if err := app.Run(); err != nil {
		log.Fatalf("failed to start app: %s", err)
	}
}
