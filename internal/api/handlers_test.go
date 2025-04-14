package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/maheshjq/web-analyzer_v1/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockAnalyzerFactory creates a mock analyzer for testing
func setupMockAnalyzer(fn func(string) (*models.AnalysisResponse, error)) {
	// Temporarily set EnableCaching to false for tests
	originalEnableCaching := EnableCaching
	EnableCaching = false

	// Override singletonAnalyzer for testing
	singletonAnalyzer = &MockAnalyzer{
		AnalyzeFn: fn,
	}

	// Restore original value after test
	defer func() {
		EnableCaching = originalEnableCaching
		singletonAnalyzer = nil // Reset singleton
	}()
}

func TestAnalyzeHandler(t *testing.T) {
	// Clear singleton analyzer
	singletonAnalyzer = nil

	// Set up mock analyzer
	singletonAnalyzer = &MockAnalyzer{
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

	// Ensure analyzer is cleared after test
	defer func() {
		singletonAnalyzer = nil
	}()

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
	assert.Equal(t, "HTML5", response.HTMLVersion, "unexpected HTML version")
	assert.Equal(t, "Test Page", response.Title, "unexpected title")
	assert.Equal(t, 1, response.Headings.H1, "unexpected H1 count")
	assert.Equal(t, 2, response.Headings.H2, "unexpected H2 count")
	assert.Equal(t, 5, response.Links.Internal, "unexpected internal links count")
	assert.Equal(t, 3, response.Links.External, "unexpected external links count")
	assert.True(t, response.ContainsLoginForm, "expected login form to be detected")
}

func TestAnalyzeHandler_InvalidRequest(t *testing.T) {
	// Test cases with invalid inputs
	testCases := []struct {
		name       string
		reqBody    string
		wantStatus int
		wantError  string
	}{
		{
			name:       "Empty request",
			reqBody:    `{}`,
			wantStatus: http.StatusBadRequest,
			wantError:  "URL is required",
		},
		{
			name:       "Invalid JSON",
			reqBody:    `{invalid json}`,
			wantStatus: http.StatusBadRequest,
			wantError:  "Invalid request body",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest("POST", "/api/analyze", strings.NewReader(tc.reqBody))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(AnalyzeHandler)
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tc.wantStatus, rr.Code, "unexpected status code")

			var errorResp models.ErrorResponse
			err = json.Unmarshal(rr.Body.Bytes(), &errorResp)
			require.NoError(t, err, "Failed to parse error response")

			assert.Contains(t, errorResp.Message, tc.wantError, "error message doesn't contain expected text")
		})
	}
}

func TestAnalyzeHandler_AnalyzerError(t *testing.T) {
	// Clear singleton analyzer
	singletonAnalyzer = nil

	// Set up mock analyzer that returns an error
	singletonAnalyzer = &MockAnalyzer{
		AnalyzeFn: func(url string) (*models.AnalysisResponse, error) {
			return nil, errors.New("analyzer error")
		},
	}

	// Ensure analyzer is cleared after test
	defer func() {
		singletonAnalyzer = nil
	}()

	reqBody := `{"url": "https://example.com"}`
	req, err := http.NewRequest("POST", "/api/analyze", strings.NewReader(reqBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(AnalyzeHandler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadGateway, rr.Code, "expected Bad Gateway status")

	var errorResp models.ErrorResponse
	err = json.Unmarshal(rr.Body.Bytes(), &errorResp)
	require.NoError(t, err, "Failed to parse error response")

	assert.Contains(t, errorResp.Message, "Failed to analyze URL", "error message doesn't match expected text")
}

func TestHealthCheckHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/health", nil)
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HealthCheckHandler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "health check should return 200 OK")
	assert.Contains(t, rr.Body.String(), "status", "response should include status field")
	assert.Contains(t, rr.Body.String(), "ok", "status should be ok")
}
