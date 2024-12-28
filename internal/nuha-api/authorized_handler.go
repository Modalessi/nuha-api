package nuha

import (
	"context"
	"net/http"
	"strings"

	"github.com/Modalessi/nuha-api/internal/auth"
)

func authorized(next http.HandlerFunc, auth *auth.AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			respondWithError(w, 401, AUTHORIZATION_HEADER_ERROR)
			return
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

		userEmail, err := auth.ValidateToken(r.Context(), tokenString)
		if err != nil {
			respondWithError(w, 401, INVALID_TOKEN_ERROR)
			return
		}

		ctx := context.WithValue(r.Context(), userEmailKey, userEmail)

		next(w, r.WithContext(ctx))
	}
}
