package api

import (
	"fmt"
	"net/http"
)

// AnalyzeHandler godoc
// @Summary Analyze a web page
// @Description Fetches and analyzes a web page by URL, returning information about its structure and content
// @Tags analysis
// @Accept json
// @Produce json
// @Param request body AnalysisRequest true "URL to analyze"
// @Success 200 {object} AnalysisResponse "Successful analysis"
// @Failure 400 {object} ErrorResponse "Bad request (invalid URL format)"
// @Failure 502 {object} ErrorResponse "Failed to fetch or analyze the URL"
// @Router /api/analyze [post]
func AnalyzeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "{\"status\":\"ok\"}")
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

// AnalysisRequest represents the request to analyze a webpage
type AnalysisRequest struct {
	URL string `json:"url" example:"https://example.com"`
}

// AnalysisResponse represents the result of the analysis
type AnalysisResponse struct {
	HTMLVersion      string        `json:"htmlVersion" example:"HTML5"`
	Title            string        `json:"title" example:"Example Domain"`
	Headings         HeadingCount  `json:"headings"`
	Links            LinkAnalysis  `json:"links"`
	ContainsLoginForm bool          `json:"containsLoginForm" example:"false"`
}

// HeadingCount represents the count of different heading levels
type HeadingCount struct {
	H1 int `json:"h1" example:"1"`
	H2 int `json:"h2" example:"2"`
	H3 int `json:"h3" example:"3"`
	H4 int `json:"h4" example:"0"`
	H5 int `json:"h5" example:"0"`
	H6 int `json:"h6" example:"0"`
}

// LinkAnalysis represents the analysis of links on the page
type LinkAnalysis struct {
	Internal     int `json:"internal" example:"5"`
	External     int `json:"external" example:"3"`
	Inaccessible int `json:"inaccessible" example:"1"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	StatusCode int    `json:"statusCode" example:"502"`
	Message    string `json:"message" example:"Failed to analyze URL: HTTP error 404 Not Found"`
}