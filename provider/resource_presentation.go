package provider

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"nz.antunovic/mattr-terraform-provider/api"
	"nz.antunovic/mattr-terraform-provider/generator"
)

func resourcePresentation() *schema.Resource {
	presentationSchema := map[string]*schema.Schema{
		"domain": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"name": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"query": &schema.Schema{
			Type:     schema.TypeList,
			Required: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"type": &schema.Schema{
						Type:     schema.TypeString,
						Required: true,
					},
					"credential_query": &schema.Schema{
						Type:     schema.TypeList,
						Required: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"required": &schema.Schema{
									Type:     schema.TypeBool,
									Required: true,
								},
								"reason": &schema.Schema{
									Type:     schema.TypeString,
									Required: true,
								},
								"trusted_issuer": &schema.Schema{
									Type:     schema.TypeList,
									Optional: true,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"required": &schema.Schema{
												Type:     schema.TypeBool,
												Required: true,
											},
											"issuer": &schema.Schema{
												Type:     schema.TypeString,
												Required: true,
											},
										},
									},
								},
								"frame": &schema.Schema{
									Type:     schema.TypeString,
									Optional: true,
								},
								"example": &schema.Schema{
									Type:     schema.TypeString,
									Optional: true,
								},
							},
						},
					},
				},
			},
		},
	}

	generator := generator.Generator{
		Path:               "/v2/credentials/web-semantic/presentations/templates",
		Client:             &api.HttpClient{},
		Schema:             presentationSchema,
		ModifyRequestBody:  modifyRequestBody,
		ModifyResponseBody: modifyResponseBody,
	}

	resource := generator.GenResource()
	return &resource
}

func modifyRequestBody(body interface{}) (interface{}, error) {
	return transformBody(body, true)
}

func modifyResponseBody(body interface{}) (interface{}, error) {
	return transformBody(body, false)
}

// TODO: fix this monstrosity
func transformBody(body interface{}, isRequest bool) (interface{}, error) {
	bodyMap, ok := body.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("Unexpected type for /v2/credentials/web-semantic/presentations/templates response: %T", body)
	}

	bodyQueryList, ok := bodyMap["query"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("Unexpected type for /v2/credentials/web-semantic/presentations/templates 'query' field: %T", bodyQueryList)
	}

	for j, bodyQuery := range bodyQueryList {
		bodyQueryMap, ok := bodyQuery.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("Unexpected type for /v2/credentials/web-semantic/presentations/templates 'query' index %d: %T", j, bodyQueryList)
		}

		credentialQuery, ok := bodyQueryMap["credentialQuery"].([]interface{})
		if !ok {
			return nil, fmt.Errorf("Unexpected type for /v2/credentials/web-semantic/presentations/templates 'credentialQuery' field: %T", credentialQuery)
		}

		for i, query := range credentialQuery {
			credQueryMap, ok := query.(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("Unexpected type for /v2/credentials/web-semantic/presentations/templates 'credentialQuery' element at index %d: %T", i, body)
			}

			if isRequest {
				if example, ok := credQueryMap["example"].(string); ok {
					var payload json.RawMessage
					err := json.Unmarshal([]byte(example), &payload)
					if err != nil {
						return nil, fmt.Errorf("Error parsing JSON for 'example' field (check your syntax): %s", err)
					}
					credQueryMap["example"] = payload
				} else if credQueryMap["example"] != nil {
					return nil, fmt.Errorf("Unexpected type for /v2/credentials/web-semantic/presentations/templates 'example' field: %T", example)
				}

				if frame, ok := credQueryMap["frame"].(string); ok {
					var payload json.RawMessage
					err := json.Unmarshal([]byte(frame), &payload)
					if err != nil {
						return nil, fmt.Errorf("Error parsing JSON for 'frame' field (check your syntax): %s", err)
					}
					credQueryMap["frame"] = payload
				} else if credQueryMap["frame"] != nil {
					return nil, fmt.Errorf("Unexpected type for /v2/credentials/web-semantic/presentations/templates 'frame' field: %T", frame)
				}
			} else {
				if example, ok := credQueryMap["example"]; ok && example != nil {
					jsonExample, err := json.Marshal(example)
					if err != nil {
						return nil, fmt.Errorf("Error serialising 'example' field from JSON (check your syntax): %s", err)
					}
					credQueryMap["example"] = string(jsonExample)
				}

				if frame, ok := credQueryMap["frame"].(interface{}); ok && frame != nil {
					jsonFrame, err := json.Marshal(frame)
					if err != nil {
						return nil, fmt.Errorf("Error serialising 'frame' field (check your syntax): %s", err)
					}
					credQueryMap["frame"] = string(jsonFrame)
				}
			}
		}
	}

	return bodyMap, nil
}
