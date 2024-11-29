package internal

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
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

func VerfiyToken(tokenString string, secretKey string) (*jwt.Token, error) {

	parseOptions := []jwt.ParserOption{
		jwt.WithValidMethods([]string{"HS256"}),
		jwt.WithExpirationRequired(),
	}

	token, err := jwt.Parse(tokenString, keyFunction(secretKey), parseOptions...)
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return token, nil
}

func keyFunction(secretKey string) jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Method.Alg())
		}

		return []byte(secretKey), nil
	}
}
