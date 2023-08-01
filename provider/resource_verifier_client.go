package provider

import (
	"fmt"
	"nz.antunovic/mattr-terraform-provider/api"
	"nz.antunovic/mattr-terraform-provider/generator"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceVerifierClient(client api.Client) *schema.Resource {
	verifierClientSchema := map[string]*schema.Schema{
		"verifier_id": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"name": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"redirect_uris": &schema.Schema{
			Type:     schema.TypeList,
			Required: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"response_types": &schema.Schema{
			Type:     schema.TypeList,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"grant_types": &schema.Schema{
			Type:     schema.TypeList,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"token_endpoint_auth_method": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"id_token_signed_response_alg": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"application_type": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"logo_uri": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"secret": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"openid_configuration_url": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"authorization_url": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
	}

	getPath := func(d *schema.ResourceData) (string, error) {
		if verifierId, ok := d.Get("verifier_id").(string); ok && len(verifierId) != 0 {
			return fmt.Sprintf("/ext/oidc/v1/verifiers/%s/clients", verifierId), nil
		}

		return "", fmt.Errorf("'verifier_id' field is required for verifier client and must be a string")
	}

	verifierClientModifyRes := func(body interface{}) (interface{}, error) {
		if bodyMap, ok := body.(map[string]interface{}); ok {
			delete(bodyMap, "verifierId")
			return bodyMap, nil
		}

		return nil, fmt.Errorf("Unexpected type for /ext/oidc/v1/issuers/{id}/clients request: %T", body)
	}

	generator := generator.Generator{
		GetPath:           getPath,
		Client:            client,
		Schema:            verifierClientSchema,
		ModifyRequestBody: verifierClientModifyRes,
	}
	resource := generator.GenResource()

	return &resource
}
