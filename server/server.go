package main

import (
	"fmt"
	"log"
	"net/http"
)

const PORT string = ":8080"
const CLIENTURL string = "http://localhost:8090"

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprint(w, "<h1>Welcome to OAuth Server</h1>")
	})

	fmt.Println("Server Listening on port", PORT)
	err := http.ListenAndServe(PORT, nil)
	log.Fatal(err)
}
