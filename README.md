# Web Page Analyzer

A web application for analyzing web pages, providing information about HTML version, page title, headings, links, and login forms.

## Project Overview

This application consists of a Go backend and a React frontend that allows users to analyze web pages by entering a URL. The analysis includes:

- HTML version detection
- Page title extraction
- Heading count by level (h1-h6)
- Internal and external link count
- Detection of inaccessible links
- Login form detection

## Technologies Used

### Backend
- Go 1.21
- Gorilla Mux (HTTP router)
- golang.org/x/net/html (HTML parsing)
- log/slog (Structured logging)

### Frontend
- React 18
- Tailwind CSS
- Axios (HTTP requests)

### DevOps
- Docker for containerization
- Make for automation

## Prerequisites

- Go 1.21 or later
- Node.js 18 or later (for frontend development)
- Docker (optional, for containerization)

## Setup Instructions

### Manual Setup

1. **Clone the repository**
   ```bash
   git clone https://github.com/yourusername/web-analyzer.git
   cd web-analyzer
   ```

2. **Build and run the backend**
   ```bash
   go mod tidy
   go run cmd/server/main.go
   ```

3. **Build the frontend**
   ```bash
   cd web
   npm install
   npm run build
   cd ..
   ```

4. **Access the application**
   Open your browser and navigate to `http://localhost:8080`

### Docker Setup

1. **Build the Docker image**
   ```bash
   docker build -t web-analyzer .
   ```

2. **Run the container**
   ```bash
   docker run -p 8080:8080 web-analyzer
   ```

3. **Access the application**
   Open your browser and navigate to `http://localhost:8080`

## Usage

1. Enter a valid URL in the input field
2. Click "Analyze" to process the URL
3. View the analysis results displayed on the page

## API Documentation

### POST /api/analyze
Analyzes a web page.

**Request Body:**
```json
{
  "url": "https://example.com"
}
```

**Success Response (200 OK):**
```json
{
  "htmlVersion": "HTML5",
  "title": "Example Page",
  "headings": {
    "h1": 1,
    "h2": 3,
    "h3": 5,
    "h4": 0,
    "h5": 0,
    "h6": 0
  },
  "links": {
    "internal": 10,
    "external": 5,
    "inaccessible": 1
  },
  "containsLoginForm": true
}
```

**Error Response (4xx/5xx):**
```json
{
  "statusCode": 502,
  "message": "Failed to analyze URL: HTTP error: 404 Not Found"
}
```

### GET /api/health
Health check endpoint.

**Success Response (200 OK):**
```json
{
  "status": "ok"
}
```

## Testing

To run the unit tests:
```bash
go test ./...
```

## Design Decisions and Assumptions

1. **HTML Version Detection**: We detect HTML versions based on the DOCTYPE declaration. If no DOCTYPE is found, we attempt to guess the version based on the presence of HTML5-specific elements.

2. **Link Classification**:
   - Internal links: Links that point to the same domain or relative links
   - External links: Links that point to different domains
   - Inaccessible links: Links that return HTTP status codes other than 2xx or 3xx

3. **Login Form Detection**: We detect login forms by looking for:
   - Forms with "login", "signin", etc. in their ID, class, or action attributes
   - Forms containing password input fields

4. **Concurrency**: We use goroutines and waitgroups to check link accessibility in parallel, improving performance when analyzing pages with many links.

5. **Error Handling**: Detailed error messages are provided when a URL cannot be accessed, including the HTTP status code and a description.

## Challenges Faced

1. **HTML Parsing**: Determining the HTML version accurately was challenging, especially for pages without proper DOCTYPE declarations.

2. **Link Accessibility**: Checking all links for accessibility could be resource-intensive and time-consuming. We implemented concurrency to mitigate this.

3. **Login Form Detection**: Login forms can vary widely, so we used a combination of heuristics to detect them.

## Possible Improvements

1. **Caching**: Implement caching to avoid re-analyzing the same URLs repeatedly.

2. **Pagination**: Add pagination for large analysis results, especially for pages with many links.

3. **Authentication**: Add user authentication to allow users to save and review their analysis history.

4. **Advanced Analysis Options**:
   - SEO analysis
   - Performance metrics
   - Accessibility evaluation
   - Mobile-friendliness check

5. **Scheduled Analysis**: Allow users to set up scheduled analysis for monitoring websites over time.

6. **API Rate Limiting**: Implement rate limiting to prevent abuse of the API.

7. **Export Functionality**: Allow users to export analysis results in different formats (CSV, PDF, etc.).

8. **Detailed Link Analysis**: Provide more detailed information about links, such as anchor text, target, etc.

9. **Security Enhancements**: Implement additional security measures such as input validation, CSRF protection, etc.

10. **Metrics and Monitoring**: Add Prometheus metrics for monitoring application performance.