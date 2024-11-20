package metrics

import (
	"context"
	"errors"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	ErrEmptyAddr = errors.New("empty addr")
)

type Metric struct {
	addr string
}

func New(addr string) Metric {
	return Metric{
		addr: addr,
	}
}

func (m *Metric) Run(ctx context.Context) error {
	if m.addr == "" {
		return ErrEmptyAddr
	}

	return http.ListenAndServe(m.addr, promhttp.Handler())
}
