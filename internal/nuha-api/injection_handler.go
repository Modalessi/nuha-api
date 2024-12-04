package nuha

import (
	"log"
	"net/http"
	"time"
)

type NuhaHandler func(*NuhaServer, http.ResponseWriter, *http.Request) error

func withServer(ns *NuhaServer, handler NuhaHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		start := time.Now()
		log.Printf("Started %s %s", r.Method, r.URL.Path)

		err := handler(ns, w, r)
		duration := time.Since(start)

		if err != nil {
			log.Printf("Error handling %s %s: %v", r.Method, r.URL.Path, err)

			if !isHeaderWritten(w) {
				respondWithError(w, http.StatusInternalServerError, SERVER_ERROR)
				log.Printf("handler did not respond with a proper error")
			}

		}

		log.Printf("Completed %s %s in %v", r.Method, r.URL.Path, duration)

	}
}

func isHeaderWritten(w http.ResponseWriter) bool {
	rw, ok := w.(interface{ Written() bool })
	return ok && rw.Written()
}
