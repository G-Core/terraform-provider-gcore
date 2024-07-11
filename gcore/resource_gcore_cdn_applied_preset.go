package gcore

import (
	"context"
	"fmt"
	"log"

	"github.com/G-Core/gcorelabscdn-go/presets"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCDNAppliedPreset() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				presetID, objectID, err := ImportAppliedPresetStringParser(d.Id())

				if err != nil {
					return nil, err
				}
				d.Set("preset_id", presetID)
				d.Set("object_id", objectID)
				d.SetId(fmt.Sprintf("%d-%d", presetID, objectID))

				return []*schema.ResourceData{d}, nil
			},
		},
		Schema: map[string]*schema.Schema{
			"preset_id": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "ID of CDN preset which will be applied to the object",
			},
			"object_id": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "ID of CDN object (resource, rule) for which preset will be applied",
			},
		},
		CreateContext: resourceCDNPresetApply,
		ReadContext:   resourceCDNAppliedPresetRead,
		DeleteContext: resourceCDNPresetUnapply,
		Description:   "Represent applied preset",
	}
}

func resourceCDNPresetApply(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start CDN Preset applying")
	config := m.(*Config)
	client := config.CDNClient

	config.CDNMutex.Lock()
	defer config.CDNMutex.Unlock()

	var req presets.ApplyRequest
	req.ObjectID = d.Get("object_id").(int)

	result, err := client.Presets().Apply(ctx, d.Get("preset_id").(int), &req)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d-%d", result.PresetID, result.ObjectID))
	resourceCDNAppliedPresetRead(ctx, d, m)

	log.Printf("[DEBUG] Finish CDN Preset (id=%d) applying to the object (id=%d)\n", result.PresetID, result.ObjectID)
	return nil
}

func resourceCDNAppliedPresetRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	appliedPresetID := d.Id()
	log.Printf("[DEBUG] Start CDN AppliedPreset reading (id=%s)\n", appliedPresetID)
	config := m.(*Config)
	client := config.CDNClient

	appliedPreset, err := client.Presets().GetAppliedPreset(ctx, d.Get("preset_id").(int), d.Get("object_id").(int))
	if err != nil {
		return diag.FromErr(err)
	}

	if appliedPreset == nil {
		log.Printf("[WARN] Preset with id %d is not applied to the object with id %d,"+
			"removing it from the state", d.Get("preset_id").(int), d.Get("object_id").(int))
		d.SetId("")
		return nil
	}

	log.Println("[DEBUG] Finish CDN AppliedPreset reading")
	return nil
}

func resourceCDNPresetUnapply(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	appliedPresetID := d.Id()
	log.Printf("[DEBUG] Start CDN Preset unapplying (id=%s)\n", appliedPresetID)

	config := m.(*Config)
	client := config.CDNClient

	config.CDNMutex.Lock()
	defer config.CDNMutex.Unlock()

	if err := client.Presets().Unapply(ctx, d.Get("preset_id").(int), d.Get("object_id").(int)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	log.Println("[DEBUG] Finish CDN Preset unapplying")
	return nil
}
