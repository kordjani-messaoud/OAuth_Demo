package auth

import (
	"fmt"
	"net/http"
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

	if params := r.form; ValidAccreditation(params)
	// Validate Auth req
	fmt.Print(r.Form)
}

func ValidAccreditation(params url.Values) (ClientId, bool){

	action := params["action"][0]
	clientId := ClientId(params["client_id"][0])
	state := params["state"][0]
	storedState := AuthzResCombinaisons[clientId].State

	if acction == "allow" &&
		state == storedState {
			return clientId, true
		}
	return clientId, false
}
