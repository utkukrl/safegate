package metrics

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/utkukrl/safegate/internal/core"
)

var (
	requestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "safegate_requests_total",
			Help: "Total number of requests processed",
		},
		[]string{"method", "endpoint", "status"},
	)

	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "safegate_request_duration_seconds",
			Help:    "Duration of HTTP requests",
			Buckets: prometheus.ExponentialBuckets(0.001, 2, 15),
		},
		[]string{"method", "endpoint"},
	)

	blockedRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "safegate_blocked_requests_total",
			Help: "Number of blocked requests",
		},
		[]string{"method", "endpoint"},
	)
)

func init() {
	prometheus.MustRegister(requestsTotal)
	prometheus.MustRegister(requestDuration)
	prometheus.MustRegister(blockedRequests)
}

func MetricsMiddleware() core.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			rw := &statusResponseWriter{ResponseWriter: w, statusCode: 200}
			next.ServeHTTP(rw, r)

			duration := time.Since(start).Seconds()
			endpoint := r.URL.Path
			method := r.Method
			status := strconv.Itoa(rw.statusCode)

			requestsTotal.WithLabelValues(method, endpoint, status).Inc()
			requestDuration.WithLabelValues(method, endpoint).Observe(duration)

			if rw.statusCode >= 400 {
				blockedRequests.WithLabelValues(method, endpoint).Inc()
			}
		})
	}
}

type statusResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *statusResponseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func MetricsHandler() http.Handler {
	return promhttp.Handler()
}
