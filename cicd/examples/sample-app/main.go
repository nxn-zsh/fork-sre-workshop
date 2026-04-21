package main

import (
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// version is set at build time via ldflags.
var version = "dev"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	reg := prometheus.NewRegistry()
	m := newMetrics(reg)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /", record(m, "/", handleRoot))
	mux.HandleFunc("GET /health", record(m, "/health", handleHealth))
	mux.HandleFunc("GET /version", record(m, "/version", handleVersion))
	mux.Handle("GET /metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}))

	slog.Info("server starting", "port", port, "version", version)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		slog.Error("server failed to start", "error", err)
		os.Exit(1)
	}
}

// record wraps a handler with Prometheus instrumentation: it increments
// the in-flight gauge, times the request with both a histogram and a
// summary, and counts the final request by path/method/status.
func record(m *metrics, path string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m.activeRequests.Inc()
		defer m.activeRequests.Dec()

		start := time.Now()
		sr := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next(sr, r)
		d := time.Since(start).Seconds()

		m.httpRequestDuration.WithLabelValues(path).Observe(d)
		m.httpRequestSummary.WithLabelValues(path).Observe(d)
		m.httpRequestsTotal.WithLabelValues(path, r.Method, strconv.Itoa(sr.status)).Inc()
	}
}

// statusRecorder captures the response status code so it can be used as
// a metric label.
type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}
