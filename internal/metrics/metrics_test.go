// internal/metrics/metrics_test.go
package metrics

import (
	"net/http/httptest"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
)

func TestMetricsInitialization(t *testing.T) {
	prometheus.DefaultRegisterer = prometheus.NewRegistry()

	// Initialize metrics without error
	Initialize()

	handler := MetricsHandler()
	assert.NotNil(t, handler, "Metrics handler should not be nil")

	// Simple test to increment counters
	RequestsTotal.WithLabelValues("/test", "200").Inc()
	AnalysisCount.Inc()
	CacheHitCount.Inc()

	RequestDuration.WithLabelValues("/test").Observe(0.5)
	AnalysisDuration.Observe(1.0)

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/metrics", nil)

	// Serve the request
	handler.ServeHTTP(recorder, request)

	assert.Equal(t, 200, recorder.Code, "Metrics endpoint should return 200 OK")
}
