.PHONY: all build build-frontend build-backend run test clean

# def target
all: build

build: build-frontend build-backend

build-frontend:
	@echo "Building frontend..."
	cd web && npm install && npm run build

build-backend:
	@echo "Building backend..."
	go mod tidy
	go build -o bin/web-analyzer ./cmd/server

# Run app
run: build
	@echo "Starting web-analyzer..."
	./bin/web-analyzer

dev:
	@echo "Starting backend in development mode..."
	go run cmd/server/main.go

# run frontend
dev-frontend:
	@echo "Starting frontend in development mode..."
	cd web && npm start

test:
	@echo "Running tests..."
	go test ./...

clean:
	@echo "Cleaning up..."
	rm -rf bin/ web/build/ web/node_modules/