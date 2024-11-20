package metrics

import (
	"context"

	"github.com/meraiku/micro/pkg/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type NotifMetric struct {
	m metrics.Metric

	counterNotif prometheus.Counter
}

func New(addr string) *NotifMetric {
	nc := promauto.NewCounter(prometheus.CounterOpts{
		Name: "user_notification_total",
		Help: "The total number of notifications",
	})

	return &NotifMetric{
		m:            metrics.New(addr),
		counterNotif: nc,
	}
}

func (n *NotifMetric) Run(ctx context.Context) error {

	return n.m.Run(ctx)
}

func (n *NotifMetric) IncrNotifications() {
	n.counterNotif.Inc()
}
