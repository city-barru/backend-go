package user

import (
	"backend-go/controllers/user"
	"backend-go/middleware"

	"github.com/gin-gonic/gin"
)

// SetupUserRoutes sets up user-related routes
func SetupUserRoutes(router *gin.RouterGroup) {
	// All user routes require authentication, some require admin role
	router.GET("/users", middleware.AuthMiddleware(), middleware.RequireRole("admin"), user.GetAll) // Only admin can get all users
	router.GET("/users/:id", middleware.AuthMiddleware(), user.GetByID)
	router.POST("/users", middleware.AuthMiddleware(), middleware.RequireRole("admin"), user.Create) // Only admin can create users
	router.PUT("/users/:id", middleware.AuthMiddleware(), user.Update)
	router.DELETE("/users/:id", middleware.AuthMiddleware(), middleware.RequireRole("admin"), user.Delete) // Only admin can delete users
}
