package api

// @title Web Page Analyzer API
// @version 1.0
// @description API for analyzing web pages, extracting HTML version, title, headings, links, and detecting login forms.
// @contact.name Web Analyzer Team

// @host localhost:8080
// @BasePath /api

// @scheme http

// AnalyzeHandler godoc
// @Summary Analyze a web page
// @Description Fetches and analyzes a web page by URL, returning information about its structure and content
// @ID analyze-web-page
// @Accept json
// @Produce json
// @Param request body models.AnalysisRequest true "URL to analyze"
// @Success 200 {object} models.AnalysisResponse "Successful analysis"
// @Failure 400 {object} models.ErrorResponse "Bad request (invalid URL format)"
// @Failure 502 {object} models.ErrorResponse "Failed to fetch or analyze the URL"
// @Router /analyze [post]
func (h *Handler) AnalyzeHandlerSwagger() {
	// This function exists only for Swagger documentation
	// The actual implementation is in AnalyzeHandler
}

// HealthCheckHandler godoc
// @Summary Health check
// @Description Returns the health status of the API
// @ID health-check
// @Produce json
// @Success 200 {object} map[string]string "Service is healthy"
// @Router /health [get]
func (h *Handler) HealthCheckHandlerSwagger() {
	// This function exists only for Swagger documentation
	// The actual implementation is in HealthCheckHandler
}