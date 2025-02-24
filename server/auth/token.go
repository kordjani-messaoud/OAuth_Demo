package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func GetAccessToken(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		http.Error(w, "Invalide Method", http.StatusMethodNotAllowed)
		return
	}
	// Validate querey params
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if clientId, isValide := ValideAccessTokenReq(r); isValide {
		GrantAccessToken(w, clientId)
		return
	}

	http.Error(w, "Unauthorized", http.StatusUnauthorized)
}

func ValideAccessTokenReq(r *http.Request) (ClientId, bool) {
	params := r.Form
	grantType := params["grant_type"][0]
	fmt.Println("Grant type:", grantType)
	code := params["code"][0]
	fmt.Println("code:", code)
	redirectURI := params["redirect_uri"][0]
	fmt.Println("redirect_uri:", redirectURI)

	clientIdstr, clientSecret, _ := r.BasicAuth()
	clientId := ClientId(clientIdstr)

	if app, prs := apps[clientId]; prs &&
		app.ClientSecret == clientSecret &&
		app.RedirectURI == redirectURI &&
		grantType == "authorization_code" {
		if authzResCombinaison, prs := AuthzResCombinaisons[clientId]; prs &&
			authzResCombinaison.Code == code {
			return clientId, true
		}
	}

	return clientId, false
}

func GrantAccessToken(w http.ResponseWriter, clientId ClientId) {
	accessToken := AuthzResCombinaisons[clientId].AccessToken

	data := map[string]string{
		"access_token": accessToken,
		"token_type":   "random string",
	}

	bytes, err := json.Marshal(data)

	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.Header().Set("Content-Length", strconv.Itoa(len(bytes)))
	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}
