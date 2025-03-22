package repositories

import (
	"gorm.io/gorm"
	"raya/models"
	"time"
)

func GetNextCategoryOrder(db *gorm.DB) (int, error) {
	var maxOrder struct {
		MaxOrder int
	}
	
	err := db.Model(&models.Category{}).
		Select("COALESCE(MAX(\"order\"), 0) as max_order").
		Scan(&maxOrder).Error
	
	if err != nil {
		return 0, err
	}
	
	return maxOrder.MaxOrder + 1, nil
}

func GetNextLinkOrder(db *gorm.DB, categoryID uint) (int, error) {
	var maxOrder struct {
		MaxOrder int
	}
	
	err := db.Model(&models.Link{}).
		Where("category_id = ?", categoryID).
		Select("COALESCE(MAX(\"order\"), 0) as max_order").
		Scan(&maxOrder).Error
	
	if err != nil {
		return 0, err
	}
	
	return maxOrder.MaxOrder + 1, nil
}

func GetAllCategories(db *gorm.DB, includeEmpty bool) ([]models.Category, error) {
	var categories []models.Category
	

	query := db.Preload("Links", "is_active = ?", true).
		Order("categories.order asc").
		Order("links.order asc")
	
	if err := query.Find(&categories).Error; err != nil {
		return nil, err
	}
	

	if !includeEmpty {
		var filteredCategories []models.Category
		for _, category := range categories {
			if len(category.Links) > 0 {
				filteredCategories = append(filteredCategories, category)
			}
		}
		return filteredCategories, nil
	}
	
	return categories, nil
}

func GetAllLinks(db *gorm.DB) ([]models.Link, error) {
	var links []models.Link
	err := db.Preload("Category").Order("category_id, order").Find(&links).Error
	return links, err
}

func GetLinkByID(db *gorm.DB, id uint) (*models.Link, error) {
	var link models.Link
	if err := db.Preload("Category").First(&link, id).Error; err != nil {
		return nil, err
	}
	return &link, nil
}

func CreateLink(db *gorm.DB, link *models.Link) error {
	if link.Order <= 0 {
		nextOrder, err := GetNextLinkOrder(db, link.CategoryID)
		if err != nil {
			return err
		}
		link.Order = nextOrder
	}
	
	return db.Create(link).Error
}

func UpdateLink(db *gorm.DB, id uint, updatedLink *models.Link) error {
	link, err := GetLinkByID(db, id)
	if err != nil {
		return err
	}

	link.Title = updatedLink.Title
	link.URL = updatedLink.URL
	link.CategoryID = updatedLink.CategoryID
	link.IsActive = updatedLink.IsActive
	link.Order = updatedLink.Order
	link.UpdatedAt = time.Now()

	return db.Save(link).Error
}

func DeleteLink(db *gorm.DB, id uint) error {
	return db.Delete(&models.Link{}, id).Error
}

func GetCategoryByID(db *gorm.DB, id uint) (*models.Category, error) {
	var category models.Category
	if err := db.Preload("Links", "is_active = ?", true).First(&category, id).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

func CreateCategory(db *gorm.DB, category *models.Category) error {
	if category.Order <= 0 {
		nextOrder, err := GetNextCategoryOrder(db)
		if err != nil {
			return err
		}
		category.Order = nextOrder
	}
	
	return db.Create(category).Error
}

func UpdateCategory(db *gorm.DB, id uint, updatedCategory *models.Category) error {
	category, err := GetCategoryByID(db, id)
	if err != nil {
		return err
	}

	category.Name = updatedCategory.Name
	category.Order = updatedCategory.Order

	return db.Save(category).Error
}

func DeleteCategory(db *gorm.DB, id uint) error {
	return db.Delete(&models.Category{}, id).Error
}