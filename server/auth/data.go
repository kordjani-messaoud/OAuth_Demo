package auth

var apps = map[ClientId]App{
	ClientId("printshop"): {
		ClientId("printshop"),
		"printsecret",
		CLIENTURL + "/callback",
		[]string{"all"},
	},
}

var AuthzResCombinaisons = make(map[ClientId]AuthzResCombinaison)
