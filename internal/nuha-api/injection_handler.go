package nuha

import (
	"log"
	"net/http"
)

type NuhaHandler func(*NuhaServer, http.ResponseWriter, *http.Request) error

func InjectNuhaServer(ns *NuhaServer, handler NuhaHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := handler(ns, w, r)
		if err != nil {
			log.Print(err.Error())
		}
	}
}
