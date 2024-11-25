package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	RequestCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "shop_service_requests_total",
			Help: "Total number of requests to the shop service",
		},
		[]string{"method", "endpoint", "status"},
	)

	RequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "shop_service_request_duration_seconds",
			Help:    "Duration of requests to the shop service",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)

	ShopOperations = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "shop_service_operations_total",
			Help: "Total number of shop operations",
		},
		[]string{"operation"},
	)

	WarehouseAssociations = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "shop_service_warehouse_associations_total",
			Help: "Total number of warehouse associations/dissociations",
		},
		[]string{"operation", "shop_id"},
	)
)
