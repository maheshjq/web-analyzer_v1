package analyzer

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/maheshjq/web-analyzer_v1/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/html"
)

// TestNewAnalyzer ensures the analyzer is created with the correct defaults
func TestNewAnalyzer(t *testing.T) {
	analyzer := NewAnalyzer()

	// Check that we have a non-nil analyzer with correct timeout
	require.NotNil(t, analyzer)
	require.NotNil(t, analyzer.client)
	assert.Equal(t, 10*time.Second, analyzer.client.Timeout)
}

// TestDetectHTMLVersion tests HTML version detection for different doctypes
func TestDetectHTMLVersion(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		expected string
	}{
		{
			name:     "HTML5 Standard Doctype",
			html:     `<!DOCTYPE html><html><head><title>Test</title></head><body></body></html>`,
			expected: "HTML5",
		},
		{
			name:     "HTML 4.01 Strict",
			html:     `<!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 4.01//EN" "http://www.w3.org/TR/html4/strict.dtd"><html><head><title>Test</title></head><body></body></html>`,
			expected: "HTML 4.01",
		},
		{
			name:     "XHTML 1.0 Strict",
			html:     `<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Strict//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-strict.dtd"><html xmlns="http://www.w3.org/1999/xhtml"><head><title>Test</title></head><body></body></html>`,
			expected: "XHTML 1.0",
		},
		{
			name:     "XHTML 1.1",
			html:     `<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.1//EN" "http://www.w3.org/TR/xhtml11/DTD/xhtml11.dtd"><html xmlns="http://www.w3.org/1999/xhtml"><head><title>Test</title></head><body></body></html>`,
			expected: "XHTML 1.1",
		},
		{
			name:     "No DOCTYPE but HTML5 elements",
			html:     `<html><head><title>Test</title></head><body><article>Content</article><section>More content</section></body></html>`,
			expected: "HTML5 (No DOCTYPE)",
		},
		{
			name:     "No DOCTYPE and no HTML5 elements",
			html:     `<html><head><title>Test</title></head><body><div>Content</div></body></html>`,
			expected: "Unknown (No DOCTYPE)",
		},
		{
			name:     "Empty document",
			html:     ``,
			expected: "Unknown (No DOCTYPE)",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			doc, err := html.Parse(strings.NewReader(tc.html))
			require.NoError(t, err)

			version := detectHTMLVersion(doc)
			assert.Equal(t, tc.expected, version)
		})
	}
}

// TestExtractTitle tests title extraction from HTML documents
func TestExtractTitle(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		expected string
	}{
		{
			name:     "Simple title",
			html:     `<html><head><title>Test Title</title></head><body></body></html>`,
			expected: "Test Title",
		},
		{
			name:     "Title with special characters",
			html:     `<html><head><title>Special &amp; Characters - Test</title></head><body></body></html>`,
			expected: "Special & Characters - Test",
		},
		{
			name:     "Empty title",
			html:     `<html><head><title></title></head><body></body></html>`,
			expected: "",
		},
		{
			name:     "No title tag",
			html:     `<html><head></head><body></body></html>`,
			expected: "",
		},
		{
			name:     "Multiple title tags (should use first)",
			html:     `<html><head><title>First Title</title><title>Second Title</title></head><body></body></html>`,
			expected: "First Title",
		},
		{
			name:     "Nested content in title",
			html:     `<html><head><title>Title <span>with span</span></title></head><body></body></html>`,
			expected: "Title ",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			doc, err := html.Parse(strings.NewReader(tc.html))
			require.NoError(t, err)

			title := extractTitle(doc)
			assert.Equal(t, tc.expected, title)
		})
	}
}

