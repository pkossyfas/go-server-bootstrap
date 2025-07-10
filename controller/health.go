package controller

import (
	"net/http"

	"github.com/pkossyfas/go-server-bootstrap/metrics"
)

var healthEndpointLatency = metrics.LatencyHistogram("health_endpoint_latency")
var healthEndpointCounter = metrics.CounterWithLabels("health_endpoint_requests_total", "code", "method")

// HealthEndpointMetrics is a metric wrapper for the health_endpoint endpoint.
var HealthEndpointMetrics = metrics.MetricDecorator(HealthEndpoint, healthEndpointLatency, healthEndpointCounter, "/health")

// HealthEndpoint implements the controller for /health endpoint.
var HealthEndpoint = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	// Allow only Get method
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)

		return
	}

	w.WriteHeader(200)
})
