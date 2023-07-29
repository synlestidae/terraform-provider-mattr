package provider

import (
	"fmt"
	"log"

	"nz.antunovic/mattr-terraform-provider/api"
	"nz.antunovic/mattr-terraform-provider/generator"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceVerifierClient(client api.Client) *schema.Resource {
	verifierClientSchema := map[string]*schema.Schema{
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
			Optional: true,
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
			Optional: true,
		},
		"logo_uri": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"secret": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"openid_configuration_url": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"authorization_url": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
	}

	getPath := func(d *schema.ResourceData) (string, error) {
		if verifierId, ok := d.Get("verifier_id").(string); ok && len(verifierId) != 0 {
			return fmt.Sprintf("/ext/oidc/v1/verifiers/%s/clients", verifierId), nil
		}

		return "", fmt.Errorf("'verifier_id' field is required for verifier client and must be a string")
	}

	verifierClientModifyRes := func(body interface{}) (interface{}, error) {
		if bodyMap, ok := body.(map[string]interface{}); ok {
			delete(bodyMap, "verifierId")
			return bodyMap, nil
		}

		return nil, fmt.Errorf("Unexpected type for /ext/oidc/v1/issuers/{id}/clients request: %T", body)
	}

	generator := generator.Generator{
		GetPath:           getPath,
		Client:            client,
		Schema:            verifierClientSchema,
		ModifyRequestBody: verifierClientModifyRes,
	}
	resource := generator.GenResource()

	return &resource
}

func resourceVerifierClientCreate(d *schema.ResourceData, m interface{}) error {
	log.Println("Creating verifier client")

	api := m.(api.ProviderConfig).Api
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
	api := m.(api.ProviderConfig).Api

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
	api := m.(api.ProviderConfig).Api

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

	api := m.(api.ProviderConfig).Api
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

func processVerifierClientData(verifier_client *api.VerifierClient, d *schema.ResourceData) error {
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

func fromTerraformVerifierClient(d *schema.ResourceData) (*api.VerifierClient, error) {
	return &api.VerifierClient{
		Name:                     d.Get("name").(string),
		RedirectUris:             castToStringSlice(d.Get("redirect_uris")),
		ResponseTypes:            castToStringSlice(d.Get("response_types")),
		GrantTypes:               castToStringSlice(d.Get("grant_types")),
		TokenEndpointAuthMethod:  d.Get("token_endpoint_auth_method").(string),
		IdTokenSignedResponseAlg: d.Get("id_token_signed_response_alg").(string),
		ApplicationType:          d.Get("application_type").(string),
	}, nil
}
