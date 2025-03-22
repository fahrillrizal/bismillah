package controllers

import (
	"net/http"
	"raya/models"
	"raya/services"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetCategoryByID(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid ID format"})
		return
	}

	category, err := services.GetCategoryByID(db, uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Category not found"})
		return
	}

	c.JSON(http.StatusOK, category)
}

func GetAllCategories(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	includeEmpty := c.Query("includeEmpty") == "true"

	categories, err := services.GetAllCategories(db, includeEmpty)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error fetching categories"})
		return
	}

	c.JSON(http.StatusOK, categories)
}

func GetAllLinks(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	links, err := services.GetAllLinks(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error fetching links"})
		return
	}

	c.JSON(http.StatusOK, links)
}

func GetLinkByID(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid ID format"})
		return
	}

	link, err := services.GetLinkByID(db, uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Link not found"})
		return
	}

	c.JSON(http.StatusOK, link)
}

func GetLinksByCategory(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	
	categoryIDStr := c.Param("category_id")
	categoryID, err := strconv.ParseUint(categoryIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "ID kategori tidak valid"})
		return
	}
	

	var category models.Category
	if err := db.First(&category, categoryID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Kategori tidak ditemukan"})
		return
	}
	

	var links []models.Link
	if err := db.Where("category_id = ? AND is_active = ?", categoryID, true).
		Order("\"order\" asc").
		Find(&links).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error mengambil link"})
		return
	}
	
	c.JSON(http.StatusOK, links)
}

func CreateLink(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
		
	categoryIDStr := c.Param("category_id")
	categoryID, err := strconv.ParseUint(categoryIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "ID kategori tidak valid"})
		return
	}
		
	var category models.Category
	if err := db.First(&category, categoryID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Kategori tidak ditemukan"})
		return
	}
		
	var link models.Link
	if err := c.ShouldBindJSON(&link); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Format input tidak valid"})
		return
	}
		
	link.CategoryID = uint(categoryID)
	

	if link.Order <= 0 {
	
		var maxOrder struct {
			MaxOrder int
		}
		if err := db.Model(&models.Link{}).
			Where("category_id = ?", categoryID).
			Select("COALESCE(MAX(\"order\"), 0) as max_order").
			Scan(&maxOrder).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error menentukan urutan"})
			return
		}
		link.Order = maxOrder.MaxOrder + 1
	}
		
	if err := services.CreateLink(db, &link); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	
	c.JSON(http.StatusCreated, link)
}

func UpdateLink(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
		
	categoryIDStr := c.Param("category_id")
	categoryID, err := strconv.ParseUint(categoryIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "ID kategori tidak valid"})
		return
	}
	
	linkIDStr := c.Param("link_id")
	linkID, err := strconv.ParseUint(linkIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "ID link tidak valid"})
		return
	}
		
	var link models.Link
	if err := db.Where("id = ? AND category_id = ?", linkID, categoryID).First(&link).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Link tidak ditemukan dalam kategori ini"})
		return
	}
		
	var updatedLink models.Link
	if err := c.ShouldBindJSON(&updatedLink); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Format input tidak valid"})
		return
	}
		
	updatedLink.CategoryID = uint(categoryID)
		
	if err := services.UpdateLink(db, uint(linkID), &updatedLink); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
		
	db.First(&link, linkID)
	c.JSON(http.StatusOK, link)
}

func DeleteLink(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
		
	categoryIDStr := c.Param("category_id")
	categoryID, err := strconv.ParseUint(categoryIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "ID kategori tidak valid"})
		return
	}
	
	linkIDStr := c.Param("link_id")
	linkID, err := strconv.ParseUint(linkIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "ID link tidak valid"})
		return
	}
		
	var link models.Link
	if err := db.Where("id = ? AND category_id = ?", linkID, categoryID).First(&link).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Link tidak ditemukan dalam kategori ini"})
		return
	}
		
	if err := db.Delete(&link).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error menghapus link"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Link berhasil dihapus"})
}

func CreateCategory(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var category models.Category
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input format"})
		return
	}


	if category.Order <= 0 {
		var maxOrder struct {
			MaxOrder int
		}
		if err := db.Model(&models.Category{}).
			Select("COALESCE(MAX(\"order\"), 0) as max_order").
			Scan(&maxOrder).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error menentukan urutan"})
			return
		}
		category.Order = maxOrder.MaxOrder + 1
	}

	if err := services.CreateCategory(db, &category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, category)
}

func UpdateCategory(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid ID format"})
		return
	}

	var category models.Category
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input format"})
		return
	}

	if err := services.UpdateCategory(db, uint(id), &category); err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"message": "Category not found"})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, category)
}

func DeleteCategory(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid ID format"})
		return
	}

	if err := services.DeleteCategory(db, uint(id)); err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"message": "Category not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error deleting category"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Category deleted successfully"})
}