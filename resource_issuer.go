package main

import (
	"fmt"

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
				Type:     schema.TypeList,
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
				Type:     schema.TypeList,
				Required: true,
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
	return nil
}

func resourceIssuerRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceIssuerUpdate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceIssuerDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}

func processIssuerData(issuerResponse *IssuerResponse, d *schema.ResourceData) error {
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
	if err := d.Set("credential_branding", flattenCredentialBranding(&issuerResponse.Credential.CredentialBranding)); err != nil {
		return fmt.Errorf("error setting 'credential_branding' field: %s", err)
	}
	if err := d.Set("federated_provider", flattenFederatedProvider(&issuerResponse.FederatedProvider)); err != nil {
		return fmt.Errorf("error setting 'federated_provider' field: %s", err)
	}
	if err := d.Set("claim_mappings", flattenClaimMappings(issuerResponse.ClaimMappings)); err != nil {
		return fmt.Errorf("error setting 'claim_mappings' field: %s", err)
	}
	if err := d.Set("static_request_parameters", &issuerResponse.StaticRequestParameters); err != nil {
		return fmt.Errorf("error setting 'static_request_parameters' field: %s", err)
	}
	if err := d.Set("forwarded_request_parameters", &issuerResponse.ForwardedRequestParameters); err != nil {
		return fmt.Errorf("error setting 'forwarded_request_parameters' field: %s", err)
	}
	return nil
}

func fromTerraformIssuer(d *schema.ResourceData) IssuerRequest {
	cred := IssuerCredential{
		IssuerDid:     d.Get("issuer_did").(string),
		IssuerLogoUrl: d.Get("issuer_logo_url").(string),
		IssuerIconUrl: d.Get("issuer_icon_url").(string),
		Name:          d.Get("name").(string),
		Description:   d.Get("description").(string),
		Context:       d.Get("context").([]string),
		Type:          d.Get("type").(string),
		CredentialBranding: CredentialBranding{
			BackgroundColor:   d.Get("background_color").(string),
			WatermarkImageUrl: d.Get("watermark_image_url").(string),
		},
	}

	fedProv := FederatedProvider{
		Url:                     d.Get("federated_provider_url").(string),
		Scope:                   d.Get("federated_provider_scope").([]string),
		ClientId:                d.Get("federated_provider_client_id").(string),
		ClientSecret:            d.Get("federated_provider_client_secret").(string),
		TokenEndpointAuthMethod: d.Get("federated_provider_token_endpoint_auth_method").(string),
		ClaimsSource:            d.Get("federated_provider_claims_source").(string),
	}

	var claimMappings []ClaimMapping
	for _, c := range d.Get("claim_mappings").([]interface{}) {
		cm := c.(map[string]interface{})
		claimMappings = append(claimMappings, ClaimMapping{
			JsonLdTerm: cm["json_ld_term"].(string),
			OidcClaim:  cm["oidc_claim"].(string),
		})
	}

	return IssuerRequest{
		Credential:                 cred,
		FederatedProvider:          fedProv,
		ClaimMappings:              claimMappings,
		StaticRequestParameters:    d.Get("static_request_parameters").(map[string]interface{}),
		ForwardedRequestParameters: d.Get("forwarded_request_parameters").([]string),
	}
}

func flattenCredentialBranding(credentialBranding *CredentialBranding) map[string]string {
	panic("Not implemented")
}

func flattenFederatedProvider(federatedProvider *FederatedProvider) map[string]string {
	panic("Not implemented")
}

func flattenClaimMappings(claimMappings []ClaimMapping) []map[string]string {
	panic("Not implemented")
}
