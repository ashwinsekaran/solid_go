// Package metrics provides Prometheus instrumentation for the REST service.
//
// It mirrors the metrics setup in worker_pool_api_processer (prom.go + the
// middleware in main.go): a Metrics struct holding the collectors, a Register
// method that builds a dedicated registry, and Middleware that wraps an
// httprouter.Handle to record per-request counts, an in-flight gauge, and a
// latency histogram. The wrapped handlers stay oblivious to metrics — the
// middleware is the only place that touches Prometheus.
package metrics

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/prometheus/client_golang/prometheus"
)

// Metrics bundles the three collectors the middleware records into.
type Metrics struct {
	TotalRequests  prometheus.Counter       // monotonic count of all handled requests
	ReqDuration    *prometheus.HistogramVec // request latency, labelled by method + path
	ActiveRequests *prometheus.GaugeVec     // in-flight requests, labelled by method + path
}

// NewMetrics constructs the collectors with their names, help text, and labels.
func NewMetrics() *Metrics {
	return &Metrics{
		TotalRequests: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "gateway_requests_total",
			Help: "Total no of Http requests",
		}),
		ReqDuration: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name:    "gateway_requests_duration_seconds",
			Help:    "Total Requests duration",
			Buckets: prometheus.DefBuckets,
		},
			[]string{"method", "path"}),
		ActiveRequests: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "gateway_active_requests_total",
			Help: "Total Active requests",
		},
			[]string{"method", "path"}),
	}
}

// Register creates a dedicated registry and registers every collector on it.
// The returned registry is served at /metrics via promhttp in main.go.
func (m *Metrics) Register() *prometheus.Registry {
	reg := prometheus.NewRegistry()
	reg.MustRegister(m.TotalRequests)
	reg.MustRegister(m.ReqDuration)
	reg.MustRegister(m.ActiveRequests)

	return reg
}

// Middleware wraps next, recording the in-flight gauge (incremented on entry,
// decremented on exit via defer), the latency histogram (observed on exit), and
// the total-request counter — all before delegating to the real handler.
func Middleware(m *Metrics, next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		m.ActiveRequests.With(prometheus.Labels{"method": r.Method, "path": r.URL.Path}).Inc()
		defer m.ActiveRequests.With(prometheus.Labels{"method": r.Method, "path": r.URL.Path}).Dec()

		timer := prometheus.NewTimer(prometheus.ObserverFunc(func(duration float64) {
			m.ReqDuration.With(prometheus.Labels{"method": r.Method, "path": r.URL.Path}).Observe(duration)
		}))
		defer timer.ObserveDuration()

		m.TotalRequests.Inc()
		next(w, r, p)
	}
}
