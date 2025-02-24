package main

import (
	"fmt"
	"log"
	"net/http"
	"server/auth"
	"server/resource"
)

const PORT string = ":8080"
const CLIENTURL string = "http://localhost:8090"

func main() {

	http.HandleFunc("/authorize", auth.Authorize)
	http.HandleFunc("/login", auth.Login)
	http.HandleFunc("/consent", auth.RequestApproval)
	http.HandleFunc("/token", auth.GetAccessToken)
	http.HandleFunc("/access", resource.DeliverData)

	fmt.Printf("Server Listening on port%s\n", PORT)
	err := http.ListenAndServe(PORT, nil)
	log.Fatal(err)
}
