package main

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceVerifier() *schema.Resource {
	return &schema.Resource{
		Create: resourceVerifierCreate,
		Read:   resourceVerifierRead,
		Update: resourceVerifierUpdate,
		Delete: resourceVerifierDelete,
		Schema: map[string]*schema.Schema{
			"verifier_did": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"presentation_template_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"claim_mapping": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"json_lq_fqn": &schema.Schema{
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
			"include_presentation": &schema.Schema{
				Type:     schema.TypeBool,
				Required: true,
			},
		},
	}
}

func resourceVerifierCreate(d *schema.ResourceData, m interface{}) error {
	log.Println("Creating verifier")

	api := m.(ProviderConfig).Api
	verifier_request, err := fromTerraformVerifier(d)
	if err != nil {
		return err
	}

	verifier_response, err := api.PostVerifier(verifier_request)
	if err != nil {
		return err
	}
	if err = processVerifierData(verifier_response, d); err != nil {
		return err
	}
	log.Println("Created verifier")
	return nil
}

func resourceVerifierRead(d *schema.ResourceData, m interface{}) error {
	log.Println("Reading credential config")

	id := d.Id()
	api := m.(ProviderConfig).Api
	verifier_response, err := api.GetVerifier(id)

	if err != nil {
		return err
	}
	if err = processVerifierData(verifier_response, d); err != nil {
		return err
	}
	log.Println("Read credential config")
	return nil
}

func resourceVerifierUpdate(d *schema.ResourceData, m interface{}) error {
	log.Println("Updating verifier")

	verifier_request, err := fromTerraformVerifier(d)
	if err != nil {
		return err
	}

	id := d.Id()
	api := m.(ProviderConfig).Api
	verifier_response, err := api.PutVerifier(id, verifier_request)
	if err != nil {
		return err
	}

	if err = processVerifierData(verifier_response, d); err != nil {
		return err
	}
	log.Println("Updated credential config")
	return nil
}

func resourceVerifierDelete(d *schema.ResourceData, m interface{}) error {
	log.Println("Deleting verifier")

	api := m.(ProviderConfig).Api
	id := d.Id()
	err := api.DeleteVerifier(id)
	if err != nil {
		return err
	}

	log.Println("Deleted verifier")
	return nil
}

func processVerifierData(verifier *Verifier, d *schema.ResourceData) error {
	log.Println("Converting verifier from REST")

	return nil
}

func fromTerraformVerifier(d *schema.ResourceData) (*Verifier, error) {
	// TODO
	return nil, nil
}
