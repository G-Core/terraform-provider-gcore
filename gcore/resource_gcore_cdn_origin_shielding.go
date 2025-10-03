package gcore

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/AlekSi/pointer"
	"github.com/G-Core/gcorelabscdn-go/originshielding"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
		CreateContext: resourceCDNOriginShieldingUpdate,
		ReadContext:   resourceCDNOriginShieldingRead,
		UpdateContext: resourceCDNOriginShieldingUpdate,
		DeleteContext: resourceCDNOriginShieldingDelete,
		Description:   "Represent origin shielding",
	}
}

func resourceCDNOriginShieldingRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	resourceID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.Errorf("Failed reading: wrong input id: %s", d.Id())
	}

	log.Printf("[DEBUG] Start CDN Origin Shielding reading (id=%d)\n", resourceID)
	config := m.(*Config)
	client := config.CDNClient

	result, err := client.OriginShielding().Get(ctx, int64(resourceID))
	if err != nil {
		if isNotFoundError(err) {
			d.SetId("") // Resource not found, remove from state
			return diag.Diagnostics{
				{
					Severity: diag.Warning,
					Summary:  fmt.Sprintf("Origin Shielding with id %s not found, removing from state", resourceID),
				},
			}
		}

		return diag.FromErr(err)
	}

	err = d.Set("shielding_pop", result.ShieldingPop)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Println("[DEBUG] Finish CDN Origin Shielding reading")
	return nil
}

func resourceCDNOriginShieldingUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	resourceID := d.Get("resource_id").(int)
	log.Printf("[DEBUG] Start CDN Origin Shielding updating (id=%d)\n", resourceID)
	config := m.(*Config)
	client := config.CDNClient

	config.CDNMutex.Lock()
	defer config.CDNMutex.Unlock()

	var req originshielding.UpdateRequest
	req.ShieldingPop = pointer.ToInt(d.Get("shielding_pop").(int))

	if _, err := client.OriginShielding().Update(ctx, int64(resourceID), &req); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d", resourceID))
	resourceCDNOriginShieldingRead(ctx, d, m)

	log.Printf("[DEBUG] Finish CDN Origin Shielding updating")
	return nil
}

func resourceCDNOriginShieldingDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	resourceID := d.Get("resource_id").(int)
	log.Printf("[DEBUG] Start CDN Origin Shielding deleting (id=%d)\n", resourceID)
	config := m.(*Config)
	client := config.CDNClient

	config.CDNMutex.Lock()
	defer config.CDNMutex.Unlock()

	var req originshielding.UpdateRequest
	var intPointer *int
	req.ShieldingPop = intPointer

	if _, err := client.OriginShielding().Update(ctx, int64(resourceID), &req); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	log.Printf("[DEBUG] Finish CDN Origin Shielding deleting")
	return nil
}
