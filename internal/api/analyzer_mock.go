package api

import (
    "github.com/maheshjq/web-analyzer_v1/internal/models"
)

// MockAnalyzer is a test implementation of the Analyzer interface
type MockAnalyzer struct {
    AnalyzeFn func(url string) (*models.AnalysisResponse, error)
}

// Analyze calls the mock implementation function
func (m *MockAnalyzer) Analyze(url string) (*models.AnalysisResponse, error) {
    return m.AnalyzeFn(url)
}