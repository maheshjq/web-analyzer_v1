package api

import (
	"net/http"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	_ "github.com/maheshjq/web-analyzer_v1/docs" // This is where the generated swagger docs are
)

// SetupSwagger adds Swagger documentation routes to the router
func SetupSwagger(router *mux.Router) {
	// Serve the Swagger UI
	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
}