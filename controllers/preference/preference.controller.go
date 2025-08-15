package preference

import (
	"backend-go/config"
	"backend-go/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetAll(c *gin.Context) {
	var preferences []models.Preference
	if err := config.DB.Find(&preferences).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve preferences",
			"message": "Could not fetch preferences from database",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Preferences retrieved successfully",
		"data":    preferences,
		"count":   len(preferences),
	})
}

func GetByID(c *gin.Context) {
	var preference models.Preference
	id := c.Param("id")

	if err := config.DB.First(&preference, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Preference not found",
			"message": "The requested preference does not exist",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Preference retrieved successfully",
		"data":    preference,
	})
}

func Create(c *gin.Context) {
	var preference models.Preference
	if err := c.ShouldBindJSON(&preference); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid input",
			"message": "Please provide valid preference data",
		})
		return
	}

	if err := config.DB.Create(&preference).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create preference",
			"message": "Could not save preference to database",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Preference created successfully",
		"data":    preference,
	})
}

func Update(c *gin.Context) {
	var preference models.Preference
	id := c.Param("id")

	if err := config.DB.First(&preference, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Preference not found",
			"message": "The requested preference does not exist",
		})
		return
	}

	if err := c.ShouldBindJSON(&preference); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid input",
			"message": "Please provide valid preference data",
		})
		return
	}

	if err := config.DB.Save(&preference).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update preference",
			"message": "Could not update preference in database",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Preference updated successfully",
		"data":    preference,
	})
}

func Delete(c *gin.Context) {
	var preference models.Preference
	id := c.Param("id")

	if err := config.DB.First(&preference, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Preference not found",
			"message": "The requested preference does not exist",
		})
		return
	}

	if err := config.DB.Delete(&preference).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to delete preference",
			"message": "Could not delete preference from database",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Preference deleted successfully",
	})
}
