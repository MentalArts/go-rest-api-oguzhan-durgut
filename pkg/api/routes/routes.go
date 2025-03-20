package routes

import (
	"mentalartsapi/pkg/api/handlers"
	"mentalartsapi/pkg/api/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all the routes for the application
func SetupRoutes(router *gin.Engine) {
	// Add telemetry middleware to all routes
	router.Use(middleware.Telemetry())

	router.GET("/ping", handlers.HandlePing)
	router.GET("/hello", handlers.HandleHello)
	router.GET("/helloWithPayload", handlers.HandleHelloWithPayload)

	router.POST("/execute", handlers.HandleExecute)

	// Author routes
	router.POST("/author", handlers.CreateAuthor)
	router.GET("/author", handlers.GetAllAuthors)
	router.GET("/author/:id", handlers.GetAuthor)
	router.PUT("/author/:id", handlers.UpdateAuthor)
	router.DELETE("/author/:id", handlers.DeleteAuthor)
}
