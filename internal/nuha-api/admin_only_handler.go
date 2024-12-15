package nuha

import "net/http"

func adminOnly(next http.HandlerFunc, adminEmail string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email := r.Context().Value(userEmailKey)
		if email == nil {
			respondWithError(w, 403, NOT_AUTHORIZED_ERROR)
			return
		}

		if emailString, ok := email.(string); ok {
			if emailString != adminEmail {
				respondWithError(w, 403, NOT_AUTHORIZED_ERROR)
				return
			}

			next(w, r)
		} else {
			respondWithError(w, 500, SERVER_ERROR)
		}
	}
}
