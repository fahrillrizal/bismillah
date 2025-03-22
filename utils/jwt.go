package utils

import (
	"errors"
	"os"
	"time"

	"raya/models"

	"github.com/golang-jwt/jwt/v5"
)

var SecretKey = []byte(os.Getenv("JWT_SECRET_KEY"))

func GenerateJWT(user models.Admin) (string, error) {
	claims := jwt.MapClaims{
		"id":       user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(SecretKey)
}

func ParseJWT(tokenString string) (*jwt.Token, *models.Admin, error) {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return SecretKey, nil
	})

	if err != nil || !token.Valid {
		return nil, nil, errors.New("token tidak valid")
	}

	idFloat, ok := claims["id"].(float64)
	if !ok {
		return nil, nil, errors.New("ID tidak valid dalam token")
	}

	username, _ := claims["username"].(string)

	user := &models.Admin{
		ID:       uint(idFloat),
		Username: username,
	}

	return token, user, nil
}