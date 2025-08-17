package preference

import (
	"backend-go/controllers/preference"
	"backend-go/middleware"

	"github.com/gin-gonic/gin"
)

func SetupPreferenceRoutes(router *gin.RouterGroup) {
	preferenceGroup := router.Group("/preferences")
	preferenceGroup.GET("/", preference.GetAll) // Get all preferences
	preferenceGroup.Use(middleware.AuthMiddleware()) // Protect all preference routes
	{
		preferenceGroup.GET("/:id", preference.GetByID)
		preferenceGroup.POST("/", middleware.RequireRole("admin"), preference.Create)      // Only admin can create preferences
		preferenceGroup.PUT("/:id", middleware.RequireRole("admin"), preference.Update)    // Only admin can update preferences
		preferenceGroup.DELETE("/:id", middleware.RequireRole("admin"), preference.Delete) // Only admin can delete preferences
		
		preferenceGroup.POST("/assign", preference.AssignPreference) // Assign preference to user
	}
}
