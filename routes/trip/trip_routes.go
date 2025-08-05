package trip

import (
	"backend-go/controllers/trip"
	"backend-go/middleware"

	"github.com/gin-gonic/gin"
)

// SetupTripRoutes sets up trip-related routes
func SetupTripRoutes(router *gin.RouterGroup) {
	tripGroup := router.Group("/trips")
	{
		// Public routes (anyone can view trips)
		tripGroup.GET("/", middleware.OptionalAuth(), trip.GetAll)
		tripGroup.GET("/:id", middleware.OptionalAuth(), trip.GetByID)

		// Protected routes (authentication required)
		protected := tripGroup.Group("/")
		protected.Use(middleware.AuthMiddleware())
		{
			protected.POST("/", trip.Create)
			protected.PUT("/:id", trip.Update)
			protected.DELETE("/:id", trip.Delete)
		}
	}
}
