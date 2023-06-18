package provider

import (
	"fmt"
	"log"
	"nz.antunovic/mattr-terraform-provider/api"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceDid() *schema.Resource {
	return &schema.Resource{
		Create: resourceDidCreate,
		Read:   resourceDidRead,
		Delete: resourceDidDelete,

		Schema: map[string]*schema.Schema{
			"method": &schema.Schema{
				Type:        schema.TypeString,
				Description: "The method (or type) of did: key, web, or ion",
				Required:    true,
				ForceNew:    true,
			},
			"url": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Domain or URL from which hostname will be extracted",
				ForceNew:    true,
			},
			"key_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"did": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"registered": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"keys": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"did_document_key_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"kms_key_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"context": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"public_key": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"controller": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"public_key_base58": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"key_agreement": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"controller": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"public_key_base58": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"authentication": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"assertion_method": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"capability_delegation": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"capability_invocation": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceDidCreate(d *schema.ResourceData, m interface{}) error {
	log.Println("Creating did")

	// TODO: check if the did exists first

	didApi := m.(ProviderConfig).Api

	// prepare the did body request
	method := d.Get("method").(string)
	options := api.DidRequestOptions{
		KeyType: d.Get("key_type").(string),
		Url:     d.Get("url").(string),
	}
	did_request := api.DidRequest{
		Method:  method,
		Options: options,
	}

	did_response, err := didApi.PostDid(did_request)
	if err != nil {
		return err
	}

	// success, process did
	err = processDidData(d, did_response)
	if err != nil {
		return err
	}
	d.SetId(did_response.Did)
	return nil
}

func resourceDidRead(d *schema.ResourceData, m interface{}) error {
	log.Println("Reading did")

	didApi := m.(ProviderConfig).Api

	did := d.Id()
	did_response, err := didApi.GetDid(did)
	if err != nil {
		return err
	}

	err = processDidData(d, did_response)
	if err != nil {
		return err
	}
	d.SetId(did)

	return nil
}

func resourceDidDelete(d *schema.ResourceData, m interface{}) error {
	log.Println("Reading did")

	didApi := m.(ProviderConfig).Api

	did := d.Id()
	err := didApi.DeleteDid(did)
	if err != nil {
		return err
	}

	d.SetId("")

	return nil

}

func processDidData(d *schema.ResourceData, response *api.DidResponse) error {
	d.Set("registration_status", response.RegistrationStatus)
	d.Set("registered", response.LocalMetadata.Registered)

	keys := []interface{}{}

	for _, k := range response.LocalMetadata.Keys {
		key := map[string]string{
			"did_document_key_id": k.DidDocumentKeyId,
			"kms_key_id":          k.KmsKeyId,
		}
		keys = append(keys, key)
	}

	d.Set("keys", keys)
	return processDidDocument(d, &response.LocalMetadata.InitialDidDocument)
}

func processDidDocument(d *schema.ResourceData, didDocument *api.DidDocument) error {
	publicKey := make([]map[string]string, len(didDocument.PublicKey))
	keyAgreement := make([]map[string]string, len(didDocument.KeyAgreement))

	for i, pubKey := range didDocument.PublicKey {
		publicKey[i] = make(map[string]string, 4)
		publicKey[i]["id"] = pubKey.Id
		publicKey[i]["type"] = pubKey.Type
		publicKey[i]["controller"] = pubKey.Controller
		publicKey[i]["public_key_base58"] = pubKey.PublicKeyBase58
	}

	for i, ka := range didDocument.KeyAgreement {
		keyAgreement[i] = make(map[string]string, 4)
		keyAgreement[i]["id"] = ka.Id
		keyAgreement[i]["type"] = ka.Type
		keyAgreement[i]["controller"] = ka.Controller
		keyAgreement[i]["public_key_base58"] = ka.PublicKeyBase58
	}

	//auth := make([]interface{}, len(didDocument.Authentication))

	var err error

	if err = d.Set("context", didDocument.Context.Uris); err != nil {
		return fmt.Errorf("error setting 'context' field: %s", err)
	}
	if err = d.Set("public_key", publicKey); err != nil {
		return fmt.Errorf("error setting 'public_key' field: %s", err)
	}
	if err = d.Set("key_agreement", keyAgreement); err != nil {
		return fmt.Errorf("error setting 'key_agreement' field: %s", err)
	}
	if err = d.Set("authentication", didDocument.Authentication); err != nil {
		return fmt.Errorf("error setting 'authentication' field: %s", err)
	}
	if err = d.Set("assertion_method", didDocument.AssertionMethod); err != nil {
		return fmt.Errorf("error setting 'assertion_method' field: %s", err)
	}
	if err = d.Set("capability_delegation", didDocument.CapabilityDelegation); err != nil {
		return fmt.Errorf("error setting 'capability_delegation' field: %s", err)
	}
	if err = d.Set("capability_invocation", didDocument.CapabilityInvocation); err != nil {
		return fmt.Errorf("error setting 'capability_invocation' field: %s", err)
	}

	return nil
}
