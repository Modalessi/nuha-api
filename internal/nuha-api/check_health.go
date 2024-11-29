package nuha

import "net/http"

func checkHealth(w http.ResponseWriter, r *http.Request) {
	respondWithText(w, 200, "this is working good")
}
