package main

import (
	"backend-go/config"
	"backend-go/models"
	"backend-go/routes"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using default values")
	}

	// Connect to database
	config.ConnectDatabase()

	// Auto-migrate the schema
	if err := models.AutoMigrate(); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Initialize Gin router
	router := gin.Default()

	// Add CORS middleware
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "OK",
			"message": "Backend Go API is running",
		})
	})

	// Setup routes
	routes.SetupRoutes(router)

	// Start server
	log.Println("Starting server on :8080...")
	router.Run(":8080")
}
