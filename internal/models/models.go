package models

// AnalysisRequest represents the request to analyze a webpage
type AnalysisRequest struct {
	URL string `json:"url"`
}

// AnalysisResponse represents the result of the analysis
type AnalysisResponse struct {
	HTMLVersion string `json:"htmlVersion"`
}
