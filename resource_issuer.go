package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceIssuer() *schema.Resource {
	return &schema.Resource{
		Create: resourceIssuerCreate,
		Read:   resourceIssuerRead,
		Update: resourceIssuerUpdate,
		Delete: resourceIssuerDelete,
		Schema: map[string]*schema.Schema{
			"issuer_did": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"issuer_logo_url": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"issuer_icon_url": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"context": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"type": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
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
			"federated_provider": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"url": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"scope": &schema.Schema{
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
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
					},
				},
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
			"claim_mappings": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"json_ld_term": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"oidc_claim": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func resourceIssuerCreate(d *schema.ResourceData, m interface{}) error {
	log.Println("Creating issuer")
	api := m.(ProviderConfig).Api
	issuer_request := fromTerraformIssuer(d)
	issuer_response, err := api.PostIssuer(&issuer_request)
	if err != nil {
		return err
	}
	processIssuerData(issuer_response, d)
	return nil
}

func resourceIssuerRead(d *schema.ResourceData, m interface{}) error {
	log.Println("Getting issuer")
	issuer_id := d.Id()
	api := m.(ProviderConfig).Api
	issuer_response, err := api.GetIssuer(issuer_id)
	if err != nil {
		return err
	}
	processIssuerData(issuer_response, d)
	return nil
}

func resourceIssuerUpdate(d *schema.ResourceData, m interface{}) error {
	log.Println("Updating issuer")
	issuer_id := d.Id()
	api := m.(ProviderConfig).Api
	issuer_request := fromTerraformIssuer(d)
	issuer_response, err := api.PutIssuer(issuer_id, &issuer_request)
	if err != nil {
		return err
	}
	processIssuerData(issuer_response, d)
	return nil
}

func resourceIssuerDelete(d *schema.ResourceData, m interface{}) error {
	log.Println("Deleting issuer")
	issuer_id := d.Id()
	api := m.(ProviderConfig).Api
	return api.DeleteIssuer(issuer_id)
}

func processIssuerData(issuerResponse *IssuerResponse, d *schema.ResourceData) error {
	log.Println("Processing issuer data")
	d.SetId(issuerResponse.Id)
	if err := d.Set("issuer_did", issuerResponse.Credential.IssuerDid); err != nil {
		return fmt.Errorf("error setting 'issuer_did' field: %s", err)
	}
	if err := d.Set("issuer_logo_url", issuerResponse.Credential.IssuerLogoUrl); err != nil {
		return fmt.Errorf("error setting 'issuer_logo_url' field: %s", err)
	}
	if err := d.Set("issuer_icon_url", issuerResponse.Credential.IssuerIconUrl); err != nil {
		return fmt.Errorf("error setting 'issuer_icon_url' field: %s", err)
	}
	if err := d.Set("name", issuerResponse.Credential.Name); err != nil {
		return fmt.Errorf("error setting 'name' field: %s", err)
	}
	if err := d.Set("description", issuerResponse.Credential.Description); err != nil {
		return fmt.Errorf("error setting 'description' field: %s", err)
	}
	if err := d.Set("context", issuerResponse.Credential.Context); err != nil {
		return fmt.Errorf("error setting 'context' field: %s", err)
	}
	if err := d.Set("type", issuerResponse.Credential.Type); err != nil {
		return fmt.Errorf("error setting 'type' field: %s", err)
	}
	if err := d.Set("credential_branding", flattenCredentialBranding(issuerResponse.Credential.CredentialBranding)); err != nil {
		return fmt.Errorf("error setting 'credential_branding' field: %s", err)
	}
	if err := d.Set("federated_provider", flattenFederatedProvider(issuerResponse.FederatedProvider)); err != nil {
		return fmt.Errorf("error setting 'federated_provider' field: %s", err)
	}
	if err := d.Set("claim_mappings", flattenClaimMappings(issuerResponse.ClaimMappings)); err != nil {
		return fmt.Errorf("error setting 'claim_mappings' field: %s", err)
	}

	static_request_parameters := issuerResponse.StaticRequestParameters
	converted_params := make(map[string]string, len(static_request_parameters))
	for k, v := range static_request_parameters {
		switch v.(type) {
		case int:
			converted_params[k] = strconv.Itoa(v.(int))
		case bool:
			converted_params[k] = strconv.FormatBool(v.(bool))
		case float32, float64:
			converted_params[k] = fmt.Sprintf("%g", v)
		default:
			converted_params[k] = fmt.Sprintf("%v", v)
		}
	}

	if err := d.Set("static_request_parameters", &converted_params); err != nil {
		return fmt.Errorf("error setting 'static_request_parameters' field: %s", err)
	}
	if err := d.Set("forwarded_request_parameters", &issuerResponse.ForwardedRequestParameters); err != nil {
		return fmt.Errorf("error setting 'forwarded_request_parameters' field: %s", err)
	}
	return nil
}

func fromTerraformIssuer(d *schema.ResourceData) IssuerRequest {
	log.Println("Preparing issuer data")

	credentialBranding := d.Get("credential_branding").(map[string]interface{})

	cred := IssuerCredential{
		IssuerDid:     castToString(d.Get("issuer_did")),
		IssuerLogoUrl: castToString(d.Get("issuer_logo_url")),
		IssuerIconUrl: castToString(d.Get("issuer_icon_url")),
		Name:          castToString(d.Get("name")),
		Description:   castToString(d.Get("description")),
		Context:       castToStringSlice(d.Get("context")),
		Type:          castToStringSlice(d.Get("type")),
		CredentialBranding: &CredentialBranding{
			BackgroundColor:   castToString(credentialBranding["background_color"]),
			WatermarkImageUrl: castToString(credentialBranding["watermark_image_url"]),
		},
	}

	federated_provider_data := d.Get("federated_provider").(map[string]interface{})
	fedProv := FederatedProvider{
		Url:                     castToString(federated_provider_data["url"]),
		Scope:                   splitList(federated_provider_data["scope"]),
		ClientId:                castToString(federated_provider_data["client_id"]),
		ClientSecret:            castToString(federated_provider_data["client_secret"]),
		TokenEndpointAuthMethod: castToString(federated_provider_data["token_endpoint_auth_method"]),
		ClaimsSource:            castToString(federated_provider_data["claims_source"]),
	}

	var claimMappings []ClaimMapping
	for _, c := range d.Get("claim_mappings").([]interface{}) {
		cm := c.(map[string]interface{})
		claimMappings = append(claimMappings, ClaimMapping{
			JsonLdTerm: cm["json_ld_term"].(string),
			OidcClaim:  cm["oidc_claim"].(string),
		})
	}

	staticRequestParameters := d.Get("static_request_parameters").(map[string]interface{})

	// TODO fix crappy conversion here
	if max_age_string, ok := staticRequestParameters["max_age"].(string); ok {
		if max_age_int, err := strconv.Atoi(max_age_string); err == nil {
			staticRequestParameters["max_age"] = max_age_int
		}
	}

	log.Println("Prepared issuer data")

	return IssuerRequest{
		Credential:                 &cred,
		FederatedProvider:          &fedProv,
		ClaimMappings:              claimMappings,
		StaticRequestParameters:    staticRequestParameters,
		ForwardedRequestParameters: castToStringSlice(d.Get("forwarded_request_parameters")),
	}
}

func flattenCredentialBranding(credentialBranding *CredentialBranding) map[string]string {
	log.Println("Brand", *credentialBranding)
	branding := make(map[string]string, 2)
	branding["background_color"] = credentialBranding.BackgroundColor
	branding["watermark_image_url"] = credentialBranding.WatermarkImageUrl
	return branding
}

func flattenFederatedProvider(federatedProvider *FederatedProvider) map[string]interface{} {
	fpProvider := make(map[string]interface{}, 6)
	fpProvider["url"] = federatedProvider.Url
	fpProvider["scope"] = strings.Join(federatedProvider.Scope, " ")
	fpProvider["client_id"] = federatedProvider.ClientId
	fpProvider["client_secret"] = federatedProvider.ClientSecret
	fpProvider["token_endpoint_auth_method"] = federatedProvider.TokenEndpointAuthMethod
	fpProvider["claims_source"] = federatedProvider.ClaimsSource
	return fpProvider
}

func flattenClaimMappings(claimMappings []ClaimMapping) []map[string]string {
	mappings := make([]map[string]string, len(claimMappings))
	for i, mapping := range claimMappings {
		mappings[i] = map[string]string{
			"json_ld_term": mapping.JsonLdTerm,
			"oidc_claim":   mapping.OidcClaim,
		}
	}
	return mappings
}

func splitList(val interface{}) []string {
	if val == nil {
		return []string{}
	}
	return strings.Split(val.(string), " ")
}

func castToString(val interface{}) string {
	if val == nil {
		return ""
	}
	return val.(string)
}

func castToStringSlice(val interface{}) []string {
	if val == nil {
		return make([]string, 0)
	}
	interfaceSlice := val.([]interface{})
	stringSlice := make([]string, len(interfaceSlice))
	for i, s := range interfaceSlice {
		stringSlice[i] = s.(string)
	}
	return stringSlice
}
