package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/maheshjq/web-analyzer_v1/internal/analyzer"
	"github.com/maheshjq/web-analyzer_v1/internal/models"
)

// AnalyzeHandler godoc
// @Summary Analyze a web page
// @Description Fetches and analyzes a web page by URL, returning information about its structure and content
// @Tags analysis
// @Accept json
// @Produce json
// @Param request body models.AnalysisRequest true "URL to analyze"
// @Success 200 {object} models.AnalysisResponse "Successful analysis"
// @Failure 400 {object} models.ErrorResponse "Bad request (invalid URL format)"
// @Failure 502 {object} models.ErrorResponse "Failed to fetch or analyze the URL"
// @Router /api/analyze [post]
func AnalyzeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req models.AnalysisRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	if req.URL == "" {
		sendErrorResponse(w, http.StatusBadRequest, "URL is required")
		return
	}

	if !strings.HasPrefix(req.URL, "http://") && !strings.HasPrefix(req.URL, "https://") {
		req.URL = "https://" + req.URL
	}

	_, err = url.ParseRequestURI(req.URL)
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "Invalid URL format: "+err.Error())
		return
	}

	analyzerInstance := analyzer.NewAnalyzer()

	// Analyze the url here
	analysisResult, err := analyzerInstance.Analyze(req.URL)
	if err != nil {
		log.Printf("Error analyzing URL %s: %v", req.URL, err)
		sendErrorResponse(w, http.StatusBadGateway, fmt.Sprintf("Failed to analyze URL: %v", err))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(analysisResult)
}

// HealthCheckHandler godoc
// @Summary Health check
// @Description Returns the health status of the API
// @Tags health
// @Produce json
// @Success 200 {object} map[string]string "Service is healthy"
// @Router /api/health [get]
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"status": "ok"}`)
}

func sendErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(models.ErrorResponse{
		StatusCode: statusCode,
		Message:    message,
	})
}
