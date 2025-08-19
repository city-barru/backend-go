package preference

import (
	"backend-go/controllers/preference"
	"backend-go/middleware"

	"github.com/gin-gonic/gin"
)

func SetupPreferenceRoutes(router *gin.RouterGroup) {
	// Public route
	router.GET("/preferences", preference.GetAll) // Get all preferences

	// Protected routes with middleware chaining
	router.GET("/preferences/:id", middleware.AuthMiddleware(), preference.GetByID)
	router.POST("/preferences", middleware.AuthMiddleware(), middleware.RequireRole("admin"), preference.Create)       // Only admin can create preferences
	router.PUT("/preferences/:id", middleware.AuthMiddleware(), middleware.RequireRole("admin"), preference.Update)    // Only admin can update preferences
	router.DELETE("/preferences/:id", middleware.AuthMiddleware(), middleware.RequireRole("admin"), preference.Delete) // Only admin can delete preferences

	router.POST("/preferences/assign", middleware.AuthMiddleware(), preference.AssignPreference) // Assign preference to user
}
