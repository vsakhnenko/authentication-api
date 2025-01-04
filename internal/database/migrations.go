package database

import (
	"authentication/internal/entities"
	"gorm.io/gorm"
	"log"
)

func AutoMigrate(db *gorm.DB) {
	err := db.AutoMigrate(&entities.User{})
	if err != nil {
		log.Fatalf("Auto migration failed: %v", err)
	}
}
