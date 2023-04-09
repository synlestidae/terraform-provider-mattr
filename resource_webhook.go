package main

import (
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
					Type: schema.TypeString,
				},
				Description: "Types of events we will look out for and send to the webhook",
				Required:    true,
			},
			"url": &schema.Schema{
				Type:        schema.TypeString,
				Description: "URL of the webhook, to which event payloads are delivered",
				Required:    true,
			},
			"disabled": &schema.Schema{
				Type:        schema.TypeBool,
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
	api := InitFromEnv()
	id := d.Id()
	webhook, err := api.GetWebhook(id)
	if err != nil {
		return err
	}
	processWebhookData(webhook, d)

	return nil
}

func resourceWebhookUpdate(d *schema.ResourceData, m interface{}) error {
	api := InitFromEnv()
	id := d.Id()
	webhook_request := fromTerraform(d)
	webhook, err := api.PutWebhook(id, &webhook_request)
	if err != nil {
		return err
	}
	processWebhookData(webhook, d)
	return nil
}

func resourceWebhookDelete(d *schema.ResourceData, m interface{}) error {
	api := InitFromEnv()
	id := d.Id()
	err := api.DeleteWebhook(id)
	if err != nil {
		return err
	}
	return nil
}

func processWebhookData(webhook_response *WebhookResponse, d *schema.ResourceData) {
	d.SetId(webhook_response.Id)
	d.Set("url", webhook_response.Url)
	d.Set("events", webhook_response.Events)
	d.Set("disabled", webhook_response.Disabled)
}

func fromTerraform(d *schema.ResourceData) WebhookRequest {
	event_list := d.Get("events").([]interface{})
	events := make([]string, len(event_list), len(event_list))

	for i, event := range event_list {
		events[i] = event.(string)
	}

	return WebhookRequest{
		Events:   events,
		Url:      d.Get("url").(string),
		Disabled: d.Get("disabled").(bool),
	}
}
