package gcore

import (
	"context"
	"github.com/G-Core/gcorelabscdn-go/originshielding"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
)

func resourceCDNOriginShielding() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"resource_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "ID of CDN resource for which shielding will be applied",
			},
			"shielding_pop": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "ID of the shielding point of present",
			},
		},
		ReadContext:   resourceCDNOriginShieldingRead,
		UpdateContext: resourceCDNOriginShieldingUpdate,
		Description:   "Represent origin shielding",
	}
}

func resourceCDNOriginShieldingRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	resourceID := d.Get("resource_id").(int)
	log.Printf("[DEBUG] Start CDN Origin Shielding reading (id=%s)\n", resourceID)
	config := m.(*Config)
	client := config.CDNClient

	result, err := client.OriginShielding().Get(ctx, int64(resourceID))
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("shielding_pop", result.ShieldingPop)

	log.Println("[DEBUG] Finish CDN Origin Shielding reading")
	return nil
}

func resourceCDNOriginShieldingUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	resourceID := d.Get("resource_id").(int)
	log.Printf("[DEBUG] Start CDN Origin Shielding updating (id=%s)\n", resourceID)
	config := m.(*Config)
	client := config.CDNClient

	var req originshielding.UpdateRequest
	req.ShieldingPop = d.Get("shielding_pop").(int)

	if _, err := client.OriginShielding().Update(ctx, int64(resourceID), &req); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] Finish CDN Origin Shielding updating")

	return resourceCDNOriginShieldingRead(ctx, d, m)
}
