package nuha

import "net/http"

func protected(ns *NuhaServer, w http.ResponseWriter, r *http.Request) error {

	msg := ""

	msg += "you are logined in and you can do this\n"

	msg += "your email is:" + r.Context().Value(USER_EMAIL_CONTEXT_KEY).(string) + "\n"
	msg += "your token is:" + r.Context().Value(USER_TOKEN_CONTEXT_KEY).(string) + "\n"

	respondWithText(w, 200, msg)
	return nil
}
