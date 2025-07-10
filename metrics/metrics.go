/*
Package metrics provides Prometheus metrics mechanisms
*/
package metrics

import (
	"net/http"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// LatencyVar holds prometheus histogram.
type LatencyVar struct {
	ReqLatency *prometheus.HistogramVec
}

// CounterVar holds prometheus counter.
type CounterVar struct {
	ReqCounter *prometheus.CounterVec
}

// Handler returns the global http.Handler that provides
// the prometheus metrics format on GET requests.
func Handler() http.Handler {
	return promhttp.Handler()
}

// NewTimer returns the global prometheus.NewTimer which creates a new Timer
func NewTimer(o prometheus.Observer) *prometheus.Timer {
	return prometheus.NewTimer(o)
}

// CounterWithLabels is a constructor that initiates the appropriate
// prometheus vector counter accepting multiple labels.
func CounterWithLabels(counterName string, labels ...string) CounterVar {
	var counter CounterVar
	counter.ReqCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: counterName,
			Help: "partitioned by " + strings.Join(labels, ", "),
		},
		labels,
	)
	prometheus.Register(counter.ReqCounter)

	return counter
}

// LatencyHistogram is a constructor that initiates
// the appropriate prometheus histogram.
func LatencyHistogram(latencyName string) LatencyVar {
	var histogram LatencyVar
	histogram.ReqLatency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    latencyName,
			Help:    "A histogram of latencies for requests.",
			Buckets: []float64{.01, .05, 0.1, 0.25, 0.5, 1},
		},
		[]string{"endpoint", "code", "method"},
	)
	prometheus.MustRegister(histogram.ReqLatency)

	return histogram
}

// MetricDecorator is a decorator function that decorates http.handlerFunc.
// Counter automatically detects the status code. If no status
// code is explicitly set, status code 200 is assumed.
func MetricDecorator(myHandler http.Handler, l LatencyVar, c CounterVar, handlerName string) http.HandlerFunc {

	return promhttp.InstrumentHandlerDuration(l.ReqLatency.MustCurryWith(
		prometheus.Labels{"endpoint": handlerName}), promhttp.InstrumentHandlerCounter(c.ReqCounter, myHandler))
}
