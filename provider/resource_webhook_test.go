package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
	resourceData := runCreate(t, resource, createData, &client)

	AssertEqual(t, createData["events"], resourceData.Get("events").(*schema.Set).List(), "Events should match")
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

	AssertEqual(t, createData["events"], resourceData.Get("events").(*schema.Set).List(), "Events should match")
	AssertEqual(t, createData["url"], resourceData.Get("url"), "URL should match")
	AssertEqual(t, createData["disabled"], resourceData.Get("disabled"), "Disabled should match")
}
