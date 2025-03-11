package provider

import (
	"fmt"

	"nz.antunovic/mattr-terraform-provider/api"
	"nz.antunovic/mattr-terraform-provider/generator"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceIssuerClient() *schema.Resource {
	issuerClientSchema := map[string]*schema.Schema{
		"issuer_id": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"name": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"redirect_uris": &schema.Schema{
			Type:     schema.TypeSet,
			Required: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"response_types": &schema.Schema{
			Type:     schema.TypeSet,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"grant_types": &schema.Schema{
			Type:     schema.TypeSet,
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
			Required: true,
		},
	}

	getPath := func(d *schema.ResourceData) (string, error) {
		if issuerId, ok := d.Get("issuer_id").(string); ok {
			return fmt.Sprintf("/ext/oidc/v1/issuers/%s/clients", issuerId), nil
		}

		return "", fmt.Errorf("'issuer_id' field is required for issuer client and must be a string")
	}

	generator := generator.Generator{
		GetPath: getPath,
		Client:  &api.HttpClient{},
		Schema:  issuerClientSchema,
	}

	resource := generator.GenResource()
	return &resource
}
