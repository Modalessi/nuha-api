package nuha

import (
	"encoding/json"
	"net/http"

	"github.com/Modalessi/nuha-api/internal"
)

func login(ns *NuhaServer, w http.ResponseWriter, r *http.Request) error {
	type loginSchema struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	defer r.Body.Close()

	loginData := loginSchema{}
	err := json.NewDecoder(r.Body).Decode(&loginData)
	if err != nil {
		respondWithError(w, 400, INVALID_JSON_ERROR)
		return err
	}

	user, err := ns.UserRepo.GetUserByEmail(r.Context(), loginData.Email)
	if err != nil {
		respondWithError(w, 500, SERVER_ERROR)
		return err
	}

	if user == nil {
		respondWithError(w, 404, EntityDoesNotExistError("User"))
		return nil
	}

	// what you wann do is use the auth service to login
	// return if error if it was not successful
	// return token if it was successful

	token, err := ns.Auth.Login(r.Context(), loginData.Email, loginData.Password)
	if err != nil {
		respondWithError(w, 400, err)
		return err
	}

	response := struct {
		Token string `json:"token"`
	}{
		Token: token,
	}

	respondWithJson(w, 200, &internal.JsonWrapper{Data: response})
	return nil
}
