package repositories

import (
	"gorm.io/gorm"
	"raya/models"
)

func GetUserByUsername(db *gorm.DB, username string) (*models.Admin, error) {
	var user models.Admin
	err := db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserByID(db *gorm.DB, userID uint) (*models.Admin, error) {
	var user models.Admin
	if err := db.First(&user, userID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func UpdatePassword(db *gorm.DB, userID uint, newPassword string) error {
	return db.Model(&models.Admin{}).Where("id = ?", userID).Update("password", newPassword).Error
}