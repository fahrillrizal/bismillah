package middleware

import (
	"net/http"
	"raya/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization header required"})
			c.Abort()
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid authorization format"})
			c.Abort()
			return
		}

		tokenString := tokenParts[1]
		_, user, err := utils.ParseJWT(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
			c.Abort()
			return
		}

		c.Set("user", user)
		c.Next()
	}
}

func DetectMobileMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userAgent := c.GetHeader("User-Agent")

		if !strings.Contains(userAgent, "Android") && !strings.Contains(userAgent, "iPhone") && !strings.Contains(userAgent, "iPad") {
			c.JSON(http.StatusForbidden, gin.H{
				"message": "Akses hanya diperbolehkan dari perangkat mobile (Android/iOS)",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
