package oauth2

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

const (
	TOKEN_ENDPOINT = "oauth2/token"
	GRANT_TYPE     = "client_credentials"
)

type ClientCredentials struct {
	lock sync.RWMutex

	clientId     string
	clientSecret string
	scopes       string
	token        *Token
	expiration   time.Time

	host   string
	client *http.Client
}

type Token struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
}

type ApplicationBuilder interface {
	New() ClientCredentials
	ClientId(clientId string)
	ClientSecret()
	Scopes()
	Host()
}

func New() ClientCredentials {

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	return ClientCredentials{
		clientId:     "",
		clientSecret: "",
		scopes:       "",
		token:        nil,
		expiration:   *new(time.Time),

		host:   "",
		client: client,
	}
}

func (m *ClientCredentials) ClientId(clientID string) {
	m.clientId = clientID
}

func (m *ClientCredentials) ClientSecret(secret string) {
	m.clientSecret = secret
}

func (m *ClientCredentials) Scopes(scope string) {
	m.scopes = scope
}

func (m *ClientCredentials) Host(host string) {
	m.host = host
}

type Application interface {
	GetToken() string
	getExpiration() string
	retriveToken() Token
}

func (m *ClientCredentials) GetToken() string {
	now := time.Now()

	if now.After(m.expiration) {

		newToken := m.retriveToken()

		m.lock.Lock()
		defer m.lock.Unlock()

		m.token = &newToken
		m.expiration = time.Now().Add(time.Duration(newToken.ExpiresIn-1) * time.Second)

		return newToken.AccessToken

	} else {
		m.lock.RLock()
		defer m.lock.RUnlock()
		return m.token.AccessToken
	}
}

func (m *ClientCredentials) retriveToken() Token {
	url := fmt.Sprintf("%s/%s", m.host, TOKEN_ENDPOINT)
	body := fmt.Sprintf("grant_type=%s&scope=%s", GRANT_TYPE, m.scopes)
	payload := strings.NewReader(body)

	request, err := http.NewRequest("POST", url, payload)

	if err != nil {
		panic("was not able to generate request")
	}

	request.SetBasicAuth(m.clientId, m.clientSecret)
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	response, err := m.client.Do(request)
	if err != nil {
		panic("was not able to generate request")
	}
	defer response.Body.Close()

	var token Token

	err = json.NewDecoder(response.Body).Decode(&token)

	if err != nil {
		panic("err")
	}

	return token
}

func (m *ClientCredentials) GetHost() string {
	return m.host
}

func (m *ClientCredentials) GetClient() *http.Client {
	return m.client
}
