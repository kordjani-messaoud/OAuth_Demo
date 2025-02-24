package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

const PORT string = ":8090"
const SERVERURL string = "http://localhost:8080"
const CLIENTURL string = "http://localhost:8090"

var clientId string
var clientSecret string
var state string
var scopes []string

type Token struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

func main() {

	http.HandleFunc("/", root)
	http.HandleFunc("/submit", Submit)
	http.HandleFunc("/callback", callback)

	fmt.Printf("Server Listening on port%s\n", PORT)
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

func callback(w http.ResponseWriter, r *http.Request) {

	params := r.URL.Query()
	token := GetAccessToken(w, params)
	secretData := GetProtectedData(w, token)

	html := fmt.Sprintf(`<!DOCTYPE html>
	<html lang="en">
	<head>
	    <meta charset="UTF-8">
	    <meta name="viewport" content="width=device-width, initial-scale=1.0">
	    <title>POST Form</title>
	</head>
	<body>
		<h1> Supper Secret Data </h1>
		<b>%s</b>
	</body>
	</html>
	`, secretData)
	fmt.Fprint(w, html)
}

func GetAccessToken(w http.ResponseWriter, params url.Values) Token {

	redirectURI := CLIENTURL + "/callback"
	form := url.Values{
		"grant_type":   []string{"authorization_code"},
		"code":         params["code"],
		"redirect_uri": []string{redirectURI},
	}

	uri := fmt.Sprintf("%s/token", SERVERURL)
	req, _ := http.NewRequest("POST", uri, strings.NewReader(form.Encode()))

	// req.SetBasicAuth(clientId, clientSecret)
	clientCredentials := fmt.Sprintf("%s:%s", clientId, clientSecret)
	clientCredentialsEncoded := base64.StdEncoding.EncodeToString([]byte(clientCredentials))
	authzHeaderValue := fmt.Sprintf("Basic %s", clientCredentialsEncoded)
	req.Header.Add("Authorization", authzHeaderValue)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
	}
	resBody, _ := io.ReadAll(res.Body)

	token := Token{}
	json.Unmarshal(resBody, &token)

	return token
}

func GetProtectedData(w http.ResponseWriter, token Token) string {
	uri := fmt.Sprintf("%s/access?access_token=%s", SERVERURL, token.AccessToken)

	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		fmt.Println("Data Access Err", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
	}

	req.SetBasicAuth(clientId, clientSecret)

	client := &http.Client{}

	res, err := client.Do(req)

	if err != nil {
		fmt.Println("Data Access Err", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
	}

	resBody, _ := io.ReadAll(res.Body)

	return string(resBody)
}
