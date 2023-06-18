package provider

import (
	"log"
	"strconv"
	"nz.antunovic/mattr-terraform-provider/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAuthentication() *schema.Resource {
	return &schema.Resource{
		Create: resourceAuthenticationCreate,
		Read:   resourceAuthenticationRead,
		Update: resourceAuthenticationUpdate,
		Delete: resourceAuthenticationDelete,
		Schema: map[string]*schema.Schema{
			"url": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"scope": &schema.Schema{
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			"client_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"client_secret": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"token_endpoint_auth_method": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"claims_source": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"static_request_parameters": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"forwarded_request_parameters": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"claims_to_sync": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceAuthenticationCreate(d *schema.ResourceData, m interface{}) error {
	log.Println("Creating authentication provider")
	api := m.(ProviderConfig).Api
	authentication_request, err := fromTerraformAuthentication(d)
	if err != nil {
		return err
	}
	authentication_response, err := api.PostAuthenticationProvider(authentication_request)
	if err != nil {
		return err
	}
	if err = processAuthenticationData(authentication_response, d); err != nil {
		return err
	}
	return nil
}

func resourceAuthenticationRead(d *schema.ResourceData, m interface{}) error {
	log.Println("Reading authentication provider")
	id := d.Id()
	api := m.(ProviderConfig).Api
	authentication_response, err := api.GetAuthenticationProvider(id)
	if err != nil {
		return err
	}
	if err = processAuthenticationData(authentication_response, d); err != nil {
		return err
	}
	return nil
}

func resourceAuthenticationUpdate(d *schema.ResourceData, m interface{}) error {
	log.Println("Updating authentication provider")
	api := m.(ProviderConfig).Api
	authentication_request, err := fromTerraformAuthentication(d)
	if err != nil {
		return err
	}
	authentication_response, err := api.PutAuthenticationProvider(authentication_request)
	if err != nil {
		return err
	}
	if err = processAuthenticationData(authentication_response, d); err != nil {
		return err
	}
	return nil
}

func resourceAuthenticationDelete(d *schema.ResourceData, m interface{}) error {
	log.Println("Deleting authentication provider")
	id := d.Id()
	api := m.(ProviderConfig).Api
	err := api.DeleteAuthenticationProvider(id)
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func processAuthenticationData(authentication *api.AuthenticationProvider, d *schema.ResourceData) error {
	d.SetId(authentication.Id)

	var err error
	if err = d.Set("url", authentication.Url); err != nil {
		return err
	}
	if err = d.Set("scope", authentication.Scope); err != nil {
		return err
	}
	if err = d.Set("client_id", authentication.ClientId); err != nil {
		return err
	}
	if err = d.Set("client_secret", authentication.ClientSecret); err != nil {
		return err
	}
	if err = d.Set("token_endpoint_auth_method", authentication.TokenEndpointAuthMethod); err != nil {
		return err
	}
	if err = d.Set("claims_source", authentication.ClaimsSource); err != nil {
		return err
	}
	if err = d.Set("static_request_parameters", authentication.StaticRequestParameters); err != nil {
		return err
	}
	if err = d.Set("forwarded_request_parameters", authentication.ForwardedRequestParameters); err != nil {
		return err
	}
	if err = d.Set("claims_to_sync", authentication.ClaimsToSync); err != nil {
		return err
	}

	return nil
}

func fromTerraformAuthentication(d *schema.ResourceData) (*api.AuthenticationProvider, error) {
	staticRequestParameters := d.Get("static_request_parameters").(map[string]any)

	forwardedParams := d.Get("forwarded_request_parameters").([]interface{})
	forwardedRequestParameters := castToStringSlice(forwardedParams)

	// TODO fix crappy conversion here
	if max_age_string, ok := staticRequestParameters["max_age"].(string); ok {
		if max_age_int, err := strconv.Atoi(max_age_string); err == nil {
			staticRequestParameters["max_age"] = max_age_int
		}
	}

	return &api.AuthenticationProvider{
		Id:                         d.Id(),
		RedirectUrl:                castToString(d.Get("redirect_url")),
		Url:                        castToString(d.Get("url")),
		Scope:                      castToStringSlice(d.Get("scope")),
		ClientId:                   castToString(d.Get("client_id")),
		ClientSecret:               castToString(d.Get("client_secret")),
		TokenEndpointAuthMethod:    castToString(d.Get("token_endpoint_auth_method")),
		ClaimsSource:               castToString(d.Get("claims_source")),
		StaticRequestParameters:    staticRequestParameters,
		ForwardedRequestParameters: forwardedRequestParameters,
		ClaimsToSync:               castToStringSlice(d.Get("claims_to_sync")),
	}, nil
}

func castToStringMap(stringMap map[string]interface{}) map[string]string {
	newStringMap := make(map[string]string, len(stringMap))
	for k, v := range stringMap {
		newStringMap[k] = v.(string)
	}
	return newStringMap
}
