package preference

import (
	"backend-go/config"
	"backend-go/models"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
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

func AssignPreference(c *gin.Context) {
	userIdStr, exists := c.Get("userID")
	if !exists {
		fmt.Println("User ID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Authentication required",
			"message": "User not authenticated",
		})
		return
	}
	userId, ok := userIdStr.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to assign preference",
			"message": "Could not save preference to database",
		})
		return
	}

	var preferences []models.Preference
	if err := c.ShouldBindJSON(&preferences); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid input",
			"message": "Please provide valid preference data",
		})
		return
	}
	// fmt.Println(preferences)

	for _, preference := range preferences {
		if err := config.DB.Model(&models.UserPreference{}).Where("user_id = ? AND preference_id = ?", userId, preference.ID).Error; err != nil {
			c.JSON(http.StatusConflict, gin.H{
				"error":   "Preference already assigned",
				"message": "This preference is already assigned to the user",
			})
			return
		}

		if err := config.DB.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "name"}},
			DoNothing: true,
		}).Create(&preference).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to assign preference",
				"message": "Could not save preference to database",
			})
			return
		}

		if err := config.DB.Create(&models.UserPreference{
			UserID:       userId,
			PreferenceID: preference.ID,
		}).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to assign preference",
				"message": "Could not save preference assignment to database",
			})
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Preference assigned successfully",
		"data":    preferences,
	})
}
