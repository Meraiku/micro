package app

import (
	"context"
	"os"

	"github.com/meraiku/micro/notification/internal/services/notification"
	"github.com/meraiku/micro/notification/pkg/metrics"
	"github.com/meraiku/micro/pkg/logging"
)

type App struct {
}

func New(ctx context.Context) *App {
	return &App{}
}

func (a *App) Run(ctx context.Context) error {

	log := logging.L(ctx)

	metricAddr := os.Getenv("METRICS_ADDR")

	m := metrics.New(metricAddr)

	log.Info(
		"metrics initialized",
		logging.String("addr", metricAddr),
	)

	go func() {
		if err := m.Run(ctx); err != nil {
			log.Error("failed to run metrics server", logging.Err(err))
		}
	}()

	notificationService, err := notification.New(ctx, m)
	if err != nil {
		return err
	}

	notificationService.Read(ctx)

	return nil
}
