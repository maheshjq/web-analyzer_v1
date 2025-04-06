# Build stage
FROM golang:1.21-alpine AS build

WORKDIR /app

# Install dependencies
RUN apk --no-cache add ca-certificates git

# Copy go.mod 
COPY go.mod ./
# Create an empty go.sum
RUN touch go.sum

# Initialize module and fetch dependencies
RUN go mod download
RUN go mod tidy

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o web-analyzer ./cmd/server

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy the binary from the build stage
COPY --from=build /app/web-analyzer .

# Create the web/build directory
RUN mkdir -p web/build

# Expose the application port
EXPOSE 8080

# Run the application
CMD ["./web-analyzer"]