package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func Login(w http.ResponseWriter, r *http.Request) {

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}
	clientId := r.Form["client_id"][0]
	redirectURI := r.Form["redirect_uri"][0]
	scopes := r.Form["scope"]
	state := r.Form["state"][0]

	html := fmt.Sprintf(`<!DOCTYPE html>
		<html lang="en">
		<head>
		    <meta charset="UTF-8">
		    <meta name="viewport" content="width=device-width, initial-scale=1.0">
		    <title>Photo Gallery</title>
		</head>
		<body style="background-color: blue">
		<h1> Photo Gallery </h1>
		<h3> Print service is requesting you permission to access you photos </h3>	
		    <form action="/consent" method="POST">
		        <input type="hidden" name="client_id" value="%s">
		        <input type="hidden" name="redirect_uri" value="%s">
		        <input type="hidden" name="state" value="%s">
		        <input type="hidden" id="scope" name="scopes" value="%s">

		        <button type="submit" name="action" value="allow">Submit</button>
		    </form>
		</body>
		</html>
		`,
		clientId,
		redirectURI,
		state,
		strings.Join(scopes, " "))
	fmt.Fprint(w, html)
}
