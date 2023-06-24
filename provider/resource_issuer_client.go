package provider

import (
	"log"
	"nz.antunovic/mattr-terraform-provider/api"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceIssuerClient() *schema.Resource {
	return &schema.Resource{
		Create: resourceIssuerClientCreate,
		Read:   resourceIssuerClientRead,
		Update: resourceIssuerClientUpdate,
		Delete: resourceIssuerClientDelete,
		Schema: map[string]*schema.Schema{
			"issuer_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"secret": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"redirect_uris": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"response_types": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"grant_types": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"token_endpoint_auth_method": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"id_token_signed_response_alg": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"application_type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceIssuerClientCreate(d *schema.ResourceData, m interface{}) error {
	log.Println("Creating issuer client")

	issuer_id := d.Get("issuer_id").(string)
	api := m.(api.ProviderConfig).Api
	issuer_client_request, err := fromTerraformIssuerClient(d)
	if err != nil {
		return err
	}

	issuer_client_response, err := api.PostIssuerClient(issuer_id, issuer_client_request)
	if err != nil {
		return err
	}
	if err = processIssuerClientData(issuer_client_response, d); err != nil {
		return err
	}
	log.Println("Created issuer client")
	return nil
}

func resourceIssuerClientRead(d *schema.ResourceData, m interface{}) error {
	log.Println("Reading credential config")

	id := d.Id()
	issuer_id := d.Get("issuer_id").(string)
	api := m.(api.ProviderConfig).Api
	issuer_client_response, err := api.GetIssuerClient(issuer_id, id)

	if err != nil {
		return err
	}
	if err = processIssuerClientData(issuer_client_response, d); err != nil {
		return err
	}
	log.Println("Read credential config")
	return nil
}

func resourceIssuerClientUpdate(d *schema.ResourceData, m interface{}) error {
	log.Println("Updating issuer client")

	issuer_client_request, err := fromTerraformIssuerClient(d)
	if err != nil {
		return err
	}

	id := d.Id()
	issuer_id := d.Get("issuer_id").(string)
	api := m.(api.ProviderConfig).Api
	issuer_client_response, err := api.PutIssuerClient(issuer_id, id, issuer_client_request)
	if err != nil {
		return err
	}

	if err = processIssuerClientData(issuer_client_response, d); err != nil {
		return err
	}
	log.Println("Updated credential config")
	return nil
}

func resourceIssuerClientDelete(d *schema.ResourceData, m interface{}) error {
	log.Println("Deleting issuer client")

	api := m.(api.ProviderConfig).Api
	issuer_id := d.Get("issuer_id").(string)
	id := d.Id()
	err := api.DeleteIssuerClient(issuer_id, id)
	if err != nil {
		return err
	}

	log.Println("Deleted issuer client")
	return nil
}

func processIssuerClientData(issuer_client *api.IssuerClientResponse, d *schema.ResourceData) error {
	log.Println("Converting issuer client from REST")

	var err error
	d.SetId(issuer_client.Id)
	if err = d.Set("name", issuer_client.Name); err != nil {
		return err
	}
	if err = d.Set("secret", issuer_client.Secret); err != nil {
		return err
	}
	if err = d.Set("redirect_uris", issuer_client.RedirectUris); err != nil {
		return err
	}
	if err = d.Set("response_types", issuer_client.ResponseTypes); err != nil {
		return err
	}
	if err = d.Set("grant_types", issuer_client.GrantTypes); err != nil {
		return err
	}
	if err = d.Set("token_endpoint_auth_method", issuer_client.TokenEndpointAuthMethod); err != nil {
		return err
	}
	if err = d.Set("id_token_signed_response_alg", issuer_client.IdTokenSignedResponseAlg); err != nil {
		return err
	}
	if err = d.Set("application_type", issuer_client.ApplicationType); err != nil {
		return err
	}

	return nil
}

func fromTerraformIssuerClient(d *schema.ResourceData) (*api.IssuerClientRequest, error) {
	return &api.IssuerClientRequest{
		Name:                     d.Get("name").(string),
		RedirectUris:             castToStringSlice(d.Get("redirect_uris")),
		ResponseTypes:            castToStringSlice(d.Get("response_types")),
		GrantTypes:               castToStringSlice(d.Get("grant_types")),
		TokenEndpointAuthMethod:  d.Get("token_endpoint_auth_method").(string),
		IdTokenSignedResponseAlg: d.Get("id_token_signed_response_alg").(string),
		ApplicationType:          d.Get("application_type").(string),
	}, nil
}
