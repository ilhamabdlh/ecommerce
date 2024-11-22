package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	OrdersTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "business_orders_total",
			Help: "Total number of orders by status",
		},
		[]string{"status"},
	)

	OrderValue = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "business_order_value",
			Help:    "Distribution of order values",
			Buckets: prometheus.LinearBuckets(10, 50, 10), // 10, 60, 110, ...
		},
		[]string{"status"},
	)

	ProductStock = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "business_product_stock",
			Help: "Current stock level by product",
		},
		[]string{"product_id", "warehouse_id"},
	)

	UserRegistrations = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "business_user_registrations_total",
			Help: "Total number of user registrations",
		},
		[]string{"source"},
	)

	ActiveUsers = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "business_active_users",
			Help: "Number of active users in the last 24 hours",
		},
	)
)
