package main

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceDid() *schema.Resource {
	return &schema.Resource{
		Create: resourceDidCreate,
		Read:   resourceDidRead,
		Delete: resourceDidDelete,

		Schema: map[string]*schema.Schema{
			"method": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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
			"initial_did_document": &schema.Schema{
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Resource{
					Schema: didDocumentSchema(),
				},
			},
		},
	}
}

func resourceDidCreate(d *schema.ResourceData, m interface{}) error {
	log.Println("Creating did")

	// TODO: check if the did exists first

	api := m.(ProviderConfig).Api

	// prepare the did body request
	method := d.Get("method").(string)
	options := DidRequestOptions{
		KeyType: d.Get("key_type").(string),
		Url:     d.Get("url").(string),
	}
	did_request := DidRequest{
		Method:  method,
		Options: options,
	}

	did_response, err := api.PostDid(did_request)
	if err != nil {
		return err
	}

	// success, process did
	processDidData(d, did_response)
	d.SetId(did_response.Did)
	return nil
}

func resourceDidRead(d *schema.ResourceData, m interface{}) error {
	log.Println("Reading did")

	api := m.(ProviderConfig).Api

	did := d.Id()
	did_response, err := api.GetDid(did)
	if err != nil {
		return err
	}

	processDidData(d, did_response)
	d.SetId(did)

	return nil
}

func resourceDidDelete(d *schema.ResourceData, m interface{}) error {
	log.Println("Reading did")

	api := m.(ProviderConfig).Api

	did := d.Id()
	err := api.DeleteDid(did)
	if err != nil {
		return err
	}

	d.SetId("")

	return nil

}

func processDidData(d *schema.ResourceData, response *DidResponse) {
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
	processDidDocument(d, &response.LocalMetadata.InitialDidDocument)
}

func processDidDocument(d *schema.ResourceData, didDocument *DidDocument) {
	publicKey := make([]map[string]interface{}, len(didDocument.PublicKey))
	keyAgreement := make([]map[string]interface{}, len(didDocument.KeyAgreement))

	for i, pubKey := range didDocument.PublicKey {
		publicKey[i] = make(map[string]interface{}, 4)
		publicKey[i]["id"] = pubKey.Id
		publicKey[i]["type"] = pubKey.Type
		publicKey[i]["controller"] = pubKey.Controller
		publicKey[i]["public_key_base58"] = pubKey.PublicKeyBase58
	}

	for i, ka := range didDocument.KeyAgreement {
		keyAgreement[i] = make(map[string]interface{}, 4)
		keyAgreement[i]["id"] = ka.Id
		keyAgreement[i]["type"] = ka.Type
		keyAgreement[i]["controller"] = ka.Controller
		keyAgreement[i]["public_key_base58"] = ka.PublicKeyBase58
	}

	auth := make([]string, len(didDocument.Authentication))
	copy(auth, didDocument.Authentication)

	assertion := make([]string, len(didDocument.AssertionMethod))
	copy(assertion, didDocument.AssertionMethod)

	delegation := make([]string, len(didDocument.CapabilityDelegation))
	copy(delegation, didDocument.CapabilityDelegation)

	invocation := make([]string, len(didDocument.CapabilityInvocation))
	copy(invocation, didDocument.CapabilityInvocation)

	didDoc := map[string]interface{}{
		"id":                    didDocument.Id,
		"@context":              didDocument.Context,
		"public_key":            publicKey,
		"key_agreement":         keyAgreement,
		"authentication":        auth,
		"assertion_method":      assertion,
		"capability_delegation": delegation,
		"capability_invocation": invocation,
	}

	d.Set("initial_did_document", didDoc)
}

func didDocumentSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"@context": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
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
					"private_key_base58": &schema.Schema{
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
	}
}
