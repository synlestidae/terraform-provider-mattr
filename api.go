package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Api struct {
	ClientId     string
	ClientSecret string
	Audience     string
	AuthUrl      string
	ApiUrl       string
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

// OK so this is good
// 12 structs need to be declared (6 requests, 6 responses)
// there will also be an error struct, with the mattr format messages

// also need to write
// * the http request thing, and it should support keepalive
// * an order of URL inference
// * inferring the URL from access token

// oh and unit tests

func (a *Api) PostDid() (*DidResponse, error) {
	return nil, fmt.Errorf("Not quite implemented yet")
}

func (a *Api) GetDid() (*DidResponse, error) {
	return nil, fmt.Errorf("Not quite implemented yet")
	//request, err := a.Request("GET", "/core/v1/dids", nil)
	//if err == nil {
	//		return nil, err
	//}
	// todo
}

func (a *Api) DeleteDid() error {
	return fmt.Errorf("Not quite implemented yet")
}

// Webhooks
func (a *Api) PostWebhook(*WebhookRequest) (*WebhookResponse, error) {
	return nil, fmt.Errorf("Not quite implemented yet")
}

func (a *Api) GetWebhook(id string) (*WebhookResponse, error) {
	return nil, fmt.Errorf("Not quite implemented yet")
}

func (a *Api) GetWebhooks() (*WebhookListResponse, error) {
	return nil, fmt.Errorf("Not quite implemented yet")
}

func (a *Api) PutWebhook(webhook *WebhookRequest) (*WebhookResponse, error) {
	return nil, fmt.Errorf("Not quite implemented yet")
}

func (a *Api) DeleteWebhook(id string) error {
	return fmt.Errorf("Not quite implemented yet")
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

func (a *Api) Request(method string, resource string, body *interface{}) (*http.Request, error) {
	url := fmt.Sprintf("%s%s", a.BaseUrl(), resource) // TODO remove trailing backslashes
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

func (a *Api) BaseUrl() string {
	return ""
}
