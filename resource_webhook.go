package main

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	_ "github.com/motemen/go-loghttp/global"
)

func resourceWebhook() *schema.Resource {
	return &schema.Resource{
		Create: resourceWebhookCreate,
		Read:   resourceWebhookRead,
		Update: resourceWebhookUpdate,
		Delete: resourceWebhookDelete,
		Schema: map[string]*schema.Schema{
			"events": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type:     schema.TypeString,
					Optional: true,
				},
			},
			"url": schema.Schema{
				Type: schema.TypeString,
			},
			"disabled": schema.Schema{
				Type: schema.TypeString,
			},
		},
	}
}

func resourceWebhookCreate(d *schema.ResourceData, m interface{}) error {
	api := InitFromEnv()

	webhook_request := WebhookRequest{
		events:   d.Get("events").([]string),
		url:      d.Get("url").(string),
		disabled: d.Get("disabled").(bool),
	}

	webhook_response, err := api.PostWebhook(&webhook_request)
	if err != nil {
		return err
	}

	d.SetId(webhook_response.Id)
	d.Set("url", webhook_response.Id)
	d.Set("events", webhook_response.events)
	d.Set("disabled", webhook_response.Disabled)

	return nil
}

func resourceWebhookRead(d *schema.ResourceData, m interface{}) error {
	return fmt.Errorf("Not quite implemented")
}

func resourceWebhookUpdate(d *schema.ResourceData, m interface{}) error {
	return fmt.Errorf("Not quite implemented")
}

func resourceWebhookDelete(d *schema.ResourceData, m interface{}) error {
	return fmt.Errorf("Not quite implemented")
}

func toTerraform(webhook_response *WebhookResponse, d *schema.ResourceData) {
	d.SetId(webhook_response.Id)
	d.Set("url", webhook_response.Id)
	d.Set("events", webhook_response.events)
	d.Set("disabled", webhook_response.Disabled)
}

func fromTerraform(d *schema.ResourceData) WebhookRequest {
	return WebhookRequest{
		events:   d.Get("events").([]string),
		url:      d.Get("url").(string),
		disabled: d.Get("disabled").(bool),
	}
}
