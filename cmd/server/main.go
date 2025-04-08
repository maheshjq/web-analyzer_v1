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
	httpSwagger "github.com/swaggo/http-swagger"
	
	"github.com/maheshjq/web-analyzer_v1/internal/api"
	// Import swagger docs when they are generated
	// _ "github.com/maheshjq/web-analyzer_v1/docs"
)

// @title Web Page Analyzer API
// @version 1.0
// @description API for analyzing web pages, extracting HTML version, title, headings, links, and detecting login forms.

// @contact.name Web Analyzer Team

// @host localhost:8080
// @BasePath /
// @schemes http

func main() {
	// Create router
	r := mux.NewRouter()
	
	// API endpoints
	apiRouter := r.PathPrefix("/api").Subrouter()
	apiRouter.HandleFunc("/analyze", api.AnalyzeHandler).Methods("POST")
	apiRouter.HandleFunc("/health", api.HealthCheckHandler).Methods("GET")
	
	// Setup Swagger
	// The URL should be the same as the swagger docs
	r.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
	))
	
	// Serve static files from the React build folder
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./web/build")))
	
	// Wrap router with CORS middleware
	corsHandler := cors.Default().Handler(r)
	
	// Create HTTP server
	srv := &http.Server{
		Addr:         ":8080",
		Handler:      corsHandler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	
	// Start server in a goroutine
	go func() {
		fmt.Println("Server starting on :8080")
		fmt.Println("Swagger UI available at http://localhost:8080/swagger/")
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