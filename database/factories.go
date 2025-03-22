package database

import (
	"fmt"
	"raya/models"
	"raya/utils"

	"gorm.io/gorm"
)

// AdminFactory creates a default admin user if none exists
func AdminFactory(db *gorm.DB) error {
	var count int64
	db.Model(&models.Admin{}).Count(&count)

	if count > 0 {
		fmt.Println("Admin users already exist, skipping factory")
		return nil
	}

	// Create default admin
	defaultPassword, err := utils.HashPassword("admin123")
	if err != nil {
		return err
	}

	admin := models.Admin{
		Username: "admin",
		Password: defaultPassword,
	}

	if err := db.Create(&admin).Error; err != nil {
		return err
	}

	fmt.Println("Default admin user created")
	return nil
}