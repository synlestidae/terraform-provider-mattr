package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"math"
	"nz.antunovic/mattr-terraform-provider/api"
	"reflect"
	"testing"
)

func AssertEqual(t *testing.T, expected, actual interface{}, msg string) {
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("%s. Expected: %v, but got: %v", msg, expected, actual)
	}
}

func runCreate(t *testing.T, resource *schema.Resource, data interface{}, client api.Client) *schema.ResourceData {
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

}
