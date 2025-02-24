package resource

import (
	"fmt"
	"net/http"
	"server/auth"
)

func DeliverData(w http.ResponseWriter, r *http.Request) {

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}

	clientId, _, ok := r.BasicAuth()
	fmt.Printf("Recived clientId from %s.\n", clientId)

	if !ok {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	token := r.Form["access_token"][0]

	if combinaison, prs := auth.AuthzResCombinaisons[auth.ClientId(clientId)]; prs &&
		combinaison.AccessToken == token {
		data := secretData[auth.ClientId(clientId)]

		fmt.Fprintf(w, "%s", data)
		return
	}

	http.Error(w, "Unauthorized", http.StatusUnauthorized)

}
