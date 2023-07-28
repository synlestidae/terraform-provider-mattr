package provider

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"nz.antunovic/mattr-terraform-provider/api"
	"nz.antunovic/mattr-terraform-provider/generator"
)

func resourceClaimSource() *schema.Resource {
	claimSourceSchema := map[string]*schema.Schema{
		"name": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"url": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"authorization_type": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"authorization_value": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"request_parameter": &schema.Schema{
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
						Required: true,
					},
					"default_value": &schema.Schema{
						Type:     schema.TypeString,
						Optional: true,
					},
				},
			},
		},
	}

	claimSourceGenerator := generator.Generator{
		Path:               "/core/v1/claimsources",
		Client:             &api.HttpClient{},
		Schema:             claimSourceSchema,
		ModifyRequestBody:  convertReqParamsBody,
		ModifyResponseBody: convertResParamsBody,
	}

	resource := claimSourceGenerator.GenResource()

	return &resource
}

func convertReqParamsBody(body interface{}) (interface{}, error) {
	reqMap, ok := body.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("Unexpected type for /v1/core/claimsources response: %T", body)
	}

	authorization := make(map[string]interface{}, 2)
	authorization["type"] = reqMap["authorizationType"]
	authorization["value"] = reqMap["authorizationValue"]
	reqMap["authorization"] = authorization
	delete(reqMap, "authorizationType")
	delete(reqMap, "authorizationValue")

	paramList, ok := reqMap["requestParameter"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("Unexpected type for /core/v1/claimsources 'requestParameter' field: %T", reqMap)
	}

	paramsMap := make(map[string]interface{}, len(paramList))
	for i, param := range paramList {
		paramElem, ok := param.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("Unexpected type for /v1/core/claimsources request parameter %d: %T", i, param)
		}

		property, ok := paramElem["name"].(string)
		if !ok {
			return nil, fmt.Errorf("Unexpected type for property field in /v1/core/claimsources request parameter: %T", paramElem["name"])
		}

		paramMap := map[string]string{
			"mapFrom":      paramElem["mapFrom"].(string),
			"defaultValue": paramElem["defaultValue"].(string),
		}

		paramsMap[property] = paramMap
	}

	delete(reqMap, "requestParameter")
	reqMap["requestParameters"] = paramsMap

	return reqMap, nil
}

func convertResParamsBody(body interface{}) (interface{}, error) {
	reqMap, ok := body.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("Unexpected type for response body: %T", body)
	}

	authorizationMap := reqMap["authorization"].(map[string]interface{})
	reqMap["authorization_type"] = authorizationMap["type"]
	reqMap["authorization_value"] = authorizationMap["value"]
	delete(reqMap, "authorization")

	paramsMap, ok := reqMap["requestParameters"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("Unexpected type for requestParameters field in response: %T", reqMap["requestParameters"])
	}

	paramList := make([]interface{}, 0, len(paramsMap))
	for property, paramMap := range paramsMap {
		paramMapTyped, ok := paramMap.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("Unexpected type for property '%s' in requestParameters: %T", property, paramMap)
		}

		param := map[string]interface{}{
			"name":         property,
			"mapFrom":      paramMapTyped["mapFrom"].(string),
			"defaultValue": paramMapTyped["defaultValue"].(string),
		}

		paramList = append(paramList, param)
	}

	reqMap["requestParameter"] = paramList
	delete(reqMap, "requestParameters")

	return reqMap, nil
}
