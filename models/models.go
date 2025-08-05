package models

import (
	"backend-go/config"
	"log"
)

// AutoMigrate runs auto-migration for all models
func AutoMigrate() error {
	err := config.DB.AutoMigrate(
		&User{},
		&Trip{},
		&Image{},
	)
	if err != nil {
		log.Printf("Failed to migrate database: %v", err)
		return err
	}
	log.Println("Database migration completed successfully!")
	return nil
}
