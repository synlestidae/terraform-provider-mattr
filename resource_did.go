package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"io/ioutil"
	"net/http"
	"os"
	"log"
)

import _ "github.com/motemen/go-loghttp/global"

const ENV_API_URL = "MATTR_API_URL"
const ENV_AUTH_URL = "MATTR_AUTH_URL"
const ENV_AUTH_AUDIENCE = "MATTR_AUDIENCE"
const ENV_CLIENT_ID = "MATTR_CLIENT_ID"
const ENV_CLIENT_SECRET = "MATTR_CLIENT_SECRET"

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

type DidRequest struct {
	Method  string            `json:"method,omitempty"`
	Options DidRequestOptions `json:"options"`
}

type DidRequestOptions struct {
	KeyType string `json:"keyType,omitempty"`
	Url     string `json:"url,omitempty"`
}

type DidResponse struct {
	Did                string `json:"did"`
	RegistrationStatus string `json:"did"`
}

func resourceDid() *schema.Resource {
	return &schema.Resource{
		Create: resourceDidCreate,
		Read:   resourceDidRead,
		Update: resourceDidUpdate,
		Delete: resourceDidDelete,

		Schema: map[string]*schema.Schema{
			"method": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"url": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Domain or URL from which hostname will be extracted",
			},
			"key_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"did": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"registration_status": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"initial_did_document": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceDidCreate(d *schema.ResourceData, m interface{}) error {
	log.Println("Creating did...")

	// prepare the did body request
	method := d.Get("method").(string)
	options := DidRequestOptions{
		KeyType: d.Get("key_type").(string),
		Url:     d.Get("url").(string),
	}
	did_request := DidRequest{
		Method:  method,
		Options: options,
	}

	req_body_json, err := json.Marshal(did_request)
	if err != nil {
		return err
	}

	// prep the request
	base_url := os.Getenv(ENV_API_URL)
	url := fmt.Sprintf("%s/core/v1/dids", base_url)
	req, err := http.NewRequest("POST", base_url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	access_token, err := getAccessToken()
	if err != nil {
		return nil
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", access_token))
	req, err = http.NewRequest("POST", url, bytes.NewBuffer(req_body_json))
	if err != nil {
		return err
	}

	// perform the request
	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode < 200 || 200 < resp.StatusCode {
		return fmt.Errorf("Got status code %d from API", resp.StatusCode)
	}

	// read raw json body
	response_body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// parse the body
	var response DidResponse
	err = json.Unmarshal(response_body, &response)
	if err != nil {
		return err
	}

	d.Set("did", response.Did)
	d.Set("registration_status", response.RegistrationStatus)
	// TODO all the other stuff

	return nil
}

func resourceDidRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceDidUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceDidRead(d, m)
}

func resourceDidDelete(d *schema.ResourceData, m interface{}) error {
	d.SetId("")
	return nil
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
