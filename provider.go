package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func Provider() *schema.Provider {
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
		},
		ResourcesMap: map[string]*schema.Resource{
			"mattr_did":                     resourceDid(),
			"mattr_webhook":                 resourceWebhook(),
			"mattr_issuer":                  resourceIssuer(),
			"mattr_credential":              resourceCredentialConfig(),
			"mattr_claim_source":            resourceClaimSource(),
			"mattr_authentication_provider": resourceAuthentication(),
			"mattr_issuer_client":           resourceIssuerClient(),
			"mattr_verifier":                resourceVerifier(),
			"mattr_verifier_client":         resourceVerifierClient(),
		},
		ConfigureFunc: ProviderConfigure,
	}
}

type ProviderConfig struct {
	Api Api
}

func ProviderConfigure(d *schema.ResourceData) (interface{}, error) {
	api := Api{
		ClientId:     getOrEmpty(d, "client_id"),
		ClientSecret: getOrEmpty(d, "client_secret"),
		Audience:     getOrEmpty(d, "audience"),
		AuthUrl:      getOrEmpty(d, "auth_url"),
		ApiUrl:       getOrEmpty(d, "api_url"),
	}

	config := ProviderConfig{
		Api: api,
	}

	return config, nil
}

func getOrEmpty(d *schema.ResourceData, key string) string {
	if value, ok := d.Get(key).(string); ok {
		return value
	}
	return ""
}
