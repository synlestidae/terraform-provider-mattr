package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type AuthRequest struct {
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Audience     string `json:"audience"`
	GrantType    string `json:"grant_type"`
}

type AuthResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

func getAccessToken() (string, error) {
	log.Println("Getting access token")

	client_id := os.Getenv(ENV_CLIENT_ID)
	client_secret := os.Getenv(ENV_CLIENT_SECRET)

	auth_url := os.Getenv(ENV_AUTH_URL)
	if len(auth_url) == 0 {
		auth_url = "https://auth.mattr.global/oauth/token"
	}
	audience := os.Getenv(ENV_AUTH_AUDIENCE)
	if len(audience) == 0 {
		audience = "https://vii.mattr.global" // TODO it should work it out from auth_url
	}

	req_body := AuthRequest{
		ClientId:     client_id,
		ClientSecret: client_secret,
		Audience:     audience,
		GrantType:    "client_credentials",
	}
	req_body_json, err := json.Marshal(req_body)
	log.Printf(string(req_body_json))

	if err != nil {
		return "", err
	}
	req, err := http.NewRequest("POST", auth_url, bytes.NewBuffer(req_body_json))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	response_body, err := ioutil.ReadAll(resp.Body)
	log.Printf(string(response_body))
	if err != nil {
		// handle error
		return "", err
	}
	var response AuthResponse
	err = json.Unmarshal(response_body, &response)
	if err != nil {
		return "", err
	}

	return response.AccessToken, nil
}
