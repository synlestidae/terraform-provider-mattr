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

	verifier_client_response, err := api.PostVerifierClient(verifier_client_request)
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
	verifier_client_response, err := api.GetVerifierClient(id)

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
	verifier_client_response, err := api.PutVerifierClient(id, verifier_client_request)
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
	err := api.DeleteVerifierClient(id)
	if err != nil {
		return err
	}

	log.Println("Deleted verifier client")
	return nil
}

func processVerifierClientData(verifier_client *VerifierClient, d *schema.ResourceData) error {
	return fmt.Errorf("Not quite implemented")
}

func fromTerraformVerifierClient(d *schema.ResourceData) (*VerifierClient, error) {
	return nil, fmt.Errorf("Not quite implemented")
}
