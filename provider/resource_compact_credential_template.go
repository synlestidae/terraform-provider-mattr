package provider

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	_ "github.com/motemen/go-loghttp/global"
	"io/ioutil"
	"nz.antunovic/mattr-terraform-provider/api"
	"nz.antunovic/mattr-terraform-provider/generator"
)

func resourceCompactCredentialTemplate() *schema.Resource {
	generator := generator.Generator{
		Path:   "/v2/credentials/compact/pdf/templates",
		Client: &api.HttpClient{},
		Schema: map[string]*schema.Schema{
			"template_path": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"font_paths": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"file_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"metadata": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
			},
			"fonts": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"file_name": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"fields": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"value": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"is_required": &schema.Schema{
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"alternative_text": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"font_name": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}

	generator.ModifyRequest = func(url *string, headers *map[string]string, body *interface{}) error {
		//config := make(map[string]interface{})
		bodyMap := (*body).(map[string]interface{})

		buffer := new(bytes.Buffer)

		// these params dont get included
		templatePath := bodyMap["templatePath"].(string)
		fontPaths := bodyMap["fontPaths"].([]string)
		delete(bodyMap, "templatePath")
		delete(bodyMap, "fontPaths")

		writer := zip.NewWriter(buffer)
		defer writer.Close()

		// read the fonts into zip fonts dir
		for _, font := range fontPaths {
			fontFile, err := ioutil.ReadFile(font)
			if err != nil {
				return err
			}

			fontWriter, err := writer.Create(fmt.Sprintf("fonts/%s", fontFile))
			if err != nil {
				return err
			}
			//defer fontWriter.Close()
			if _, err := fontWriter.Write(fontFile); err != nil {
				return err
			}
		}

		// write template.pdf
		templateBytes, err := ioutil.ReadFile(templatePath)
		if err != nil {
			return err
		}
		templateWriter, err := writer.Create("template.pdf")
		if err != nil {
			return err
		}
		//defer templateWriter.close()
		if _, err := templateWriter.Write(templateBytes); err != nil {
			return err
		}

		// now serialize the JSON
		bodyJson, err := json.Marshal(bodyMap)
		if err != nil {
			return err
		}
		configWriter, err := writer.Create("config.json")
		if err != nil {
			return err
		}
		//defer configWriter.Close()
		if _, err := configWriter.Write(bodyJson); err != nil {
			return err
		}

		(*headers)["Content-Type"] = "application/zip"
		*body = buffer.Bytes()

		return nil
	}

	resource := generator.GenResource()

	return &resource
}
