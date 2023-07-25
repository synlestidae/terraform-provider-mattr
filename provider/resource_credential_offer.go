package provider

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	_ "github.com/motemen/go-loghttp/global"
	"nz.antunovic/mattr-terraform-provider/api"
	"nz.antunovic/mattr-terraform-provider/generator"
)

func credOfferModifyRes(res interface{}) (interface{}, error) {
	resMap, ok := res.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("Unexpected type for /core/v1/openid/offers: %T", resMap)
	}

	resMap["id"] = resMap["uri"]

	return res, nil
}

func resourceCredentialOffer() *schema.Resource {
	generator := generator.Generator{
		Path:               "/core/v1/openid/offers",
		Client:             &api.HttpClient{},
		Immutable:          true,
		ModifyResponseBody: credOfferModifyRes,
		Schema: map[string]*schema.Schema{
			"credentials": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "List of IDs of credential configurations",
				Required:    true,
				ForceNew:    true,
			},
			"uri": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}

	resource := generator.GenResource()
	resource.Read = func(*schema.ResourceData, interface{}) error {
		return nil
	}
	resource.Delete = func(*schema.ResourceData, interface{}) error {
		return nil
	}

	return &resource
}
