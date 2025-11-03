package main

import (
	"axion/lib/youShallNotPass/apps"
)

func main() {
	authApp := apps.New()

	authApp.ClientId("1Tpwubp_LOrKefYt_23oT5sXOlwa")
	authApp.ClientSecret("1fNcZMLPt2Gnkfx6EBLcOzeVRnmbntjyNi0FM_GVPYga")
	authApp.Scopes("")
	authApp.Host("https://localhost:9444")

	token := authApp.GetToken()

	println(token)
}
