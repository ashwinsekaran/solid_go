package main

import (
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Metrics struct {
	TotalRequests  prometheus.Counter
	ReqDuration    *prometheus.HistogramVec
	ActiveRequests *prometheus.GaugeVec
}

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

func (m *Metrics) Register() *prometheus.Registry {
	reg := prometheus.NewRegistry()
	reg.MustRegister(m.TotalRequests)
	reg.MustRegister(m.ReqDuration)
	reg.MustRegister(m.ActiveRequests)

	return reg
}

// "github.com/prometheus/client_golang/prometheus"
// "github.com/prometheus/client_golang/prometheus/promhttp"
type XXX struct {
	HistCount *prometheus.HistogramVec
}

func NewXXX() *XXX {
	return &XXX{
		HistCount: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name:    "app_hist_count",
			Help:    "Describe Histogram count",
			Buckets: prometheus.DefBuckets,
		},
			[]string{"xxx"},
		),
	}
}
func (x *XXX) Register() *prometheus.Registry {
	reg := prometheus.NewRegistry()
	reg.MustRegister(x.HistCount)
	return reg
}
func mains() {
	start := time.Now()
	m := NewXXX()
	reg := m.Register()
	r := httprouter.New()
	r.Handler("GET", "/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
	//m.ReqCount.With(prometheus.Labels{"aaa": "enter_value_here"}).Inc()
	defer func() {
		duration := time.Since(start).Seconds()
		m.HistCount.With(prometheus.Labels{"xxx": "enter_value_here"}).Observe(duration)
	}()
}
