package gcore

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataClient() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataClientRead,
		Description: "Represent CDN account data",
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
				Description: "ID of your CDN account.",
				Computed:    true,
			},
			"cname": {
				Type:        schema.TypeString,
				Description: "Domain zone to which a CNAME record of your CDN resources should be pointed.",
				Computed:    true,
			},
		},
	}
}

func dataClientRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start reading client")

	config := m.(*Config)
	client := config.CDNClient

	result, err := client.ClientsMe().Get(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] Client received: %v", *result)

	d.SetId(fmt.Sprintf("%d", result.ID))
	d.Set("cname", result.Cname)

	log.Println("[DEBUG] Finish reading client")
	return nil
}
