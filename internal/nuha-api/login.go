package nuha

import (
	"encoding/json"
	"net/http"

	"github.com/Modalessi/nuha-api/internal"
	"github.com/Modalessi/nuha-api/internal/repositories"
	"golang.org/x/crypto/bcrypt"
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

	ur := repositories.NewUserRespository(r.Context(), ns.DBQueries)
	user, err := ur.GetUserByEmail(loginData.Email)
	if err != nil {
		respondWithError(w, 500, SERVER_ERROR)
		return err
	}

	if user == nil {
		respondWithError(w, 404, EntityDoesNotExistError("User"))
		return nil
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginData.Password))
	if err != nil {
		respondWithError(w, 401, WRONG_CREDINTALS_ERROR)
		return err
	}

	// create jwt token
	token, err := internal.NewJWTTokenWithClaims(user.Name, user.Email, ns.JWTSecret)
	if err != nil {
		respondWithError(w, 500, SERVER_ERROR)
		return err
	}

	respone := struct {
		Token string `json:"token"`
		Email string `json:"email"`
		Name  string `json:"name"`
	}{
		Token: token,
		Email: user.Email,
		Name:  user.Name,
	}

	respondWithJson(w, 200, &internal.JsonWrapper{Data: respone})
	return nil
}
