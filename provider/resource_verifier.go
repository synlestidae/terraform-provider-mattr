package provider

import (
	"log"
	"nz.antunovic/mattr-terraform-provider/api"

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
						"json_ld_fqn": &schema.Schema{
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

func processVerifierData(verifier *api.Verifier, d *schema.ResourceData) error {
	log.Println("Converting verifier from REST")

	d.SetId(verifier.Id)

	if err := d.Set("verifier_did", verifier.VerifierDid); err != nil {
		return err
	}

	if err := d.Set("presentation_template_id", verifier.PresentationTemplateId); err != nil {
		return err
	}

	if err := d.Set("claim_mapping", flattenVerifierClaimMappings(verifier.ClaimMappings)); err != nil {
		return err
	}

	if err := d.Set("include_presentation", verifier.IncludePresentation); err != nil {
		return err
	}

	return nil
}

func fromTerraformVerifier(d *schema.ResourceData) (*api.Verifier, error) {
	var claimMappings []api.VerifierClaimMapping
	for _, c := range d.Get("claim_mapping").([]interface{}) {
		cm := c.(map[string]interface{})

		claimMappings = append(claimMappings, api.VerifierClaimMapping{
			JsonLdFqn: cm["json_ld_fqn"].(string),
			OidcClaim: cm["oidc_claim"].(string),
		})
	}

	return &api.Verifier{
		Id:                     d.Id(),
		VerifierDid:            d.Get("verifier_did").(string),
		PresentationTemplateId: d.Get("presentation_template_id").(string),
		ClaimMappings:          claimMappings,
		IncludePresentation:    d.Get("include_presentation").(bool),
	}, nil
}

func flattenVerifierClaimMappings(claimMappings []api.VerifierClaimMapping) []map[string]string {
	mappings := make([]map[string]string, len(claimMappings))
	for i, mapping := range claimMappings {
		mappings[i] = map[string]string{
			"json_ld_fqn": mapping.JsonLdFqn,
			"oidc_claim":  mapping.OidcClaim,
		}
	}
	return mappings
}
