package trip

import (
	"backend-go/config"
	"backend-go/models"
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CreateTripRequest represents the request structure for creating a trip
type CreateTripRequest struct {
	Name           string             `json:"name" binding:"required"`
	Description    string             `json:"description"`
	CoverImage     string             `json:"cover_image"`
	Price          float64            `json:"price" binding:"required,min=0"`
	Duration       int                `json:"duration" binding:"required,min=1"`
	StartLatitude  float64            `json:"start_latitude" binding:"required"`
	StartLongitude float64            `json:"start_longitude" binding:"required"`
	EndLatitude    float64            `json:"end_latitude" binding:"required"`
	EndLongitude   float64            `json:"end_longitude" binding:"required"`
	PreferenceIDs  []uint             `json:"preference_ids"`
	Points         []models.TripPoint `json:"points"`
}

// UpdateTripRequest represents the request structure for updating a trip
type UpdateTripRequest struct {
	Name           *string            `json:"name"`
	Description    *string            `json:"description"`
	CoverImage     *string            `json:"cover_image"`
	Price          *float64           `json:"price"`
	Duration       *int               `json:"duration"`
	StartLatitude  *float64           `json:"start_latitude"`
	StartLongitude *float64           `json:"start_longitude"`
	EndLatitude    *float64           `json:"end_latitude"`
	EndLongitude   *float64           `json:"end_longitude"`
	PreferenceIDs  []uint             `json:"preference_ids"`
	Points         []models.TripPoint `json:"points"`
}

// GetAll retrieves all trips with optional filtering and pagination
func GetAll(c *gin.Context) {
	var trips []models.Trip
	query := config.DB.Preload("User").Preload("Images").Preload("Preferences").Preload("Points")

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

	if err := config.DB.Preload("User").Preload("Images").Preload("Preferences").Preload("Points").First(&trip, id).Error; err != nil {
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

	// Find preferences
	var preferences []models.Preference
	if len(req.PreferenceIDs) > 0 {
		if err := config.DB.Find(&preferences, req.PreferenceIDs).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid preference IDs"})
			return
		}
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
		Preferences:    preferences,
		Points:         req.Points,
	}

	if err := config.DB.Create(&trip).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create trip",
			"message": "Could not save trip to database",
		})
		return
	}

	// Load associations for the response
	config.DB.Preload("User").Preload("Preferences").Preload("Points").First(&trip, trip.ID)

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

	// Use a transaction for atomic updates
	tx := config.DB.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
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

	if err := tx.Save(&trip).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update trip",
			"message": "Could not save changes to database",
		})
		return
	}

	// Update preferences if provided
	if req.PreferenceIDs != nil {
		var preferences []models.Preference
		if len(req.PreferenceIDs) > 0 {
			if err := tx.Find(&preferences, req.PreferenceIDs).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid preference IDs"})
				return
			}
		}
		if err := tx.Model(&trip).Association("Preferences").Replace(preferences); err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update preferences"})
			return
		}
	}

	// Update points if provided
	if req.Points != nil {
		if err := tx.Model(&trip).Association("Points").Replace(req.Points); err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update points"})
			return
		}
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	// Load the user data for the response
	config.DB.Preload("User").Preload("Preferences").Preload("Points").First(&trip, trip.ID)

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

