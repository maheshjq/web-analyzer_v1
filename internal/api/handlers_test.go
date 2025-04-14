package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/maheshjq/web-analyzer_v1/internal/models"
)

func TestAnalyzeHandler(t *testing.T) {
	// Save original factory and restore it after test
	originalFactory := NewAnalyzerFunc
	defer func() { NewAnalyzerFunc = originalFactory }()

	// Create mock analyzer with predictable behavior
	NewAnalyzerFunc = func() Analyzer {
		return &MockAnalyzer{
			AnalyzeFn: func(url string) (*models.AnalysisResponse, error) {
				// Return predetermined test data
				return &models.AnalysisResponse{
					HTMLVersion:       "HTML5",
					Title:             "Test Page",
					Headings:          models.HeadingCount{H1: 1, H2: 2},
					Links:             models.LinkAnalysis{Internal: 5, External: 3},
					ContainsLoginForm: true,
				}, nil
			},
		}
	}

	// Create test request
	reqBody := `{"url": "https://example.com"}`
	req, err := http.NewRequest("POST", "/api/analyze", strings.NewReader(reqBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	rr := httptest.NewRecorder()

	// Call the handler
	handler := http.HandlerFunc(AnalyzeHandler)
	handler.ServeHTTP(rr, req)

	// Check response
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Parse response
	var response models.AnalysisResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatal("Failed to parse response:", err)
	}

	// Verify response data
	if response.HTMLVersion != "HTML5" {
		t.Errorf("unexpected HTML version: got %v, want %v", response.HTMLVersion, "HTML5")
	}
	if response.Title != "Test Page" {
		t.Errorf("unexpected title: got %v, want %v", response.Title, "Test Page")
	}
	// Add more assertions as needed
}
