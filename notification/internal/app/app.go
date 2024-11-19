package app

import (
	"context"
	"log"
)

type App struct {
}

func New(ctx context.Context) *App {
	return &App{}
}

func (a *App) Run(ctx context.Context) error {

	log.Println("Starting notification service")
	return nil
}
