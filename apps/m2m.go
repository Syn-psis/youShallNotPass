package apps

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

const (
	TOKEN_ENDPOINT = "/oauth2/token"
	GRANT_TYPE     = "grant_type=client_credentials"
)

type M2M struct {
	host string

	clientId     string
	clientSecret string
	scopes       string
	token        Token
	expiration   time.Time

	client http.Client
}

type Token struct {
	access_token string
	token_type   string
	expires_in   int
}

type Application interface {
	GetToken() string
	getExpiration() string
	retriveToken() Token
}

func getToken(m M2M) string {
	return ""
}

func retriveToken(m M2M) Token {
	url := fmt.Sprintf("%s/%s", m.host, TOKEN_ENDPOINT)
	body := fmt.Sprintf("%s&scope=%s", GRANT_TYPE, m.scopes)
	payload := strings.NewReader(body)

	request, err := http.NewRequest("POST", url, payload)

	if err != nil {
		panic("was not able to generate request")
	}

	request.SetBasicAuth(m.clientId, m.clientSecret)
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	m.client.Do(request)

	return Token{
		token_type:   "null",
		expires_in:   34,
		access_token: "null",
	}
}
