package provider

import (
	"fmt"
	"reflect"
	"log"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"math"
	"nz.antunovic/mattr-terraform-provider/api"
	"testing"
)

type Request struct {
	url     string
	headers map[string]string
	body    interface{}
}

type TestClient struct {
	logs      []Request
	responses map[string]interface{}
}

func (client *TestClient) Post(url string, headers map[string]string, body interface{}) (interface{}, error) {
	return client.respond("POST", url, headers, body)
}

func (client *TestClient) Get(url string, headers map[string]string) (interface{}, error) {
	return client.respond("GET", url, headers, nil)
}

func (client *TestClient) Put(url string, headers map[string]string, body interface{}) (interface{}, error) {
	return client.respond("PUT", url, headers, body)
}

func (client *TestClient) Delete(url string, headers map[string]string) error {
	_, err := client.respond("DELETE", url, headers, nil)
	return err
}

func (client *TestClient) respond(method string, url string, headers map[string]string, body interface{}) (interface{}, error) {
	endpoint := fmt.Sprintf("%s %s", method, url)
	log.Printf("Locating response for %s", endpoint)
	response, ok := client.responses[endpoint]
	if !ok {
		log.Printf("Failed to find response for %s", endpoint)
		return nil, fmt.Errorf("Unable to find response for %s", endpoint)
	}
	log.Printf("Successfully located response for %s", endpoint)
	return response, nil
}

func TestResourceWebhookCreate(t *testing.T) {
	client := TestClient{
		responses: map[string]interface{}{
			"POST https://test.api/core/v1/webhooks": map[string]interface{} {
				"id": "8e485582-6ef6-49bc-80fa-25a1b36a8322",
				"events": []interface{}{ "OidcIssuerCredentialIssued" },
				"url": "https://test.api/webhook",
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
		"events": []interface{} { "OidcIssuerCredentialIssued" },
		"url": "https://test.api/webhook",
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
