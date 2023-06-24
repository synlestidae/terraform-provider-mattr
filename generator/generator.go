package generator

import (
	"nz.antunovic/mattr-terraform-provider/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Generator struct {
	Path string
	Immutable bool
	Singleton bool
	Schema map[string]*schema.Schema
	Client api.Client
}

func (generator *Generator) GenResource() schema.Resource {

	create := func (d *schema.ResourceData, m interface{}) error {
		//var reqVisitor RequestVisitor

		panic("Not quite implemented")
	}

	read := func (d *schema.ResourceData, m interface{}) error {
		//var reqVisitor RequestVisitor

		panic("Not quite implemented")
	}

	update := func (d *schema.ResourceData, m interface{}) error {
		//var reqVisitor RequestVisitor

		panic("Not quite implemented")
	}

	deleteResource := func (d *schema.ResourceData, m interface{}) error {
		panic("Not quite implemented")
	}

	resource := schema.Resource {
		Create: create,
		Update: update,
		Read: read,
	}

	if !generator.Immutable {
		resource.Delete = deleteResource
	}

	return resource
}
