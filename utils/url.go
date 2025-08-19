package utils

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
)

// GetBaseURL returns the base URL for the application
func GetBaseURL(c *gin.Context) string {
	// Check if BASE_URL is set in environment
	if baseURL := os.Getenv("BASE_URL"); baseURL != "" {
		return baseURL
	}

	// Fallback to constructing from request
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}

	// Check for X-Forwarded-Proto header (common in reverse proxy setups)
	if proto := c.GetHeader("X-Forwarded-Proto"); proto != "" {
		scheme = proto
	}

	host := c.Request.Host
	return fmt.Sprintf("%s://%s", scheme, host)
}

// GetFullImageURL returns the full URL for an image path
func GetFullImageURL(c *gin.Context, imagePath string) string {
	baseURL := GetBaseURL(c)
	return fmt.Sprintf("%s%s", baseURL, imagePath)
}
