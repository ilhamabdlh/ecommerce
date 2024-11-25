package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	RequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "warehouse_requests_total",
			Help: "The total number of requests",
		},
		[]string{"method", "endpoint", "status"},
	)

	RequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "warehouse_request_duration_seconds",
			Help:    "Request duration in seconds",
			Buckets: []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
		},
		[]string{"method", "endpoint"},
	)

	StockLevels = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "warehouse_stock_levels",
			Help: "Current stock levels by warehouse and product",
		},
		[]string{"warehouse_id", "product_id"},
	)

	TransferOperations = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "warehouse_transfer_operations_total",
			Help: "The total number of stock transfer operations",
		},
		[]string{"status"},
	)

	ActiveWarehouses = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "warehouse_active_total",
			Help: "The total number of active warehouses",
		},
	)
)
