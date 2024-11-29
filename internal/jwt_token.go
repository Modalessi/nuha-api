package internal

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func NewJWTTokenWithClaims(name string, email string, secretKey string) (string, error) {
	claims := jwt.MapClaims{
		"sub":  email,
		"name": name,
		"exp":  time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secretKey))
}
