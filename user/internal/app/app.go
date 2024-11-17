package app

import (
	"context"
	"os"

	"github.com/meraiku/micro/pkg/logging"
	"github.com/meraiku/micro/user/internal/config"
	"github.com/meraiku/micro/user/internal/domain/user/memory"
	"github.com/meraiku/micro/user/internal/service/user"
)

type API interface {
	Run(ctx context.Context) error
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
		a.initLogger,
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

func (a *App) initLogger(_ context.Context) error {

	log := logging.NewLogger(
		logging.WithLevel(logging.LevelDebug),
		logging.WithSource(false),
	)

	log.Info("logger initialized")

	return nil
}

func (a *App) initConfig(_ context.Context) error {

	config.Load()

	return nil
}

func (a *App) initUserService(ctx context.Context) error {
	var repo user.Repository

	repoEnv := os.Getenv("USER_REPO")
	if repoEnv == "" {
		repoEnv = "memory"
	}

	switch repoEnv {
	case "memory":
		memoryRepo := memory.New()
		repo = memoryRepo
	}

	logging.L(ctx).Info(
		"user service initialized",
		logging.String("repo", repoEnv),
	)

	a.userService = user.New(repo)
	return nil
}

func (a *App) initAPI(ctx context.Context) error {

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

	logging.L(ctx).Info(
		"api initialized",
		logging.String("transport", transport),
	)

	return nil
}

func (a *App) Run(ctx context.Context) error {
	return a.api.Run(ctx)
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
