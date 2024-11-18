package app

import (
	"context"

	v1 "github.com/meraiku/micro/websocket/intrenal/controllers/http/v1"
)

type App struct {
}

func New(ctx context.Context) *App {
	return &App{}
}

func (a *App) Run(ctx context.Context) error {
	chatService := v1.NewChatServiceAPI(":8080")

	return chatService.Run(ctx)
}
