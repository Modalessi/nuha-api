package nuha

import (
	"encoding/json"
	"net/http"
	"regexp"

	"github.com/Modalessi/nuha-api/internal"
	"github.com/Modalessi/nuha-api/internal/factories"
)

func register(ns *NuhaServer, w http.ResponseWriter, r *http.Request) error {
	type registerSchema struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	defer r.Body.Close()

	registerData := registerSchema{}
	err := json.NewDecoder(r.Body).Decode(&registerData)
	if err != nil {
		respondWithError(w, 400, INVALID_JSON_ERROR)
		return err
	}

	if !isValidCredentials(registerData.Email, registerData.Password) {
		respondWithError(w, 400, INVALID_CREDINTALS_ERROR)
		return err
	}

	uf := factories.NewUserFactory(r.Context(), ns.DB)
	// check if user already exist

	exist, err := uf.DoesUserExistWithEmail(registerData.Email)
	if err != nil {
		respondWithError(w, 500, SERVER_ERROR)
		return err
	}

	if exist {
		respondWithError(w, 400, USER_ALREADY_EXIST_ERROR)
		return err
	}

	// add user
	user, err := uf.CreateNewUser(registerData.Name, registerData.Email, registerData.Password)
	if err != nil {
		respondWithError(w, 500, err)
	}

	// return the data

	response := struct {
		UserName  string `json:"user_name"`
		UserEmail string `json:"user_email"`
		CreatedAt string `json:"created_at"`
	}{
		UserName:  user.Name,
		UserEmail: user.Email,
		CreatedAt: user.CreatedAt.String(),
	}

	respondWithJson(w, 201, &internal.JsonWrapper{Data: response})
	return nil
}

func isValidCredentials(email string, pw string) bool {
	// check the email is valid email for registration
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	isValidEmail, err := regexp.MatchString(emailRegex, email)
	if err != nil {
		return false
	}
	if !isValidEmail {
		return false
	}

	// check password with regex (8 chars)
	validPasswordRegex := `^[A-Za-z0-9]{8,}$`

	isValidPassword, err := regexp.MatchString(validPasswordRegex, pw)
	if err != nil {
		return false
	}

	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(pw)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(pw)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(pw)

	return isValidPassword && hasUpper && hasLower && hasNumber
}
