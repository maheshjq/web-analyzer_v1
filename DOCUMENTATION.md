# Web Page Analyzer - Documentation

## Overview

The Web Page Analyzer is a web application that allows users to analyze web pages by providing a URL. The application returns detailed information about the web page, including its HTML version, page title, headings, links, and login form detection.

## Architecture

The application follows a client-server architecture:

- **Backend**: Written in Go, responsible for fetching and analyzing web pages
- **Frontend**: React application that provides a user interface for requesting analysis and displaying results

### Backend

The backend is structured using a modular approach:

```
web-analyzer/
├── cmd/server/         # Application entry point
├── internal/           # Internal packages
│   ├── analyzer/       # Core analysis logic
│   ├── api/            # HTTP handlers and middleware
│   └── models/         # Data models
```

#### Key Components:

1. **WebAnalyzer**: The core component responsible for analyzing web pages
2. **API Handlers**: HTTP handlers for processing requests and returning responses
3. **Middleware**: Logging, error recovery, and CORS support

### Frontend

The React frontend is organized as follows:

```
web/
├── public/           # Static assets
└── src/              # Source code
    ├── components/   # React components
    └── services/     # API services
```

#### Key Components:

1. **AnalysisForm**: Form for submitting URLs
2. **AnalysisResult**: Component for displaying analysis results
3. **ErrorDisplay**: Component for displaying error messages
4. **API Service**: Service for communicating with the backend

## Implementation Details

### HTML Version Detection

The application detects the HTML version by examining the DOCTYPE declaration:

- `<!DOCTYPE html>` indicates HTML5
- `<!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 4.01//EN"...>` indicates HTML 4.01
- `<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0...>` indicates XHTML 1.0

If no DOCTYPE is found, the application attempts to identify HTML5 by looking for HTML5-specific elements.

### Link Analysis

Links are categorized as:

- **Internal**: Links pointing to the same domain or relative URLs
- **External**: Links pointing to different domains
- **Inaccessible**: Links that return HTTP status codes other than 2xx or 3xx

### Login Form Detection

The application uses heuristics to detect login forms:

1. Forms with "login", "signin", or similar terms in their attributes
2. Forms containing password input fields

### Concurrency

The application uses Go's concurrency features to improve performance:

- Goroutines for parallel processing of link accessibility checks
- Context for managing timeouts and cancellations
- Waitgroups for coordinating goroutines

## API Documentation

### `POST /api/analyze`

Analyzes a web page.

**Request:**

```json
{
  "url": "https://example.com"
}
```

**Response:**

```json
{
  "htmlVersion": "HTML5",
  "title": "Example Domain",
  "headings": {
    "h1": 1,
    "h2": 0,
    "h3": 0,
    "h4": 0,
    "h5": 0,
    "h6": 0
  },
  "links": {
    "internal": 0,
    "external": 1,
    "inaccessible": 0
  },
  "containsLoginForm": false
}
```

**Error Response:**

```json
{
  "statusCode": 502,
  "message": "Failed to analyze URL: HTTP error: 404 Not Found"
}
```

## Testing

The application includes unit tests for key components:

- HTML version detection
- Title extraction
- Heading counting
- Link classification
- Login form detection

## Building and Running

### Prerequisites

- Go 1.21 or later
- Node.js 18 or later
- npm 9 or later

### With Make

```bash
# Build the application
make build

# Run the application
make run

# Run tests
make test

# Clean build artifacts
make clean
```

### With Docker

```bash
# Build Docker image
docker build -t web-analyzer .

# Run Docker container
docker run -p 8080:8080 web-analyzer
```

## Design Decisions and Assumptions

1. **HTTP Client Timeouts**: The application uses a 10-second timeout for HTTP requests to prevent hanging on slow websites.

2. **Link Accessibility Checks**: To avoid overloading the server with requests, link accessibility checks are performed in parallel with a shorter timeout.

3. **Error Handling**: Detailed error messages are provided to help users understand why a URL could not be analyzed.

4. **CORS Support**: Cross-Origin Resource Sharing (CORS) is enabled to allow the frontend to communicate with the backend from different origins during development.

5. **Structured Logging**: The application uses structured logging to facilitate log analysis and monitoring.

## Potential Improvements

1. **Caching**: Implement caching to avoid re-analyzing the same URLs.

2. **Rate Limiting**: Add rate limiting to prevent abuse of the API.

3. **Authentication**: Add user authentication to allow users to save and review analysis history.

4. **More Analysis Options**: Expand the analysis to include performance metrics, SEO analysis, and accessibility evaluation.

5. **Improved Login Form Detection**: Enhance login form detection with more sophisticated heuristics or machine learning.

6. **Pagination**: Add pagination for large analysis results, especially for pages with many links.

7. **Metrics and Monitoring**: Add Prometheus metrics for monitoring application performance.

8. **Export Functionality**: Allow users to export analysis results in different formats.