package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

type Api struct {
	ClientId             string
	ClientSecret         string
	Audience             string
	AuthUrl              string
	ApiUrl               string
	AccessToken          string
	AccessTokenExpiresAt int64
}

func (e ApiError) Error() string {
	var sb strings.Builder

	// Format the main error message
	if e.StatusCode != 0 && e.Method != "" && e.Url != "" {
		sb.WriteString(fmt.Sprintf("Got status code %d from %s %s\n", e.StatusCode, e.Method, e.Url))
	}

	if e.Code != "" && e.Message != "" {
		sb.WriteString(fmt.Sprintf("Got '%s' error: %s\n", e.Code, e.Message))
	}

	// Format the error details, if any
	if len(e.Details) > 0 {
		for _, detail := range e.Details {
			if detail.Param != "" && detail.Msg != "" {
				sb.WriteString(fmt.Sprintf("Error in %s with '%s': %s\n", detail.Location, detail.Param, detail.Msg))
			}
		}
	}

	return strings.TrimSpace(sb.String())
}

type ErrorDetail struct {
	Msg      string `json:"msg,omitempty"`
	Param    string `json:"param,omitempty"`
	Location string `json:"location,omitempty"`
}

type ApiError struct {
	Method     string        `json:"method,omitempty"`
	Url        string        `json:"url,omitempty"`
	Code       string        `json:"code,omitempty"`
	StatusCode int           `json:"statusCode,omitempty"`
	Message    string        `json:"message,omitempty"`
	Details    []ErrorDetail `json:"details,omitempty"`
}

func ParseError(responseBody []byte) (ApiError, error) {
	var apiError ApiError
	err := json.Unmarshal(responseBody, &apiError)
	return apiError, err
}

type ProviderConfig struct {
	Api Api
}

func (a *Api) Init() {
	if len(a.AuthUrl) == 0 {
		a.AuthUrl = "https://auth.mattr.global/oauth/token"
	}
	if len(a.Audience) == 0 {
		a.Audience = "https://vii.mattr.global" // TODO it should work it out from auth_url
	}
}

func (a *Api) GetUrl(path string) (string, error) {
	return fmt.Sprintf("%s%s", a.ApiUrl, path), nil
}

func (a *Api) GetAccessToken() (string, error) {
	timeStarted := time.Now().Unix()

	var expireTolerance int64 = 15 // get a new token 15 seconds before it expires

	if len(a.AccessToken) != 0 && a.AccessTokenExpiresAt+expireTolerance < time.Now().Unix() {
		log.Printf("Using cached access token")
		return a.AccessToken, nil
	}

	log.Printf("Getting new access token")

	auth_url := a.AuthUrl
	if len(auth_url) == 0 {
		auth_url = "https://auth.mattr.global/oauth/token"
	}
	audience := a.Audience
	if len(audience) == 0 {
		audience = "https://vii.mattr.global" // TODO it should work it out from auth_url
	}

	req_body := AuthRequest{
		ClientId:     a.ClientId,
		ClientSecret: a.ClientSecret,
		Audience:     audience,
		GrantType:    "client_credentials",
	}
	log.Println(req_body)
	req_body_json, err := json.Marshal(req_body)

	if err != nil {
		return "", err
	}
	req, err := http.NewRequest("POST", auth_url, bytes.NewBuffer(req_body_json))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	if resp == nil {
		return "", fmt.Errorf("Response unavailable. Not sure why.")
	}
	if resp.StatusCode < 200 || 299 < resp.StatusCode {
		return "", fmt.Errorf("Invalid status code while retrieving token: %d", resp.StatusCode)
	}
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	response_body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
		return "", err
	}
	var response AuthResponse
	err = json.Unmarshal(response_body, &response)
	if err != nil {
		return "", err
	}

	timeElapsed := time.Now().Unix() - timeStarted
	a.AccessToken = response.AccessToken
	a.AccessTokenExpiresAt = time.Now().Unix() + int64(response.ExpiresIn) - timeElapsed

	return response.AccessToken, nil
}

func (a *Api) Request(method string, url string, body interface{}) (*http.Request, error) {
	log.Printf("Preparing %s request to %s", method, url)

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	if body != nil {
		log.Printf("Will upload JSON to %s", url)
		req_body_json, err := json.Marshal(body)
		log.Printf("JSON: %s", string(req_body_json))
		if err != nil {
			return nil, err
		}
		req_body := bytes.NewBuffer(req_body_json)
		req, err = http.NewRequest(method, url, req_body)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", "application/json")
	}

	req.Header.Set("Accept", "application/json")

	log.Printf("Getting access token")

	access_token, err := a.GetAccessToken()
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", access_token))
	return req, err
}

func Get[T any](a *Api, path string) (*T, error) {
	url, _ := a.GetUrl(path) // TODO error handling
	log.Printf("GET from %s", url)
	client := http.DefaultClient
	request, err := a.Request("GET", url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || 299 < resp.StatusCode {
		return nil, fmt.Errorf("Got status code %d from API", resp.StatusCode)
	}

	return processResponse[T](resp)
}

func Post[T any](a *Api, path string, body interface{}) (*T, error) {
	return Send[T](a, "POST", path, body)
}

func Send[T any](a *Api, method string, path string, body interface{}) (*T, error) {
	url, _ := a.GetUrl(path) // TODO error handling
	log.Printf("%s to %s", method, url)
	client := http.DefaultClient
	request, err := a.Request(method, url, body)
	if err != nil {
		return nil, err
	}
	log.Printf("Doing request")
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	result, err := processResponse[T](resp)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func unmarshal[T any](b []byte) (v T, err error) {
	return v, json.Unmarshal(b, &v)
}
