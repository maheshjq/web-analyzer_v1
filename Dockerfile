FROM golang:1.21-alpine AS build

WORKDIR /app

RUN apk --no-cache add ca-certificates git

COPY go.mod ./
COPY go.sum ./

RUN go mod download
RUN go mod verify

COPY . .

RUN go mod tidy

RUN go get -d -v ./...

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o web-analyzer ./cmd/server

FROM alpine:latest

WORKDIR /app

COPY --from=build /app/web-analyzer .

RUN mkdir -p web/build

EXPOSE 8080

CMD ["./web-analyzer"]