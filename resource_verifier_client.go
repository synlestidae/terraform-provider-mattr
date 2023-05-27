package main

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceVerifierClient() *schema.Resource {
	return &schema.Resource{
		Create: resourceVerifierClientCreate,
		Read:   resourceVerifierClientRead,
		Update: resourceVerifierClientUpdate,
		Delete: resourceVerifierClientDelete,
		Schema: map[string]*schema.Schema{
			"verifier_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
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
			"logo_url": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceVerifierClientCreate(d *schema.ResourceData, m interface{}) error {
	log.Println("Creating verifier client")

	api := m.(ProviderConfig).Api
	verifier_client_request, err := fromTerraformVerifierClient(d)
	if err != nil {
		return err
	}

	verifier_id, ok := d.Get("verifier_id").(string)
	if !ok {
		return fmt.Errorf("verifier_id is required for creating verifier client")
	}

	verifier_client_response, err := api.PostVerifierClient(verifier_id, verifier_client_request)
	if err != nil {
		return err
	}
	if err = processVerifierClientData(verifier_client_response, d); err != nil {
		return err
	}
	log.Println("Created verifier client")
	return nil
}

func resourceVerifierClientRead(d *schema.ResourceData, m interface{}) error {
	log.Println("Reading credential config")

	id := d.Id()
	api := m.(ProviderConfig).Api

	verifier_id, ok := d.Get("verifier_id").(string)
	if !ok {
		return fmt.Errorf("verifier_id is required for creating verifier client")
	}

	verifier_client_response, err := api.GetVerifierClient(verifier_id, id)

	if err != nil {
		return err
	}

	if err = processVerifierClientData(verifier_client_response, d); err != nil {
		return err
	}
	log.Println("Read credential config")
	return nil
}

func resourceVerifierClientUpdate(d *schema.ResourceData, m interface{}) error {
	log.Println("Updating verifier client")

	verifier_client_request, err := fromTerraformVerifierClient(d)
	if err != nil {
		return err
	}

	id := d.Id()
	api := m.(ProviderConfig).Api

	verifier_id, ok := d.Get("verifier_id").(string)
	if !ok {
		return fmt.Errorf("verifier_id is required for creating verifier client")
	}

	verifier_client_response, err := api.PutVerifierClient(verifier_id, id, verifier_client_request)
	if err != nil {
		return err
	}

	if err = processVerifierClientData(verifier_client_response, d); err != nil {
		return err
	}
	log.Println("Updated credential config")
	return nil
}

func resourceVerifierClientDelete(d *schema.ResourceData, m interface{}) error {
	log.Println("Deleting verifier client")

	api := m.(ProviderConfig).Api
	id := d.Id()

	verifier_id, ok := d.Get("verifier_id").(string)
	if !ok {
		return fmt.Errorf("verifier_id is required for creating verifier client")
	}

	err := api.DeleteVerifierClient(verifier_id, id)
	if err != nil {
		return err
	}

	log.Println("Deleted verifier client")
	return nil
}

func processVerifierClientData(verifier_client *VerifierClient, d *schema.ResourceData) error {
	log.Println("Converting verifier client from REST")

	d.SetId(verifier_client.Id)

	var err error
	if err = d.Set("name", verifier_client.Name); err != nil {
		return err
	}
	if err = d.Set("redirect_uris", verifier_client.RedirectUris); err != nil {
		return err
	}
	if err = d.Set("response_types", verifier_client.ResponseTypes); err != nil {
		return err
	}
	if err = d.Set("grant_types", verifier_client.GrantTypes); err != nil {
		return err
	}
	if err = d.Set("token_endpoint_auth_method", verifier_client.TokenEndpointAuthMethod); err != nil {
		return err
	}
	if err = d.Set("id_token_signed_response_alg", verifier_client.IdTokenSignedResponseAlg); err != nil {
		return err
	}
	if err = d.Set("application_type", verifier_client.ApplicationType); err != nil {
		return err
	}

	return nil
}

func fromTerraformVerifierClient(d *schema.ResourceData) (*VerifierClient, error) {
	return &VerifierClient{
		Name:                     d.Get("name").(string),
		RedirectUris:             castToStringSlice(d.Get("redirect_uris")),
		ResponseTypes:            castToStringSlice(d.Get("response_Types")),
		GrantTypes:               castToStringSlice(d.Get("grant_types")),
		TokenEndpointAuthMethod:  d.Get("token_endpoint_auth_method").(string),
		IdTokenSignedResponseAlg: d.Get("id_token_signed_response_alg").(string),
		ApplicationType:          d.Get("application_type").(string),
	}, nil
}
