package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

// @title Web Page Analyzer API
// @version 1.0
// @description API for analyzing web pages, extracting HTML version, title, headings, links, and detecting login forms.

// @contact.name Web Analyzer Team

// @host localhost:8080
// @BasePath /api

// @schemes http

func main() {
	// Create router
	r := mux.NewRouter()
	
	// API endpoints (you'll implement these later)
	apiRouter := r.PathPrefix("/api").Subrouter()
	apiRouter.HandleFunc("/analyze", analyzeHandler).Methods("POST")
	apiRouter.HandleFunc("/health", healthCheckHandler).Methods("GET")
	
	// Serve static files from the React build folder
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./web/build")))
	
	// Wrap router with CORS middleware
	handler := cors.Default().Handler(r)
	
	// Create HTTP server
	srv := &http.Server{
		Addr:         ":8080",
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	
	// Start server in a goroutine
	go func() {
		fmt.Println("Server starting on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()
	
	// Wait for interrupt signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	
	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
	fmt.Println("Server gracefully stopped")
}

// Placeholder API handlers (implement these with your analyzer logic)
func analyzeHandler(w http.ResponseWriter, r *http.Request) {
	// Your analyze implementation will go here
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"message": "Analysis endpoint placeholder"}`)
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"status": "ok"}`)
}