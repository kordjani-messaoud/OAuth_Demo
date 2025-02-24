package auth

import (
	"fmt"
	"net/http"
	"net/url"
)

func RequestApproval(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		http.Error(w, "Invalide Method", http.StatusMethodNotAllowed)
		return

	}
	// Validate querey params
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	params := r.Form

	if clientId, isValid := ValidAccreditation(params); isValid {
		SendCode(w, r, clientId)
	}

	// Validate Auth req
	fmt.Print(r.Form)
}

func ValidAccreditation(params url.Values) (ClientId, bool) {

	action := params["action"][0]
	clientId := ClientId(params["client_id"][0])
	state := params["state"][0]
	storedState := AuthzResCombinaisons[clientId].State

	if action == "allow" &&
		state == storedState {
		return clientId, true
	}
	return clientId, false
}

func SendCode(w http.ResponseWriter, r *http.Request, clientId ClientId) {

	code := AuthzResCombinaisons[clientId].Code
	state := AuthzResCombinaisons[clientId].State
	redirectURI := apps[clientId].RedirectURI

	url := fmt.Sprintf("%s?code=%s&state=%s", redirectURI, code, state)

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}
