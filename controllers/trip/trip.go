package trip

import (
	"backend-go/config"
	"backend-go/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetAll retrieves all trips
func GetAll(c *gin.Context) {
	var trips []models.Trip
	config.DB.Preload("User").Find(&trips)
	c.JSON(http.StatusOK, gin.H{"data": trips})
}

// GetByID retrieves a single trip by ID
func GetByID(c *gin.Context) {
	var trip models.Trip
	id := c.Param("id")

	if err := config.DB.Preload("User").First(&trip, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Trip not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": trip})
}

// Create creates a new trip
func Create(c *gin.Context) {
	var trip models.Trip

	if err := c.ShouldBindJSON(&trip); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	trip.UserID = userID.(uint)

	config.DB.Create(&trip)
	c.JSON(http.StatusCreated, gin.H{"data": trip})
}

// Update updates an existing trip
func Update(c *gin.Context) {
	var trip models.Trip
	id := c.Param("id")

	if err := config.DB.First(&trip, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Trip not found"})
		return
	}

	// Check if user owns the trip or is admin
	userID, _ := c.Get("userID")
	userRole, _ := c.Get("role")

	if trip.UserID != userID.(uint) && userRole.(string) != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only update your own trips"})
		return
	}

	if err := c.ShouldBindJSON(&trip); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	config.DB.Save(&trip)
	c.JSON(http.StatusOK, gin.H{"data": trip})
}

// Delete deletes a trip
func Delete(c *gin.Context) {
	var trip models.Trip
	id := c.Param("id")

	if err := config.DB.First(&trip, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Trip not found"})
		return
	}

	// Check if user owns the trip or is admin
	userID, _ := c.Get("userID")
	userRole, _ := c.Get("role")

	if trip.UserID != userID.(uint) && userRole.(string) != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only delete your own trips"})
		return
	}

	config.DB.Delete(&trip)
	c.JSON(http.StatusOK, gin.H{"message": "Trip deleted successfully"})
}
