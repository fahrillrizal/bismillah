package controllers

import (
	"net/http"
	"raya/services"
	"raya/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// LoginUser godoc
// @Summary Login user
// @Description Authenticate user and return a token
// @Tags auth
// @Accept json
// @Produce json
// @Param login body object true "Login Credentials" 
// @Success 200 {object} map[string]string "token: JWT Token"
// @Failure 400 {object} map[string]string "message: Username and password are required"
// @Failure 401 {object} map[string]string "message: Invalid credentials"
// @Router /api/login [post]
func LoginUser(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var input struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Username and password are required"})
		return
	}

	token, err := services.LoginUser(db, input.Username, input.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid credentials"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

// ChangePassword godoc
// @Summary Change user password
// @Description Change the password of an authenticated user
// @Tags auth
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param password body object true "Password Data"
// @Success 200 {object} map[string]string "message: Password berhasil diubah"
// @Failure 400 {object} map[string]string "message: Input tidak valid"
// @Failure 401 {object} map[string]string "message: Unauthorized"
// @Router /api/change-password [post]
func ChangePassword(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	user := c.MustGet("user").(*models.Admin)

	var input struct {
		CurrentPassword string `json:"current_password" binding:"required"`
		NewPassword     string `json:"new_password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Input tidak valid"})
		return
	}

	if err := services.UpdatePassword(db, user.ID, input.CurrentPassword, input.NewPassword); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password berhasil diubah"})
}

// LogoutUser godoc
// @Summary Logout user
// @Description Logout a user from the system
// @Tags auth
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]string "message: Logout successful"
// @Router /api/logout [post]
func LogoutUser(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}