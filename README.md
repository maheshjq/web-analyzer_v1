# Web Page Analyzer

A web application that analyzes web pages and provides detailed information about their structure and content.


## Features

- **HTML Version Detection**: Identifies the HTML version used (HTML5, HTML4.01, XHTML, etc.)
- **Page Title Extraction**: Retrieves the title of the analyzed page
- **Heading Analysis**: Counts headings by level (h1-h6)
- **Link Analysis**: Categorizes links as internal or external and identifies inaccessible links
- **Login Form Detection**: Determines if the page contains a login form
- **Modern, Responsive UI**: Clean interface built with React and Tailwind CSS
- **API Documentation**: Interactive Swagger documentation for the API endpoints
- **Dockerized Deployment**: Easy containerized deployment with Docker and docker-compose

## Technology Stack

### Backend
- Go (Golang) 1.21+
- Gorilla Mux for routing
- Swagger for API documentation
- Structured logging with slog
- Graceful shutdown handling
- Concurrent link processing

### Frontend
- React 18
- Tailwind CSS
- Axios for API communication
- Responsive design

## Getting Started

### Prerequisites
- Go 1.21 or later
- Node.js 18 or later
- npm 9 or later
- Docker and docker-compose (optional, for containerized deployment)

### Installation and Setup

#### Clone the repository
```bash
git clone https://github.com/yourusername/web-analyzer.git
cd web-analyzer
```

#### Run with Make (Recommended)

The project includes a Makefile with common commands:

```bash
# Install dependencies and build both frontend and backend
make build

# Run the application (builds first if necessary)
make run

# Run backend in development mode
make dev

# Run frontend in development mode (in a separate terminal)
make dev-frontend

# Run tests
make test

# Clean build artifacts
make clean
```

#### Manual Setup

If you prefer to run commands manually:

1. Build the frontend:
```bash
cd web
npm install
npm run build
cd ..
```

2. Build and run the backend:
```bash
go mod tidy
go build -o bin/web-analyzer ./cmd/server
./bin/web-analyzer
```

#### Docker Setup

```bash
# Build and run with docker-compose
docker-compose up --build
```

### Accessing the Application

- Web Interface: [http://localhost:8080](http://localhost:8080)
- Swagger API Documentation: [http://localhost:8080/swagger/](http://localhost:8080/swagger/)
- API Endpoint: [http://localhost:8080/api/analyze](http://localhost:8080/api/analyze)

## API Documentation

### POST /api/analyze

Analyzes a web page by URL.

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
    "h2": 2,
    "h3": 0,
    "h4": 0,
    "h5": 0,
    "h6": 0
  },
  "links": {
    "internal": 5,
    "external": 3,
    "inaccessible": 1
  },
  "containsLoginForm": false
}
```

### GET /api/health

Health check endpoint to verify the API is running.

**Response:**
```json
{
  "status": "ok"
}
```

## Design Decisions and Implementation Details

### HTML Version Detection
The application detects HTML versions by examining:
1. DOCTYPE declarations
2. Presence of HTML5-specific elements (for pages without a DOCTYPE)

### Link Classification
- **Internal Links**: Same domain or relative URLs
- **External Links**: Different domains
- **Accessibility**: Links are checked for accessibility (2xx/3xx status codes)

### Login Form Detection
Login forms are detected through:
- Form attributes containing "login", "signin", etc.
- Presence of password input fields

### Concurrency
The application uses Go's concurrency features:
- Goroutines for parallel link accessibility checking
- Mutex for thread-safe updates
- WaitGroups for synchronization

## Testing

The project includes unit tests for key components:
- HTML version detection
- Title extraction
- Heading counting
- Link classification
- Login form detection

To run tests:
```bash
go test ./...
```

## Challenges and Solutions

- **Link Accessibility Checking**: Implemented concurrent checking with short timeouts to avoid blocking the main analysis.
- **HTML Parsing**: Used golang.org/x/net/html library for reliable HTML parsing regardless of malformed content.
- **Frontend-Backend Integration**: Used structured responses and proper error handling to ensure smooth communication.

## Future Improvements

- Caching of analysis results to improve performance
- More detailed link analysis with metadata
- Performance metrics for page load times
- SEO analysis features
- Accessibility evaluation
- Export functionality for analysis results

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

- The Go team for the excellent net/html package
- The React and Tailwind CSS communities for building great frontend tools
- All contributors and testers who have helped improve this project