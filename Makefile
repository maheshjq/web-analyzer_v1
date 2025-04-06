.PHONY: all build build-frontend build-backend run test clean

# Default target
all: build

# Build everything
build: build-frontend build-backend

# Build frontend
build-frontend:
	@echo "Building frontend..."
	cd web && npm install && npm run build

# Build backend
build-backend:
	@echo "Building backend..."
	go mod tidy
	go build -o bin/web-analyzer ./cmd/server

# Run the application
run: build
	@echo "Starting web-analyzer..."
	./bin/web-analyzer

# Run in development mode
dev:
	@echo "Starting backend in development mode..."
	go run cmd/server/main.go

# Run frontend in development mode
dev-frontend:
	@echo "Starting frontend in development mode..."
	cd web && npm start

# Run tests
test:
	@echo "Running tests..."
	go test ./...

# Clean build artifacts
clean:
	@echo "Cleaning up..."
	rm -rf bin/ web/build/ web/node_modules/