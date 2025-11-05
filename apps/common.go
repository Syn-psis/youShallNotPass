package apps

import (
	"net/http"
)

type AuthApplication interface {
	GetToken() string
	GetHost() string
	GetClient() *http.Client
}
