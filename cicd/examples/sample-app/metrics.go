package main

import "github.com/prometheus/client_golang/prometheus"

// metrics groups the four Prometheus metric types used by the workshop:
// Counter, Gauge, Histogram, and Summary
type metrics struct {
	httpRequestsTotal   *prometheus.CounterVec
	activeRequests      prometheus.Gauge
	httpRequestDuration *prometheus.HistogramVec
	httpRequestSummary  *prometheus.SummaryVec
}

func newMetrics(reg prometheus.Registerer) *metrics {
	m := &metrics{
		httpRequestsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "Total HTTP requests by path, method, and status.",
			},
			[]string{"path", "method", "status"},
		),
		activeRequests: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "http_active_requests",
			Help: "Number of in-flight HTTP requests.",
		}),
		httpRequestDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_duration_seconds",
				Help:    "HTTP request duration in seconds (histogram).",
				Buckets: []float64{0.01, 0.05, 0.1, 0.5, 1, 5},
			},
			[]string{"path"},
		),
		httpRequestSummary: prometheus.NewSummaryVec(
			prometheus.SummaryOpts{
				Name:       "http_request_summary_seconds",
				Help:       "HTTP request duration in seconds (summary with client-side quantiles).",
				Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
			},
			[]string{"path"},
		),
	}
	reg.MustRegister(m.httpRequestsTotal, m.activeRequests, m.httpRequestDuration, m.httpRequestSummary)
	return m
}
