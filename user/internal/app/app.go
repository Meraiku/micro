package app

import (
	v1 "github.com/meraiku/micro/user/internal/controller/rest/v1"
	"github.com/meraiku/micro/user/internal/domain/user/memory"
	"github.com/meraiku/micro/user/internal/service/user"
)

type App struct {
}

func New() (*App, error) {
	return &App{}, nil
}

func (a *App) Run() error {
	memoryRepo := memory.New()
	userService := user.New(memoryRepo)
	api := v1.New(userService)

	return api.Run()
}
