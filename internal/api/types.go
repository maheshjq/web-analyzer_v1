package api

import (
	"log/slog"
	"net/http"
)

// Handler encapsulates the dependencies for the API handlers
type Handler struct {
	logger *slog.Logger
}

// NewHandler creates a new Handler instance
func NewHandler(logger *slog.Logger) *Handler {
	return &Handler{
		logger: logger,
	}
}

// sendError sends an error response with the given status code and message
func (h *Handler) sendError(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	w.Write([]byte(`{"statusCode": ` + string(rune(statusCode)) + `, "message": "` + message + `"}`))
}