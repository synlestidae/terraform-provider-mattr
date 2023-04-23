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
					Schema: map[string]*schema.Schema{},
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
	//api := m.(ProviderConfig).Api
	return nil
}

func resourceClaimSourceRead(d *schema.ResourceData, m interface{}) error {
	log.Println("Reading claim source")
	return nil
}

func resourceClaimSourceUpdate(d *schema.ResourceData, m interface{}) error {
	log.Println("Updating claim source")
	return nil
}

func resourceClaimSourceDelete(d *schema.ResourceData, m interface{}) error {
	log.Println("Deleting claim source")
	return nil
}

func processClaimSourceData(config *ClaimSource, d *schema.ResourceData) {
}

func fromTerraformClaimSource(d *schema.ResourceData) ClaimSource {
	panic("Not quite implemented yet")
}
