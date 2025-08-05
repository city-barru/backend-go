package models

import "gorm.io/gorm"

type Image struct {
	gorm.Model

	TripID *uint `json:"trip_id"`
}
