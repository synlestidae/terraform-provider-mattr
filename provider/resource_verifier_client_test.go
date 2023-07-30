package provider

import (
	"encoding/json"
	"testing"
)

func TestResourceVerifierClientCreate(t *testing.T) {
	jsonData := `{
	"id": "da9bb6e4-c9ae-4468-b6ac-72b90d6efd5d",
	"secret": "H2epdcmNJ46hXJo5opdzvhbZK9W2ZGPkQh.E",
	"name": "OIDC Client for the verifier",
	"redirectUris": [
		"https://example.com/callback"
	],
	"responseTypes": [
		"code"
	],
	"grantTypes": [
		"authorization_code"
	],
	"tokenEndpointAuthMethod": "client_secret_post",
	"idTokenSignedResponseAlg": "ES256",
	"applicationType": "web",
	"logoUri": "https://example.com/logo.png"
}`
	var data map[string]interface{}
	err := json.Unmarshal([]byte(jsonData), &data)

	if err != nil {
		t.Fatalf("Error in test data: %s", err)
	}

	client := TestClient{
		responses: map[string]interface{}{
			"POST https://test.api/ext/oidc/v1/verifiers/402c65eb-48e9-4a4c-b5e9-1ea615baccee/clients": data,
		},
	}

	// Test Create
	createData := map[string]interface{}{
		"name":                         "OIDC Client for the verifier",
		"verifier_id":                  "402c65eb-48e9-4a4c-b5e9-1ea615baccee",
		"redirect_uris":                []interface{}{"https://example.com/callback"},
		"response_types":               []interface{}{"code"},
		"grant_types":                  []interface{}{"authorization_code"},
		"token_endpoint_auth_method":   "client_secret_post",
		"id_token_signed_response_alg": "ES256",
		"application_type":             "web",
		"logo_uri":                     "https://example.com/logo.png",
	}

	resource := resourceVerifierClient(&client)
	resourceData := runCreate(t, resource, createData, &client)
	AssertEqual(t, "OIDC Client for the verifier", resourceData.Get("name"), "Name should match")
	AssertEqual(t, "402c65eb-48e9-4a4c-b5e9-1ea615baccee", resourceData.Get("verifier_id"), "Verifier ID should be correct")
	AssertEqual(t, "H2epdcmNJ46hXJo5opdzvhbZK9W2ZGPkQh.E", resourceData.Get("secret"), "Secret should be correct")
	AssertEqual(t, []interface{}{"https://example.com/callback"}, resourceData.Get("redirect_uris"), "Redirect URIs should be correct")
	AssertEqual(t, []interface{}{"code"}, resourceData.Get("response_types"), "Response types should be correct")
	AssertEqual(t, []interface{}{"authorization_code"}, resourceData.Get("grant_types"), "Grant types should be correct")
	AssertEqual(t, "client_secret_post", resourceData.Get("token_endpoint_auth_method"), "Auth method should be correct")
	AssertEqual(t, "ES256", resourceData.Get("id_token_signed_response_alg"), "Algorithm should be correct")
	AssertEqual(t, "web", resourceData.Get("application_type"), "Application type should be correct")
	AssertEqual(t, "https://example.com/logo.png", resourceData.Get("logo_uri"), "Logo should be correct")
}
