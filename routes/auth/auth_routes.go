package auth

import (
	"backend-go/controllers/auth"
	"backend-go/middleware"

	"github.com/gin-gonic/gin"
)

// SetupAuthRoutes sets up authentication routes
func SetupAuthRoutes(router *gin.RouterGroup) {
	// Public routes
	router.POST("/auth/register", auth.Register)
	router.POST("/auth/login", auth.Login)
	router.GET("/auth/roles", auth.GetRoles)

	// Protected routes with middleware chaining
	router.GET("/auth/profile", middleware.AuthMiddleware(), auth.GetProfile)
	router.PUT("/auth/profile", middleware.AuthMiddleware(), auth.UpdateProfile)
}
