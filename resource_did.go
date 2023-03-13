package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"net/http"
	//"net/http"

	//"log"
	"os"
)

const ENV_API_URL = "MATTR_API_URL"
const ENV_AUTH_URL = "MATTR_AUTH_URL"
const ENV_AUTH_AUDIENCE = "MATTR_AUTH_AUDIENCE"
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
			"registered": &schema.Schema{
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
	base_url := os.Getenv(ENV_API_URL)
	url := fmt.Sprintf("%s/core/v1/dids")

	req, err := http.NewRequest("POST", "base_url", nil)
	if err != nil {
		return err
	}

	access_token, err := getAccessToken()

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
	client_id := os.Getenv(ENV_CLIENT_ID)
	client_secret := os.Getenv(ENV_CLIENT_ID)

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

  response_body , err := ioutil.ReadAll(resp.Body)
  if err != nil {
        // handle error
  }
  json.Unmarshal
  return "", err
}
