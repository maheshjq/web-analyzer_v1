package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/maheshjq/web-analyzer_v1/internal/api"
	"github.com/maheshjq/web-analyzer_v1/internal/metrics"

	// Uncomment when you have generated swagger docs
	_ "github.com/maheshjq/web-analyzer_v1/docs"
)

func init() {
	metrics.Initialize()
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
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	router := mux.NewRouter()

	// Define all API endpoints
	apiRouter := router.PathPrefix("/api").Subrouter()
	apiRouter.HandleFunc("/analyze", api.AnalyzeHandler).Methods("POST", "OPTIONS")
	apiRouter.HandleFunc("/health", api.HealthCheckHandler).Methods("GET")

	router.Handle("/metrics", metrics.MetricsHandler())

	// Add pprof endpoints
	// These handlers are used for performance profiling
	router.HandleFunc("/debug/pprof/", pprof.Index)
	router.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	router.HandleFunc("/debug/pprof/profile", pprof.Profile)
	router.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	router.HandleFunc("/debug/pprof/trace", pprof.Trace)
	router.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
	router.Handle("/debug/pprof/heap", pprof.Handler("heap"))
	router.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))
	router.Handle("/debug/pprof/block", pprof.Handler("block"))

	// Setup Swagger UI for API docs
	router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
	))

	// Serve static files from the web build directory
	spa := http.FileServer(http.Dir("./web/build"))
	router.PathPrefix("/").Handler(spa)

	// Add middleware for metrics, logging and recovery
	router.Use(api.MetricsMiddleware)
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
		fmt.Println("Prometheus metrics available at http://localhost:" + port + "/metrics")
		fmt.Println("pprof debugging available at http://localhost:" + port + "/debug/pprof/")

		// logger.Info("Server started", "port", port) // this is for debugging
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
