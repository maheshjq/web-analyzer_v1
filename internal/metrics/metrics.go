// internal/metrics/metrics.go
package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	RequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "web_analyzer_requests_total",
			Help: "Total number of requests by endpoint and status",
		},
		[]string{"endpoint", "status"},
	)

	RequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "web_analyzer_request_duration_seconds",
			Help:    "Request duration in seconds by endpoint",
			Buckets: prometheus.LinearBuckets(0.1, 0.1, 10), // 0.1s to 1s
		},
		[]string{"endpoint"},
	)

	AnalysisCount = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "web_analyzer_analysis_count_total",
			Help: "Total number of web pages analyzed",
		},
	)

	AnalysisDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "web_analyzer_analysis_duration_seconds",
			Help:    "Duration of web page analysis in seconds",
			Buckets: prometheus.LinearBuckets(0.5, 0.5, 10), // 0.5s to 5s
		},
	)

	// CacheHitCount tracks cache hits
	CacheHitCount = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "web_analyzer_cache_hits_total",
			Help: "Total number of cache hits",
		},
	)

	// CacheMissCount tracks cache misses
	CacheMissCount = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "web_analyzer_cache_misses_total",
			Help: "Total number of cache misses",
		},
	)
)

func Initialize() {
	prometheus.MustRegister(
		RequestsTotal,
		RequestDuration,
		AnalysisCount,
		AnalysisDuration,
		CacheHitCount,
		CacheMissCount,
	)
}

func MetricsHandler() http.Handler {
	return promhttp.Handler()
}
