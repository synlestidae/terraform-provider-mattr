package main

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	_ "github.com/motemen/go-loghttp/global"
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
		},
	}
}

func resourceDidCreate(d *schema.ResourceData, m interface{}) error {
	log.Println("Creating did")

	// TODO: check if the did exists first

	api := InitFromEnv()

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

	api := InitFromEnv()

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

	api := InitFromEnv()

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
}
