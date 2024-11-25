package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	RequestCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "product_service_requests_total",
			Help: "Total number of requests to the product service",
		},
		[]string{"method", "endpoint", "status"},
	)

	RequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "product_service_request_duration_seconds",
			Help:    "Duration of requests to the product service",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)

	ProductOperations = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "product_service_operations_total",
			Help: "Total number of product operations",
		},
		[]string{"operation"},
	)
)
