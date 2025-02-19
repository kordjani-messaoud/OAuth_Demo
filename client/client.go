package main

import (
	"crypto/rand"
	"fmt"
	"log"
	"net/http"
	"strings"
)

const PORT string = ":8090"
const SERVERURL string = "http://localhost:8080"
const CLIENTURL string = "http://localhost:8090"

var clientId string
var clientSecret string
var state string
var scopes []string

func main() {

	http.HandleFunc("/", root)
	http.HandleFunc("/submit", Submit)

	fmt.Printf("Server Listening on port%s.\n", PORT)
	err := http.ListenAndServe(PORT, nil)
	log.Fatal(err)
}

func root(w http.ResponseWriter, r *http.Request) {
	state = rand.Text()
	responseType := "code"

	html := fmt.Sprintf(`<!DOCTYPE html>
	<html lang="en">
	<head>
	    <meta charset="UTF-8">
	    <meta name="viewport" content="width=device-width, initial-scale=1.0">
	    <title>POST Form</title>
	</head>
	<body>
	    <form action="/submit" method="GET">
	        <input type="hidden" name="state" value="%s">
	        <input type="hidden" name="redirect_uri" value="%s">
	        <input type="hidden" name="response_type" value="%s">

	        <label for="client_id">Client ID:</label>
	        <input type="text" id="client_id" name="client_id" required><br><br>

	        <label for="client_secret">Client Secret:</label>
	        <input type="text" id="client_secret" name="client_secret" required><br><br>

	        <label for="scopes">Scopes:</label>
	        <input type="text" id="scope" name="scope" required><br><br>

	        <button type="submit" name="action" value="allow">Submit</button>
	    </form>
	</body>
	</html>
	`, state, CLIENTURL, responseType)
	fmt.Fprint(w, html)
}

func Submit(w http.ResponseWriter, r *http.Request) {

	if r.Method != "GET" {
		http.Error(w, "Invalide Method", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	clientId = r.Form["client_id"][0]
	clientSecret = r.Form["client_secret"][0]
	scopes = r.Form["scope"]
	responseType := r.Form["response_type"][0]

	url := fmt.Sprintf("%s/authorize?response_type=%s&client_id=%s&client_secret=%s&redirect_uri=%s&scope=%s&state=%s",
		SERVERURL,
		responseType,
		clientId,
		clientSecret,
		CLIENTURL+"/callback",
		strings.Join(scopes, " "),
		state,
	)

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}
