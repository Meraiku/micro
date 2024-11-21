package metrics

import (
	"github.com/meraiku/micro/pkg/metrics"
)

type NotifMetric struct {
	metrics.Metric
}

func New(addr string) *NotifMetric {

	return &NotifMetric{
		metrics.New(addr),
	}
}
