package main

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceCustomDomain() *schema.Resource {
	return &schema.Resource{
		Create: resourceCustomDomainCreate,
		Read:   resourceCustomDomainRead,
		Update: resourceCustomDomainUpdate,
		Delete: resourceCustomDomainDelete,
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"logo_url": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"domain": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"homepage": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"verification_token": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_verified": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"verified_at": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceCustomDomainCreate(d *schema.ResourceData, m interface{}) error {
	log.Println("Creating custom domain")

	api := m.(ProviderConfig).Api
	custom_domain_request, err := fromTerraformCustomDomain(d)
	if err != nil {
		return err
	}

	custom_domain_response, err := api.PostCustomDomain(custom_domain_request)
	if err != nil {
		return err
	}
	if err = processCustomDomainData(custom_domain_response, d); err != nil {
		return err
	}

	log.Println("Created custom domain")
	return nil
}

func resourceCustomDomainRead(d *schema.ResourceData, m interface{}) error {
	log.Println("Reading credential config")

	api := m.(ProviderConfig).Api

	custom_domain_response, err := api.GetCustomDomain()
	if err != nil {
		return err
	}
	if err = processCustomDomainData(custom_domain_response, d); err != nil {
		return err
	}

	log.Println("Read credential config")
	return nil
}

func resourceCustomDomainUpdate(d *schema.ResourceData, m interface{}) error {
	log.Println("Updating custom domain")

	custom_domain_request, err := fromTerraformCustomDomain(d)
	if err != nil {
		return err
	}

	id := d.Id()
	api := m.(ProviderConfig).Api

	custom_domain_response, err := api.PutCustomDomain(id, custom_domain_request)
	if err != nil {
		return err
	}
	if err = processCustomDomainData(custom_domain_response, d); err != nil {
		return err
	}

	log.Println("Updated credential config")
	return nil
}

func resourceCustomDomainDelete(d *schema.ResourceData, m interface{}) error {
	log.Println("Deleting custom domain")

	api := m.(ProviderConfig).Api

	err := api.DeleteCustomDomain()
	if err != nil {
		return err
	}

	log.Println("Deleted custom domain")
	return nil
}

func processCustomDomainData(custom_domain *CustomDomainResponse, d *schema.ResourceData) error {
	log.Println("Converting custom domain from REST")

	d.SetId(custom_domain.Domain)

	d.Set("name", custom_domain.Name)
	d.Set("logo_url", custom_domain.LogoUrl)
	d.Set("domain", custom_domain.Domain)
	d.Set("homepage", custom_domain.Homepage)
	d.Set("verification_token", custom_domain.VerificationToken)
	d.Set("is_verified", custom_domain.IsVerified)
	d.Set("verified_at", custom_domain.VerifiedAt)

	return nil
}

func fromTerraformCustomDomain(d *schema.ResourceData) (*CustomDomainRequest, error) {
	return &CustomDomainRequest{
		Name:     d.Get("name").(string),
		LogoUrl:  d.Get("logo_url").(string),
		Domain:   d.Get("domain").(string),
		Homepage: d.Get("homepage").(string),
	}, nil
}
