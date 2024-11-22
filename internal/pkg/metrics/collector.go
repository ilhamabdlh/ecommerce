package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type MetricsCollector struct {
	requestDuration *prometheus.HistogramVec
	errorCounter    *prometheus.CounterVec
	activeRequests  prometheus.Gauge
	totalRequests   *prometheus.CounterVec
}

func NewMetricsCollector(serviceName string) *MetricsCollector {
	return &MetricsCollector{
		requestDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name: serviceName + "_request_duration_seconds",
				Help: "Duration of requests in seconds",
			},
			[]string{"method", "endpoint", "status"},
		),
		errorCounter: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: serviceName + "_errors_total",
				Help: "Total number of errors",
			},
			[]string{"type"},
		),
		activeRequests: promauto.NewGauge(prometheus.GaugeOpts{
			Name: serviceName + "_active_requests",
			Help: "Number of active requests",
		}),
		totalRequests: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: serviceName + "_requests_total",
				Help: "Total number of requests",
			},
			[]string{"method", "endpoint"},
		),
	}
}
