package auth

const (
	PORT      = ":8080"
	SERVERURL = "http://localhost:8080"
	CLIENTURL = "http://localhost:8090"
)

type ClientId string

type App struct {
	ClientId     ClientId
	ClientSecret string
	RedirectURI  string
	Scopes       []string
}

type AuthzResCombinason struct {
	Code        string
	State       string
	AccessToken string
}
