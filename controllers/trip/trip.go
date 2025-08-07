package trip

import (
	"backend-go/config"
	"backend-go/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateTripRequest represents the request structure for creating a trip
type CreateTripRequest struct {
	Name           string  `json:"name" binding:"required"`
	Description    string  `json:"description"`
	CoverImage     string  `json:"cover_image"`
	Price          float64 `json:"price" binding:"required,min=0"`
	Duration       int     `json:"duration" binding:"required,min=1"`
	StartLatitude  float64 `json:"start_latitude" binding:"required"`
	StartLongitude float64 `json:"start_longitude" binding:"required"`
	EndLatitude    float64 `json:"end_latitude" binding:"required"`
	EndLongitude   float64 `json:"end_longitude" binding:"required"`
}

// UpdateTripRequest represents the request structure for updating a trip
type UpdateTripRequest struct {
	Name           *string  `json:"name"`
	Description    *string  `json:"description"`
	CoverImage     *string  `json:"cover_image"`
	Price          *float64 `json:"price"`
	Duration       *int     `json:"duration"`
	StartLatitude  *float64 `json:"start_latitude"`
	StartLongitude *float64 `json:"start_longitude"`
	EndLatitude    *float64 `json:"end_latitude"`
	EndLongitude   *float64 `json:"end_longitude"`
}

// GetAll retrieves all trips with optional filtering and pagination
func GetAll(c *gin.Context) {
	var trips []models.Trip
	query := config.DB.Preload("User").Preload("Images")

	// Optional filtering by user ID
	if userID := c.Query("user_id"); userID != "" {
		query = query.Where("user_id = ?", userID)
	}

	// Optional filtering by price range
	if minPrice := c.Query("min_price"); minPrice != "" {
		query = query.Where("price >= ?", minPrice)
	}
	if maxPrice := c.Query("max_price"); maxPrice != "" {
		query = query.Where("price <= ?", maxPrice)
	}

	// Execute query
	if err := query.Find(&trips).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve trips",
			"message": "Could not fetch trips from database",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Trips retrieved successfully",
		"data":    trips,
		"count":   len(trips),
	})
}

// GetByID retrieves a single trip by ID
func GetByID(c *gin.Context) {
	var trip models.Trip
	id := c.Param("id")

	if err := config.DB.Preload("User").Preload("Images").First(&trip, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Trip not found",
			"message": "The requested trip does not exist",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Trip retrieved successfully",
		"data":    trip,
	})
}

// Create creates a new trip (only trip_owner users)
func Create(c *gin.Context) {
	// Check if user is authenticated
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Authentication required",
			"message": "You must be logged in to create a trip",
		})
		return
	}

	// Check if user has trip_owner role
	userRole, exists := c.Get("role")
	if !exists || userRole.(string) != "trip_owner" {
		c.JSON(http.StatusForbidden, gin.H{
			"error":   "Insufficient permissions",
			"message": "Only trip owners can create trips",
		})
		return
	}

	var req CreateTripRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	// Create trip
	trip := models.Trip{
		Name:           req.Name,
		Description:    req.Description,
		CoverImage:     req.CoverImage,
		Price:          req.Price,
		Duration:       req.Duration,
		StartLatitude:  req.StartLatitude,
		StartLongitude: req.StartLongitude,
		EndLatitude:    req.EndLatitude,
		EndLongitude:   req.EndLongitude,
		UserID:         userID.(uint),
	}

	if err := config.DB.Create(&trip).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create trip",
			"message": "Could not save trip to database",
		})
		return
	}

	// Load the user data for the response
	config.DB.Preload("User").First(&trip, trip.ID)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Trip created successfully",
		"data":    trip,
	})
}

// Update updates an existing trip (only the trip owner)
func Update(c *gin.Context) {
	var trip models.Trip
	id := c.Param("id")

	// Find the trip
	if err := config.DB.First(&trip, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Trip not found",
			"message": "The requested trip does not exist",
		})
		return
	}

	// Check if user is authenticated
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Authentication required",
			"message": "You must be logged in to update a trip",
		})
		return
	}

	// Check if user has trip_owner role
	userRole, exists := c.Get("role")
	if !exists || userRole.(string) != "trip_owner" {
		c.JSON(http.StatusForbidden, gin.H{
			"error":   "Insufficient permissions",
			"message": "Only trip owners can update trips",
		})
		return
	}

	// Check if user owns the trip
	if trip.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{
			"error":   "Access denied",
			"message": "You can only update your own trips",
		})
		return
	}

	var req UpdateTripRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	// Update only provided fields
	if req.Name != nil {
		trip.Name = *req.Name
	}
	if req.Description != nil {
		trip.Description = *req.Description
	}
	if req.CoverImage != nil {
		trip.CoverImage = *req.CoverImage
	}
	if req.Price != nil {
		trip.Price = *req.Price
	}
	if req.Duration != nil {
		trip.Duration = *req.Duration
	}
	if req.StartLatitude != nil {
		trip.StartLatitude = *req.StartLatitude
	}
	if req.StartLongitude != nil {
		trip.StartLongitude = *req.StartLongitude
	}
	if req.EndLatitude != nil {
		trip.EndLatitude = *req.EndLatitude
	}
	if req.EndLongitude != nil {
		trip.EndLongitude = *req.EndLongitude
	}

	if err := config.DB.Save(&trip).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update trip",
			"message": "Could not save changes to database",
		})
		return
	}

	// Load the user data for the response
	config.DB.Preload("User").First(&trip, trip.ID)

	c.JSON(http.StatusOK, gin.H{
		"message": "Trip updated successfully",
		"data":    trip,
	})
}

// Delete deletes a trip (only the trip owner)
func Delete(c *gin.Context) {
	var trip models.Trip
	id := c.Param("id")

	// Find the trip
	if err := config.DB.First(&trip, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Trip not found",
			"message": "The requested trip does not exist",
		})
		return
	}

	// Check if user is authenticated
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Authentication required",
			"message": "You must be logged in to delete a trip",
		})
		return
	}

	// Check if user has trip_owner role
	userRole, exists := c.Get("role")
	if !exists || userRole.(string) != "trip_owner" {
		c.JSON(http.StatusForbidden, gin.H{
			"error":   "Insufficient permissions",
			"message": "Only trip owners can delete trips",
		})
		return
	}

	// Check if user owns the trip
	if trip.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{
			"error":   "Access denied",
			"message": "You can only delete your own trips",
		})
		return
	}

	if err := config.DB.Delete(&trip).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to delete trip",
			"message": "Could not remove trip from database",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Trip deleted successfully",
		"data": gin.H{
			"id":   trip.ID,
			"name": trip.Name,
		},
	})
}

// GetMyTrips retrieves all trips created by the current user
func GetMyTrips(c *gin.Context) {
	// Check if user is authenticated
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Authentication required",
			"message": "You must be logged in to view your trips",
		})
		return
	}

	// Check if user has trip_owner role
	userRole, exists := c.Get("role")
	if !exists || userRole.(string) != "trip_owner" {
		c.JSON(http.StatusForbidden, gin.H{
			"error":   "Insufficient permissions",
			"message": "Only trip owners can view their trips",
		})
		return
	}

	var trips []models.Trip
	if err := config.DB.Preload("User").Preload("Images").Where("user_id = ?", userID).Find(&trips).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve trips",
			"message": "Could not fetch your trips from database",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Your trips retrieved successfully",
		"data":    trips,
		"count":   len(trips),
	})
}
