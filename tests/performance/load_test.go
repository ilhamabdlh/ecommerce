package performance

import (
	"testing"
	"time"
)

func TestAPIEndpointsPerformance(t *testing.T) {
	rate := vegeta.Rate{Freq: 100, Per: time.Second}
	duration := 5 * time.Second
	targeter := vegeta.NewStaticTargeter(vegeta.Target{
		Method: "GET",
		URL:    "http://localhost:8080/api/v1/products",
	})
	attacker := vegeta.NewAttacker()

	var metrics vegeta.Metrics
	for res := range attacker.Attack(targeter, rate, duration, "Load Test") {
		metrics.Add(res)
	}
	metrics.Close()

	// Assert performance requirements
	if metrics.Latencies.P95 > 500*time.Millisecond {
		t.Errorf("P95 latency too high: %s", metrics.Latencies.P95)
	}
	if metrics.Success < 0.95 {
		t.Errorf("Success rate too low: %.2f", metrics.Success)
	}
}
