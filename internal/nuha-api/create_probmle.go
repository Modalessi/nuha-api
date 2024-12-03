package nuha

import "net/http"

func createProblem(ns *NuhaServer, w http.ResponseWriter, r *http.Request) error {
	respondWithText(w, 200, "you can do this")
	return nil
}
