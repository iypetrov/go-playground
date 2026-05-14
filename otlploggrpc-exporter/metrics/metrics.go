package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type PluginMetrics struct {
	OutputClientLogs *prometheus.CounterVec
}

func NewPluginMetrics(reg prometheus.Registerer) *PluginMetrics {
	namespace := "fluentbit_example"
	m := &PluginMetrics{
		// https://github.com/prometheus/client_golang/pull/713
		OutputClientLogs: promauto.With(reg).NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "output_client_logs_total",
			Help:      "Total number of the forwarded logs to the output client",
		}, []string{"host"}),
	}
	return m
}
