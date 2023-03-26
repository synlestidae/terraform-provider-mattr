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
	Did                string        `json:"did"` // TODO not available on GET
	RegistrationStatus string        `json:"registrationStatus"`
	LocalMetadata      LocalMetadata `json:"localMetadata"`
}

type LocalMetadata struct {
	Keys          []KeyMetadata   `json:"keys"`
	Registered    int64           `json:"registered"`
	InitialDidDoc json.RawMessage `json:"initialDidDocument"`
}

type KeyMetadata struct {
	DidDocumentKeyId string `json:"didDocumentKeyId"`
	KmsKeyId         string `json:"kmsKeyId"`
}

func resourceDid() *schema.Resource {
	return &schema.Resource{
		Create: resourceDidCreate,
		Read:   resourceDidRead,
		Delete: resourceDidDelete,

		Schema: map[string]*schema.Schema{
			"method": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"url": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Domain or URL from which hostname will be extracted",
				ForceNew: true,
			},
			"key_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"did": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"registered": &schema.Schema{
				Type:     schema.TypeInt,
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
	log.Println("Creating did")

	// TODO: check if the did exists first

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
	base_url, err := getBaseUrl()
	if err != nil {
		return err
	}
	url := fmt.Sprintf("%s/core/v1/dids", base_url)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(req_body_json))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	access_token, err := getAccessToken()
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", access_token))

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
	d.SetId(response.Did)

	return nil
}

func processDidData(d *schema.ResourceData, response *DidResponse) {
	d.Set("registration_status", response.RegistrationStatus)
	d.Set("registered", response.LocalMetadata.Registered)

	keys := []interface{}{}

	for _, k := range response.LocalMetadata.Keys {
		key := map[string]string{
			"did_document_key_id": k.DidDocumentKeyId,
			"kms_key_id":          k.KmsKeyId,
		}
		keys = append(keys, key)
	}

	d.Set("keys", keys)
}

func resourceDidRead(d *schema.ResourceData, m interface{}) error {
	// TODO: handle 404 - 404 means SetId("") i think
	log.Println("Reading did resource...")
	var err error
	did := d.Id()
	log.Printf("Did is %s\n", did)
	base_url, err := getBaseUrl()
	if err != nil {
		return err
	}
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
	d.SetId(did)
	return err
}

func resourceDidDelete(d *schema.ResourceData, m interface{}) error {
	log.Println("Deleting did resource...")
	var err error
	did := d.Id()
	base_url, err := getBaseUrl()
	if err != nil {
		return err
	}
	url := fmt.Sprintf("%s/core/v1/dids/%s", base_url, did)

	// formulate request
	request, err := didRequest("DELETE", url, nil)

	// perform the request
	client := http.DefaultClient
	resp, err := client.Do(request)
	if err != nil {
		return err
	}

	if resp.StatusCode < 200 || 299 < resp.StatusCode {
		return fmt.Errorf("Got status code %d from API", resp.StatusCode)
	}

	d.SetId("")
	return err
}

func processDidResponse(resp *http.Response) (DidResponse, error) {
	var response DidResponse

	if resp.StatusCode < 200 || 299 < resp.StatusCode {
		response_body, _ := ioutil.ReadAll(resp.Body)
		log.Println("Response body: ", string(response_body))
		return response, fmt.Errorf("Got status code %d from API", resp.StatusCode)
	}

	// read raw json body
	response_body, err := ioutil.ReadAll(resp.Body)
	log.Println("Response body: ", string(response_body))
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
	var req *http.Request
	access_token, err := getAccessToken()
	if err != nil {
		return req, nil
	}

	if did_request != nil {
		req.Header.Set("Content-Type", "application/json")
		req_body_json, err := json.Marshal(did_request)
		if err != nil {
			return req, err
		}

		req, err = http.NewRequest(method, url, bytes.NewBuffer(req_body_json))
	} else {
		req, err = http.NewRequest(method, url, nil)
		req.Header.Set("Accept", "application/json")
	}
	if err != nil {
		return req, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", access_token))
	return req, err
}

func getBaseUrl() (string, error) {
	var err error
	base_url := os.Getenv(ENV_API_URL)
	if len(base_url) == 0 {
		err = fmt.Errorf("%s environment variable not set", ENV_API_URL)
	}
	return base_url, err
}
