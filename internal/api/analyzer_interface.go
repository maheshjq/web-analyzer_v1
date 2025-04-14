package api

import (
	"sync"
	"time"

	"github.com/maheshjq/web-analyzer_v1/internal/analyzer"
	"github.com/maheshjq/web-analyzer_v1/internal/models"
)

// toggle flag -> enable/disable caching
var EnableCaching = true

// Analyzer interface defines the behavior for a web page analyzer
type Analyzer interface {
	Analyze(url string) (*models.AnalysisResponse, error)
}

// Global singleton
var singletonAnalyzer Analyzer
var once sync.Once

func GetAnalyzer() Analyzer {
	once.Do(func() {
		realAnalyzer := &DefaultAnalyzer{}

		if EnableCaching {
			// Cache results if caching is enabled
			singletonAnalyzer = NewCachedAnalyzer(realAnalyzer, 15*time.Minute)
		} else {
			// Use analyzer directly if caching is disabled
			singletonAnalyzer = realAnalyzer
		}
	})
	return singletonAnalyzer
}

// DefaultAnalyzer is a wrapper for the actual analyzer imple
type DefaultAnalyzer struct{}

// Analyze implements the Analyzer interface by calling the actual analyzer
func (da *DefaultAnalyzer) Analyze(url string) (*models.AnalysisResponse, error) {
	// Create an instance of actual analyzer
	realAnalyzer := analyzer.NewAnalyzer()

	// Call the actual analyze method
	return realAnalyzer.Analyze(url)
}
