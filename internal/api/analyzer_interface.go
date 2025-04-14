package api

import (
	"time"

	"github.com/maheshjq/web-analyzer_v1/internal/analyzer"
	"github.com/maheshjq/web-analyzer_v1/internal/models"
)

// Analyzer interface defines the behavior for a web page analyzer
type Analyzer interface {
	Analyze(url string) (*models.AnalysisResponse, error)
}

// analyzer factory function type
type AnalyzerFactory func() Analyzer

// Default analyzer factory variable that can be replaced in tests
var NewAnalyzerFunc AnalyzerFactory

// Initialize our factory.
func init() {
	// Create a cached analyzer that wraps the actual analyzer
	NewAnalyzerFunc = func() Analyzer {
		realAnalyzer := &DefaultAnalyzer{}
		// Cache results for 15 minutes
		return NewCachedAnalyzer(realAnalyzer, 15*time.Minute)
	}
}

// DefaultAnalyzer is a wrapper around the actual analyzer implementation
type DefaultAnalyzer struct{}

// Analyze implements the Analyzer interface by calling the actual analyzer
func (da *DefaultAnalyzer) Analyze(url string) (*models.AnalysisResponse, error) {
	// Create an instance of your actual analyzer
	realAnalyzer := analyzer.NewAnalyzer()

	// Call the actual analyze method
	return realAnalyzer.Analyze(url)
}
