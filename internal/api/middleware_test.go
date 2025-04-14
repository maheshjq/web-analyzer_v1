package api

import (
	"bytes"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoggingMiddleware(t *testing.T) {
	var logBuffer bytes.Buffer
	testLogger := slog.New(slog.NewJSONHandler(&logBuffer, nil))

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	loggingHandler := LoggingMiddleware(testLogger)(nextHandler)

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("User-Agent", "test-agent")

	recorder := httptest.NewRecorder()

	loggingHandler.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)

	logOutput := logBuffer.String()
	assert.Contains(t, logOutput, "Request started")
	assert.Contains(t, logOutput, "Request completed")
	assert.Contains(t, logOutput, "GET")
	assert.Contains(t, logOutput, "/test")
	assert.Contains(t, logOutput, "test-agent")
}

func TestRecoverMiddleware(t *testing.T) {
	var logBuffer bytes.Buffer
	testLogger := slog.New(slog.NewJSONHandler(&logBuffer, nil))

	panicHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("test panic")
	})

	safeHandler := RecoverMiddleware(testLogger)(panicHandler)

	req := httptest.NewRequest("GET", "/panic", nil)

	recorder := httptest.NewRecorder()

	safeHandler.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "Internal Server Error")

	logOutput := logBuffer.String()
	assert.Contains(t, logOutput, "Panic recovered")
	assert.Contains(t, logOutput, "test panic")
}

func TestCorsMiddleware(t *testing.T) {
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	corsHandler := CorsMiddleware(nextHandler)

	t.Run("Normal Request", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test", nil)
		recorder := httptest.NewRecorder()

		corsHandler.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Equal(t, "OK", recorder.Body.String())
		assert.Equal(t, "*", recorder.Header().Get("Access-Control-Allow-Origin"))
		assert.Equal(t, "GET, POST, OPTIONS", recorder.Header().Get("Access-Control-Allow-Methods"))
	})

	t.Run("OPTIONS Request", func(t *testing.T) {
		req := httptest.NewRequest("OPTIONS", "/test", nil)
		recorder := httptest.NewRecorder()

		corsHandler.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Equal(t, "", recorder.Body.String())
		assert.Equal(t, "*", recorder.Header().Get("Access-Control-Allow-Origin"))
	})
}
