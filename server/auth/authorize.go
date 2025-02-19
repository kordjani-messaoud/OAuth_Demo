package auth

import (
	"crypto/rand"
	"fmt"
	"net/http"
	"net/url"
	"slices"
	"strings"
)

// I need to think more about the scopes
var scopes []string
var isScopesUpdated bool = false

func Authorize(w http.ResponseWriter, r *http.Request) {
	// Validate method
	if r.Method != "GET" {
		http.Error(w, "Invalide Method", http.StatusMethodNotAllowed)
		return

	}
	// Validate querey params
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	// Validate Auth req
	if params := r.Form; ValidAuthzReq(params) {
		clientId := ClientId(params["client_id"][0])
		state := params["state"][0]
		AutzResParams := NewAuthzResCombinaison(state)
		AuthzResCombinaisons[clientId] = *AutzResParams
		// Ask Resource Owner permission
		RequestUserPermission(w, r, params)

	} else {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}
}

// Validate Authorize requests
// Returning bool
// Get sent scopes and drop every one that is not on the registered scope list
// Check if scopes had been update
func ValidAuthzReq(params url.Values) bool {
	// Get params
	responseType := params["response_type"][0]
	redirectURI := params["redirect_uri"][0]
	clientId := ClientId(params["client_id"][0])
	clientSecret := params["client_secret"][0]
	scopes = params["scope"]

	// Validate params
	if app, present := apps[clientId]; present &&
		redirectURI == app.RedirectURI &&
		responseType == "code" &&
		clientSecret == app.ClientSecret {
		for _, scope := range scopes {
			if !slices.Contains(app.Scopes, scope) {
				scopes = slices.DeleteFunc(scopes, func(v string) bool {
					return v == scope
				})
				isScopesUpdated = true
			}
		}
		if len(scopes) != 0 {
			return true
		}
	}
	return false
}

func NewAuthzResCombinaison(state string) *AuthzResCombinason {
	return &AuthzResCombinason{
		rand.Text(),
		state,
		rand.Text(),
	}
}

// Q:Why do we parse the req form again couldn't we just ust the data stored. Res: Because the data is not stored
func RequestUserPermission(w http.ResponseWriter, r *http.Request, params url.Values) {
	clientId := ClientId(params["client_id"][0])
	redirectURI := params["redirect_uri"][0]
	state := params["state"][0]

	url := fmt.Sprintf("%s/login?client_id=%s&redirect_uri=%s&scope=%s&state=%s",
		SERVERURL,
		clientId,
		redirectURI,
		strings.Join(scopes, " "),
		state,
	)

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}
