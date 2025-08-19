package image

import (
	"backend-go/controllers/image"
	"backend-go/middleware"

	"github.com/gin-gonic/gin"
)

// SetupImageRoutes sets up image-related routes
func SetupImageRoutes(router *gin.RouterGroup) {
	// Public routes to serve images
	router.GET("/images/:filename", image.GetImage)
	router.GET("/covers/:filename", image.GetCoverImage)

	// Protected routes with middleware chaining
	router.POST("/images/upload", middleware.AuthMiddleware(), image.Upload)
	router.POST("/images/upload-cover", middleware.AuthMiddleware(), image.UploadCoverImage)
	router.GET("/images/my-images", middleware.AuthMiddleware(), image.GetMyImages)
	router.GET("/images/trip/:trip_id", middleware.OptionalAuth(), image.GetImagesByTrip)
	router.DELETE("/images/:id", middleware.AuthMiddleware(), image.DeleteImage)
}
