package api

import (
	"fmt"
	"net/http"
)

// AnalyzeHandler is a minimal implementation
func AnalyzeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "{\"status\":\"ok\"}")
}
