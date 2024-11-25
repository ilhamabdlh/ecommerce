package utils

import (
	"time"

	"github.com/sony/gobreaker"
)

var WarehouseBreaker *gobreaker.CircuitBreaker

func InitCircuitBreaker() {
	settings := gobreaker.Settings{
		Name:        "warehouse-service",
		MaxRequests: 3,
		Interval:    10 * time.Second,
		Timeout:     60 * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			return counts.Requests >= 3 && failureRatio >= 0.6
		},
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			Logger.Infof("Circuit Breaker state changed from %v to %v", from, to)
		},
	}

	WarehouseBreaker = gobreaker.NewCircuitBreaker(settings)
}
