package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	RequestCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "order_service_requests_total",
			Help: "Total number of requests to the order service",
		},
		[]string{"method", "endpoint", "status"},
	)

	RequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "order_service_request_duration_seconds",
			Help:    "Duration of requests to the order service",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)

	OrderOperations = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "order_service_operations_total",
			Help: "Total number of order operations",
		},
		[]string{"operation", "status"},
	)

	OrderStatusTransitions = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "order_service_status_transitions_total",
			Help: "Total number of order status transitions",
		},
		[]string{"from_status", "to_status"},
	)

	OrderTotalAmount = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "order_service_total_amount",
			Help:    "Distribution of order total amounts",
			Buckets: []float64{10, 50, 100, 500, 1000, 5000},
		},
		[]string{"status"},
	)

	StockCheckErrors = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "order_service_stock_check_errors_total",
			Help: "Total number of stock check errors",
		},
	)
)
