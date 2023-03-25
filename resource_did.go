package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

import _ "github.com/motemen/go-loghttp/global"

type DidRequest struct {
	Method  string            `json:"method"`
	Options DidRequestOptions `json:"options,omitempty"`
}

type DidRequestOptions struct {
	KeyType string `json:"keyType,omitempty"`
	Url     string `json:"url,omitempty"`
}

type DidResponse struct {
	Did                string `json:"did"`
	RegistrationStatus string `json:"registrationStatus"`
	LocalMetadata      LocalMetadata    `json:"localMetadata"`
}

type LocalMetadata struct {
    Keys           []KeyMetadata `json:"keys"`
    Registered     int64           `json:"registered"`
    InitialDidDoc  json.RawMessage `json:"initialDidDocument"`
}

type KeyMetadata struct {
    DidDocumentKeyId string `json:"didDocumentKeyId"`
    KmsKeyId         string `json:"kmsKeyId"`
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
			"registered": &schema.Schema {
				Type: schema.TypeInt,
				Computed: true,
			},
			"keys": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
								"did_document_key_id": {
										Type:     schema.TypeString,
										Optional: true,
								},
								"kms_key_id": {
										Type:     schema.TypeString,
										Optional: true,
								},
						},
				},
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

	response, err := processDidResponse(resp)
	if err != nil {
		return err
	}

	processDidData(d, &response)

	return nil
}

func processDidData(d *schema.ResourceData, response *DidResponse) {
	d.Set("id", response.Did)
	d.Set("did", response.Did)
	d.Set("registration_status", response.RegistrationStatus)

	keys := []interface{}{}

	for _, k := range response.LocalMetadata.Keys {
		key := map[string]string{
			"did_document_key_id": k.DidDocumentKeyId,
			"kms_key_id": k.KmsKeyId,
		}
		keys = append(keys, key)
	}
}

func resourceDidRead(d *schema.ResourceData, m interface{}) error {
	var err error
	did := d.Get("id").(string)
	base_url := os.Getenv(ENV_API_URL)
	url := fmt.Sprintf("%s/core/v1/dids/%s", base_url, did)

	// formulate request
	request, err := didRequest("GET", url, nil)
	
	// perform the request
	client := http.DefaultClient
	resp, err := client.Do(request)
	if err != nil {
		return err
	}

	// consume the request
	did_response, err := processDidResponse(resp)
	if err != nil {
		return err
	}

	processDidData(d, &did_response)
	return err
}

func resourceDidUpdate(d *schema.ResourceData, m interface{}) error {
	// TODO this method seems a bit dodgy
	err := resourceDidDelete(d, m) // cannot update did in place. need to delete and re-create
	if err != nil {
		return err
	}
	
	return resourceDidCreate(d, m)
}

func resourceDidDelete(d *schema.ResourceData, m interface{}) error {
	var err error
	did := d.Get("id").(string)
	base_url := os.Getenv(ENV_API_URL)
	url := fmt.Sprintf("%s/core/v1/dids/%s", base_url, did)

	// formulate request
	request, err := didRequest("DELETE", url, nil)
	
	// perform the request
	client := http.DefaultClient
	resp, err := client.Do(request)
	if err != nil {
		return err
	}

	if resp.StatusCode < 200 || 200 < resp.StatusCode{
		return fmt.Errorf("Got status code %d from API", resp.StatusCode)
	}

	d.SetId("")
	return err
}

func processDidResponse(resp *http.Response) (DidResponse, error) {
	var response DidResponse

	if resp.StatusCode < 200 || 200 < resp.StatusCode {
		return response, fmt.Errorf("Got status code %d from API", resp.StatusCode)
	}

	// read raw json body
	response_body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return response, err
	}

	// parse the body
	err = json.Unmarshal(response_body, &response)
	if err != nil {
		return response, err
	}

	return response, nil
}

func didRequest(method string, url string, did_request *DidRequest) (*http.Request, error) {
	// prep the request
	base_url := os.Getenv(ENV_API_URL)
	req, err := http.NewRequest(method, base_url, nil)
	if err != nil {
		return req, err
	}
	req.Header.Set("Content-Type", "application/json")
	access_token, err := getAccessToken()
	if err != nil {
		return req, nil
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", access_token))

	if did_request != nil {
		req_body_json, err := json.Marshal(did_request)
		if err != nil {
			return req, err
		}
	
		return http.NewRequest(method, url, bytes.NewBuffer(req_body_json))
	}

	return http.NewRequest(method, url, nil)
}
