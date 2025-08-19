package trip

import (
	"backend-go/controllers/trip"
	"backend-go/middleware"

	"github.com/gin-gonic/gin"
)

// SetupTripRoutes sets up trip-related routes
func SetupTripRoutes(router *gin.RouterGroup) {
	// Public routes (anyone can view trips)
	router.GET("/trips", middleware.OptionalAuth(), trip.GetAll)
	router.GET("/trips/:id", middleware.OptionalAuth(), trip.GetByID)

	// Protected routes with middleware chaining
	// Trip owner only routes - require authentication + trip_owner role
	router.POST("/trips", middleware.AuthMiddleware(), middleware.RequireRole("trip_owner"), trip.Create)
	router.PUT("/trips/:id", middleware.AuthMiddleware(), middleware.RequireRole("trip_owner"), trip.Update)
	router.DELETE("/trips/:id", middleware.AuthMiddleware(), middleware.RequireRole("trip_owner"), trip.Delete)
	router.GET("/trips/my-trips", middleware.AuthMiddleware(), middleware.RequireRole("trip_owner"), trip.GetMyTrips)
}
