package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"nz.antunovic/mattr-terraform-provider/api"
	"nz.antunovic/mattr-terraform-provider/generator"
)

func resourceDid() *schema.Resource {
	schema := map[string]*schema.Schema{
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
		"keys": &schema.Schema{
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"did_document_key_id": {
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
					"kms_key_id": {
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
				},
			},
		},
	}

	modifyResponseBody := func(responseBody interface{}) (interface{}, error) {
		if orig, ok := (responseBody).(map[string]interface{}); ok {
			if localMetadata, ok := orig["localMetadata"].(map[string]interface{}); ok {
				newResponse := make(map[string]interface{})
				newResponse["keys"] = localMetadata["keys"]

				did, ok := orig["did"].(string)
				if ok && len(did) != 0 {
					newResponse["id"] = did
					return newResponse, nil
				}

				initialDidDocument, ok := localMetadata["initialDidDocument"].(map[string]interface{})
				if !ok {
					return nil, fmt.Errorf("Internal error: unable to determine did")
				}
				did, ok = initialDidDocument["id"].(string)
				if !ok {
					return nil, fmt.Errorf("Internal error: unable to determine did")
				}
				newResponse["id"] = did
				return newResponse, nil
			} else {
				return nil, fmt.Errorf("Unexpected type for DID `localMetadata` field: %T", orig["localMetadata"])
			}
		} else {
			return nil, fmt.Errorf("Unexpected type for DID response: %T", responseBody)
		}
	}

	generator := generator.Generator{
		Path:               "/core/v1/dids",
		Immutable:          true,
		Client:             &api.HttpClient{},
		Schema:             schema,
		ModifyResponseBody: modifyResponseBody,
	}

	resource := generator.GenResource()
	return &resource
}
