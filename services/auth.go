package services

import (
	"errors"
	"gorm.io/gorm"

	"raya/repositories"
	"raya/utils"
)

func LoginUser(db *gorm.DB, username, password string) (string, error) {
	user, err := repositories.GetUserByUsername(db, username)
	if err != nil {
		return "", errors.New("user tidak ditemukan")
	}

	if !utils.CheckPassword(password, user.Password) {
		return "", errors.New("password salah")
	}

	token, err := utils.GenerateJWT(*user)
	if err != nil {
		return "", err
	}

	return token, nil
}

func UpdatePassword(db *gorm.DB, userID uint, currentPassword, newPassword string) error {
	admin, err := repositories.GetUserByID(db, userID)
	if err != nil {
		return errors.New("admin tidak ditemukan")
	}

	if !utils.CheckPassword(currentPassword, admin.Password) {
		return errors.New("password saat ini salah")
	}

	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return errors.New("gagal membuat password baru")
	}

	return repositories.UpdatePassword(db, userID, hashedPassword)
}

func LogoutUser(token string) error {
	return nil
}
