package controller

import (
	"context"
	"fmt"
	"net/http"

	db "github.com/pkossyfas/go-server-bootstrap/dbconnector"
	"github.com/pkossyfas/go-server-bootstrap/logger"
	"github.com/pkossyfas/go-server-bootstrap/metrics"
)

var readyEndpointLatency = metrics.LatencyHistogram("ready_endpoint_latency")
var readyEndpointCounter = metrics.CounterWithLabels("ready_endpoint_requests_total", "code", "method")

// ReadyEndpointMetrics is a metric wrapper for the ready_endpoint endpoint.
var ReadyEndpointMetrics = metrics.MetricDecorator(ReadyEndpoint, readyEndpointLatency, readyEndpointCounter, "/ready")

// ReadyEndpoint implements the controller for /ready endpoint.
var ReadyEndpoint = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	// Allow only Get method
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)

		return
	}

	if db.DBPool == nil {
		logger.Error(fmt.Errorf("db connection initialized has failed"), "readiness probe failed")
		w.WriteHeader(503)

		return
	}

	err := db.PingDB(context.Background(), db.DBPool)

	if err != nil {
		logger.Error(err, "readiness probe failed")
		w.WriteHeader(503)

		return
	}

	w.WriteHeader(200)
})
