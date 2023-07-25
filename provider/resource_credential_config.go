package provider

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"nz.antunovic/mattr-terraform-provider/api"
	"nz.antunovic/mattr-terraform-provider/generator"
)

func resourceCredentialConfig() *schema.Resource {
	credentialConfigSchema := map[string]*schema.Schema{
		"name": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"description": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"type": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"additional_types": &schema.Schema{
			Type:     schema.TypeList,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"contexts": &schema.Schema{
			Type:     schema.TypeList,
			Required: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"issuer_name": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"issuer_logo_url": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"issuer_icon_url": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"proof_type": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"background_color": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"watermark_image_url": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"claim_mapping": &schema.Schema{
			Type:     schema.TypeSet,
			Required: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": &schema.Schema{
						Type:     schema.TypeString,
						Required: true,
					},
					"map_from": &schema.Schema{
						Type:     schema.TypeString,
						Optional: true,
					},
					"default_value": &schema.Schema{
						Type:     schema.TypeString,
						Optional: true,
					},
					"required": &schema.Schema{
						Type:     schema.TypeBool,
						Optional: true,
					},
				},
			},
		},
		"persist": &schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
		},
		"revocable": &schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
		},
		"include_id": &schema.Schema{
			Type:     schema.TypeBool,
			Computed: true,
		},
		"claim_source_id": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"years": &schema.Schema{
			Type:     schema.TypeInt,
			Optional: true,
		},
		"months": &schema.Schema{
			Type:     schema.TypeInt,
			Optional: true,
		},
		"weeks": &schema.Schema{
			Type:     schema.TypeInt,
			Optional: true,
		},
		"days": &schema.Schema{
			Type:     schema.TypeInt,
			Optional: true,
		},
		"hours": &schema.Schema{
			Type:     schema.TypeInt,
			Optional: true,
		},
		"minutes": &schema.Schema{
			Type:     schema.TypeInt,
			Optional: true,
		},
		"seconds": &schema.Schema{
			Type:     schema.TypeInt,
			Optional: true,
		},
	}

	claimSourceGenerator := generator.Generator{
		Path:               "/core/v2/credentials/web-semantic/configurations",
		Client:             &api.HttpClient{},
		Schema:             credentialConfigSchema,
		ModifyRequestBody:  convertCredentialConfigReq,
		ModifyResponseBody: convertCredentialConfigRes,
	}

	provider := claimSourceGenerator.GenResource()
	return &provider
}

func convertCredentialConfigReq(body interface{}) (interface{}, error) {
	bodyMap, ok := body.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("Unexpected type for /core/v2/credentials/web-semantic/configurations: %T", body)
	}

	// issuer
	issuer := make(map[string]interface{})
	issuer["name"] = bodyMap["issuerName"]
	issuer["logoUrl"] = bodyMap["issuerLogoUrl"]
	issuer["iconUrl"] = bodyMap["issuerIconUrl"]
	bodyMap["issuer"] = issuer
	delete(bodyMap, "issuerName")
	delete(bodyMap, "issuerLogoUrl")
	delete(bodyMap, "issuerIconUrl")

	// brandingMap
	brandingMap := make(map[string]interface{}, 2)
	brandingMap["backgroundColor"] = bodyMap["backgroundColor"]
	brandingMap["watermarkImageUrl"] = bodyMap["watermarkImageUrl"]
	bodyMap["credentialBranding"] = brandingMap
	delete(bodyMap, "backgroundColor")
	delete(bodyMap, "watermarkImageUrl")

	// claim mappings
	claimMappingList, ok := bodyMap["claimMapping"].([]interface{})
	claimMappingApi := make(map[string]interface{})
	for _, claimMapping := range claimMappingList {
		claimMappingMap, ok := claimMapping.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("Unexpected type for /core/v2/credentials/web-semantic/configurations `claim_mapping`: %T", body)
		}
		name := claimMappingMap["name"].(string) // TODO cast this safely, return error if error
		delete(claimMappingMap, "name")
		claimMappingApi[name] = claimMappingMap
	}
	bodyMap["claimMappings"] = claimMappingApi
	delete(bodyMap, "claimMapping")

	// expires
	expiresMap := make(map[string]interface{}, 7)
	expiresMap["years"] = bodyMap["years"]
	expiresMap["months"] = bodyMap["months"]
	expiresMap["weeks"] = bodyMap["weeks"]
	expiresMap["days"] = bodyMap["days"]
	expiresMap["hours"] = bodyMap["hours"]
	expiresMap["minutes"] = bodyMap["minutes"]
	expiresMap["seconds"] = bodyMap["seconds"]
	bodyMap["expiresIn"] = expiresMap
	delete(bodyMap, "years")
	delete(bodyMap, "months")
	delete(bodyMap, "weeks")
	delete(bodyMap, "days")
	delete(bodyMap, "hours")
	delete(bodyMap, "minutes")
	delete(bodyMap, "seconds")

	return bodyMap, nil
}

func convertCredentialConfigRes(body interface{}) (interface{}, error) {
	bodyMap, ok := body.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("Unexpected type for response body: %T", body)
	}

	// issuer
	if bodyIssuer, ok := bodyMap["issuer"].(map[string]interface{}); ok {
		bodyMap["issuerName"] = bodyIssuer["name"]
		bodyMap["issuerLogoUrl"] = bodyIssuer["logoUrl"]
		bodyMap["issuerIconUrl"] = bodyIssuer["iconUrl"]
	}
	delete(bodyMap, "issuer")

	// brandingMap
	brandingMap, ok := bodyMap["credentialBranding"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("Unexpected type for response body field 'credentialBranding': %T", bodyMap["credentialBranding"])
	}
	bodyMap["backgroundColor"] = brandingMap["backgroundColor"]
	bodyMap["watermarkImageUrl"] = brandingMap["watermarkImageUrl"]
	delete(bodyMap, "credentialBranding")

	// claim mappings
	claimMappingApi, ok := bodyMap["claimMappings"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("Unexpected type for response body field 'claimMappings': %T", bodyMap["claimMappings"])
	}
	claimMappingList := make([]interface{}, 0, len(claimMappingApi))
	for name, claimMappingMap := range claimMappingApi {
		claimMapping, ok := claimMappingMap.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("Unexpected type for response body field 'claimMappings'[%s]: %T", name, claimMappingMap)
		}
		claimMapping["name"] = name
		claimMappingList = append(claimMappingList, claimMapping)
	}
	bodyMap["claimMapping"] = claimMappingList
	delete(bodyMap, "claimMappings")

	// expires
	expiresMap, ok := bodyMap["expiresIn"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("Unexpected type for response body field 'expiresIn': %T", bodyMap["expiresIn"])
	}
	bodyMap["years"] = expiresMap["years"]
	bodyMap["months"] = expiresMap["months"]
	bodyMap["weeks"] = expiresMap["weeks"]
	bodyMap["days"] = expiresMap["days"]
	bodyMap["hours"] = expiresMap["hours"]
	bodyMap["minutes"] = expiresMap["minutes"]
	bodyMap["seconds"] = expiresMap["seconds"]
	delete(bodyMap, "expiresIn")

	return bodyMap, nil
}
