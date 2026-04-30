package domain

import (
	"log"

	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	log.Println("Running auto-migration from entities...")

	err := db.AutoMigrate(
		&User{},
		&Token{},
		&Product{},
		// Add all your entities here
	)

	if err != nil {
		return err
	}

	log.Println("✅ Auto-migration completed")
	return nil
}
