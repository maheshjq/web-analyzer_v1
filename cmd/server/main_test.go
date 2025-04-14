package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestHealthEndpoint(t *testing.T) {
	originalPort := os.Getenv("PORT")
	defer os.Setenv("PORT", originalPort)

	os.Setenv("PORT", "8081")

	router := mux.NewRouter()
	apiRouter := router.PathPrefix("/api").Subrouter()

	apiRouter.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	}).Methods("GET")

	server := httptest.NewServer(router)
	defer server.Close()

	// test health
	resp, err := http.Get(server.URL + "/api/health")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// test JSON
	contentType := resp.Header.Get("Content-Type")
	assert.Equal(t, "application/json", contentType)
}
