package nuha

import "net/http"

func protected(ns *NuhaServer, w http.ResponseWriter, r *http.Request) error {

	msg := ""

	msg += "you are logined in and you can do this\n"

	msg += "your email is:" + r.Context().Value(userEmailKey).(string) + "\n"
	msg += "your name is:" + r.Context().Value(userNameKey).(string) + "\n"

	respondWithText(w, 200, msg)
	return nil
}
