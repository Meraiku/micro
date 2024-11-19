package app

import (
	"context"
	"log"

	"github.com/meraiku/micro/notification/internal/services/notification"
)

type App struct {
}

func New(ctx context.Context) *App {
	return &App{}
}

func (a *App) Run(ctx context.Context) error {

	log.Println("Starting notification service")

	notificationService, err := notification.New(ctx)
	if err != nil {
		return err
	}

	notificationService.Read()

	return nil
}
