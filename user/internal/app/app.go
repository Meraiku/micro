package app

import (
	"context"
	"os"

	"github.com/meraiku/micro/user/internal/config"
	"github.com/meraiku/micro/user/internal/domain/user/memory"
	"github.com/meraiku/micro/user/internal/service/user"
)

type API interface {
	Run() error
}

type App struct {
	userService *user.Service
	api         API
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
		a.initConfig,
		a.initUserService,
		a.initAPI,
	}

	for _, dep := range deps {
		if err := dep(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initConfig(_ context.Context) error {

	config.Load()

	return nil
}

func (a *App) initUserService(_ context.Context) error {
	var repo user.Repository

	switch os.Getenv("USER_REPO") {
	default:
		memoryRepo := memory.New()
		repo = memoryRepo
	}

	a.userService = user.New(repo)
	return nil
}

func (a *App) initAPI(_ context.Context) error {

	transport := os.Getenv("API")
	if transport == "" {
		transport = "REST"
	}

	switch transport {
	case "REST":
		a.api = newRestService()
	case "GRPC":
		a.api = newGRPCService()
	}

	return nil
}

func (a *App) Run() error {
	return a.api.Run()
}

func setupUserRepository(repoType string) user.Repository {
	var repo user.Repository

	switch repoType {
	default:
		memoryRepo := memory.New()
		repo = memoryRepo
	}

	return repo
}
