package routes

import (
	"backend-go/routes/auth"
	"backend-go/routes/image"
	"backend-go/routes/preference"
	"backend-go/routes/trip"
	"backend-go/routes/user"

	"github.com/gin-gonic/gin"
)

// SetupRoutes initializes all route groups
func SetupRoutes(router *gin.Engine) {
	// API version 1
	v1 := router.Group("/api/v1")
	{
		// Authentication routes (public)
		auth.SetupAuthRoutes(v1)

		// User routes (protected)
		user.SetupUserRoutes(v1)

		// Trip routes (protected)
		trip.SetupTripRoutes(v1)

		// Preference routes (public & protected)
		preference.SetupPreferenceRoutes(v1)

		// Image routes (public & protected)
		image.SetupImageRoutes(v1)
	}
}
