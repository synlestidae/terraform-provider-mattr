package generator

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"nz.antunovic/mattr-terraform-provider/api"
)

type Generator struct {
	Path      string
	GetPath   func(*schema.ResourceData) (string, error)
	Immutable bool
	Singleton bool
	Schema    map[string]*schema.Schema
	Client    api.Client

	ModifyRequestBody  func(requestBody interface{}) (interface{}, error)
	ModifyResponseBody func(responseBody interface{}) (interface{}, error)

	ModifyRequest      func(url *string, headers *map[string]string, body *interface{}) error
	ModifyResponse     func(headers *map[string]string, body *interface{}) error
	ModifyResourceData func(resourceData *schema.ResourceData) error
	GetId              func(requestBody *interface{}, responseBody *interface{}) string
}

func (generator *Generator) GenResource() schema.Resource {
	create := func(d *schema.ResourceData, m interface{}) error {
		return generator.sendRequestAndProcessResponse(d, m, "create")
	}

	read := func(d *schema.ResourceData, m interface{}) error {
		return generator.sendRequestAndProcessResponse(d, m, "read")
	}

	update := func(d *schema.ResourceData, m interface{}) error {
		return generator.sendRequestAndProcessResponse(d, m, "update")
	}

	deleteResource := func(d *schema.ResourceData, m interface{}) error {
		return generator.sendRequestAndProcessResponse(d, m, "delete")
	}

	resource := schema.Resource{
		Create: create,
		Read:   read,
		Delete: deleteResource,
		Schema: generator.Schema,
	}

	if !generator.Immutable {
		resource.Update = update
	}

	return resource
}

// Helper function for sending requests and processing responses
func (generator *Generator) sendRequestAndProcessResponse(d *schema.ResourceData, m interface{}, method string) error {
	api := m.(api.ProviderConfig).Api
	requestVisitor := RequestVisitor{
		schema: generator.Schema,
	}

	var path string
	var err error
	if generator.GetPath != nil {
		path, err = generator.GetPath(d)
		if err != nil {
			return err
		}
	} else {
		path = generator.Path
	}

	log.Printf("Path is %s", path)

	url, err := api.GetUrl(path)
	if err != nil {
		return err
	}

	log.Printf("Url is '%s'", url)

	fullUrl := url
	if !generator.Singleton && method != "create" {
		fullUrl = fmt.Sprintf("%s/%s", url, d.Id())
	}

	accessToken, err := api.GetAccessToken()
	if err != nil {
		return err
	}
	headers := map[string]string{
		"Authorization": "Bearer " + accessToken,
	}

	var body interface{}
	if method != "read" {
		body, err = requestVisitor.accept(d)
		if err != nil {
			return err
		}
	}

	// modify request
	if generator.ModifyRequestBody != nil {
		body, err = generator.ModifyRequestBody(body)
		if err != nil {
			return err
		}
	}

	if generator.ModifyRequest != nil {
		err = generator.ModifyRequest(&url, &headers, &body)
		if err != nil {
			return err
		}
	}

	log.Printf("Full URL '%s'", fullUrl)

	// send request
	var response interface{}
	switch method {
	case "create":
		response, err = generator.Client.Post(fullUrl, headers, body)
	case "read":
		response, err = generator.Client.Get(fullUrl, headers)
	case "update":
		response, err = generator.Client.Put(fullUrl, headers, body)
	case "delete":
		err = generator.Client.Delete(fullUrl, headers)
	default:
		return fmt.Errorf("unknown method: %s", method)
	}

	if err != nil {
		return err
	}

	// modify response
	if generator.ModifyResponseBody != nil {
		response, err = generator.ModifyResponseBody(response)
		if err != nil {
			return err
		}
	}
	if generator.ModifyResponse != nil {
		generator.ModifyResponse(&map[string]string{}, &response) // TODO response headers
	}

	// process response
	responseVisitor := ResponseVisitor{}
	transformedResponse, err := responseVisitor.accept(response)
	if err != nil {
		return err
	}

	var id string
	if generator.GetId != nil {
		id = generator.GetId(&body, &response)
	} else {
		id = responseVisitor.id
	}

	d.SetId(id)
	if data, ok := transformedResponse.(map[string]interface{}); ok {
		for key, val := range data {
			err := d.Set(key, val)
			if err != nil {
				log.Printf("Unable to set '%s' = '%s'. Ignoring.", key, err)
			}
		}
	}

	return nil
}
