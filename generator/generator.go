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

func (generator *Generator) sendRequestAndProcessResponse(d *schema.ResourceData, m interface{}, operation string) error {
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

	log.Printf("Going to send request for resource: %s", path)

	url, err := api.GetUrl(path)
	if err != nil {
		return err
	}

	fullUrl := url
	if !generator.Singleton && operation != "create" {
		fullUrl = fmt.Sprintf("%s/%s", url, d.Id())
	}


	log.Printf("Full resource URL is: %s", fullUrl)
	log.Printf("Getting access token for %s", fullUrl)

	accessToken, err := api.GetAccessToken()
	if err != nil {
		return err
	}
	headers := map[string]string{
		"Authorization": "Bearer " + accessToken,
	}

	var body interface{}
	if operation == "create" || operation == "update" {
		log.Printf("Operation for %s is create or update, generating request body", fullUrl)
		body, err = requestVisitor.accept(d)
		if err != nil {
			return err
		}
	}

	// modify request
	if generator.ModifyRequestBody != nil && (operation == "create" || operation == "update") {
		log.Printf("Operation for %s is create or update, modifying request body", fullUrl)
		body, err = generator.ModifyRequestBody(body)
		if err != nil {
			return err
		}
	}

	if generator.ModifyRequest != nil && (operation == "create" || operation == "update") {
		log.Printf("Operation for %s is create or update, modifying request", fullUrl)
		err = generator.ModifyRequest(&url, &headers, &body)
		if err != nil {
			return err
		}
	}

	// send request
	var response interface{}
	switch operation {
	case "create":
		response, err = generator.Client.Post(fullUrl, headers, body)
	case "read":
		response, err = generator.Client.Get(fullUrl, headers)
	case "update":
		response, err = generator.Client.Put(fullUrl, headers, body)
	case "delete":
		err = generator.Client.Delete(fullUrl, headers)
	default:
		return fmt.Errorf("unknown operation: %s", operation)
	}

	if err != nil {
		return err
	}

	// on successful delete, exit early
	if operation == "delete" {
		log.Printf("Delete for %s was successful", fullUrl)
		return nil
	}

	// modify response
	if generator.ModifyResponseBody != nil {
		log.Printf("Modifying response body for %s", fullUrl)
		response, err = generator.ModifyResponseBody(response)
		if err != nil {
			return err
		}
	}
	if generator.ModifyResponse != nil {
		log.Printf("Modifying response for %s", fullUrl)
		err = generator.ModifyResponse(&map[string]string{}, &response) // TODO response headers
		if err != nil {
			return err
		}
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
	// TODO: move this to visitor
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
