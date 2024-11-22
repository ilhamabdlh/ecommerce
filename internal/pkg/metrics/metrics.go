package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	RequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "http_request_duration_seconds",
			Help: "Duration of HTTP requests",
		},
		[]string{"path", "method", "status"},
	)

	ActiveConnections = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "active_connections",
		Help: "Number of active connections",
	})

	TotalOrders = promauto.NewCounter(prometheus.CounterOpts{
		Name: "total_orders",
		Help: "Total number of orders created",
	})
)
