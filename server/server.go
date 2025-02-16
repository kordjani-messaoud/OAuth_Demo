package main

import (
	"fmt"
	"log"
	"net/http"
	"server/auth"
)

const PORT string = ":8080"
const CLIENTURL string = "http://localhost:8090"

func main() {

	http.HandleFunc("/authorize", auth.Authorize)

	fmt.Println("Server Listening on port", PORT)
	err := http.ListenAndServe(PORT, nil)
	log.Fatal(err)
}