// TestCountHeadings tests heading count functionality
func TestCountHeadings(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		expected models.HeadingCount
	}{
		{
			name: "All heading levels",
			html: `
				<html><body>
					<h1>Heading 1</h1>
					<h2>Heading 2a</h2>
					<h2>Heading 2b</h2>
					<h3>Heading 3</h3>
					<h4>Heading 4</h4>
					<h5>Heading 5</h5>
					<h6>Heading 6</h6>
				</body></html>
			`,
			expected: models.HeadingCount{H1: 1, H2: 2, H3: 1, H4: 1, H5: 1, H6: 1},
		},
		{
			name:     "No headings",
			html:     `<html><body><p>No headings here</p></body></html>`,
			expected: models.HeadingCount{H1: 0, H2: 0, H3: 0, H4: 0, H5: 0, H6: 0},
		},
		{
			name: "Nested headings",
			html: `
				<html><body>
					<div>
						<h1>Nested <span>Heading</span> 1</h1>
					</div>
				</body></html>
			`,
			expected: models.HeadingCount{H1: 1, H2: 0, H3: 0, H4: 0, H5: 0, H6: 0},
		},
		{
			name: "Multiple h1s",
			html: `
				<html><body>
					<h1>First H1</h1>
					<section>
						<h1>Second H1</h1>
					</section>
				</body></html>
			`,
			expected: models.HeadingCount{H1: 2, H2: 0, H3: 0, H4: 0, H5: 0, H6: 0},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			doc, err := html.Parse(strings.NewReader(tc.html))
			require.NoError(t, err)

			headings := models.HeadingCount{}
			countHeadings(doc, &headings)

			assert.Equal(t, tc.expected.H1, headings.H1)
			assert.Equal(t, tc.expected.H2, headings.H2)
			assert.Equal(t, tc.expected.H3, headings.H3)
			assert.Equal(t, tc.expected.H4, headings.H4)
			assert.Equal(t, tc.expected.H5, headings.H5)
			assert.Equal(t, tc.expected.H6, headings.H6)
		})
	}
}

// TestIsInternalLink tests the isInternalLink function
func TestIsInternalLink(t *testing.T) {
	tests := []struct {
		name     string
		href     string
		host     string
		expected bool
	}{
		{
			name:     "Absolute URL same host",
			href:     "https://example.com/page",
			host:     "example.com",
			expected: true,
		},
		{
			name:     "Absolute URL different host",
			href:     "https://other.com/page",
			host:     "example.com",
			expected: false,
		},
		{
			name:     "Absolute URL with www subdomain",
			href:     "https://www.example.com/page",
			host:     "example.com",
			expected: false, // Different by exact match
		},
		{
			name:     "Relative URL root path",
			href:     "/page",
			host:     "example.com",
			expected: true,
		},
		{
			name:     "Relative URL with dot",
			href:     "./page",
			host:     "example.com",
			expected: true,
		},
		{
			name:     "Relative URL with parent",
			href:     "../page",
			host:     "example.com",
			expected: true,
		},
		{
			name:     "Fragment URL",
			href:     "#section",
			host:     "example.com",
			expected: true,
		},
		{
			name:     "Empty URL",
			href:     "",
			host:     "example.com",
			expected: true,
		},
		{
			name:     "JavaScript URL",
			href:     "javascript:void(0)",
			host:     "example.com",
			expected: false, // Consider js links as external
		},
		{
			name:     "Mailto URL",
			href:     "mailto:test@example.com",
			host:     "example.com",
			expected: false, // Mailto links are external
		},
		{
			name:     "Invalid URL",
			href:     "::::invalid",
			host:     "example.com",
			expected: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := isInternalLink(tc.href, tc.host)
			assert.Equal(t, tc.expected, result)
		})
	}
}

// TestIsAccessibleLink tests the isAccessibleLink function
func TestIsAccessibleLink(t *testing.T) {
	// Set up a custom client with a mock transport
	mockTransport := &mockRoundTripper{
		responses: map[string]*http.Response{
			"https://example.com/ok": {
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString("")),
			},
			"https://example.com/redirect": {
				StatusCode: http.StatusFound,
				Body:       io.NopCloser(bytes.NewBufferString("")),
			},
			"https://example.com/forbidden": {
				StatusCode: http.StatusForbidden,
				Body:       io.NopCloser(bytes.NewBufferString("")),
			},
			"https://example.com/notfound": {
				StatusCode: http.StatusNotFound,
				Body:       io.NopCloser(bytes.NewBufferString("")),
			},
		},
		errors: map[string]error{
			"https://example.com/timeout": &url.Error{
				Err: &timeoutError{},
			},
		},
	}

	client := &http.Client{
		Transport: mockTransport,
	}

	tests := []struct {
		name     string
		link     string
		expected bool
	}{
		{
			name:     "Status 200 OK",
			link:     "https://example.com/ok",
			expected: true,
		},
		{
			name:     "Status 302 Found (redirect)",
			link:     "https://example.com/redirect",
			expected: true,
		},
		{
			name:     "Status 403 Forbidden",
			link:     "https://example.com/forbidden",
			expected: false, // 4xx is not accessible
		},
		{
			name:     "Status 404 Not Found",
			link:     "https://example.com/notfound",
			expected: false, // 4xx is not accessible
		},
		{
			name:     "Timeout error",
			link:     "https://example.com/timeout",
			expected: false,
		},
		{
			name:     "Non-existent in test mapping",
			link:     "https://example.com/nonexistent",
			expected: false,
		},
		{
			name:     "Fragment URL",
			link:     "#section",
			expected: true, // Fragment links are always accessible
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := isAccessibleLink(tc.link, client)
			assert.Equal(t, tc.expected, result)
		})
	}
}

