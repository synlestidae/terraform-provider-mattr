package provider

import (
	"archive/zip"
	"bytes"
	"fmt"
	"encoding/json"
	"log"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	_ "github.com/motemen/go-loghttp/global"
	"io/ioutil"
	"nz.antunovic/mattr-terraform-provider/api"
	"nz.antunovic/mattr-terraform-provider/generator"
	"path/filepath"
)

type ZipCreator struct {
	writer *zip.Writer
}

func (z *ZipCreator) writeTemplate(path string) error {
	templateBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	templateWriter, err := z.writer.Create("template.pdf")
	if err != nil {
		return err
	}
	if _, err := templateWriter.Write(templateBytes); err != nil {
		return err
	}
	return nil
}

func (z *ZipCreator) writeFontPaths(fonts []interface{}) error {
	for _, font := range fonts {
		fontMap := font.(map[string]interface{})
		fontPath := fontMap["fileName"].(string)
		log.Printf("Reading font file: %s", fontPath)
		fontFile, err := ioutil.ReadFile(fontPath)
		if err != nil {
			log.Printf("Error reading font file.")
			return err
		}

		dstPath := filepath.Base(fontPath)
		fontWriter, err := z.writer.Create(fmt.Sprintf("fonts/%s", dstPath))

		if err != nil {
			log.Printf("Error creating font file %s in ZIP.", fontPath)
			return err
		}
		if _, err := fontWriter.Write(fontFile); err != nil {
			log.Printf("Error writign font file %s to ZIP.", fontPath)
			return err
		}
	}

	return nil
}

func (z *ZipCreator) writeConfig(config interface{}) error {
	configMap := config.(map[string]interface{})
	fontList := configMap["fonts"].([]interface{})

	for _, font := range fontList {
		font := font.(map[string]interface{})
		fontFile := filepath.Base(font["fileName"].(string))
		font["fileName"] = fontFile
	}

	bodyJson, err := json.Marshal(config)
	if err != nil {
		return err
	}
	configWriter, err := z.writer.Create("config.json")
	if err != nil {
		return err
	}
	
	if _, err := configWriter.Write(bodyJson); err != nil {
		return err
	}

	return nil
}


func templateGenerator() generator.Generator {
	generator := generator.Generator{
		Client: &api.HttpClient{},
		Schema: map[string]*schema.Schema{
			"template_path": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
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
		log.Printf("Generating ZIP file for %s", generator.Path)
		bodyMap := (*body).(map[string]interface{})

		buffer := new(bytes.Buffer)

		// these params dont get included
		templatePath := bodyMap["templatePath"].(string)
		delete(bodyMap, "templatePath")
		delete(bodyMap, "fontPaths")

		writer := zip.NewWriter(buffer)
		zipCreator := ZipCreator {
			writer: writer,
		}
		if err := zipCreator.writeTemplate(templatePath); err != nil {
			return err
		}
		if err := zipCreator.writeFontPaths(bodyMap["fonts"].([]interface{})); err != nil {
			return err
		}
		if err := zipCreator.writeConfig(bodyMap); err != nil {
			return err
		}
		if err := writer.Close(); err != nil {
			return err
		}

		bytes := buffer.Bytes()

		log.Printf("Produced a ZIP file of %d byte(s)", len(bytes))
		*body = bytes

		return nil
	}

	return generator
}

func resourceSemanticCompactCredentialTemplate() *schema.Resource {
	generator := templateGenerator()
	generator.Path = "/v2/credentials/compact-semantic/pdf/templates"
	resource := generator.GenResource()

	return &resource
}

func resourceCompactCredentialTemplate() *schema.Resource {
	generator := templateGenerator()
	generator.Path = "/v2/credentials/compact/pdf/templates"
	resource := generator.GenResource()

	return &resource
}
