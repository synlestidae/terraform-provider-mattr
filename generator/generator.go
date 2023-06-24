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
}

func (generator *Generator) GenResource() schema.Resource {
	create := func(d *schema.ResourceData, m interface{}) error {
		api := m.(api.ProviderConfig).Api
		var requestVisitor RequestVisitor

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

		response, err := generator.Client.Post(url, headers, body)
		if err != nil {
			return err
		}
		responseVisitor := ResponseVisitor{
			resourceData: d,
		}
		if _, err := responseVisitor.accept(response); err != nil {
			return err
		}

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

		response, err := generator.Client.Get(fullUrl, headers)
		if err != nil {
			return err
		}
		responseVisitor := ResponseVisitor{
			resourceData: d,
		}
		if _, err := responseVisitor.accept(response); err != nil {
			return err
		}

		return nil
	}

	update := func(d *schema.ResourceData, m interface{}) error {
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

		response, err := generator.Client.Put(fullUrl, headers, body)
		if err != nil {
			return err
		}
		responseVisitor := ResponseVisitor{
			resourceData: d,
		}
		if _, err := responseVisitor.accept(response); err != nil {
			return err
		}

		return nil
	}

	deleteResource := func(d *schema.ResourceData, m interface{}) error {
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

		return generator.Client.Delete(fullUrl, headers)
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
