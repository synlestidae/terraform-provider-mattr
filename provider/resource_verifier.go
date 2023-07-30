package provider

import (
	"nz.antunovic/mattr-terraform-provider/api"
	"nz.antunovic/mattr-terraform-provider/generator"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceVerifier() *schema.Resource {
	verifierSchema := map[string]*schema.Schema{
		"verifier_did": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"presentation_template_id": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"claim_mapping": &schema.Schema{
			Type:     schema.TypeList,
			Required: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"json_ld_fqn": &schema.Schema{
						Type:     schema.TypeString,
						Required: true,
					},
					"oidc_claim": &schema.Schema{
						Type:     schema.TypeString,
						Required: true,
					},
				},
			},
		},
		"include_presentation": &schema.Schema{
			Type:     schema.TypeBool,
			Required: true,
		},
	}

	verifierGenerator := generator.Generator{
		Path:   "/ext/oidc/v1/verifiers",
		Client: &api.HttpClient{},
		Schema: verifierSchema,
	}

	provider := verifierGenerator.GenResource()
	return &provider
}

