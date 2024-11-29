package nuha

import (
	"net/http"

	"github.com/Modalessi/nuha-api/internal"
)

func respondWithError(w http.ResponseWriter, code int, err error) {
	w.Header().Add("Content-Type", "application/json")
	data := err.Error()

	w.WriteHeader(code)
	w.Write([]byte(data))
}

func respondWithJson(w http.ResponseWriter, code int, payload internal.Jsonable) {
	w.Header().Add("Content-Type", "application/json")
	data := payload.JSON()

	w.WriteHeader(code)
	w.Write(data)
}

func respondWithText(w http.ResponseWriter, code int, text string) {
	w.Header().Add("Content-Type", "text")

	w.WriteHeader(code)
	w.Write([]byte(text))
}
