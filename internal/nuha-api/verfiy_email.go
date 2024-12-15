package nuha

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

func verifyEmail(ns *NuhaServer, w http.ResponseWriter, r *http.Request) error {
	token, err := getVerficationTokenFromUrl(r.URL)
	if err != nil {
		respondWithError(w, 400, err)
		return err
	}

	err = ns.Auth.VerifyUser(r.Context(), token)
	if err != nil {
		respondWithError(w, 400, err)
		return err
	}

	respondWithSuccess(w, 200, "user is verified")
	return nil
}

func getVerficationTokenFromUrl(url *url.URL) (string, error) {
	// ex: www.nuha.com/verify/{token}
	urlPath := url.Path
	parts := strings.Split(urlPath, "/")
	if len(parts) < 3 {
		return "", fmt.Errorf("invalid url format")
	}

	token := parts[len(parts)-1]
	if token == "" {
		return "", fmt.Errorf("token is missing")
	}

	return token, nil
}
