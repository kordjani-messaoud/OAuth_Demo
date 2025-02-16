package auth

import (
	"net/http"
	"net/url"
)

func Authorize(w http.ResponseWriter, r *http.Request) {

	if r.Method != "GET" {
		http.Error(w, "Invalide Method", http.StatusMethodNotAllowed)
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}

	if params := r.Form; ValidateAuthReq(params) {

	} else {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}
}

func ValidateAuthReq(params url.Values) bool {

	responseType := params["responseType"][0]
	redirectURI := params["redirectURI"][0]
	clientId := params["clientID"][0]
	clientSecret := params["clientSecret"][0]
	scopes = params["scopes"]
	state := params["state"][0]

}
