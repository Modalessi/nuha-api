package nuha

import "net/http"

func checkHealth(w http.ResponseWriter, r *http.Request) {
	respondWithSuccess(w, 200, "this is working good")
}
