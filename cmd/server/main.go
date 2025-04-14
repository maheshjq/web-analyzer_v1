package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/maheshjq/web-analyzer_v1/internal/api"
	// Uncomment when you have generated swagger docs
	_ "github.com/maheshjq/web-analyzer_v1/docs"
)

func init() {
	// api.NewAnalyzerFunc = func() api.Analyzer {
	// 	// Create default analyzer
	// 	defaultAnalyzer := &api.DefaultAnalyzer{}

	// 	// Wrap with caching (cache results for 5 minutes)
	// 	return api.NewCachedAnalyzer(defaultAnalyzer, 5*time.Minute)
	// }
}

// @title Web Page Analyzer API
// @version 1.0
// @description API for analyzing web pages, extracting HTML version, title, headings, links, and detecting login forms.

// @contact.name Web Analyzer Team

// @host localhost:8080
// @BasePath /
// @schemes http

func main() {
	api.EnableCaching = true
	// Setup logger for the app
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	// Create the router for handling routes
	router := mux.NewRouter()

	// Define all API endpoints here
	apiRouter := router.PathPrefix("/api").Subrouter()
	apiRouter.HandleFunc("/analyze", api.AnalyzeHandler).Methods("POST", "OPTIONS")
	apiRouter.HandleFunc("/health", api.HealthCheckHandler).Methods("GET")

	// Setup Swagger UI for API docs
	router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
	))

	// Serve static files from the web build directory
	spa := http.FileServer(http.Dir("./web/build"))
	router.PathPrefix("/").Handler(spa)

	// Add middleware for logging and recovery
	router.Use(api.LoggingMiddleware(logger))
	router.Use(api.RecoverMiddleware(logger))

	// Setup CORS so the frontend can talk to the backend
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}).Handler(router)

	// Get the port from env or use 8080 as default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Configure the HTTP server
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      corsHandler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Start the server in a goroutine
	go func() {
		fmt.Printf("Server starting on :%s\n", port)
		fmt.Println("Swagger UI available at http://localhost:" + port + "/swagger/")
		// logger.Info("Server started", "port", port) // Left this commented out from debugging
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server crashed: %v", err)
		}
	}()

	// Wait for an interrupt signal to shut down
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	// Gracefully shutdown the server
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logger.Error("Shutdown failed", "error", err)
	}
	fmt.Println("Server gracefully stopped")
}
