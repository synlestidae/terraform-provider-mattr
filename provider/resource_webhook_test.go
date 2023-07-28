package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"math"
	"nz.antunovic/mattr-terraform-provider/api"
	"reflect"
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
	resource := resourceWebhook(&client)

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

	// Assert that the resource data matches the input data
	AssertEqual(t, createData["events"], createCtx.Get("events"), "Events should match")
	AssertEqual(t, createData["url"], createCtx.Get("url"), "URL should match")
	AssertEqual(t, createData["disabled"], createCtx.Get("disabled"), "Disabled should match")
}

// Helper function to compare values and fail the test if not equal

func AssertEqual(t *testing.T, expected, actual interface{}, msg string) {
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("%s. Expected: %v, but got: %v", msg, expected, actual)
	}
}
