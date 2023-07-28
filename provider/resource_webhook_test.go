package provider

import (
	"testing"
)

func TestResourceWebhookCreate(t *testing.T) {
	client := TestClient{
		responses: map[string]interface{}{
			"POST https://test.api/core/v1/webhooks": map[string]interface{}{
				"id":       "8e485582-6ef6-49bc-80fa-25a1b36a8322",
				"events":   []interface{}{"OidcIssuerCredentialIssued"},
				"url":      "https://test.api/webhook",
				"disabled": false,
			},
		},
	}

	// Test Create
	createData := map[string]interface{}{
		"events":   []interface{}{"OidcIssuerCredentialIssued"},
		"url":      "https://test.api/webhook",
		"disabled": false,
	}

	resource := resourceWebhook(&client)
	resourceData := runCreate(t, resource, &createData, &client)

	AssertEqual(t, createData["events"], resourceData.Get("events"), "Events should match")
	AssertEqual(t, createData["url"], resourceData.Get("url"), "URL should match")
	AssertEqual(t, createData["disabled"], resourceData.Get("disabled"), "Disabled should match")
}

func TestResourceWebhookCreateExtraneousFields(t *testing.T) {
	client := TestClient{
		responses: map[string]interface{}{
			"POST https://test.api/core/v1/webhooks": map[string]interface{}{
				"id":          "8e485582-6ef6-49bc-80fa-25a1b36a8322",
				"events":      []interface{}{"OidcIssuerCredentialIssued"},
				"url":         "https://test.api/webhook",
				"disabled":    false,
				"extraneous":  true,
				"unnecessary": "absolutely",
			},
		},
	}
	createData := map[string]interface{}{
		"events":   []interface{}{"OidcIssuerCredentialIssued"},
		"url":      "https://test.api/webhook",
		"disabled": false,
	}

	resource := resourceWebhook(&client)
	resourceData := runCreate(t, resource, createData, &client)

	AssertEqual(t, createData["events"], resourceData.Get("events"), "Events should match")
	AssertEqual(t, createData["url"], resourceData.Get("url"), "URL should match")
	AssertEqual(t, createData["disabled"], resourceData.Get("disabled"), "Disabled should match")
}

/*func runCreate(t *testing.T, resource *schema.Resource, data interface{}, client api.Client) *schema.ResourceData {
	provider_config := api.ProviderConfig{
		Api: api.Api{
			ClientId:             "test-id",
			ClientSecret:         "test-scret",
			Audience:             "test",
			AuthUrl:              "https://test.api/auth",
			ApiUrl:               "https://test.api",
			AccessToken:          "test-token",
			AccessTokenExpiresAt: math.MaxInt,
		},
	}

	// Test Create
	createData := map[string]interface{}{
		"events":   []interface{}{"OidcIssuerCredentialIssued"},
		"url":      "https://test.api/webhook",
		"disabled": false,
	}

	createCtx := schema.TestResourceDataRaw(t, resource.Schema, createData)

	// Set the createData in the resource context
	for k, v := range createData {
		createCtx.Set(k, v)
	}

	err := resource.Create(createCtx, provider_config)

	// Assert that Create succeeded without errors
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	// Assert that the resource has an ID after creation
	if id := createCtx.Id(); id == "" {
		t.Fatal("Create should set the ID")
	}

	return createCtx

}*/
