package m2m

import (
	"axion/lib/youShallNotPass/apps"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

const (
	SCIM_USER = "scim2/Users"
)

type ScimResource[T any] struct {
	TotalResults int      `json:"totalResults"`
	StartIndex   int      `json:"startIndex"`
	ItemsPerPage int      `json:"itemsPerPage"`
	Schemas      []string `json:"schemas"`
	Resources    []T      `json:"Resources"`
}

type ScimSearch struct {
	Schemas    []string `json:"schemas"`
	Attributes []string `json:"attributes"`
	Filter     string   `json:"filter"`
	Domain     string   `json:"domain"`
	StartIndex int      `json:"startIndex"`
	Count      int      `json:"count"`
}

func FilterUser[T any](app apps.AuthApplication, filters *string) (ScimResource[T], error) {
	var user ScimResource[T]
	defaultFilters := ""

	if filters == nil {
		filters = &defaultFilters
	}

	url := fmt.Sprintf("%s/%s?%s", app.GetHost(), SCIM_USER, *filters)
	request, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return user, err
	}

	request.Header.Set("Authorization", "Bearer "+app.GetToken())
	request.Header.Set("accept", "application/scim+json")

	response, err := app.GetClient().Do(request)

	if err != nil {
		return user, err
	}
	defer response.Body.Close()

	err = json.NewDecoder(response.Body).Decode(&user)
	if err != nil {
		return user, err
	}

	return user, nil

}

func CreateUser[T any](app apps.AuthApplication, user T) error {

	url := fmt.Sprintf("%s/%s", app.GetHost(), SCIM_USER)
	body, err := json.Marshal(user)

	if err != nil {
		return err
	}

	payload := bytes.NewReader(body)

	request, err := http.NewRequest("POST", url, payload)

	if err != nil {
		return err
	}

	request.Header.Set("Authorization", "Bearer "+app.GetToken())
	request.Header.Set("accept", "application/scim+json")

	response, err := app.GetClient().Do(request)

	if err != nil {
		return err
	}

	if response.StatusCode != 201 {
		errorMessage := fmt.Sprintf("user not created due to server response: %d", response.StatusCode)
		return errors.New(errorMessage)
	}

	return nil
}

func SearchUser[T any](app apps.AuthApplication, search ScimSearch) (ScimResource[T], error) {

	var users ScimResource[T]

	url := fmt.Sprintf("%s/%s/.search", app.GetHost(), SCIM_USER)
	body, err := json.Marshal(search)

	if err != nil {
		return users, err
	}

	payload := bytes.NewReader(body)

	request, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return users, err
	}

	request.Header.Set("Authorization", "Bearer "+app.GetToken())
	request.Header.Set("accept", "application/scim+json")

	response, err := app.GetClient().Do(request)
	if err != nil {
		return users, err
	}

	if response.StatusCode != 200 {
		errorMessage := fmt.Sprintf("user not created due to server response: %d", response.StatusCode)
		return users, errors.New(errorMessage)
	}

	defer response.Body.Close()

	err = json.NewDecoder(response.Body).Decode(&users)
	if err != nil {
		return users, err
	}

	return users, nil
}

func GetUserByID[T any](app apps.AuthApplication, UserID string) (T, error) {

	var user T

	url := fmt.Sprintf("%s/%s/%s", app.GetHost(), SCIM_USER, UserID)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return user, err
	}

	request.Header.Set("Authorization", "Bearer "+app.GetToken())
	request.Header.Set("accept", "application/scim+json")

	response, err := app.GetClient().Do(request)
	if err != nil {
		return user, err
	}

	if response.StatusCode != 200 {
		errorMessage := fmt.Sprintf("could not retrive user due to server response: %d", response.StatusCode)
		return user, errors.New(errorMessage)
	}

	defer response.Body.Close()

	err = json.NewDecoder(response.Body).Decode(&user)
	if err != nil {
		return user, err
	}

	return user, nil

}

func PutUpdateUser[T any, U any](app apps.AuthApplication, UserID string, Update U) (T, error) {

	var user T

	url := fmt.Sprintf("%s/%s/%s", app.GetHost(), SCIM_USER, UserID)
	body, err := json.Marshal(user)
	if err != nil {
		return user, err
	}

	payload := bytes.NewReader(body)

	request, err := http.NewRequest("PUT", url, payload)
	if err != nil {
		return user, err
	}

	request.Header.Set("Authorization", "Bearer "+app.GetToken())
	request.Header.Set("accept", "application/scim+json")

	response, err := app.GetClient().Do(request)
	if err != nil {
		return user, err
	}

	if response.StatusCode != 200 {
		errorMessage := fmt.Sprintf("could not update user due to server response: %d", response.StatusCode)
		return user, errors.New(errorMessage)
	}

	defer response.Body.Close()

	err = json.NewDecoder(response.Body).Decode(&user)
	if err != nil {
		return user, err
	}

	return user, nil

}

func DeleteUser(app apps.AuthApplication, UserID string) error {
	url := fmt.Sprintf("%s/%s/%s", app.GetHost(), SCIM_USER, UserID)

	request, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	request.Header.Set("Authorization", "Bearer "+app.GetToken())
	request.Header.Set("accept", "application/scim+json")

	response, err := app.GetClient().Do(request)
	if err != nil {
		return err
	}

	if response.StatusCode != 204 {
		errorMessage := fmt.Sprintf("could not delete user due to server response: %d", response.StatusCode)
		return errors.New(errorMessage)
	}

	return nil

}

func PatchUpdateUser[T any]() (T, error) {

}
