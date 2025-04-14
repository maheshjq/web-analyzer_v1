// internal/api/cached_analyzer.go
package api

import (
	"sync"
	"time"

	"github.com/maheshjq/web-analyzer_v1/internal/models"
)

// CachedResult holds a cached analysis result and its expiration time
type CachedResult struct {
	Result    *models.AnalysisResponse
	ExpiresAt time.Time
}

// CachedAnalyzer implements Analyzer with caching
type CachedAnalyzer struct {
	delegate Analyzer
	cache    map[string]CachedResult
	ttl      time.Duration
	mu       sync.RWMutex
}

// NewCachedAnalyzer creates a new cached analyzer
func NewCachedAnalyzer(delegate Analyzer, ttl time.Duration) *CachedAnalyzer {
	return &CachedAnalyzer{
		delegate: delegate,
		cache:    make(map[string]CachedResult),
		ttl:      ttl,
	}
}

// Analyze implements the Analyzer interface with caching
func (ca *CachedAnalyzer) Analyze(url string) (*models.AnalysisResponse, error) {
	// Check cache first
	ca.mu.RLock()
	cached, found := ca.cache[url]
	ca.mu.RUnlock()

	now := time.Now()

	// If found and not expired, return the cached result
	if found && now.Before(cached.ExpiresAt) {
		return cached.Result, nil
	}

	// Not in cache or expired, call the delegate
	result, err := ca.delegate.Analyze(url)
	if err != nil {
		return nil, err
	}

	// Store in cache
	ca.mu.Lock()
	ca.cache[url] = CachedResult{
		Result:    result,
		ExpiresAt: now.Add(ca.ttl),
	}
	ca.mu.Unlock()

	return result, nil
}
