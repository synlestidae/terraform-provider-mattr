package main

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceClaimSource() *schema.Resource {
	return &schema.Resource{
		Create: resourceClaimSourceCreate,
		Read:   resourceClaimSourceRead,
		Update: resourceClaimSourceUpdate,
		Delete: resourceClaimSourceDelete,
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"url": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"authorization": &schema.Schema{
				Type:     schema.TypeMap,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"value": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"request_parameter": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"property": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"map_from": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"default_value": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func resourceClaimSourceCreate(d *schema.ResourceData, m interface{}) error {
	log.Println("Creating claim source")
	api := m.(ProviderConfig).Api
	claimSource := fromTerraformClaimSource(d)
	createdClaimSource, err := api.PostClaimSource(&claimSource)
	if err != nil {
		return err
	}
	processClaimSourceData(createdClaimSource, d)
	return nil
}

func resourceClaimSourceRead(d *schema.ResourceData, m interface{}) error {
	log.Println("Reading claim source")
	api := m.(ProviderConfig).Api
	claimSource, err := api.GetClaimSource(d.Id())
	if err != nil {
		return err
	}
	processClaimSourceData(claimSource, d)
	return nil
}

func resourceClaimSourceUpdate(d *schema.ResourceData, m interface{}) error {
	log.Println("Updating claim source")
	api := m.(ProviderConfig).Api
	claimSource := fromTerraformClaimSource(d)
	updatedClaimSource, err := api.PutClaimSource(d.Id(), &claimSource)
	if err != nil {
		return err
	}
	processClaimSourceData(updatedClaimSource, d)
	return nil
}

func resourceClaimSourceDelete(d *schema.ResourceData, m interface{}) error {
	log.Println("Deleting claim source")
	api := m.(ProviderConfig).Api
	return api.DeleteClaimSource(d.Id())
}

func processClaimSourceData(claimSource *ClaimSource, d *schema.ResourceData) {
	authorization := make(map[string]string, 2)
	authorization["type"] = claimSource.Authorization.Type
	authorization["value"] = claimSource.Authorization.Value

	requestParameters := make([]map[string]string, len(claimSource.RequestParameters))
	i := 0
	for k, parameter := range claimSource.RequestParameters {
		requestParameters[i] = map[string]string{
			"name":          k,
			"map_from":      parameter.MapFrom,
			"default_value": parameter.DefaultValue,
		}
		i++
	}

	d.Set("name", claimSource.Name)
	d.Set("url", claimSource.Url)
	d.Set("authorization", authorization)
	d.Set("request_parameter", requestParameters)
	d.SetId(claimSource.Id)
}

func fromTerraformClaimSource(d *schema.ResourceData) ClaimSource {
	authorizationMap := d.Get("authorization").(map[string]interface{})
	requestParametersList := d.Get("request_parameter").([]interface{})

	authorization := ClaimSourceAuthorization{
		Type:  authorizationMap["type"].(string),
		Value: authorizationMap["value"].(string),
	}

	requestParametersMap := make(map[string]ClaimSourceRequestParameter, len(requestParametersList))
	for _, param := range requestParametersList {
		paramMap := param.(map[string]interface{})
		requestParametersMap[paramMap["property"].(string)] = ClaimSourceRequestParameter{
			MapFrom:      paramMap["map_from"].(string),
			DefaultValue: paramMap["default_value"].(string),
		}
	}

	return ClaimSource{
		Id:                d.Id(),
		Name:              d.Get("name").(string),
		Url:               d.Get("url").(string),
		Authorization:     authorization,
		RequestParameters: requestParametersMap,
	}
}
