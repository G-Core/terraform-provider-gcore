package gcore

import (
	"context"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataPreset() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataPresetRead,
		Description: "Represent presets data",
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
				Description: "ID of the CDN preset",
				Required:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "Name of the CDN preset",
				Required:    true,
			},
		},
	}
}

func dataPresetRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start reading presets")

	config := m.(*Config)
	client := config.CDNClient

	result, err := client.Presets().Get(ctx, d.Get("id").(int))
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] Preset received: %v", *result)

	if result.Name != d.Get("name").(string) {
		return diag.Errorf("The provided preset name \"%s\" doesn't match the received name \"%s\". "+
			"Please check if the specified preset ID is correct or fix the name in your config.", d.Get("name").(string), result.Name)
	}

	d.SetId(strconv.Itoa(result.ID))

	log.Println("[DEBUG] Finish reading presets")
	return nil
}