// TestDetectLoginForm tests the detectLoginForm function
func TestDetectLoginForm(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		expected bool
	}{
		{
			name: "Form with password input",
			html: `
				<html><body>
					<form>
						<input type="text" name="username">
						<input type="password" name="password">
						<button type="submit">Submit</button>
					</form>
				</body></html>
			`,
			expected: true,
		},
		{
			name: "Form with login in action attribute",
			html: `
				<html><body>
					<form action="/login">
						<input type="text" name="username">
						<input type="text" name="pass">
						<button type="submit">Submit</button>
					</form>
				</body></html>
			`,
			expected: true,
		},
		{
			name: "Form with signin in id attribute",
			html: `
				<html><body>
					<form id="signin-form">
						<input type="text" name="username">
						<input type="text" name="pass">
						<button type="submit">Submit</button>
					</form>
				</body></html>
			`,
			expected: true,
		},
		{
			name: "Form with login in class attribute",
			html: `
				<html><body>
					<form class="login-form">
						<input type="text" name="username">
						<input type="text" name="pass">
						<button type="submit">Submit</button>
					</form>
				</body></html>
			`,
			expected: true,
		},
		{
			name: "Regular form, no login indicators",
			html: `
				<html><body>
					<form action="/submit">
						<input type="text" name="name">
						<input type="text" name="email">
						<button type="submit">Submit</button>
					</form>
				</body></html>
			`,
			expected: false,
		},
		{
			name: "Contact form with email",
			html: `
				<html><body>
					<form action="/contact">
						<input type="text" name="name">
						<input type="email" name="email">
						<textarea name="message"></textarea>
						<button type="submit">Send</button>
					</form>
				</body></html>
			`,
			expected: false,
		},
		{
			name: "No form",
			html: `
				<html><body>
					<div>No form here</div>
				</body></html>
			`,
			expected: false,
		},
		{
			name:     "Empty document",
			html:     ``,
			expected: false,
		},
		{
			name: "Nested form with password",
			html: `
				<html><body>
					<div>
						<form>
							<div>
								<input type="password" name="pass">
							</div>
						</form>
					</div>
				</body></html>
			`,
			expected: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			doc, err := html.Parse(strings.NewReader(tc.html))
			require.NoError(t, err)

			result := detectLoginForm(doc)
			assert.Equal(t, tc.expected, result)
		})
	}
}

// TestAnalyzeLinks tests the analyzeLinks function
func TestAnalyzeLinks(t *testing.T) {
	// Simple test for link extraction and categorization
	htmlStr := `
		<html><body>
			<a href="/">Home</a>
			<a href="/about">About</a>
			<a href="https://example.com/contact">Contact</a>
			<a href="https://external.com">External</a>
			<a href="mailto:test@example.com">Email</a>
			<a href="#section">Section</a>
			<a href="javascript:void(0)">JS Link</a>
		</body></html>
	`

	doc, err := html.Parse(strings.NewReader(htmlStr))
	require.NoError(t, err)

	// Create a test client with mocked responses
	client := &http.Client{
		Transport: &mockRoundTripper{
			responses: map[string]*http.Response{
				"https://example.com/contact": {
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString("")),
				},
				"https://external.com": {
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString("")),
				},
			},
		},
	}

	// Analyze links
	result := analyzeLinks(doc, "example.com", client)

	// Check the results
	assert.GreaterOrEqual(t, result.Internal, 3) // Home, About, Section should be internal
	assert.GreaterOrEqual(t, result.External, 2) // External and Email should be external
	assert.Equal(t, 0, result.Inaccessible)      // All links are accessible in our mock
}

