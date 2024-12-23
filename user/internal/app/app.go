package app

import (
	"context"
	"os"

	"github.com/meraiku/micro/pkg/logging"
	"github.com/meraiku/micro/pkg/metrics"
	"github.com/meraiku/micro/user/internal/config"
	"github.com/meraiku/micro/user/internal/domain/user/memory"
	"github.com/meraiku/micro/user/internal/service/user"
)

type API interface {
	Run(ctx context.Context) error
}

type App struct {
	api     API
	metrics *metrics.Metric
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
		a.initLogger,
		a.initMetrics,
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

	var enabled bool

	logAddr := os.Getenv("LOGSTASH_ADDR")
	if logAddr == "" {
		enabled = false
	}

	log := logging.NewLogger(
		logging.WithLevel(logging.LevelDebug),
		logging.WithSource(false),
		logging.WithLogstash(enabled, logAddr),
	)

	log.Info("logger initialized")

	return nil
}

func (a *App) initMetrics(_ context.Context) error {
	metricsAddr := os.Getenv("METRICS_ADDR")

	m := metrics.New(metricsAddr)

	a.metrics = &m
	return nil
}

func (a *App) initConfig(_ context.Context) error {

	config.Load()

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
	go func() {
		if err := a.metrics.Run(ctx); err != nil {
			logging.L(ctx).Info("starting service without metrics", logging.Err(err))
		}
	}()

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
