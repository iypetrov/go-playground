package metrics

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
)

var (
	reg     *prometheus.Registry
	regOnce sync.Once
)

func NewRegistry() *prometheus.Registry {
	regOnce.Do(func() {
		promreg := prometheus.NewRegistry()
		promreg.MustRegister(
			collectors.NewGoCollector(),
			collectors.NewProcessCollector(
				collectors.ProcessCollectorOpts{},
			),
		)
		reg = promreg
	})

	return reg
}
