package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"nz.antunovic/mattr-terraform-provider/api"
	"nz.antunovic/mattr-terraform-provider/generator"
)

func resourceAuthentication() *schema.Resource {
	schema := map[string]*schema.Schema{
		"url": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"scope": &schema.Schema{
			Type:     schema.TypeList,
			Elem:     &schema.Schema{Type: schema.TypeString},
			Optional: true,
		},
		"client_id": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"client_secret": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"token_endpoint_auth_method": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"claims_source": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"static_request_parameters": &schema.Schema{
			Type:     schema.TypeMap,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"forwarded_request_parameters": &schema.Schema{
			Type:     schema.TypeList,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"claims_to_sync": &schema.Schema{
			Type:     schema.TypeList,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"redirect_url": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
	}

	providerGen := generator.Generator{
		Path:   "/core/v1/users/authenticationproviders",
		Schema: schema,
		Client: &api.HttpClient{},
	}

	provider := providerGen.GenResource()
	return &provider
}
