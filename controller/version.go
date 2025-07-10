package controller

import (
	"encoding/json"
	"net/http"

	"github.com/pkossyfas/go-server-bootstrap/metrics"
)

var (
	// AppVersion defines the version of the application. By default
	// is set to unset and it should be overwritten during compilation.
	AppVersion = "unset"
)

// GetVersion defines the response for the version request.
type GetVersion struct {
	Version string `json:"version"`
}

var versionLatency = metrics.LatencyHistogram("version_latency")
var versionCounter = metrics.CounterWithLabels("version_requests_total", "code", "method")

// VersionEndpointMetrics is a metric wrapper for the version endpoint.
var VersionEndpointMetrics = metrics.MetricDecorator(VersionEndpoint, versionLatency, versionCounter, "/version")

// VersionEndpoint responds with the app version.
var VersionEndpoint = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	version := GetVersion{}
	version.Version = AppVersion
	json.NewEncoder(w).Encode(version)
})
