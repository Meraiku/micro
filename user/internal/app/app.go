package app

import (
	"context"
	"os"

	v1 "github.com/meraiku/micro/user/internal/controller/rest/v1"
	"github.com/meraiku/micro/user/internal/domain/user/memory"
	"github.com/meraiku/micro/user/internal/service/user"
)

type App struct {
	userService *user.Service
}

func New(ctx context.Context) (*App, error) {
	a := &App{}

	if err := a.initDeps(ctx); err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) initDeps(ctx context.Context) error {

	deps := []func(ctx context.Context) error{
		a.initUserService,
	}

	for _, dep := range deps {
		if err := dep(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initUserService(ctx context.Context) error {
	var repo user.Repository

	switch os.Getenv("USER_REPO") {
	default:
		memoryRepo := memory.New()
		repo = memoryRepo
	}

	a.userService = user.New(repo)
	return nil
}

func (a *App) Run() error {

	transport := os.Getenv("API")

	switch transport {
	case "REST":
		return a.runRestAPI()
	case "GRPC":
		return nil
	default:
		return a.runRestAPI()
	}

}

func (a *App) runRestAPI() error {
	api := v1.New(a.userService)
	return api.Run()
}
