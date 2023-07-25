package provider

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"nz.antunovic/mattr-terraform-provider/api"
	"nz.antunovic/mattr-terraform-provider/generator"
)

func resourceCustomDomain() *schema.Resource {
	schema := map[string]*schema.Schema{
		"name": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"logo_url": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"domain": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"homepage": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"verification_token": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"is_verified": &schema.Schema{
			Type:     schema.TypeBool,
			Computed: true,
		},
		"verified_at": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
	}

	path := "/core/v1/config/domain"

	modifyCustomDomainRes := func(response interface{}) (interface{}, error) {
		responseMap, ok := response.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("Unexpected type for %s: %T", path, response)
		}
		responseMap["id"] = responseMap["domain"]
		return responseMap, nil
	}

	custDomain := generator.Generator{
		Path:               path,
		Client:             &api.HttpClient{},
		Schema:             schema,
		Singleton:          true,
		ModifyResponseBody: modifyCustomDomainRes,
	}

	resource := custDomain.GenResource()
	return &resource
}
