package models

import "gorm.io/gorm"

type Trip struct {
	gorm.Model
	Name        string  `json:"name" gorm:"not null"`
	Description string  `json:"description"`
	CoverImage  string  `json:"cover_image"`
	Price       float64 `json:"price" gorm:"not null"`
	Duration    int     `json:"duration" gorm:"not null"`

	StartLatitude  float64 `json:"start_latitude" gorm:"not null"`
	StartLongitude float64 `json:"start_longitude" gorm:"not null"`
	EndLatitude    float64 `json:"end_latitude" gorm:"not null"`
	EndLongitude   float64 `json:"end_longitude" gorm:"not null"`

	// User association
	UserID uint `json:"user_id" gorm:"not null"`
	User   User `json:"user" gorm:"foreignKey:UserID"`

	Images      []Image      `json:"images,omitempty" gorm:"foreignKey:TripID"`
	Preferences []Preference `json:"preferences,omitempty" gorm:"many2many:trip_preferences;"`
	Points      []TripPoint  `json:"points,omitempty" gorm:"foreignKey:TripID"`
}

type TripPoint struct {
	gorm.Model
	TripID    uint    `json:"trip_id" gorm:"not null"`
	Latitude  float64 `json:"latitude" gorm:"not null"`
	Longitude float64 `json:"longitude" gorm:"not null"`
}
