package nuha

import (
	"fmt"
	"net/http"
)

func logout(ns *NuhaServer, w http.ResponseWriter, r *http.Request) error {
	token, ok := r.Context().Value(USER_TOKEN_CONTEXT_KEY).(string)
	defer r.Body.Close()

	if !ok {
		respondWithError(w, 500, SERVER_ERROR)
		return fmt.Errorf("error: casting USER_TOKEN_CONTEXT_KEY value to string")
	}

	err := ns.Auth.Logout(r.Context(), token)
	if err != nil {
		respondWithError(w, 400, fmt.Errorf("you already logged out"))
		return err
	}

	respondWithSuccess(w, 200, "user logged out successfuly")
	return nil
}
