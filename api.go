package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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

func (a *Api) Init() {
	if len(a.AuthUrl) == 0 {
		a.AuthUrl = "https://auth.mattr.global/oauth/token"
	}
	if len(a.Audience) == 0 {
		a.Audience = "https://vii.mattr.global" // TODO it should work it out from auth_url
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

type PublicKey struct {
	Id              string `json:"id"`
	Type            string `json:"type"`
	Controller      string `json:"controller"`
	PublicKeyBase58 string `json:"publicKeyBase58"`
}

type KeyAgreement struct {
	Id              string `json:"id"`
	Type            string `json:"type"`
	Controller      string `json:"controller"`
	PublicKeyBase58 string `json:"publicKeyBase58"`
}

type DidDocument struct {
	Id                   string         `json:"id"`
	Context              string         `json:"@context"`
	PublicKey            []PublicKey    `json:"publicKey"`
	KeyAgreement         []KeyAgreement `json:"keyAgreement"`
	Authentication       []string       `json:"authentication"`
	AssertionMethod      []string       `json:"assertionMethod"`
	CapabilityDelegation []string       `json:"capabilityDelegation"`
	CapabilityInvocation []string       `json:"capabilityInvocation"`
}

type LocalMetadata struct {
	Keys               []KeyMetadata `json:"keys"`
	Registered         int64         `json:"registered"`
	InitialDidDocument DidDocument   `json:"initialDidDocument"`
}

type KeyMetadata struct {
	DidDocumentKeyId string `json:"didDocumentKeyId"`
	KmsKeyId         string `json:"kmsKeyId"`
}

type WebhookRequest struct {
	Events   []string `json:"events"`
	Url      string   `json:"url"`
	Disabled bool     `json:"disabled,omitempty"`
}

type WebhookResponse struct {
	Id       string   `json:"id"`
	Events   []string `json:"events"`
	Url      string   `json:"url"`
	Disabled bool     `json:"disabled,omitempty"`
}

type WebhookListResponse struct {
	NextCursor string `json:"nextCursor"`
}

type IssuerCredential struct {
	IssuerDid          string              `json:"issuerDid,omitempty"`
	IssuerLogoUrl      string              `json:"issuerLogoUrl,omitempty"`
	IssuerIconUrl      string              `json:"issuerIconUrl,omitempty"`
	Name               string              `json:"name,omitempty"`
	Description        string              `json:"description,omitempty"`
	Context            []string            `json:"context,omitempty"`
	Type               []string            `json:"type,omitempty"`
	CredentialBranding *CredentialBranding `json:"credentialBranding,omitempty"`
	FederatedProvider  *FederatedProvider  `json:"federatedProvider,omitempty"`
}

type FederatedProvider struct {
	Url                     string   `json:"url,omitempty"`
	Scope                   []string `json:"scope"`
	ClientId                string   `json:"clientId,omitempty"`
	ClientSecret            string   `json:"clientSecret,omitempty"`
	TokenEndpointAuthMethod string   `json:"tokenEndpointAuthMethod,omitempty"`
	ClaimsSource            string   `json:"claimsSource,omitempty"`
}

type ClaimMapping struct {
	JsonLdTerm string `json:"jsonLdTerm"`
	OidcClaim  string `json:"oidcClaim"`
}

type CredentialBranding struct {
	BackgroundColor   string `json:"backgroundColor,omitempty"`
	WatermarkImageUrl string `json:"watermarkImageUrl,omitempty"`
}

type IssuerRequest struct {
	Credential                 *IssuerCredential  `json:"credential,omitempty"`
	FederatedProvider          *FederatedProvider `json:"federatedProvider,omitempty"`
	StaticRequestParameters    map[string]any     `json:"staticRequestParameters,omitempty"`
	ForwardedRequestParameters []string           `json:"forwardedRequestParameters"`
	ClaimMappings              []ClaimMapping     `json:"claimMappings"`
}

type IssuerResponse struct {
	Id                         string             `json:"id"`
	Credential                 *IssuerCredential  `json:"credential"`
	FederatedProvider          *FederatedProvider `json:"federatedProvider"`
	StaticRequestParameters    map[string]any     `json:"staticRequestParameters"`
	ForwardedRequestParameters []string           `json:"forwardedRequestParameters"`
	ClaimMappings              []ClaimMapping     `json:"claimMappings"`
}

type IssuerListResponse struct {
	Data       []IssuerResponse `json:"data"`
	NextCursor string           `json:"nextCursor"`
}

type IssuerInfo struct {
	Name    string `json:"name"`
	LogoUrl string `json:"logoUrl,omitempty"`
	IconUrl string `json:"iconUrl,omitempty"`
}

type Context string

type ClaimMappingConfig struct {
	MapFrom      string `json:"mapFrom"`
	Required     bool   `json:"required,omitempty"`
	DefaultValue string `json:"defaultValue,omitempty"`
}

type ExpiresIn struct {
	Years   int `json:"years,omitempty"`
	Months  int `json:"months,omitempty"`
	Weeks   int `json:"weeks,omitempty"`
	Days    int `json:"days,omitempty"`
	Hours   int `json:"hours,omitempty"`
	Minutes int `json:"minutes,omitempty"`
	Seconds int `json:"seconds,omitempty"`
}

type CredentialConfig struct {
	Id                 string                        `json:"id,omitempty"`
	Name               string                        `json:"name"`
	Description        string                        `json:"description,omitempty"`
	Type               string                        `json:"type"`
	AdditionalTypes    []string                      `json:"additionalTypes,omitempty"`
	Contexts           []string                      `json:"contexts,omitempty"`
	Issuer             IssuerInfo                    `json:"issuer"`
	CredentialBranding CredentialBranding            `json:"credentialBranding,omitempty"`
	ClaimMappings      map[string]ClaimMappingConfig `json:"claimMappings,omitempty"`
	Persist            bool                          `json:"persist,omitempty"`
	Revocable          bool                          `json:"revocable,omitempty"`
	ClaimSourceId      string                        `json:"claimSourceId,omitempty"`
	ExpiresIn          ExpiresIn                     `json:"expiresIn,omitempty"`
}

type WellKnownResponse struct {
}

type IssuerClientRequest struct {
	Name                     string   `json:"string"`
	RedirectUris             []string `json:"redirectUris"`
	ResponseTypes            []string `json:"responseTypes"`
	GrantTypes               []string `json:"grantTypes"`
	TokenEndpointAuthMethod  string   `json:"tokenEndpointAuthMethod"`
	IdTokenSignedResponseAlg string   `json:"idTokenSignedResponseAlg"`
	ApplicationType          string   `json:"applicationType"`
}

type IssuerClientResponse struct {
	Id                       string   `json:"id"`
	Secret                   string   `json:"secret"`
	Name                     string   `json:"string"`
	RedirectUris             []string `json:"redirectUris"`
	ResponseTypes            []string `json:"responseTypes"`
	GrantTypes               []string `json:"grantTypes"`
	TokenEndpointAuthMethod  string   `json:"tokenEndpointAuthMethod"`
	IdTokenSignedResponseAlg string   `json:"idTokenSignedResponseAlg"`
	ApplicationType          string   `json:"applicationType"`
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

type ClaimSourceAuthorization struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type ClaimSourceRequestParameter struct {
	MapFrom      string `json:"mapFrom"`
	DefaultValue string `json:"defaultValue"`
}

type ClaimSource struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Url  string `json:"url"`

	Authorization     ClaimSourceAuthorization               `json:"authorization"`
	RequestParameters map[string]ClaimSourceRequestParameter `json:"requestParameters"`
}

func (a *Api) PostDid(did DidRequest) (*DidResponse, error) {
	return Post[DidResponse](a, "/core/v1/dids", did)
}

func (a *Api) GetDid(id string) (*DidResponse, error) {
	return Get[DidResponse](a, "/core/v1/dids")
}

func (a *Api) DeleteDid(id string) error {
	return Delete(a, fmt.Sprintf("/core/v1/dids/%s", id))
}

// Webhooks
func (a *Api) PostWebhook(webhook *WebhookRequest) (*WebhookResponse, error) {
	return Post[WebhookResponse](a, "/core/v1/webhooks", webhook)
}

func (a *Api) GetWebhook(id string) (*WebhookResponse, error) {
	return Get[WebhookResponse](a, fmt.Sprintf("/core/v1/webhooks/%s", id))
}

func (a *Api) GetWebhooks() (*WebhookListResponse, error) {
	return Get[WebhookListResponse](a, "/core/v1/webhooks")
}

func (a *Api) PutWebhook(id string, webhook *WebhookRequest) (*WebhookResponse, error) {
	return Put[WebhookResponse](a, fmt.Sprintf("/core/v1/webhooks/%s", id), webhook)
}

func (a *Api) DeleteWebhook(id string) error {
	return Delete(a, fmt.Sprintf("/core/v1/webhooks/%s", id))
}

// .well-known
func (a *Api) GetWellKnown() (*WellKnownResponse, error) {
	return nil, fmt.Errorf("Not quite implemented yet")
}

// OIDC Issuers
func (a *Api) PostIssuer(issuer *IssuerRequest) (*IssuerResponse, error) {
	return Post[IssuerResponse](a, "/ext/oidc/v1/issuers", issuer)
}

func (a *Api) GetIssuer(id string) (*IssuerResponse, error) {
	return Get[IssuerResponse](a, fmt.Sprintf("/ext/oidc/v1/issuers/%s", id))
}

func (a *Api) GetIssuers() (*IssuerListResponse, error) {
	return Get[IssuerListResponse](a, fmt.Sprintf("/ext/oidc/v1/issuers"))
}

func (a *Api) PutIssuer(id string, issuer *IssuerRequest) (*IssuerResponse, error) {
	return Put[IssuerResponse](a, fmt.Sprintf("/ext/oidc/v1/issuers/%s", id), issuer)
}

func (a *Api) DeleteIssuer(id string) error {
	return Delete(a, fmt.Sprintf("/ext/oidc/v1/issuers/%s", id))
}

// Credential config

func (a *Api) PostCredentialConfig(config *CredentialConfig) (*CredentialConfig, error) {
	return Post[CredentialConfig](a, "/v2/credentials/web-semantic/configurations", config)
}

func (a *Api) GetCredentialConfig(id string) (*CredentialConfig, error) {
	return Get[CredentialConfig](a, fmt.Sprintf("/v2/credentials/web-semantic/configurations/%s", id))
}

func (a *Api) PutCredentialConfig(id string, config *CredentialConfig) (*CredentialConfig, error) {
	return Put[CredentialConfig](a, fmt.Sprintf("/v2/credentials/web-semantic/configurations/%s", id), config)
}

func (a *Api) DeleteCredentialConfig(id string) error {
	return Delete(a, fmt.Sprintf("/v2/credentials/web-semantic/configurations/%s", id))
}

// Claim source

func (a *Api) PostClaimSource(claimSource *ClaimSource) (*ClaimSource, error) {
	return Post[ClaimSource](a, "/v1/claimsources", claimSource)
}

func (a *Api) GetClaimSource(id string) (*ClaimSource, error) {
	return Get[ClaimSource](a, fmt.Sprintf("/v1/claimsources/%s", id))
}

func (a *Api) PutClaimSource(id string, claimSource *ClaimSource) (*ClaimSource, error) {
	return Put[ClaimSource](a, fmt.Sprintf("/v1/claimsources/%s", id), claimSource)
}

func (a *Api) DeleteClaimSource(id string) error {
	return Delete(a, fmt.Sprintf("/v1/claimsources/%s", id))
}

// End claim source

// Issuer Clients
func (a *Api) PostIssuerClient(issuerId string, request *IssuerClientRequest) (*IssuerClientResponse, error) {
	return Post[IssuerClientResponse](a, fmt.Sprintf("/ext/oidc/v1/issuers/%s/clients", issuerId), request)
}

func (a *Api) GetIssuerClients(issuerId string) (*IssuerClientListResponse, error) {
	return Get[IssuerClientListResponse](a, fmt.Sprintf("/ext/oidc/v1/issuers/%s/clients", issuerId))
}

func (a *Api) GetIssuerClient(issuerId string, issuerClientId string) (*IssuerClientResponse, error) {
	return Get[IssuerClientResponse](a, fmt.Sprintf("/ext/oidc/v1/issuers/%s/clients/%s", issuerId, issuerClientId))
}

func (a *Api) PutIssuerClient(issuerId string, clientId string, request *IssuerClientRequest) (*IssuerClientResponse, error) {
	return Post[IssuerClientResponse](a, fmt.Sprintf("/ext/oidc/v1/issuers/%s/clients/%s", issuerId, clientId), request)
}

func (a *Api) DeleteIssuerClient(issuerId string, clientId string) error {
	return Delete(a, fmt.Sprintf("/ext/oidc/v1/issuers/%s/clients/%s", issuerId, clientId))
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

	log.Println("Access token is %s", a.AccessToken)

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

	//result, err := processResponse[T](resp)
	return processResponse[T](resp)
}

func Post[T any](a *Api, path string, body interface{}) (*T, error) {
	return Send[T](a, "POST", path, body)
}

func Put[T any](a *Api, path string, body interface{}) (*T, error) {
	return Send[T](a, "PUT", path, body)
}

func Delete(a *Api, path string) error {
	url, _ := a.GetUrl(path) // TODO error handling
	log.Printf("DELETE %s", url)
	client := http.DefaultClient
	request, err := a.Request("DELETE", url, nil)
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

func processResponse[T any](resp *http.Response) (*T, error) {
	defer resp.Body.Close()

	if resp.StatusCode < 200 || 299 < resp.StatusCode {
		response_body, _ := ioutil.ReadAll(resp.Body)
		log.Println("Response body: ", string(response_body))
		return nil, fmt.Errorf("Got status code %d from API", resp.StatusCode)
	}

	// read raw json body
	response_body, err := ioutil.ReadAll(resp.Body)
	log.Println("Response body: ", string(response_body))
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
