package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type Api struct {
	ClientId     string
	ClientSecret string
	Audience     string
	AuthUrl      string
	ApiUrl       string
}

func (a *Api) Init() {
	if len(a.AuthUrl) == 0 {
		a.AuthUrl = "https://auth.mattr.global/oauth/token"
	}
	if len(a.Audience) == 0 {
		a.Audience = "https://vii.mattr.global" // TODO it should work it out from auth_url
	}
}

func InitFromEnv() Api {
	clientId := os.Getenv(ENV_CLIENT_ID)
	clientSecret := os.Getenv(ENV_CLIENT_SECRET)
	authUrl := os.Getenv(ENV_AUTH_URL)
	audience := os.Getenv(ENV_AUTH_AUDIENCE)
	apiUrl := os.Getenv(ENV_API_URL)

	return Api{
		ClientId:     clientId,
		ClientSecret: clientSecret,
		AuthUrl:      authUrl,
		Audience:     audience,
		ApiUrl:       apiUrl,
	}
}

type DidRequest struct {
	Method  string            `json:"method"`
	Options DidRequestOptions `json:"options,omitempty"`
}

type DidRequestOptions struct {
	KeyType string `json:"keyType,omitempty"`
	Url     string `json:"url,omitempty"`
}

type DidResponse struct {
	Did                string        `json:"did,omitempty"` // Not available on GET
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

type WebhookRequest struct {
	Events   []string `json:"events"`
	Url      string
	disabled bool
}

type WebhookResponse struct {
	Id       string   `json:"id"`
	Events   []string `json:"events"`
	Url      string
	disabled bool
}

type WebhookListResponse struct {
	NextCursor string `json:"nextCursor"`
}

type IssuerCredential struct {
	IssuerDid          string               `json:"issuerDid"`
	IssuerLogoUrl      string               `json:"issuerLogoUrl"`
	IssuerIconUrl      string               `json:"issuerIconUrl"`
	Name               string               `json:"name"`
	Description        string               `json:"description"`
	Context            []string             `json:"string"`
	Type               string               `json:"string"`
	CredentialBranding []CredentialBranding `json:"credentialBranding"`
	FederatedProvider  FederatedProvider    `json:"federatedProvider"`
}

type FederatedProvider struct {
	Url                     string
	Scope                   []string `json:"scope"`
	ClientId                string   `json:"clientId"`
	ClientSecret            string   `json:"clientSecret"`
	TokenEndpointAuthMethod string   `json:"tokenEndpointAuthMethod"`
	ClaimsSource            string   `json:"claimsSource"`
}

type ClaimMapping struct {
	JsonLdTerm string `json:"jsonLdTerm"`
	OidcClaim  string `json:"oidcClaim"`
}

type CredentialBranding struct {
	BackgroundColor   string `json:"backgroundColor"`
	WatermarkImageUrl string `json:"watermarkImageUrl"`
}

type IssuerRequest struct {
	Credential                 IssuerCredential  `json:"credential"`
	FederatedProvider          FederatedProvider `json:"federatedProvider"`
	StaticRequestParameters    map[string]any    `json:"staticRequestParameters"`
	ForwardedRequestParameters []string          `json:"forwardedRequestParameters"`
	ClaimMappings              []ClaimMapping    `json:"claimMappings"`
}

type IssuerResponse struct {
	Id                         string            `json:"id"`
	Credential                 IssuerCredential  `json:"credential"`
	FederatedProvider          FederatedProvider `json:"federatedProvider"`
	StaticRequestParameters    map[string]any    `json:"staticRequestParameters"`
	ForwardedRequestParameters []string          `json:"forwardedRequestParameters"`
	ClaimMappings              []ClaimMapping    `json:"claimMappings"`
}

type IssuerListResponse struct {
	Data       []IssuerResponse `json:"data"`
	NextCursor string           `json:"nextCursor"`
}

type WellKnownResponse struct {
}

type IssuerClientRequest struct {
	Name                     string `json:"string"`
	RedirectUris             []string
	ResponseTypes            []string
	GrantTypes               []string
	TokenEndpointAuthMethod  string `json:"tokenEndpointAuthMethod"`
	IdTokenSignedResponseAlg string `json:"idTokenSignedResponseAlg"`
	ApplicationType          string `json:"applicationType"`
}

type IssuerClientResponse struct {
	Id                       string `json:"id"`
	Secret                   string `json:"secret"`
	Name                     string `json:"string"`
	RedirectUris             []string
	ResponseTypes            []string
	GrantTypes               []string
	TokenEndpointAuthMethod  string `json:"tokenEndpointAuthMethod"`
	IdTokenSignedResponseAlg string `json:"idTokenSignedResponseAlg"`
	ApplicationType          string `json:"applicationType"`
}

type IssuerClientListResponse struct {
	NextCursor string                 `json:"nextCursor"`
	Data       []IssuerClientResponse `json:"data"`
}

type VerifierRequest struct {
}

type VerifierResponse struct {
}

type VerifierListResponse struct {
	NextCursor string             `json:"nextCursor"`
	Data       []VerifierResponse `json:"data"`
}

type VerifierClientRequest struct {
	VerifierDid            string         `json:"verifierDid"`
	PresentationTemplateId string         `json:"presentationTemplateId"`
	ClaimMappings          []ClaimMapping `json:"claimMappings"`
}

type VerifierClientResponse struct {
	Id                     string         `json:"id"`
	VerifierDid            string         `json:"verifierDid"`
	PresentationTemplateId string         `json:"presentationTemplateId"`
	ClaimMappings          []ClaimMapping `json:"claimMappings"`
}

type VerifierClientListResponse struct {
	NextCursor string             `json:"nextCursor"`
	Data       []VerifierResponse `json:"data"`
}

func (a *Api) PostDid(did DidRequest) (*DidResponse, error) {
	return Post[DidResponse](a, "/core/v1/dids", did)
}

func (a *Api) GetDid(id string) (*DidResponse, error) {
	return Get[DidResponse](a, "/core/v1/dids")
}

func (a *Api) DeleteDid(id string) error {
	return Delete[DidResponse](a, "/core/v1/dids")
}

// Webhooks
func (a *Api) PostWebhook(webhook *WebhookRequest) (*WebhookResponse, error) {
	//return nil, fmt.Errorf("Not quite implemented yet")
	return Post[WebhookResponse](a, "/core/v1/webhooks", webhook)
}

func (a *Api) GetWebhook(id string) (*WebhookResponse, error) {
	return Get[WebhookResponse](a, fmt.Sprintf("/core/v1/webhooks/%s", id))
}

func (a *Api) GetWebhooks() (*WebhookListResponse, error) {
	return Get[WebhookListResponse](a, "/core/v1/webhooks")
}

func (a *Api) PutWebhook(webhook *WebhookResponse) (*WebhookResponse, error) {
	return Put[WebhookResponse](a, fmt.Sprintf("/core/v1/webhooks/%s", webhook.Id), webhook)
}

func (a *Api) DeleteWebhook(id string) error {
	return Delete[WebhookResponse](a, fmt.Sprintf("/core/v1/webhooks/%s", id))
}

// .well-known
func (a *Api) GetWellKnown() (*WellKnownResponse, error) {
	return nil, fmt.Errorf("Not quite implemented yet")
}

// OIDC Issuers
func (a *Api) PostIssuer() (*IssuerResponse, error) {
	return nil, fmt.Errorf("Not quite implemented yet")
}

func (a *Api) GetIssuer() (*IssuerResponse, error) {
	return nil, fmt.Errorf("Not quite implemented yet")
}

func (a *Api) GetIssuers() (*IssuerListResponse, error) {
	return nil, fmt.Errorf("Not quite implemented yet")
}

func (a *Api) PutIssuer() (*IssuerResponse, error) {
	return nil, fmt.Errorf("Not quite implemented yet")
}

func (a *Api) DeleteIssuer() error {
	return fmt.Errorf("Not quite implemented yet")
}

// Issuer Clients
func (a *Api) PostIssuerClient() (*IssuerClientResponse, error) {
	return nil, fmt.Errorf("Not quite implemented yet")
}

func (a *Api) GetIssuerClients() (*IssuerClientListResponse, error) {
	return nil, fmt.Errorf("Not quite implemented yet")
}

func (a *Api) GetIssuerClient() (*IssuerClientResponse, error) {
	return nil, fmt.Errorf("Not quite implemented yet")
}

func (a *Api) PutIssuerClient() (*IssuerClientResponse, error) {
	return nil, fmt.Errorf("Not quite implemented yet")
}

func (a *Api) DeleteIssuerClient() error {
	return fmt.Errorf("Not quite implemented yet")
}

// Verifiers
func (a *Api) CreateVerifier(verifier *VerifierRequest) (*VerifierResponse, error) {
	return nil, fmt.Errorf("Not quite implemented yet")
}

func (a *Api) GetVerifiers(cursor string) (*VerifierListResponse, error) {
	return nil, fmt.Errorf("Not quite implemented yet")
}

func (a *Api) GetVerifier(id string) (*VerifierResponse, error) {
	return nil, fmt.Errorf("Not quite implemented yet")
}

func (a *Api) PutVerifier(id string, verifier *VerifierRequest) (*VerifierResponse, error) {
	return nil, fmt.Errorf("Not quite implemented yet")
}

func (a *Api) DeleteVerifier(id string) error {
	return fmt.Errorf("Not quite implemented yet")
}

func (a *Api) GetUrl(path string) (string, error) {
	return "", fmt.Errorf("Not quite implemented yet")
}

func (a *Api) GetAccessToken() (string, error) {
	var auth_url string
	if len(a.AuthUrl) == 0 {
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

	return response.AccessToken, nil
}

func (a *Api) Request(method string, resource string, body interface{}) (*http.Request, error) {
	url := fmt.Sprintf("%s%s", a.ApiUrl, resource) // TODO remove trailing backslashes

	var req *http.Request
	access_token, err := a.GetAccessToken()
	if err != nil {
		return req, nil
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
		req_body_json, err := json.Marshal(body)
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

/*func (a *Api) BaseUrl() (string, error) {
	if len(a.ApiUrl) > 0 {
		return a.ApiUrl
	}
	var err error
	base_url := os.Getenv(ENV_API_URL)
	if len(base_url) == 0 {
		err = fmt.Errorf("%s environment variable not set", ENV_API_URL)
	}
	return base_url, err
}*/

func Get[T any](a *Api, path string) (*T, error) {
	client := http.DefaultClient
	request, err := a.Request("GET", path, nil)
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

	//result, err := processResponse[T](resp)
	return processResponse[T](resp)
}

func Post[T any](a *Api, path string, body interface{}) (*T, error) {
	return Send[T](a, "POST", path, body)
}

func Put[T any](a *Api, path string, body interface{}) (*T, error) {
	return Send[T](a, "PUT", path, body)
}

func Delete[T any](a *Api, path string) error {
	client := http.DefaultClient
	request, err := a.Request("DELETE", path, nil)
	if err != nil {
		return err
	}
	resp, err := client.Do(request)
	if err != nil {
		return err
	}

	if resp.StatusCode < 200 || 299 < resp.StatusCode {
		return fmt.Errorf("Got status code %d from API", resp.StatusCode)
	}

	return nil
}

func Send[T any](a *Api, method string, path string, body interface{}) (*T, error) {
	client := http.DefaultClient
	request, err := a.Request("POST", path, body)
	if err != nil {
		return nil, err
	}
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

func processResponse[T any](resp *http.Response) (*T, error) {
	defer resp.Body.Close()

	if resp.StatusCode < 200 || 299 < resp.StatusCode {
		response_body, _ := ioutil.ReadAll(resp.Body)
		log.Println("Response body: ", string(response_body))
		return nil, fmt.Errorf("Got status code %d from API", resp.StatusCode)
	}

	// read raw json body
	response_body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// parse the body
	var response T
	err = json.Unmarshal(response_body, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func unmarshal[T any](b []byte) (v T, err error) {
	return v, json.Unmarshal(b, &v)
}
