package image

import (
	"backend-go/config"
	"backend-go/models"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Upload handles regular image uploads
func Upload(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Parse multipart form
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse form"})
		return
	}

	files := form.File["images"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No images provided"})
		return
	}

	// Get trip ID if provided
	var tripID *uint
	if tripIDStr := c.PostForm("trip_id"); tripIDStr != "" {
		if id, err := strconv.ParseUint(tripIDStr, 10, 32); err == nil {
			tripIDUint := uint(id)
			tripID = &tripIDUint
		}
	}

	var uploadedImages []models.Image
	uploadDir := "uploads/images"

	// Create upload directory if it doesn't exist
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create upload directory"})
		return
	}

	for _, file := range files {
		// Validate file type
		if !isValidImageType(file.Header.Get("Content-Type")) {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid file type for %s", file.Filename)})
			return
		}

		// Generate unique filename
		ext := filepath.Ext(file.Filename)
		fileName := fmt.Sprintf("image_%s_%s%s",
			time.Now().Format("20060102_150405"),
			uuid.New().String(),
			ext)

		filePath := filepath.Join(uploadDir, fileName)

		// Save file
		if err := c.SaveUploadedFile(file, filePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
			return
		}

		// Create image record
		userIDUint := userID.(uint)
		image := models.Image{
			TripID:       tripID,
			URL:          fmt.Sprintf("http://localhost:8080/api/v1/images/%s", fileName),
			FileName:     fileName,
			OriginalName: file.Filename,
			FileSize:     file.Size,
			MimeType:     file.Header.Get("Content-Type"),
			UploadedBy:   &userIDUint,
		}

		if err := config.DB.Create(&image).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image record"})
			return
		}

		uploadedImages = append(uploadedImages, image)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Images uploaded successfully",
		"images":  uploadedImages,
	})
}

// UploadCoverImage handles cover image uploads for trips
func UploadCoverImage(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	file, err := c.FormFile("cover_image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No cover image provided"})
		return
	}

	// Validate file type
	if !isValidImageType(file.Header.Get("Content-Type")) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type"})
		return
	}

	// Get trip ID if provided
	var tripID *uint
	if tripIDStr := c.PostForm("trip_id"); tripIDStr != "" {
		if id, err := strconv.ParseUint(tripIDStr, 10, 32); err == nil {
			tripIDUint := uint(id)
			tripID = &tripIDUint
		}
	}

	uploadDir := "uploads/covers"

	// Create upload directory if it doesn't exist
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create upload directory"})
		return
	}

	// Generate unique filename
	ext := filepath.Ext(file.Filename)
	fileName := fmt.Sprintf("cover_%s_%s%s",
		time.Now().Format("20060102_150405"),
		uuid.New().String(),
		ext)

	filePath := filepath.Join(uploadDir, fileName)

	// Save file
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	// Create image record
	userIDUint := userID.(uint)
	image := models.Image{
		TripID:       tripID,
		URL:          fmt.Sprintf("http://localhost:8080/api/v1/covers/%s", fileName),
		FileName:     fileName,
		OriginalName: file.Filename,
		FileSize:     file.Size,
		MimeType:     file.Header.Get("Content-Type"),
		UploadedBy:   &userIDUint,
	}

	if err := config.DB.Create(&image).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image record"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Cover image uploaded successfully",
		"image":   image,
	})
}

// GetImage serves regular images
func GetImage(c *gin.Context) {
	filename := c.Param("filename")
	if filename == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Filename is required"})
		return
	}

	filePath := filepath.Join("uploads/images", filename)

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Image not found"})
		return
	}

	// Open and serve the file
	file, err := os.Open(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open image"})
		return
	}
	defer file.Close()

	// Get file info for content type
	var image models.Image
	if err := config.DB.Where("file_name = ?", filename).First(&image).Error; err == nil {
		c.Header("Content-Type", image.MimeType)
	} else {
		// Fallback to detecting content type from file extension
		ext := strings.ToLower(filepath.Ext(filename))
		switch ext {
		case ".jpg", ".jpeg":
			c.Header("Content-Type", "image/jpeg")
		case ".png":
			c.Header("Content-Type", "image/png")
		case ".gif":
			c.Header("Content-Type", "image/gif")
		case ".webp":
			c.Header("Content-Type", "image/webp")
		default:
			c.Header("Content-Type", "application/octet-stream")
		}
	}

	// Copy file content to response
	_, err = io.Copy(c.Writer, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to serve image"})
		return
	}
}

// GetCoverImage serves cover images
func GetCoverImage(c *gin.Context) {
	filename := c.Param("filename")
	if filename == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Filename is required"})
		return
	}

	filePath := filepath.Join("uploads/covers", filename)

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cover image not found"})
		return
	}

	// Open and serve the file
	file, err := os.Open(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open cover image"})
		return
	}
	defer file.Close()

	// Get file info for content type
	var image models.Image
	if err := config.DB.Where("file_name = ?", filename).First(&image).Error; err == nil {
		c.Header("Content-Type", image.MimeType)
	} else {
		// Fallback to detecting content type from file extension
		ext := strings.ToLower(filepath.Ext(filename))
		switch ext {
		case ".jpg", ".jpeg":
			c.Header("Content-Type", "image/jpeg")
		case ".png":
			c.Header("Content-Type", "image/png")
		case ".gif":
			c.Header("Content-Type", "image/gif")
		case ".webp":
			c.Header("Content-Type", "image/webp")
		default:
			c.Header("Content-Type", "application/octet-stream")
		}
	}

	// Copy file content to response
	_, err = io.Copy(c.Writer, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to serve cover image"})
		return
	}
}

// GetMyImages returns images uploaded by the authenticated user
func GetMyImages(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var images []models.Image
	if err := config.DB.Where("uploaded_by = ?", userID).Find(&images).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch images"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"images": images,
	})
}

// GetImagesByTrip returns images for a specific trip
func GetImagesByTrip(c *gin.Context) {
	tripIDStr := c.Param("trip_id")
	tripID, err := strconv.ParseUint(tripIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid trip ID"})
		return
	}

	var images []models.Image
	if err := config.DB.Where("trip_id = ?", uint(tripID)).Find(&images).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch images"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"images": images,
	})
}

// DeleteImage deletes an image (only by the uploader)
func DeleteImage(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	imageIDStr := c.Param("id")
	imageID, err := strconv.ParseUint(imageIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid image ID"})
		return
	}

	// Find the image and verify ownership
	var image models.Image
	if err := config.DB.Where("id = ? AND uploaded_by = ?", uint(imageID), userID).First(&image).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Image not found or unauthorized"})
		return
	}

	// Determine file path based on URL
	var filePath string
	if strings.Contains(image.URL, "/covers/") {
		filePath = filepath.Join("uploads/covers", image.FileName)
	} else {
		filePath = filepath.Join("uploads/images", image.FileName)
	}

	// Delete the file from filesystem
	if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete file"})
		return
	}

	// Delete the database record
	if err := config.DB.Delete(&image).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete image record"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Image deleted successfully",
	})
}

// isValidImageType checks if the MIME type is a valid image type
func isValidImageType(mimeType string) bool {
	validTypes := []string{
		"image/jpeg",
		"image/jpg",
		"image/png",
		"image/gif",
		"image/webp",
	}

	for _, validType := range validTypes {
		if mimeType == validType {
			return true
		}
	}
	return false
}
