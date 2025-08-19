package models

import "gorm.io/gorm"

type Image struct {
	gorm.Model

	TripID       *uint  `json:"trip_id"`
	URL          string `json:"url" gorm:"not null"`
	FileName     string `json:"file_name" gorm:"not null"`
	OriginalName string `json:"original_name"`
	FileSize     int64  `json:"file_size"`
	MimeType     string `json:"mime_type"`
	UploadedBy   *uint  `json:"uploaded_by"`
}
