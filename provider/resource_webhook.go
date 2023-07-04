package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	_ "github.com/motemen/go-loghttp/global"
	"nz.antunovic/mattr-terraform-provider/api"
	"nz.antunovic/mattr-terraform-provider/generator"
)

func resourceWebhook() *schema.Resource {
	generator := generator.Generator{
		Path:   "/core/v1/webhooks",
		Client: &api.HttpClient{},
		Schema: map[string]*schema.Schema{
			"events": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Types of events we will look out for and send to the webhook",
				Required:    true,
			},
			"url": &schema.Schema{
				Type:        schema.TypeString,
				Description: "URL of the webhook, to which event payloads are delivered",
				Required:    true,
			},
			"disabled": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "If true, the webhook is disabled.",
			},
		},
	}

	resource := generator.GenResource()

	return &resource
}
