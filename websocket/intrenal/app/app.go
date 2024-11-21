package app

import (
	"context"
	"os"

	"github.com/meraiku/micro/pkg/logging"
	"github.com/meraiku/micro/pkg/metrics"
	"github.com/meraiku/micro/websocket/intrenal/config"
	v1 "github.com/meraiku/micro/websocket/intrenal/controllers/http/v1"
	"github.com/meraiku/micro/websocket/intrenal/repo/chatRepo/memory"
	"github.com/meraiku/micro/websocket/intrenal/services/auth"
	"github.com/meraiku/micro/websocket/intrenal/services/chat"
)

type App struct {
	cfg         config.Config
	metrics     *metrics.Metric
	chatService v1.ChatService
	authService v1.AuthService
}

func New(ctx context.Context) (*App, error) {
	a := &App{}

	if err := a.initDeps(ctx); err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) initDeps(ctx context.Context) error {

	deps := []func(context.Context) error{
		a.initConfig,
		a.initChatService,
		a.initAuthService,
		a.initMetrics,
	}

	for _, dep := range deps {
		if err := dep(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initConfig(ctx context.Context) error {
	config.Load(".env")

	cfg := config.New()

	a.cfg = cfg

	return nil
}

func (a *App) initChatService(ctx context.Context) error {

	repo := memory.NewRepository()

	cs := chat.NewService(ctx, repo)

	a.chatService = cs

	return nil
}

func (a *App) initAuthService(ctx context.Context) error {

	au, err := auth.New(ctx, a.cfg.Service.Auth.Addr)
	if err != nil {
		return err
	}

	a.authService = au

	return nil
}

func (a *App) initMetrics(ctx context.Context) error {
	metricsAddr := os.Getenv("METRICS_ADDR")

	m := metrics.New(metricsAddr)

	a.metrics = &m

	go func() {
		if err := a.metrics.Run(ctx); err != nil {
			logging.L(ctx).Info("starting service without metrics", logging.Err(err))
		}
	}()

	return nil
}

func (a *App) Run(ctx context.Context) error {

	chatService := v1.NewChatServiceAPI(
		ctx,
		a.cfg.HTTP.Addr(),
		a.chatService,
		a.authService,
	)

	return chatService.Run(ctx)
}
