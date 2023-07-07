package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	_ "github.com/motemen/go-loghttp/global"
	"nz.antunovic/mattr-terraform-provider/api"
	"nz.antunovic/mattr-terraform-provider/generator"
	"bytes"
	"encoding/json"
  "archive/zip"
	"io/ioutil"
	"fmt"
	"log"
	"path"
	"github.com/k0kubun/pp"
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
            Type:     schema.TypeString,
        },
			},
      "font_paths": &schema.Schema{ 
				Type:     schema.TypeList,
				Elem: &schema.Schema {
					Type: schema.TypeString,
				},
				Optional: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"file_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"title": &schema.Schema{
				Type:     schema.TypeString,
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
				Optional: true,
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

	generator.ModifyRequestBody = func(body *interface{}) error {
    bodyMap := (*body).(map[string]interface{})

		if title, ok := bodyMap["title"].(interface{}); ok {
			delete(bodyMap, "title")
			bodyMap["metadata"] = map[string]string {
				"title": title.(string),
			}
		}

		return nil
	}

	generator.ModifyRequest = func(url *string, headers *map[string]string, body *interface{}) error {
    bodyMap := (*body).(map[string]interface{})

		d := pp.Sprint(bodyMap)
		log.Println("Here is: %s", d);
		log.Println("Here is another: %s", map[string]string{"hello": "world"});

		// these params dont get included
    templatePath := bodyMap["templatePath"].(string)
    fontPaths := castToStringSlice(bodyMap["fontPaths"].([]interface{}))
    delete(bodyMap, "templatePath")
    delete(bodyMap, "fontPaths")


		buffer := new(bytes.Buffer)
    writer := zip.NewWriter(buffer)

    // read the fonts into zip fonts dir
		for _, font := range fontPaths {
			log.Printf("Reading font %s", font)
			fontFile, err := ioutil.ReadFile(font)
			if err != nil {
				return err
			}

			fontWriter, err := writer.Create(fmt.Sprintf("fonts/%s", path.Base(font)))
			if err != nil {
				return err
			}
			if _, err := fontWriter.Write(fontFile); err != nil {
				return err
			}
		}

		// write template.pdf
		log.Printf("Reading template %s", templatePath)
		templateBytes, err := ioutil.ReadFile(templatePath)
		if err != nil {
			return err
		}
		templateWriter, err := writer.Create("template.pdf")
		if err != nil {
			return err
		}
		if _, err := templateWriter.Write(templateBytes); err != nil {
			return err
		}

		bodyJson, err := json.Marshal(bodyMap)
		if err != nil {
			return err
		}
		configWriter, err := writer.Create("config.json")
		if err != nil {
			return err
		}
		if _, err := configWriter.Write(bodyJson); err != nil {
			return err
		}

		err = writer.Close()
		if err != nil {
				return err
		}
		
		(*headers)["Content-Type"] = "application/zip"
		zip := buffer.Bytes()

		if err = ioutil.WriteFile("/tmp/payload.zip", zip, 755); err != nil {
			return err
		}
		log.Printf("Going to upload zip file of %d bytes", len(zip))
		*body = zip 
		return nil
	}


	generator.ModifyResponseBody = func(body *interface{}) error {
		log.Printf("Body is %T", body)
    bodyMap1 := (*body).(*interface{})
    bodyMap := (*bodyMap1).(map[string]interface{})

		metadata := bodyMap["metadata"].(interface{})
		if metadata != nil {
			metadataMap := metadata.(map[string]interface{})
			if metadataMap["title"] != nil {
				title := metadataMap["title"].(string)
				bodyMap["title"] = metadataMap["title"]
				bodyMap["metadata"] = map[string]string {
					"title": title,
				}
				delete(bodyMap, "metadata")
			}
		}

		return nil
	}

	resource := generator.GenResource()

	return &resource
}
