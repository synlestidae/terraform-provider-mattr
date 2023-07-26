package provider

import (
	"fmt"

	"nz.antunovic/mattr-terraform-provider/api"
	"nz.antunovic/mattr-terraform-provider/generator"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceIssuer() *schema.Resource {
	schema := map[string]*schema.Schema{
		"issuer_did": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"issuer_logo_url": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"issuer_icon_url": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"name": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"issuer_name": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"description": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"context": &schema.Schema{
			Type:     schema.TypeList,
			Required: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"type": &schema.Schema{
			Type:     schema.TypeList,
			Required: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
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
		"url": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"scope": &schema.Schema{
			Type:     schema.TypeList,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"client_id": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"client_secret": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"token_endpoint_auth_method": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"claims_source": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"static_request_parameters": &schema.Schema{
			Type:     schema.TypeMap,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"forwarded_request_parameters": &schema.Schema{
			Type:     schema.TypeList,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"claim_mappings": &schema.Schema{
			Type:     schema.TypeList,
			Required: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"json_ld_term": &schema.Schema{
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
	}

	issuerGenerator := generator.Generator{
		Path:               "/ext/oidc/v1/issuers",
		Schema:             schema,
		Client:             &api.HttpClient{},
		ModifyRequestBody:  issuerConvertReq,
		ModifyResponseBody: issuerConvertRes,
	}

	issuerResource := issuerGenerator.GenResource()
	return &issuerResource
}

func issuerConvertReq(body interface{}) (interface{}, error) {
	bodyMap, ok := body.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("Unexpected type for /ext/oidc/v1/issuers response: %T", body)
	}

	newBody := make(map[string]interface{})

	credential := make(map[string]interface{})
	credentialBranding := make(map[string]interface{})

	credential["name"] = bodyMap["name"]
	credential["issuerDid"] = bodyMap["issuerDid"]
	credential["issuerName"] = bodyMap["issuerName"]
	credential["issuerLogoUrl"] = bodyMap["issuerLogoUrl"]
	credential["issuerIconUrl"] = bodyMap["issuerIconUrl"]
	credential["issuerName"] = bodyMap["issuerName"]
	credential["description"] = bodyMap["description"]
	credential["context"] = bodyMap["context"]
	credential["type"] = bodyMap["type"]
	credential["proofType"] = bodyMap["proofType"]
	credential["credentialBranding"] = credentialBranding
	credentialBranding["backgroundColor"] = bodyMap["backgroundColor"]
	credentialBranding["watermarkImageUrl"] = bodyMap["watermarkImageUrl"]
	credential["credentialBranding"] = credentialBranding
	newBody["credential"] = credential

	federatedProvider := make(map[string]interface{})
	federatedProvider["url"] = bodyMap["url"]
	federatedProvider["scope"] = bodyMap["scope"]
	federatedProvider["clientId"] = bodyMap["clientId"]
	federatedProvider["clientSecret"] = bodyMap["clientSecret"]
	federatedProvider["tokenEndpointAuthMethod"] = bodyMap["tokenEndpointAuthMethod"]
	federatedProvider["claimsSource"] = bodyMap["claimsSource"]
	newBody["federatedProvider"] = federatedProvider

	newBody["staticRequestParameters"] = bodyMap["staticRequestParameters"]
	newBody["forwardedRequestParameters"] = bodyMap["forwardedRequestParameters"]
	newBody["claimMappings"] = bodyMap["claimMappings"]

	return newBody, nil
}

func issuerConvertRes(body interface{}) (interface{}, error) {
	return flattenMap(body), nil
}

func flattenMap(input interface{}) interface{} {
	if inputMap, ok := input.(map[string]interface{}); ok {
		flattenedMap := make(map[string]interface{})

		for key, value := range inputMap {
			if subValue, ok := value.(map[string]interface{}); ok {
				for skey, sval := range subValue {
					if skey != "" {
						flattenedMap[skey] = flattenMap(sval)
					}
				}
			} else if key != "" {
				flattenedMap[key] = value
			}
		}

		return flattenedMap
	}

	return input
}

func castToStringSlice(val interface{}) []string {
	if val == nil {
		return make([]string, 0)
	}
	interfaceSlice := val.([]interface{})
	stringSlice := make([]string, len(interfaceSlice))
	for i, s := range interfaceSlice {
		stringSlice[i] = s.(string)
	}
	return stringSlice
}
