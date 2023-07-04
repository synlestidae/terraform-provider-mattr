package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
	"nz.antunovic/mattr-terraform-provider/api"
	"strconv"
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
			"secret": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
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
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"contexts": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"issuer": &schema.Schema{
				Type:     schema.TypeMap,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeMap,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"name": &schema.Schema{
								Type:     schema.TypeString,
								Required: true,
							},
							"logo_url": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"icon_url": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
						},
					},
				},
			},
			"credential_branding": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeMap,
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
			},
			"claim_mapping": &schema.Schema{ //TODO this needs to be changed
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
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
				Elem: &schema.Schema{
					Type: schema.TypeMap,
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
		},
	}
}

func resourceCredentialConfigCreate(d *schema.ResourceData, m interface{}) error {
	log.Println("Creating credential config")
	api := m.(api.ProviderConfig).Api
	config_request := fromTerraformCredentialConfig(d)
	config_response, err := api.PostCredentialConfig(&config_request)
	if err != nil {
		return err
	}
	processCredentialConfigData(config_response, d)
	log.Println("Created credential config")
	return nil
}

func resourceCredentialConfigRead(d *schema.ResourceData, m interface{}) error {
	log.Println("Reading credential config")
	id := d.Id()
	api := m.(api.ProviderConfig).Api
	config_response, err := api.GetCredentialConfig(id)
	if err != nil {
		return err
	}
	processCredentialConfigData(config_response, d)
	log.Println("Read credential config")
	return nil
}

func resourceCredentialConfigUpdate(d *schema.ResourceData, m interface{}) error {
	log.Println("Updating credential config")
	api := m.(api.ProviderConfig).Api
	id := d.Id()
	config_request := fromTerraformCredentialConfig(d)
	config_response, err := api.PutCredentialConfig(id, &config_request)
	if err != nil {
		return err
	}
	processCredentialConfigData(config_response, d)
	log.Println("Updated credential config")
	return nil
}

func resourceCredentialConfigDelete(d *schema.ResourceData, m interface{}) error {
	log.Println("Deleting credential config")
	api := m.(api.ProviderConfig).Api
	err := api.DeleteCredentialConfig(d.Id())
	if err != nil {
		return err
	}
	log.Println("Deleted credential config")
	return nil
}

func processCredentialConfigData(config *api.CredentialConfig, d *schema.ResourceData) {
	log.Println("Processing credential config data")
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
	i := 0
	claimList := make([]map[string]any, len(config.ClaimMappings))
	for k, claim := range config.ClaimMappings {
		claimList[i] = make(map[string]any)
		claimList[i]["name"] = k
		claimList[i]["map_from"] = claim.MapFrom
		claimList[i]["required"] = claim.Required
		claimList[i]["default_value"] = claim.DefaultValue
		i++
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
	d.Set("claim_mapping", d.Get("claim_mapping")) // TODO see the cheat I'm using
	d.Set("expires_in", expiresIn)

	d.Set("persist", config.Persist)
	d.Set("revocable", config.Revocable)
	d.Set("claim_source_id", config.ClaimSourceId)

	d.SetId(config.Id)
	log.Println("Processed credential config")
}

func fromTerraformCredentialConfig(d *schema.ResourceData) api.CredentialConfig {
	log.Println("Converting credential config from REST")
	configIssuerMap := d.Get("issuer").(map[string]interface{})
	configBrandingMap := d.Get("credential_branding").(map[string]interface{})
	claimMappingsList := d.Get("claim_mapping").([]interface{})
	configExpiresInMap := d.Get("expires_in").(map[string]interface{})

	configIssuer := api.IssuerInfo{
		Name:    configIssuerMap["name"].(string),
		LogoUrl: configIssuerMap["logo_url"].(string),
		IconUrl: configIssuerMap["icon_url"].(string),
	}

	configBranding := api.CredentialBranding{
		BackgroundColor:   configBrandingMap["background_color"].(string),
		WatermarkImageUrl: configBrandingMap["watermark_image_url"].(string),
	}

	claimMappings := make(map[string]api.ClaimMappingConfig, len(claimMappingsList))
	for _, claim := range claimMappingsList {
		claimObj := claim.(map[string]any)
		claimMappings[claimObj["name"].(string)] = api.ClaimMappingConfig{
			MapFrom:      claimObj["map_from"].(string),
			Required:     claimObj["required"].(bool),
			DefaultValue: claimObj["default_value"].(string),
		}
	}

	var yearStr, monthStr, weekStr, dayStr, hourStr, minuteStr, secondStr = configExpiresInMap["years"].(string),
		configExpiresInMap["months"].(string),
		configExpiresInMap["weeks"].(string),
		configExpiresInMap["days"].(string),
		configExpiresInMap["hours"].(string),
		configExpiresInMap["minutes"].(string),
		configExpiresInMap["seconds"].(string)

	years, err := strconv.Atoi(yearStr)
	if err != nil {
		panic("Failed to convert years") // TODO error handling
	}
	months, err := strconv.Atoi(monthStr)
	if err != nil {
		panic("Failed to convert months") // TODO error handling
	}
	weeks, err := strconv.Atoi(weekStr)
	if err != nil {
		panic("Failed to convert weeks") // TODO error handling
	}
	days, err := strconv.Atoi(dayStr)
	if err != nil {
		panic("Failed to convert days") // TODO error handling
	}
	hours, err := strconv.Atoi(hourStr)
	if err != nil {
		panic("Failed to convert hours") // TODO error handling
	}
	minutes, err := strconv.Atoi(minuteStr)
	if err != nil {
		panic("Failed to convert minutes") // TODO error handling
	}
	seconds, err := strconv.Atoi(secondStr)
	if err != nil {
		panic("Failed to convert seconds") // TODO error handling
	}

	configExpiresIn := api.ExpiresIn{
		Years:   years,
		Months:  months,
		Weeks:   weeks,
		Days:    days,
		Hours:   hours,
		Minutes: minutes,
		Seconds: seconds,
	}

	log.Println("Converted from resource data")

	return api.CredentialConfig{
		Name:               d.Get("name").(string),
		Description:        d.Get("description").(string),
		Type:               d.Get("type").(string),
		AdditionalTypes:    castToStringSlice(d.Get("additional_types").([]interface{})),
		Contexts:           castToStringSlice(d.Get("contexts").([]interface{})),
		Issuer:             configIssuer,
		CredentialBranding: configBranding,
		ClaimMappings:      claimMappings,
		Persist:            d.Get("persist").(bool),
		Revocable:          d.Get("revocable").(bool),
		ClaimSourceId:      d.Get("claim_source_id").(string),
		ExpiresIn:          configExpiresIn,
	}
}
