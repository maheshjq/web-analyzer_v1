package analyzer

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/html"
)

func TestDetectHTMLVersion(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		expected string
	}{
		{
			name:     "HTML5",
			html:     `<!DOCTYPE html><html><head><title>Test</title></head><body></body></html>`,
			expected: "HTML5",
		},
		{
			name:     "HTML 4.01",
			html:     `<!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 4.01//EN" "http://www.w3.org/TR/html4/strict.dtd"><html><head><title>Test</title></head><body></body></html>`,
			expected: "HTML 4.01",
		},
		{
			name:     "XHTML 1.0",
			html:     `<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Strict//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-strict.dtd"><html xmlns="http://www.w3.org/1999/xhtml"><head><title>Test</title></head><body></body></html>`,
			expected: "XHTML 1.0",
		},
		{
			name:     "No DOCTYPE but HTML5 elements",
			html:     `<html><head><title>Test</title></head><body><nav>Menu</nav><section>Content</section></body></html>`,
			expected: "HTML5 (No DOCTYPE)",
		},
		{
			name:     "Unknown",
			html:     `<html><head><title>Test</title></head><body></body></html>`,
			expected: "Unknown (No DOCTYPE)",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			doc, err := html.Parse(strings.NewReader(tc.html))
			if err != nil {
				t.Fatalf("Failed to parse HTML: %v", err)
			}

			version := detectHTMLVersion(doc)
			assert.Equal(t, tc.expected, version)
		})
	}
}

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
			name:     "Empty title",
			html:     `<html><head><title></title></head><body></body></html>`,
			expected: "",
		},
		{
			name:     "No title",
			html:     `<html><head></head><body></body></html>`,
			expected: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			doc, err := html.Parse(strings.NewReader(tc.html))
			if err != nil {
				t.Fatalf("Failed to parse HTML: %v", err)
			}

			title := extractTitle(doc)
			assert.Equal(t, tc.expected, title)
		})
	}
}

func TestCountHeadings(t *testing.T) {
	html := `
		<html>
			<head><title>Test</title></head>
			<body>
				<h1>Heading 1</h1>
				<h2>Heading 2</h2>
				<h2>Another H2</h2>
				<h3>Heading 3</h3>
				<h4>Heading 4</h4>
				<h5>Heading 5</h5>
				<h6>Heading 6</h6>
			</body>
		</html>
	`

	doc, err := html.Parse(strings.NewReader(html))
	if err != nil {
		t.Fatalf("Failed to parse HTML: %v", err)
	}

	headings := struct {
		H1, H2, H3, H4, H5, H6 int
	}{}

	countHeadings(doc, &headings)

	assert.Equal(t, 1, headings.H1)
	assert.Equal(t, 2, headings.H2)
	assert.Equal(t, 1, headings.H3)
	assert.Equal(t, 1, headings.H4)
	assert.Equal(t, 1, headings.H5)
	assert.Equal(t, 1, headings.H6)
}

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
			name:     "Relative URL",
			href:     "/page",
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

func TestDetectLoginForm(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		expected bool
	}{
		{
			name: "Login form with password input",
			html: `
				<html><body>
					<form id="login-form">
						<input type="text" name="username">
						<input type="password" name="password">
						<button type="submit">Login</button>
					</form>
				</body></html>
			`,
			expected: true,
		},
		{
			name: "Login form with login in action",
			html: `
				<html><body>
					<form action="/login">
						<input type="text" name="username">
						<input type="text" name="password">
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
			name: "No form",
			html: `
				<html><body>
					<div>No form here</div>
				</body></html>
			`,
			expected: false,
		},
	]

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			doc, err := html.Parse(strings.NewReader(tc.html))
			if err != nil {
				t.Fatalf("Failed to parse HTML: %v", err)
			}

			result := detectLoginForm(doc)
			assert.Equal(t, tc.expected, result)
		})
	}
}

// Mock HTTP server for testing
type mockTransport struct {
	responses map[string]*http.Response
}

func (m *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	resp, ok := m.responses[req.URL.String()]
	if !ok {
		return &http.Response{
			StatusCode: http.StatusNotFound,
			Body:       http.NoBody,
		}, nil
	}
	return resp, nil
}

func TestIsAccessibleLink(t *testing.T) {
	// Set up mock responses
	mockResp := map[string]*http.Response{
		"https://example.com/good": {
			StatusCode: http.StatusOK,
			Body:       http.NoBody,
		},
		"https://example.com/redirect": {
			StatusCode: http.StatusFound,
			Body:       http.NoBody,
		},
		"https://example.com/bad": {
			StatusCode: http.StatusNotFound,
			Body:       http.NoBody,
		},
	}

	client := &http.Client{
		Transport: &mockTransport{responses: mockResp},
	}

	tests := []struct {
		name     string
		link     string
		expected bool
	}{
		{
			name:     "Accessible link",
			link:     "https://example.com/good",
			expected: true,
		},
		{
			name:     "Redirect link",
			link:     "https://example.com/redirect",
			expected: true,
		},
		{
			name:     "Not found link",
			link:     "https://example.com/bad",
			expected: false,
		},
		{
			name:     "Fragment link",
			link:     "#section",
			expected: true,
		},
		{
			name:     "Non-existent link",
			link:     "https://example.com/nonexistent",
			expected: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := isAccessibleLink(tc.link, client)
			assert.Equal(t, tc.expected, result)
		})
	}
}