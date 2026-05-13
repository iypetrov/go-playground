package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

func AllCollectors() []prometheus.Collector {
	return []prometheus.Collector{
		OutputClientLogs,
	}
}

var (
	namespace = "fluentbit_example"

	// OutputClientLogs is a prometheus metric which keeps logs to the Output Client
	OutputClientLogs = promauto.With(prometheus.DefaultRegisterer).NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Name:      "output_client_logs_total",
		Help:      "Total number of the forwarded logs to the output client",
	}, []string{"host"})
)
