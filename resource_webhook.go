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
				Description: "Types of events we will look out for and send to the webhook",
			},
			"url": &schema.Schema{
				Type:        schema.TypeString,
				Description: "URL of the webhook, to which event payloads are delivered",
			},
			"disabled": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "If true, the webhook is disabled.",
			},
		},
	}
}

func resourceWebhookCreate(d *schema.ResourceData, m interface{}) error {
	api := InitFromEnv()
	webhook_request := fromTerraform(d)
	webhook_response, err := api.PostWebhook(&webhook_request)
	if err != nil {
		return err
	}
	processWebhookData(webhook_response, d)
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

func processWebhookData(webhook_response *WebhookResponse, d *schema.ResourceData) {
	d.SetId(webhook_response.Id)
	d.Set("url", webhook_response.Id)
	d.Set("events", webhook_response.Events)
	d.Set("disabled", webhook_response.Disabled)
}

func fromTerraform(d *schema.ResourceData) WebhookRequest {
	return WebhookRequest{
		Events:   d.Get("events").([]string),
		Url:      d.Get("url").(string),
		Disabled: d.Get("disabled").(bool),
	}
}
