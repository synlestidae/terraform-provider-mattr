package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"nz.antunovic/mattr-terraform-provider/api"
	"nz.antunovic/mattr-terraform-provider/generator"
)

func resourceIssuer() *schema.Resource {
	issuerSchema := map[string]*schema.Schema{
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
			Type:     schema.TypeSet,
			Required: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"type": &schema.Schema{
			Type:     schema.TypeSet,
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
			Type:     schema.TypeSet,
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
		"callback_url": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
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
			Type:     schema.TypeSet,
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
		"openid_configuration_url": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
	}

	issuerGenerator := generator.Generator{
		Path:               "/ext/oidc/v1/issuers",
		Schema:             issuerSchema,
		Client:             &api.HttpClient{},
		ModifyRequestBody:  issuerConvertReq,
		ModifyResponseBody: issuerConvertRes,
	}

	issuerResource := issuerGenerator.GenResource()

	// we need to integrate data from generator and resource data
	// to compute openid-configuration url
	// TODO: this is complicated because the generator "modify" functions aren't
	// powerful enough.

	createOrig := issuerResource.Create
	readOrig := issuerResource.Read
	updateOrig := issuerResource.Update

	issuerResource.Create = func(d *schema.ResourceData, m interface{}) error {
		err := createOrig(d, m)
		if err != nil {
			return err
		}
		return setOpenIdConfigurationUrl(d, m)
	}

	issuerResource.Read = func(d *schema.ResourceData, m interface{}) error {
		err := readOrig(d, m)
		if err != nil {
			return err
		}
		return setOpenIdConfigurationUrl(d, m)
	}

	issuerResource.Update = func(d *schema.ResourceData, m interface{}) error {
		err := updateOrig(d, m)
		if err != nil {
			return err
		}
		return setOpenIdConfigurationUrl(d, m)
	}

	return &issuerResource
}

func setOpenIdConfigurationUrl(d *schema.ResourceData, m interface{}) error {
	// e.g. GET https://YOUR_TENANT_URL/ext/oidc/v1/issuers/983c0a86-204f-4431-9371-f5a22e506599/.well-known/openid-configuration
	api := m.(api.ProviderConfig).Api
	id := d.Id()
	path := fmt.Sprintf("%s/.well-known/openid-configuration", id)
	openIdConfigurationUrl, err := api.GetUrl(path)
	if err != nil {
		return d.Set("openid_configuration_url", openIdConfigurationUrl)
	}
	return err
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
	proofType := bodyMap["proofType"]
	if proofType != nil {
		credential["proofType"] = proofType
	}
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
	bodyMap, ok := body.(map[string]interface{}) // TODO cast safely
	if !ok {
		return nil, fmt.Errorf("Unexpected response type for issuer: %T", body)
	}
	staticRequestParameters := bodyMap["staticRequestParameters"]
	delete(bodyMap, "staticRequestParameters")
	flattenedMap := flattenMap(body).(map[string]interface{})
	flattenedMap["staticRequestParameters"] = staticRequestParameters

	// type needs to be converted to array of strings
	if issuerType, ok := flattenedMap["type"].(string); ok {
		flattenedMap["type"] = []interface{}{issuerType}
	}

	return flattenedMap, nil
}

func flattenMap(input interface{}) interface{} {
	if inputMap, ok := input.(map[string]interface{}); ok {
		outputMap := make(map[string]interface{})

		for key, value := range inputMap {
			flattened := flattenMap(value)
			if flattenedMap, ok := flattened.(map[string]interface{}); ok {
				for subKey, subValue := range flattenedMap {
					outputMap[subKey] = subValue
				}
			} else {
				valueStr = fmt.Sprintf("%v", value)
				outputMap[key] = valueStr
			}
		}

		return outputMap
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
