package auth

import (
	"backend-go/controllers/auth"
	"backend-go/middleware"

	"github.com/gin-gonic/gin"
)

// SetupAuthRoutes sets up authentication routes
func SetupAuthRoutes(router *gin.RouterGroup) {
	authGroup := router.Group("/auth")
	{
		// Public routes
		authGroup.POST("/register", auth.Register)
		authGroup.POST("/login", auth.Login)

		// Protected routes
		protected := authGroup.Group("/")
		protected.Use(middleware.AuthMiddleware())
		{
			protected.GET("/profile", auth.GetProfile)
			protected.PUT("/profile", auth.UpdateProfile)
		}
	}
}
