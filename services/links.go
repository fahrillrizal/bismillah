package services

import (
	"errors"
	"gorm.io/gorm"
	"raya/models"
	"raya/repositories"
)

func GetAllCategories(db *gorm.DB, includeEmpty bool) ([]models.Category, error) {
	return repositories.GetAllCategories(db, includeEmpty)
}

func GetAllLinks(db *gorm.DB) ([]models.Link, error) {
	return repositories.GetAllLinks(db)
}

func GetLinkByID(db *gorm.DB, id uint) (*models.Link, error) {
	return repositories.GetLinkByID(db, id)
}

func CreateLink(db *gorm.DB, link *models.Link) error {
	if link.Title == "" || link.URL == "" || link.CategoryID == 0 {
		return errors.New("title, URL, and category are required")
	}
	
	return repositories.CreateLink(db, link)
}

func UpdateLink(db *gorm.DB, id uint, link *models.Link) error {
	if link.Title == "" || link.URL == "" || link.CategoryID == 0 {
		return errors.New("title, URL, and category are required")
	}
	return repositories.UpdateLink(db, id, link)
}

func DeleteLink(db *gorm.DB, id uint) error {
	return repositories.DeleteLink(db, id)
}

func GetCategoryByID(db *gorm.DB, id uint) (*models.Category, error) {
	return repositories.GetCategoryByID(db, id)
}

func CreateCategory(db *gorm.DB, category *models.Category) error {
	if category.Name == "" {
		return errors.New("category name is required")
	}
	
	return repositories.CreateCategory(db, category)
}

func UpdateCategory(db *gorm.DB, id uint, category *models.Category) error {
	if category.Name == "" {
		return errors.New("category name is required")
	}
	return repositories.UpdateCategory(db, id, category)
}

func DeleteCategory(db *gorm.DB, id uint) error {
	return repositories.DeleteCategory(db, id)
}