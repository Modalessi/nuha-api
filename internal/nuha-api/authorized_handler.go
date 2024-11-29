package nuha

import (
	"context"
	"net/http"
	"strings"

	"github.com/Modalessi/nuha-api/internal"
	"github.com/golang-jwt/jwt/v5"
)

func Authorized(next http.HandlerFunc, jwtSecret string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			respondWithError(w, 401, AUTHORIZATION_HEADER_ERROR)
			return
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

		token, err := internal.VerfiyToken(tokenString, jwtSecret)
		if err != nil {
			respondWithError(w, 401, INVALID_TOKEN_ERROR)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			respondWithError(w, 401, INVALID_TOKEN_ERROR)
			return
		}

		ctx := context.WithValue(r.Context(), userEmailKey, claims["sub"].(string))
		ctx = context.WithValue(ctx, userNameKey, claims["name"].(string))

		next(w, r.WithContext(ctx))
	}
}
