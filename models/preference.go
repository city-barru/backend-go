package models

import "gorm.io/gorm"

type Preference struct {
	gorm.Model
	Name string `json:"name" gorm:"not null;unique"`
}

type UserPreference struct {
	gorm.Model
	UserID       uint `json:"user_id" gorm:"not null"`       // Foreign key to User
	PreferenceID uint `json:"preference_id" gorm:"not null"` // Foreign key to Preference
}

type TripPreference struct {
	gorm.Model
	TripID       uint `json:"trip_id" gorm:"not null"`       // Foreign key to Trip
	PreferenceID uint `json:"preference_id" gorm:"not null"` // Foreign key to Preference
}
