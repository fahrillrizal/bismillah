package controllers

import (
	"net/http"
	"raya/models"
	"raya/services"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetCategoryByID godoc
// @Summary Get a category by ID
// @Description Get details of a specific category by its ID
// @Tags categories
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "Category ID"
// @Success 200 {object} models.Category
// @Failure 400 {object} map[string]string "message: Invalid ID format"
// @Failure 404 {object} map[string]string "message: Category not found"
// @Router /api/categories/{id} [get]
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

// GetAllCategories godoc
// @Summary Get all categories
// @Description Get a list of all categories
// @Tags categories
// @Produce json
// @Param includeEmpty query bool false "Include categories with no links"
// @Success 200 {array} models.Category
// @Failure 500 {object} map[string]string "message: Error fetching categories"
// @Router /api/categories [get]
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

// GetAllLinks godoc
// @Summary Get all links
// @Description Get a list of all links
// @Tags links
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {array} models.Link
// @Failure 500 {object} map[string]string "message: Error fetching links"
// @Router /api/links/all [get]
func GetAllLinks(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	links, err := services.GetAllLinks(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error fetching links"})
		return
	}

	c.JSON(http.StatusOK, links)
}

// GetLinkByID godoc
// @Summary Get a link by ID
// @Description Get details of a specific link by its ID
// @Tags links
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "Link ID"
// @Success 200 {object} models.Link
// @Failure 400 {object} map[string]string "message: Invalid ID format"
// @Failure 404 {object} map[string]string "message: Link not found"
// @Router /api/links/{id} [get]
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

// GetLinksByCategory godoc
// @Summary Get links by category
// @Description Get all links belonging to a specific category
// @Tags links
// @Produce json
// @Security ApiKeyAuth
// @Param category_id path int true "Category ID"
// @Success 200 {array} models.Link
// @Failure 400 {object} map[string]string "message: ID kategori tidak valid"
// @Failure 404 {object} map[string]string "message: Kategori tidak ditemukan"
// @Failure 500 {object} map[string]string "message: Error mengambil link"
// @Router /api/categories/{category_id}/links [get]
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

// CreateLink godoc
// @Summary Create a new link
// @Description Create a new link in a specific category
// @Tags links
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param category_id path int true "Category ID"
// @Param link body models.Link true "Link Data"
// @Success 201 {object} models.Link
// @Failure 400 {object} map[string]string "message: Format input tidak valid"
// @Failure 404 {object} map[string]string "message: Kategori tidak ditemukan"
// @Failure 500 {object} map[string]string "message: Error menentukan urutan"
// @Router /api/categories/{category_id}/links [post]
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

// UpdateLink godoc
// @Summary Update a link
// @Description Update a link in a specific category
// @Tags links
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param category_id path int true "Category ID"
// @Param link_id path int true "Link ID"
// @Param link body models.Link true "Link Data"
// @Success 200 {object} models.Link
// @Failure 400 {object} map[string]string "message: Format input tidak valid"
// @Failure 404 {object} map[string]string "message: Link tidak ditemukan dalam kategori ini"
// @Router /api/categories/{category_id}/links/{link_id} [patch]
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

// DeleteLink godoc
// @Summary Delete a link
// @Description Delete a link from a specific category
// @Tags links
// @Produce json
// @Security ApiKeyAuth
// @Param category_id path int true "Category ID"
// @Param link_id path int true "Link ID"
// @Success 200 {object} map[string]string "message: Link berhasil dihapus"
// @Failure 400 {object} map[string]string "message: ID link tidak valid"
// @Failure 404 {object} map[string]string "message: Link tidak ditemukan dalam kategori ini"
// @Failure 500 {object} map[string]string "message: Error menghapus link"
// @Router /api/categories/{category_id}/links/{link_id} [delete]
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

// CreateCategory godoc
// @Summary Create a new category
// @Description Create a new category with the input data
// @Tags categories
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param category body models.Category true "Category Data"
// @Success 201 {object} models.Category
// @Failure 400 {object} map[string]string "message: Invalid input format"
// @Failure 500 {object} map[string]string "message: Error menentukan urutan"
// @Router /api/categories [post]
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

// UpdateCategory godoc
// @Summary Update a category
// @Description Update a category with the input data
// @Tags categories
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "Category ID"
// @Param category body models.Category true "Category Data"
// @Success 200 {object} models.Category
// @Failure 400 {object} map[string]string "message: Invalid input format"
// @Failure 404 {object} map[string]string "message: Category not found"
// @Router /api/categories/{id} [patch]
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

// DeleteCategory godoc
// @Summary Delete a category
// @Description Delete a category by its ID
// @Tags categories
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "Category ID"
// @Success 200 {object} map[string]string "message: Category deleted successfully"
// @Failure 400 {object} map[string]string "message: Invalid ID format"
// @Failure 404 {object} map[string]string "message: Category not found"
// @Failure 500 {object} map[string]string "message: Error deleting category"
// @Router /api/category/{id} [delete]
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