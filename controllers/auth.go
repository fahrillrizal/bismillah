package controllers

import (
	"net/http"
	"raya/services"
	"raya/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

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

func LogoutUser(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}