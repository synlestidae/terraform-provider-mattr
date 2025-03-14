package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"nz.antunovic/mattr-terraform-provider/api"
)

func Provider() *schema.Provider {
	client := api.HttpClient{}
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"client_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"client_secret": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"audience": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"auth_url": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"api_url": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"access_token": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"mattr_did":                                  resourceDid(),
			"mattr_webhook":                              resourceWebhook(&client),
			"mattr_issuer":                               resourceIssuer(),
			"mattr_credential_web":                       resourceCredentialConfig(),
			"mattr_claim_source":                         resourceClaimSource(),
			"mattr_authentication_provider":              resourceAuthentication(),
			"mattr_issuer_client":                        resourceIssuerClient(),
			"mattr_verifier":                             resourceVerifier(),
			"mattr_verifier_client":                      resourceVerifierClient(&client),
			"mattr_custom_domain":                        resourceCustomDomain(),
			"mattr_compact_credential_template":          resourceCompactCredentialTemplate(),
			"mattr_semantic_compact_credential_template": resourceSemanticCompactCredentialTemplate(),
			"mattr_credential_offer":                     resourceCredentialOffer(),
			"mattr_presentation":                         resourcePresentation(),
		},
		ConfigureFunc: ProviderConfigure,
	}
}

func ProviderConfigure(d *schema.ResourceData) (interface{}, error) {
	a := api.Api{
		ClientId:     getOrEmpty(d, "client_id"),
		ClientSecret: getOrEmpty(d, "client_secret"),
		Audience:     getOrEmpty(d, "audience"),
		AuthUrl:      getOrEmpty(d, "auth_url"),
		ApiUrl:       getOrEmpty(d, "api_url"),
		AccessToken:  getOrEmpty(d, "access_token"),
	}

	config := api.ProviderConfig{
		Api: a,
	}

	return config, nil
}

func getOrEmpty(d *schema.ResourceData, key string) string {
	if value, ok := d.Get(key).(string); ok {
		return value
	}
	return ""
}
