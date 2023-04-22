// provider.go

package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"client_id": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc(ENV_CLIENT_ID, nil),
			},
			"client_secret": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc(ENV_CLIENT_SECRET, nil),
			},
			"audience": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc(ENV_AUTH_AUDIENCE, nil),
			},
			"auth_url": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc(ENV_AUTH_URL, nil),
			},
			"api_url": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc(ENV_API_URL, nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"mattr_did":        resourceDid(),
			"mattr_webhook":    resourceWebhook(),
			"mattr_issuer":     resourceIssuer(),
			"mattr_credential": resourceCredentialConfig(),
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
