package database

import "gorm.io/gorm"

func SeedDatabase(db *gorm.DB) error {
	return AdminFactory(db)
}