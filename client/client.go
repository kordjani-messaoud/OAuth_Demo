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
	http.HandleFunc("/submit", submit)

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
	        <input type="hidden" name="redirectURI" value="%s">
	        <input type="hidden" name="responseType" value="%s">

	        <label for="client_id">Client ID:</label>
	        <input type="text" id="clientID" name="clientID" required><br><br>

	        <label for="client_secret">Client Secret:</label>
	        <input type="text" id="clientSecret" name="clientSecret" required><br><br>

	        <label for="scopes">Scopes:</label>
	        <input type="text" id="scopes" name="scopes" required><br><br>

	        <button type="submit" name="action" value="allow">Submit</button>
	    </form>
	</body>
	</html>
	`, state, CLIENTURL, responseType)
	fmt.Fprint(w, html)
}

func submit(w http.ResponseWriter, r *http.Request) {

	if r.Method != "GET" {
		http.Error(w, "Invalide Method", http.StatusMethodNotAllowed)
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}

	clientId = r.Form["clientID"][0]
	clientSecret = r.Form["clientSecret"][0]
	scopes = r.Form["scopes"]
	responseType := r.Form["responseType"]

	url := fmt.Sprintf("%s/authorize?response_type=%s&client_id=%s&client_secret=%s&redirect_uri=%sscope%s&state=%s",
		SERVERURL,
		responseType,
		clientId,
		clientSecret,
		CLIENTURL,
		strings.Join(scopes, " "),
		state,
	)

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}
