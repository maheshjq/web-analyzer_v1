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

// @title Web Page Analyzer API
// @version 1.0
// @description API for analyzing web pages, extracting HTML version, title, headings, links, and detecting login forms.

// @contact.name Mahesh Nanayakkara

// @host localhost:8080
// @BasePath /
// @schemes http

func main() {
	// Setup logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	// router
	r := mux.NewRouter()

	// endpoints
	apiRouter := r.PathPrefix("/api").Subrouter()
	apiRouter.HandleFunc("/analyze", api.AnalyzeHandler).Methods("POST", "OPTIONS")
	apiRouter.HandleFunc("/health", api.HealthCheckHandler).Methods("GET")

	// Swagger
	r.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
	))

	// static files
	spa := http.FileServer(http.Dir("./web/build"))
	r.PathPrefix("/").Handler(spa)

	// middleware
	r.Use(api.LoggingMiddleware(logger))
	r.Use(api.RecoverMiddleware(logger))

	// CORS
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}).Handler(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// server
	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      corsHandler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		fmt.Printf("Server starting on :%s\n", port)
		fmt.Println("Swagger UI available at http://localhost:" + port + "/swagger/")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
	fmt.Println("Server gracefully stopped")
}
