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
			// Trip owner only routes - require trip_owner role
			tripOwnerRoutes := protected.Group("/")
			tripOwnerRoutes.Use(middleware.RequireRole("trip_owner"))
			{
				tripOwnerRoutes.POST("/", trip.Create)
				tripOwnerRoutes.PUT("/:id", trip.Update)
				tripOwnerRoutes.DELETE("/:id", trip.Delete)
				tripOwnerRoutes.GET("/my-trips", trip.GetMyTrips)
			}
		}
	}
}
