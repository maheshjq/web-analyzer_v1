package api

import (
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

// SetupSwagger adds Swagger documentation routes to the router
func SetupSwagger(router *mux.Router) {
	// Serve the Swagger UI
	router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"), // The URL of the API docs
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	))
}