func SeedTrips(c *gin.Context) {
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

	// Query OpenStreetMap for Jakarta tourism attractions
	osmQuery := `[out:json][timeout:25];
	(
		way["tourism"~"attraction|museum|viewpoint|gallery|theme_park|zoo|aquarium|artwork"]
		   ["name"]
		   (bbox:-6.3713,-106.9758,-6.0835,-106.6486);
		relation["tourism"~"attraction|museum|viewpoint|gallery|theme_park|zoo|aquarium"]
		        ["name"]
		        (bbox:-6.3713,-106.9758,-6.0835,-106.6486);
		node["tourism"~"attraction|museum|viewpoint|gallery|theme_park|zoo|aquarium|artwork"]
		    ["name"]
		    (bbox:-6.3713,-106.9758,-6.0835,-106.6486);
	);
	out center meta;`

	osmResp, err := http.Post("https://overpass-api.de/api/interpreter",
		"application/x-www-form-urlencoded",
		strings.NewReader("data="+url.QueryEscape(osmQuery)))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to query OpenStreetMap",
		})
		return
	}
	defer osmResp.Body.Close()

	var osmData struct {
		Elements []struct {
			Type   string  `json:"type"`
			ID     int64   `json:"id"`
			Lat    float64 `json:"lat"`
			Lon    float64 `json:"lon"`
			Center *struct {
				Lat float64 `json:"lat"`
				Lon float64 `json:"lon"`
			} `json:"center,omitempty"`
			Tags map[string]string `json:"tags"`
		} `json:"elements"`
	}

	if err := json.NewDecoder(osmResp.Body).Decode(&osmData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to parse OSM data",
		})
		return
	}

	// Preference mapping
	preferenceMap := map[string][]string{
		"attraction": {"Sightseeing", "Photography", "Culture", "History"},
		"museum":     {"Education", "Culture", "History", "Art"},
		"viewpoint":  {"Photography", "Sightseeing", "Nature"},
		"gallery":    {"Art", "Culture", "Photography"},
		"theme_park": {"Family", "Entertainment", "Adventure"},
		"zoo":        {"Family", "Education", "Nature"},
		"aquarium":   {"Family", "Education", "Marine Life"},
		"artwork":    {"Art", "Culture", "Photography"},
	}

	var trips []models.Trip
	createdPrefs := make(map[string]uint) // Cache preferences

	for i, element := range osmData.Elements {
		if i >= 20 { // Limit to 20 attractions to avoid too many requests
			break
		}

		name := element.Tags["name"]
		if name == "" {
			continue
		}

		// Get coordinates
		lat, lon := element.Lat, element.Lon
		if element.Center != nil {
			lat, lon = element.Center.Lat, element.Center.Lon
		}

		// Generate description and additional info
		description := generateDescription(element.Tags)

		// Create trip points (simulated popular spots)
		var tripPoints []models.TripPoint
		for j := 0; j < 3; j++ {
			// Generate nearby points for photo spots
			offsetLat := lat + (rand.Float64()-0.5)*0.001
			offsetLon := lon + (rand.Float64()-0.5)*0.001

			tripPoints = append(tripPoints, models.TripPoint{
				Latitude:  offsetLat,
				Longitude: offsetLon,
			})
		}

		// Get preferences for this attraction
		tourismType := element.Tags["tourism"]
		var preferences []models.Preference

		if prefNames, exists := preferenceMap[tourismType]; exists {
			for _, prefName := range prefNames {
				if prefID, cached := createdPrefs[prefName]; cached {
					preferences = append(preferences, models.Preference{Model: gorm.Model{ID: prefID}})
				} else {
					// Create new preference
					pref := models.Preference{Name: prefName}
					if err := config.DB.FirstOrCreate(&pref, models.Preference{Name: prefName}).Error; err == nil {
						preferences = append(preferences, pref)
						createdPrefs[prefName] = pref.ID
					}
				}
			}
		}

		// Generate realistic price and duration based on type
		price, duration := generatePriceAndDuration(tourismType)

		trip := models.Trip{
			Name:           name,
			Description:    description,
			CoverImage:     generateCoverImage(tourismType),
			Price:          price,
			Duration:       duration,
			StartLatitude:  lat,
			StartLongitude: lon,
			EndLatitude:    lat,
			EndLongitude:   lon,
			UserID:         userID.(uint),
			Points:         tripPoints,
			Preferences:    preferences,
		}

		trips = append(trips, trip)
	}

	if len(trips) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"message": "No suitable attractions found",
		})
		return
	}

	if err := config.DB.Create(&trips).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to seed trips",
			"message": "Could not add trips to database",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("%d trips seeded successfully", len(trips)),
		"trips":   len(trips),
	})
}

func generateDescription(tags map[string]string) string {
	desc := "Explore this amazing attraction in Jakarta. "

	if addr := tags["addr:full"]; addr != "" {
		desc += "Located at " + addr + ". "
	}

	if website := tags["website"]; website != "" {
		desc += "Visit their website for more information. "
	}

	if phone := tags["phone"]; phone != "" {
		desc += "Contact: " + phone + ". "
	}

	desc += "Perfect for photography and sightseeing!"
	return desc
}

func generatePriceAndDuration(tourismType string) (float64, int) {
	priceMap := map[string][2]float64{
		"museum":     {15000, 50000},
		"theme_park": {100000, 300000},
		"zoo":        {30000, 80000},
		"aquarium":   {50000, 150000},
		"attraction": {0, 25000},
		"viewpoint":  {0, 15000},
		"gallery":    {10000, 40000},
		"artwork":    {0, 10000},
	}

	durationMap := map[string][2]int{
		"museum":     {60, 180},
		"theme_park": {240, 480},
		"zoo":        {120, 300},
		"aquarium":   {90, 240},
		"attraction": {30, 120},
		"viewpoint":  {20, 60},
		"gallery":    {45, 120},
		"artwork":    {10, 30},
	}

	priceRange := priceMap["attraction"]       // default
	durationRange := durationMap["attraction"] // default

	if pr, exists := priceMap[tourismType]; exists {
		priceRange = pr
	}
	if dr, exists := durationMap[tourismType]; exists {
		durationRange = dr
	}

	price := priceRange[0] + rand.Float64()*(priceRange[1]-priceRange[0])
	duration := int(float64(durationRange[0]) + rand.Float64()*float64(durationRange[1]-durationRange[0]))

	return price, duration
}

func generateCoverImage(tourismType string) string {
	imageMap := map[string]string{
		"museum":     "https://images.unsplash.com/photo-1566127992631-137a642a90f4?w=800",
		"theme_park": "https://images.unsplash.com/photo-1544552866-d3ed42536cfd?w=800",
		"zoo":        "https://images.unsplash.com/photo-1564760055775-d63b17a55c44?w=800",
		"aquarium":   "https://images.unsplash.com/photo-1544551763-46a013bb70d5?w=800",
		"viewpoint":  "https://images.unsplash.com/photo-1477959858617-67f85cf4f1df?w=800",
		"gallery":    "https://images.unsplash.com/photo-1578662996442-48f60103fc96?w=800",
		"artwork":    "https://images.unsplash.com/photo-1578662996442-48f60103fc96?w=800",
		"attraction": "https://images.unsplash.com/photo-1539650116574-75c0c6d73f6e?w=800",
	}

	if img, exists := imageMap[tourismType]; exists {
		return img
	}
	return imageMap["attraction"]
}
