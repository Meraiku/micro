package app

import (
	"context"
	"net"
	"os"

	v1 "github.com/meraiku/micro/websocket/intrenal/controllers/http/v1"
)

type App struct {
}

func New(ctx context.Context) *App {
	return &App{}
}

func (a *App) Run(ctx context.Context) error {
	port := os.Getenv("PORT")
	if port == "" {
		port = "2000"
	}
	host := os.Getenv("HOST")

	chatService := v1.NewChatServiceAPI(ctx, net.JoinHostPort(host, port))

	return chatService.Run(ctx)
}
