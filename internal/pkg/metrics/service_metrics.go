package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	OrdersCreated = promauto.NewCounter(prometheus.CounterOpts{
		Name: "orders_created_total",
		Help: "The total number of created orders",
	})

	StockTransfers = promauto.NewCounter(prometheus.CounterOpts{
		Name: "stock_transfers_total",
		Help: "The total number of stock transfers between warehouses",
	})

	ProductStockLevel = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "product_stock_level",
			Help: "Current stock level for products",
		},
		[]string{"product_id", "warehouse_id"},
	)

	APIRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "api_request_duration_seconds",
			Help:    "Duration of API requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint", "status"},
	)
)