// Test helpers
// ------------

// MockRoundTripper for testing HTTP client behavior
type mockRoundTripper struct {
	responses map[string]*http.Response
	errors    map[string]error
}

func (m *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	// Check if we have an error for this URL
	if err, ok := m.errors[req.URL.String()]; ok {
		return nil, err
	}

	// Check if we have a response for this URL
	if resp, ok := m.responses[req.URL.String()]; ok {
		return resp, nil
	}

	// Default to 404 Not Found
	return &http.Response{
		StatusCode: http.StatusNotFound,
		Body:       io.NopCloser(bytes.NewBufferString("Not Found")),
	}, nil
}

// timeoutError implements the net.Error interface for timeout testing
type timeoutError struct{}

func (e *timeoutError) Error() string   { return "timeout error" }
func (e *timeoutError) Timeout() bool   { return true }
func (e *timeoutError) Temporary() bool { return true }

func TestFindElement(t *testing.T) {
	htmlStr := `
		<html>
			<head><title>Test</title></head>
			<body>
				<header>Header</header>
				<nav>Navigation</nav>
				<article>
					<section>Section content</section>
				</article>
				<footer>Footer</footer>
			</body>
		</html>
	`

	doc, err := html.Parse(strings.NewReader(htmlStr))
	require.NoError(t, err)

	// Test for elements that exist
	assert.True(t, findElement(doc, "header"))
	assert.True(t, findElement(doc, "nav"))
	assert.True(t, findElement(doc, "article"))
	assert.True(t, findElement(doc, "section"))
	assert.True(t, findElement(doc, "footer"))

	// Test for elements that don’t exist
	assert.False(t, findElement(doc, "aside"))
	assert.False(t, findElement(doc, "canvas"))
	assert.False(t, findElement(doc, "video"))
}

// TestAnalyze tests the entire Analyze function with a mock client
func TestAnalyze(t *testing.T) {
	// Create a test HTML document
	testHTML := `
		<!DOCTYPE html>
		<html>
		<head>
			<title>Test Page</title>
		</head>
		<body>
			<h1>Main Heading</h1>
			<h2>Subheading 1</h2>
			<p>Some content</p>
			<h2>Subheading 2</h2>
			<p>More content</p>
			<a href="/">Home link</a>
			<a href="/about">About link</a>
			<a href="https://external.com">External link</a>
			
			<form id="login-form">
				<input type="text" name="username">
				<input type="password" name="password">
				<button type="submit">Login</button>
			</form>
		</body>
		</html>
	`

	// Create a mock HTTP client
	mockTransport := &mockRoundTripper{
		responses: map[string]*http.Response{
			"https://test.example.com": {
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(testHTML)),
				Header:     http.Header{"Content-Type": []string{"text/html"}},
			},
			"https://error.example.com": {
				StatusCode: http.StatusNotFound,
				Body:       io.NopCloser(strings.NewReader("Not Found")),
			},
		},
		errors: map[string]error{
			"https://example.com/timeout": &url.Error{
				Err: &timeoutError{},
			},
		},
	}

	// Create an analyzer with the mock client
	analyzer := &Analyzer{
		client: &http.Client{
			Transport: mockTransport,
		},
	}

	// Test successful analysis
	t.Run("Successful Analysis", func(t *testing.T) {
		result, err := analyzer.Analyze("https://test.example.com")

		require.NoError(t, err)
		assert.Equal(t, "HTML5", result.HTMLVersion)
		assert.Equal(t, "Test Page", result.Title)
		assert.Equal(t, 1, result.Headings.H1)
		assert.Equal(t, 2, result.Headings.H2)
		assert.True(t, result.ContainsLoginForm)
		assert.GreaterOrEqual(t, result.Links.Internal, 2) // At least 2 internal links
		assert.GreaterOrEqual(t, result.Links.External, 1) // At least 1 external link
	})

	// Test HTTP error
	t.Run("HTTP Error", func(t *testing.T) {
		_, err := analyzer.Analyze("https://error.example.com")

		require.Error(t, err)
		assert.Contains(t, err.Error(), "HTTP error: 404")
	})

	// Test network error
	t.Run("Network Error", func(t *testing.T) {
		_, err := analyzer.Analyze("https://timeout.example.com")

		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to fetch URL")
	})
}