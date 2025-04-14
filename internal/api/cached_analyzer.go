// internal/api/cached_analyzer.go
package api

import (
	"log"
	"sync"
	"time"

	"github.com/maheshjq/web-analyzer_v1/internal/models"
)

// CachedResult holds a cached analysis result and its expiration time
type CachedResult struct {
	Result    *models.AnalysisResponse
	ExpiresAt time.Time
}

// NewCachedAnalyzer creates a new cached analyzer
func NewCachedAnalyzer(delegate Analyzer, ttl time.Duration) *CachedAnalyzer {
	ca := &CachedAnalyzer{
		delegate: delegate,
		cache:    make(map[string]CachedResult),
		ttl:      ttl,
	}

	// Start a goroutine to periodically clean the cache
	go ca.cleanupCache(ttl)

	return ca
}

// Add this method
func (ca *CachedAnalyzer) cleanupCache(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		now := time.Now()
		ca.mu.Lock()
		for url, cached := range ca.cache {
			if now.After(cached.ExpiresAt) {
				delete(ca.cache, url)
			}
		}
		ca.mu.Unlock()
	}
}

// Analyze implements the Analyzer interface with caching
func (ca *CachedAnalyzer) Analyze(url string) (*models.AnalysisResponse, error) {
	start := time.Now()

	// Check cache first
	ca.mu.RLock()
	cached, found := ca.cache[url]
	ca.mu.RUnlock()

	now := time.Now()

	// If found and not expired, return the cached result
	if found && now.Before(cached.ExpiresAt) {
		ca.mu.Lock()
		ca.cacheHits++
		ca.mu.Unlock()
		log.Printf("CACHE HIT for %s - returning cached result (%.2f ms)",
			url, float64(time.Since(start).Microseconds())/1000)
		return cached.Result, nil
	}

	ca.mu.Lock()
	ca.cacheMisses++
	ca.mu.Unlock()
	log.Printf("CACHE MISS for %s - performing analysis", url)

	// Not in cache or expired, call the delegate
	analysisStart := time.Now()
	result, err := ca.delegate.Analyze(url)
	analysisDuration := time.Since(analysisStart)

	if err != nil {
		log.Printf("ERROR analyzing %s: %v", url, err)
		return nil, err
	}

	log.Printf("Analysis completed for %s (took %.2f sec)",
		url, analysisDuration.Seconds())

	// Store in cache
	ca.mu.Lock()
	ca.cache[url] = CachedResult{
		Result:    result,
		ExpiresAt: now.Add(ca.ttl),
	}
	ca.mu.Unlock()
	log.Printf("Cached result for %s (expires in %v)", url, ca.ttl)

	return result, nil
}

type CachedAnalyzer struct {
	delegate    Analyzer
	cache       map[string]CachedResult
	ttl         time.Duration
	mu          sync.RWMutex
	cacheHits   int
	cacheMisses int
}

// Add methods to get metrics
func (ca *CachedAnalyzer) CacheHits() int {
	ca.mu.RLock()
	defer ca.mu.RUnlock()
	return ca.cacheHits
}

func (ca *CachedAnalyzer) CacheMisses() int {
	ca.mu.RLock()
	defer ca.mu.RUnlock()
	return ca.cacheMisses
}
