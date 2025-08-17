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