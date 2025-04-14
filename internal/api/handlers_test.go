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

// Add at the top of your handlers_test.go file:
var mockAnalyzeFunc func(url string) (*models.AnalysisResponse, error)

// Then create a mock analyzer type that uses this function
type testMockAnalyzer struct{}

func (m *testMockAnalyzer) Analyze(url string) (*models.AnalysisResponse, error) {
	return mockAnalyzeFunc(url)
}

func TestAnalyzeHandler(t *testing.T) {
	// Save original and restore after test
	originalAnalyzer := singletonAnalyzer
	defer func() { singletonAnalyzer = originalAnalyzer }()

	// Create mock response
	mockResponse := &models.AnalysisResponse{
		HTMLVersion:       "HTML5",
		Title:             "Example Domain",
		Headings:          models.HeadingCount{H1: 1, H2: 0},
		Links:             models.LinkAnalysis{Internal: 0, External: 1, Inaccessible: 0},
		ContainsLoginForm: false,
	}

	// Set up mock function
	mockAnalyzeFunc = func(url string) (*models.AnalysisResponse, error) {
		return mockResponse, nil
	}

	// Replace the analyzer with our test mock
	singletonAnalyzer = &testMockAnalyzer{}

	// Create test request
	reqBody := `{"url": "https://example.com"}`
	req, err := http.NewRequest("POST", "/api/analyze", strings.NewReader(reqBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	rr := httptest.NewRecorder()

	// Call the handler
	handler := http.HandlerFunc(AnalyzeHandler)
	handler.ServeHTTP(rr, req)

	// Check response
	require.Equal(t, http.StatusOK, rr.Code)

	// Parse response
	var response models.AnalysisResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err, "Failed to parse response")

	// Verify response data
	assert.Equal(t, mockResponse.HTMLVersion, response.HTMLVersion, "unexpected HTML version")
	assert.Equal(t, mockResponse.Title, response.Title, "unexpected title")
	assert.Equal(t, mockResponse.Headings.H1, response.Headings.H1, "unexpected H1 count")
	assert.Equal(t, mockResponse.Headings.H2, response.Headings.H2, "unexpected H2 count")
	assert.Equal(t, mockResponse.Links.Internal, response.Links.Internal, "unexpected internal links count")
	assert.Equal(t, mockResponse.Links.External, response.Links.External, "unexpected external links count")
	assert.Equal(t, mockResponse.ContainsLoginForm, response.ContainsLoginForm, "unexpected login form detection")
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