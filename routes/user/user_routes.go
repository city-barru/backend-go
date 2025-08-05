package user

import (
	"backend-go/controllers/user"
	"backend-go/middleware"

	"github.com/gin-gonic/gin"
)

// SetupUserRoutes sets up user-related routes
func SetupUserRoutes(router *gin.RouterGroup) {
	userGroup := router.Group("/users")
	userGroup.Use(middleware.AuthMiddleware()) // Protect all user routes
	{
		userGroup.GET("/", middleware.RequireRole("admin"), user.GetAll) // Only admin can get all users
		userGroup.GET("/:id", user.GetByID)
		userGroup.POST("/", middleware.RequireRole("admin"), user.Create) // Only admin can create users
		userGroup.PUT("/:id", user.Update)
		userGroup.DELETE("/:id", middleware.RequireRole("admin"), user.Delete) // Only admin can delete users
	}
}
