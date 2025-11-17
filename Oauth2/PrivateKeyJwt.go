package oauth2

import (
	"crypto/rsa"
	"crypto/tls"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	CLIENT_ASSERTION_TYPE = "urn:ietf:params:oauth:client-assertion-type:jwt-bearer"
)

type PrivateKeyJwt struct {
	lock sync.RWMutex

	clientId   string
	privateKey rsa.PrivateKey
	scopes     string
	token      *Token
	expiration time.Time

	host   string
	client *http.Client
}

func NewPrivateKeyJwt() PrivateKeyJwt {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	return PrivateKeyJwt{
		clientId:   "",
		privateKey: rsa.PrivateKey{},
		scopes:     "",
		token:      nil,
		expiration: *new(time.Time),

		host:   "",
		client: client,
	}
}

func (m *PrivateKeyJwt) ClientId(clientID string) {
	m.clientId = clientID
}

func (m *PrivateKeyJwt) PrivateKey(key rsa.PrivateKey) {
	m.privateKey = key
}

func (m *PrivateKeyJwt) Scopes(scope string) {
	m.scopes = scope
}

func (m *PrivateKeyJwt) Host(host string) {
	m.host = host
}

func (m *PrivateKeyJwt) GetToken() string {
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

func (m *PrivateKeyJwt) BuildClientAssertion() (string, error) {
	assertion := jwt.MapClaims{
		"iss": m.clientId,
		"sub": m.clientId,
		"exp": time.Now().Add(5 * time.Minute).Unix(),
		"iat": time.Now().Unix(),
		"jti": uuid.New(),
		"aud": m.host,
	}
	jwtAssertion := jwt.NewWithClaims(jwt.SigningMethodPS256, assertion)
	token, err := jwtAssertion.SignedString(m.privateKey)

	if err != nil {
		return "", err
	}
	return token, nil
}

func (m *PrivateKeyJwt) retriveToken() Token {

	jwtAssertion, err := m.BuildClientAssertion()
	if err != nil {
		panic("failed to generate jwt assertion")
	}

	body := url.Values{}
	body.Set("grant_type", GRANT_TYPE)
	body.Set("scope", m.scopes)
	body.Set("client_assertion_type", CLIENT_ASSERTION_TYPE)
	body.Set("client_assertion", jwtAssertion)

	payload := strings.NewReader(body.Encode())

	request, err := http.NewRequest("POST", m.host, payload)

	if err != nil {
		panic("was not able to generate request")
	}

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

func (m *PrivateKeyJwt) GetHost() string {
	return m.host
}

func (m *PrivateKeyJwt) GetClient() *http.Client {
	return m.client
}
