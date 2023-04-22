package main

import (
	"log"

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
				Optional: true,
			},
			"type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"additional_types": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem:     schema.TypeString,
			},
			"contexts": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem:     schema.TypeString,
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
			"claim_mappings": &schema.Schema{ //TODO this needs to be changed
				Type:     schema.TypeMap,
				Required: true,
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
			},
			"persist": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"revocable": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"claim_source_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"expires_in": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Resource{
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
	log.Println("Creating credential config")
	api := m.(ProviderConfig).Api
	config_request := fromTerraformCredentialConfig(d)
	config_response, err := api.PostCredentialConfig(&config_request)
	if err != nil {
		return err
	}
	processCredentialConfigData(config_response, d)
	d.SetId(config_response.Id)
	return nil
}

func resourceCredentialConfigRead(d *schema.ResourceData, m interface{}) error {
	log.Println("Reading credential config")
	id := d.Id()
	api := m.(ProviderConfig).Api
	config_response, err := api.GetCredentialConfig(id)
	if err != nil {
		return err
	}
	processCredentialConfigData(config_response, d)
	return nil
}

func resourceCredentialConfigUpdate(d *schema.ResourceData, m interface{}) error {
	log.Println("Updating credential config")
	api := m.(ProviderConfig).Api
	id := d.Id()
	config_request := fromTerraformCredentialConfig(d)
	config_response, err := api.PutCredentialConfig(id, &config_request)
	if err != nil {
		return err
	}
	processCredentialConfigData(config_response, d)
	return nil
}

func resourceCredentialConfigDelete(d *schema.ResourceData, m interface{}) error {
	api := m.(ProviderConfig).Api
	err := api.DeleteCredentialConfig(d.Id())
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func processCredentialConfigData(config *CredentialConfig, d *schema.ResourceData) {
	// issuerMap
	issuerMap := make(map[string]string, 3)
	issuerMap["name"] = config.Issuer.Name
	issuerMap["logo_url"] = config.Issuer.LogoUrl
	issuerMap["icon_url"] = config.Issuer.IconUrl

	// brandingMap
	brandingMap := make(map[string]string, 2)
	brandingMap["background_color"] = config.CredentialBranding.BackgroundColor
	brandingMap["watermark_image_url"] = config.CredentialBranding.WatermarkImageUrl

	// claimMap
	claimMap := make(map[string]map[string]any)
	for k, claim := range config.ClaimMappings {
		claimMap[k] = make(map[string]any)
		claimMap[k]["map_from"] = claim.MapFrom
		claimMap[k]["required"] = claim.Required
		claimMap[k]["default_value"] = claim.DefaultValue
	}

	// expiresInMap
	expiresIn := make(map[string]int)
	expiresIn["years"] = config.ExpiresIn.Years
	expiresIn["months"] = config.ExpiresIn.Months
	expiresIn["weeks"] = config.ExpiresIn.Weeks
	expiresIn["days"] = config.ExpiresIn.Days
	expiresIn["hours"] = config.ExpiresIn.Hours
	expiresIn["minutes"] = config.ExpiresIn.Minutes
	expiresIn["seconds"] = config.ExpiresIn.Seconds

	d.Set("name", config.Name)
	d.Set("description", config.Description)
	d.Set("type", config.Type)
	d.Set("additional_types", config.AdditionalTypes)
	d.Set("contexts", config.Contexts)

	d.Set("issuer", issuerMap)
	d.Set("credential_branding", brandingMap)
	d.Set("claim_mappings", claimMap)
	d.Set("expires_in", expiresIn)

	d.Set("persist", config.Persist)
	d.Set("revocable", config.Revocable)
	d.Set("claim_source_id", config.ClaimSourceId)
}

func fromTerraformCredentialConfig(d *schema.ResourceData) CredentialConfig {
	configIssuerMap := d.Get("issuer").(map[string]string)
	configBrandingMap := d.Get("credential_branding").(map[string]string)
	claimMappingsMap := d.Get("claim_mappingss").(map[string]map[string]any)
	configExpiresInMap := d.Get("expires_in").(map[string]int)

	configIssuer := IssuerInfo{
		Name:    configIssuerMap["name"],
		LogoUrl: configIssuerMap["logo_url"],
		IconUrl: configIssuerMap["name"],
	}

	configBranding := CredentialBranding{
		BackgroundColor:   configBrandingMap["background_color"],
		WatermarkImageUrl: configBrandingMap["watermark_image_url"],
	}

	claimMappings := make(map[string]ClaimMappingConfig, len(claimMappingsMap))
	for k, claim := range claimMappingsMap {
		claimMappings[k] = ClaimMappingConfig{
			MapFrom:      claim["map_from"].(string),
			Required:     claim["required"].(bool),
			DefaultValue: claim["default_value"].(string),
		}
	}

	configExpiresIn := ExpiresIn{
		Years:   configExpiresInMap["years"],
		Months:  configExpiresInMap["months"],
		Weeks:   configExpiresInMap["weeks"],
		Days:    configExpiresInMap["days"],
		Hours:   configExpiresInMap["hours"],
		Minutes: configExpiresInMap["minutes"],
		Seconds: configExpiresInMap["seconds"],
	}

	return CredentialConfig{
		Name:               d.Get("name").(string),
		Description:        d.Get("description").(string),
		Type:               d.Get("type").(string),
		AdditionalTypes:    d.Get("additional_types").([]string),
		Contexts:           d.Get("contexts").([]string),
		Issuer:             configIssuer,
		CredentialBranding: configBranding,
		ClaimMappings:      claimMappings,
		Persist:            d.Get("persist").(bool),
		Revocable:          d.Get("revocable").(bool),
		ClaimSourceId:      d.Get("claim_source_id").(string),
		ExpiresIn:          configExpiresIn,
	}

	panic("Not implemented")
}
