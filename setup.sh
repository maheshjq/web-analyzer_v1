#!/bin/bash

# Create project directories
mkdir -p cmd/server
mkdir -p internal/analyzer
mkdir -p internal/api
mkdir -p internal/models
mkdir -p web/build

# Create go.mod file
cat > go.mod << 'EOF'
module github.com/mahesh/web-analyzer

go 1.21

require (
	github.com/gorilla/mux v1.8.1
	github.com/rs/cors v1.10.1
	golang.org/x/net v0.19.0
	golang.org/x/sync v0.5.0
)
EOF

# Create minimal main.go
cat > cmd/server/main.go << 'EOF'
package main

import (
	"fmt"
	"net/http"
	
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	r := mux.NewRouter()
	
	// Just a placeholder endpoint
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Web Analyzer")
	})
	
	// Wrap router with CORS middleware
	handler := cors.Default().Handler(r)
	
	fmt.Println("Server starting on :8080")
	http.ListenAndServe(":8080", handler)
}
EOF

# Create minimal analyzer.go
cat > internal/analyzer/analyzer.go << 'EOF'
package analyzer

import (
	"fmt"
	"strings"
	
	"golang.org/x/net/html"
	"golang.org/x/sync/errgroup"
)

// Analyze is a minimal implementation to satisfy imports
func Analyze(url string) (string, error) {
	// This is just a placeholder to satisfy the imports
	g := new(errgroup.Group)
	g.Go(func() error {
		return nil
	})
	
	// Using the html package to satisfy import
	_, err := html.Parse(strings.NewReader("<html><body>Test</body></html>"))
	if err != nil {
		return "", err
	}
	
	return "Analysis complete", nil
}
EOF

# Create minimal models.go
cat > internal/models/models.go << 'EOF'
package models

// AnalysisRequest represents the request to analyze a webpage
type AnalysisRequest struct {
	URL string `json:"url"`
}

// AnalysisResponse represents the result of the analysis
type AnalysisResponse struct {
	HTMLVersion string `json:"htmlVersion"`
}
EOF

# Create minimal handlers.go
cat > internal/api/handlers.go << 'EOF'
package api

import (
	"fmt"
	"net/http"
)

// AnalyzeHandler is a minimal implementation
func AnalyzeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "{\"status\":\"ok\"}")
}
EOF

# Create empty index.html in web/build
cat > web/build/index.html << 'EOF'
<!DOCTYPE html>
<html>
<head>
    <title>Web Analyzer</title>
</head>
<body>
    <h1>Web Analyzer</h1>
</body>
</html>
EOF

echo "Setup complete. You can now build the Docker image."