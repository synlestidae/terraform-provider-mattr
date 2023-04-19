package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceCredentialConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceCredentialConfigCreate,
		Read:   resourceCredentialConfigRead,
		Update: resourceCredentialConfigUpdate,
		Delete: resourceCredentialConfigDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"additional_types": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"contexts": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"issuer": &schema.Schema{
				Type:     schema.TypeMap,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"icon_url": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"iconUrl": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"credential_branding": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"background_color": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"watermark_image_url": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"claim_mapping": &schema.Schema{
				Type: schema.TypeMap,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"map_from": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"default_value": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"required": &schema.Schema{
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
				Optional: true,
			},
			"persist": &schema.Schema{
				Type:     schema.TypeBool,
				Required: false,
			},
			"revocable": &schema.Schema{
				Type:     schema.TypeBool,
				Required: false,
			},
			"claim_source_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: false,
			},
			"expires_in": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
				Elem: schema.Resource{
					Schema: map[string]*schema.Schema{
						"years": &schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
						},
						"months": &schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
						},
						"weeks": &schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
						},
						"days": &schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
						},
						"hours": &schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
						},
						"minutes": &schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
						},
						"seconds": &schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceCredentialConfigCreate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceCredentialConfigRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceCredentialConfigUpdate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceCredentialConfigDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}
