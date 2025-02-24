package resource

import "server/auth"

var secretData = map[auth.ClientId]string{
	auth.ClientId("printshop"): "super secret information",
}
