package generator

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"nz.antunovic/mattr-terraform-provider/api"
)

type Generator struct {
	Path      string
	Immutable bool
	Singleton bool
	Schema    map[string]*schema.Schema
	Client    api.Client
	ModifyRequestBody func (*interface{}) error
	ModifyRequest func (url *string, headers *map[string]string, body *interface{}) error
	ModifyResponseBody func (*interface{}) error
	ModifyResponse func (headers *map[string]string, body *interface{}) error
	ModifyResourceData func (*schema.ResourceData) error
	GetId func(*interface{}, *interface{}) string
}

func (generator *Generator) GenResource() schema.Resource {
	create := func(d *schema.ResourceData, m interface{}) error {
		// prepare request
		api := m.(api.ProviderConfig).Api
		requestVisitor := RequestVisitor{
			schema: generator.Schema,
		}

		url, err := api.GetUrl(generator.Path)
		if err != nil {
			return err
		}
		accessToken, err := api.GetAccessToken()
		if err != nil {
			return err
		}
		headers := map[string]string{
			"Authorization": "Bearer " + accessToken,
		}

		body, err := requestVisitor.accept(d)
		if err != nil {
			return err
		}

		// modify request

		if generator.ModifyRequestBody != nil {
			generator.ModifyRequestBody(&body)
		}

		if generator.ModifyRequest != nil {
			generator.ModifyRequest(&url, &headers, &body)
		}

		// send

		response, err := generator.Client.Post(url, headers, body)
		if err != nil {
			return err
		}

		// modify response

		if generator.ModifyResponseBody != nil {
			generator.ModifyResponseBody(&response)
		}
		if generator.ModifyResponse != nil {
			generator.ModifyResponse(&map[string]string{}, &response) // TODO response headers
		}

		// process response

		responseVisitor := ResponseVisitor{
			resourceData: d,
		}
		if _, err := responseVisitor.accept(response); err != nil {
			return err
		}

		// consume response

		var id string
		if generator.GetId != nil {
			id = generator.GetId(&body, &response)
		} else {
			id = responseVisitor.id
		}

		d.SetId(id)

		return nil
	}

	read := func(d *schema.ResourceData, m interface{}) error {
		api := m.(api.ProviderConfig).Api

		url, err := api.GetUrl(generator.Path)
		if err != nil {
			return err
		}
		fullUrl := fmt.Sprintf("%s/%s", url, d.Id())
		accessToken, err := api.GetAccessToken()
		if err != nil {
			return err
		}
		headers := map[string]string{
			"Authorization": "Bearer " + accessToken,
		}

		// send request

		response, err := generator.Client.Get(fullUrl, headers)
		if err != nil {
			return err
		}

		// modify response

		if generator.ModifyResponseBody != nil {
			generator.ModifyResponseBody(&response)
		}
		if generator.ModifyResponse != nil {
			generator.ModifyResponse(&map[string]string{}, &response) // TODO response headers
		}

		// process response
		responseVisitor := ResponseVisitor{
			resourceData: d,
		}
		if _, err := responseVisitor.accept(response); err != nil {
			return err
		}

		return nil
	}

	update := func(d *schema.ResourceData, m interface{}) error {
		// prepare request

		api := m.(api.ProviderConfig).Api
		var requestVisitor RequestVisitor

		url, err := api.GetUrl(generator.Path)
		if err != nil {
			return err
		}
		fullUrl := fmt.Sprintf("%s/%s", url, d.Id())
		accessToken, err := api.GetAccessToken()
		if err != nil {
			return err
		}
		headers := map[string]string{
			"Authorization": "Bearer " + accessToken,
		}

		body, err := requestVisitor.accept(d)
		if err != nil {
			return err
		}

		// modify request

		if generator.ModifyRequestBody != nil {
			generator.ModifyRequestBody(&body)
		}

		if generator.ModifyRequest != nil {
			generator.ModifyRequest(&url, &headers, &body)
		}

		// send

		response, err := generator.Client.Put(fullUrl, headers, body)
		if err != nil {
			return err
		}

		// modify response

		if generator.ModifyResponseBody != nil {
			generator.ModifyResponseBody(&response)
		}
		if generator.ModifyResponse != nil {
			generator.ModifyResponse(&map[string]string{}, &response) // TODO response headers
		}

		// process response

		responseVisitor := ResponseVisitor{
			resourceData: d,
		}
		if _, err := responseVisitor.accept(response); err != nil {
			return err
		}

		var id string
		if generator.GetId != nil {
			id = generator.GetId(&body, &response)
		} else {
			id = responseVisitor.id
		}

		d.SetId(id)

		return nil
	}

	deleteResource := func(d *schema.ResourceData, m interface{}) error {
		// prepare request

		api := m.(api.ProviderConfig).Api
		url, err := api.GetUrl(generator.Path)
		if err != nil {
			return err
		}
		fullUrl := fmt.Sprintf("%s/%s", url, d.Id())
		accessToken, err := api.GetAccessToken()
		if err != nil {
			return err
		}
		headers := map[string]string{
			"Authorization": "Bearer " + accessToken,
		}

		// send

		err = generator.Client.Delete(fullUrl, headers)

		// consume response

		if err != nil {
			d.SetId("")
		}

		return err
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
