#!/bin/bash

# Install Swag if not already installed
if ! command -v swag &> /dev/null; then
    echo "Installing swag..."
    go install github.com/swaggo/swag/cmd/swag@latest
fi

# Generate Swagger docs
echo "Generating Swagger documentation..."
swag init -g cmd/server/main.go -o docs

echo "Swagger documentation generated in docs/ directory"
echo "Run the server and access the Swagger UI at http://localhost:8080/swagger/index.html"