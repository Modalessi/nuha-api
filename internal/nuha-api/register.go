package nuha

import (
	"encoding/json"
	"net/http"
)

func register(ns *NuhaServer, w http.ResponseWriter, r *http.Request) error {
	type registerSchema struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
		Password  string `json:"password"`
	}
	defer r.Body.Close()

	registerData := registerSchema{}
	err := json.NewDecoder(r.Body).Decode(&registerData)
	if err != nil {
		respondWithError(w, 400, INVALID_JSON_ERROR)
		return err
	}

	userID, err := ns.Auth.Register(r.Context(), registerData.Email, registerData.Password)
	if err != nil || userID == nil {
		respondWithError(w, 400, INVALID_CREDINTALS_ERROR)
		return err
	}

	user, err := ns.UserRepo.StoreNewUserData(r.Context(), *userID, registerData.FirstName, registerData.LastName)
	if err != nil {
		respondWithError(w, 500, SERVER_ERROR)
		return err
	}

	respondWithJson(w, 201, user)
	return nil
}
