package handlers

import "net/http"

func CheckHealth(w http.ResponseWriter, r *http.Request) {
	respondWithText(w, 200, "this is working good")
}
