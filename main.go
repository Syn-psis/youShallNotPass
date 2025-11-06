package main

import (
	"fmt"
	m2m "github.com/Syn-psis/youShallNotPass/apps/M2M"
)

type User struct {
	Meta struct {
		Location     string `json:"location"`
		ResourceType string `json:"resourceType"`
	} `json:"meta"`

	Roles []struct {
		AudienceValue   string `json:"audienceValue"`
		Display         string `json:"display"`
		AudienceType    string `json:"audienceType"`
		Value           string `json:"value"`
		Ref             string `json:"$ref"`
		AudienceDisplay string `json:"audienceDisplay"`
	} `json:"roles"`

	Groups []struct {
		Display string `json:"display"`
		Value   string `json:"value"`
		Ref     string `json:"$ref"`
	} `json:"groups"`

	ID       string `json:"id"`
	UserName string `json:"userName"`
}

func main() {
	authApp := m2m.New()

	authApp.ClientId("1Tpwubp_LOrKefYt_23oT5sXOlwa")
	authApp.ClientSecret("1fNcZMLPt2Gnkfx6EBLcOzeVRnmbntjyNi0FM_GVPYga")
	authApp.Scopes("internal_user_mgt_create internal_user_mgt_update internal_user_mgt_delete internal_user_mgt_view internal_user_mgt_list")
	authApp.Host("https://localhost:9444")

	response, err := m2m.FilterUser[User](&authApp, nil)

	if err != nil {
		panic("failure")
	}
	fmt.Printf("%s", response.Resources[0].ID)
}
