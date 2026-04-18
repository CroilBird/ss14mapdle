// Package endpoints contains all the endpointGenerators of the API
package endpoints

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// endpointGenerators is a list of functions that return an Endpoint
// this is done so that the database can be passed into each handler
// separately. This is useful for testing.
var endpointGenerators []func(db *gorm.DB) Endpoint

// RegisterEndpointGenerator adds an endpoint generator function
func RegisterEndpointGenerator(generator func(db *gorm.DB) Endpoint) {
	endpointGenerators = append(endpointGenerators, generator)
}

// Endpoint represents a collection of routes for a given endpoint,
// each with a method, path, and handler
type Endpoint struct {
	Routes []Route
}

// RegisterEndpointHandlers registers all the routes from all modules of the API
// this is done after all generators have been added
// since the generators are added in init functions, this function should be called after all
// modules have been imported
// Typically this is done in the main function
func RegisterEndpointHandlers(router *gin.Engine, db *gorm.DB) {
	for _, generator := range endpointGenerators {
		endpoint := generator(db)
		registerRouteHandlers(router, endpoint.Routes)
	}
}

// Route represents a single route with a method, path, and handler
type Route struct {
	Method  string
	Path    string
	Handler func(c *gin.Context)
}

func registerRouteHandlers(router *gin.Engine, routes []Route) {
	for _, route := range routes {
		router.Handle(route.Method, route.Path, route.Handler)
	}
}
