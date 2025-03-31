package gcore

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/G-Core/gcorelabscdn-go/logsuploader"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCDNLogsUploaderConfig() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Enables or disables the config. Default value is true.",
				Default:     true,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the config.",
			},
			"policy": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "ID of the policy that should be assigned to given config.",
			},
			"target": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "ID of the target to which logs should be uploaded.",
			},
			"for_all_resources": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "If set to true, the config will be applied to all CDN resources. If set to false, the config will be applied to the resources specified in the resources field. Default value is false.",
				Default:     false,
			},
			"resources": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Description: "List of resource IDs to which the config should be applied.",
			},
		},
		CreateContext: resourceCDNLogsUploaderConfigCreate,
		ReadContext:   resourceCDNLogsUploaderConfigRead,
		UpdateContext: resourceCDNLogsUploaderConfigUpdate,
		DeleteContext: resourceCDNLogsUploaderConfigDelete,
		Description:   "Represent CDN logs uploader config",
	}
}

func resourceCDNLogsUploaderConfigCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start CDN Logs Uploader Config creating")
	config := m.(*Config)
	client := config.CDNClient

	var req logsuploader.ConfigCreateRequest

	if v, ok := d.GetOk("enabled"); ok {
		req.Enabled = v.(bool)
	}

	if v, ok := d.GetOk("name"); ok {
		req.Name = v.(string)
	}

	if v, ok := d.GetOk("policy"); ok {
		req.Policy = int64(v.(int))
	}

	if v, ok := d.GetOk("target"); ok {
		req.Target = int64(v.(int))
	}

	if v, ok := d.GetOk("for_all_resources"); ok {
		req.ForAllResources = v.(bool)
	}

	if v, ok := d.GetOk("resources"); ok {
		resources := make([]int64, len(v.([]interface{})))
		for i, v := range v.([]interface{}) {
			resources[i] = int64(v.(int))
		}
		req.Resources = resources
	} else {
		req.Resources = make([]int64, 0)
	}

	result, err := client.LogsUploader().ConfigCreate(ctx, &req)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d", result.ID))
	resourceCDNLogsUploaderConfigRead(ctx, d, m)

	log.Printf("[DEBUG] Finish CDN Logs Uploader Config creating (id=%d)\n", result.ID)
	return nil
}

func resourceCDNLogsUploaderConfigRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	configID := d.Id()
	log.Printf("[DEBUG] Start CDN Logs Uploader Config reading (id=%s)\n", configID)
	config := m.(*Config)
	client := config.CDNClient

	id, err := strconv.ParseInt(configID, 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}

	result, err := client.LogsUploader().ConfigGet(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("enabled", result.Enabled)
	d.Set("name", result.Name)
	d.Set("policy", result.Policy)
	d.Set("target", result.Target)
	d.Set("for_all_resources", result.ForAllResources)
	d.Set("resources", result.Resources)

	log.Println("[DEBUG] Finish CDN Logs Uploader Config reading")
	return nil
}

func resourceCDNLogsUploaderConfigUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	configID := d.Id()
	log.Printf("[DEBUG] Start CDN Logs Uploader Config updating (id=%s)\n", configID)
	config := m.(*Config)
	client := config.CDNClient

	id, err := strconv.ParseInt(configID, 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}

	var req logsuploader.ConfigUpdateRequest

	if v, ok := d.GetOk("enabled"); ok {
		req.Enabled = v.(bool)
	}

	if v, ok := d.GetOk("name"); ok {
		req.Name = v.(string)
	}

	if v, ok := d.GetOk("policy"); ok {
		req.Policy = int64(v.(int))
	}

	if v, ok := d.GetOk("target"); ok {
		req.Target = int64(v.(int))
	}

	if v, ok := d.GetOk("for_all_resources"); ok {
		req.ForAllResources = v.(bool)
	}

	if v, ok := d.GetOk("resources"); ok {
		resources := make([]int64, len(v.([]interface{})))
		for i, v := range v.([]interface{}) {
			resources[i] = int64(v.(int))
		}
		req.Resources = resources
	} else {
		req.Resources = make([]int64, 0)
	}

	if _, err := client.LogsUploader().ConfigUpdate(ctx, id, &req); err != nil {
		return diag.FromErr(err)
	}

	log.Println("[DEBUG] Finish CDN Logs Uploader Config updating")
	return resourceCDNLogsUploaderConfigRead(ctx, d, m)
}

func resourceCDNLogsUploaderConfigDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	configID := d.Id()
	log.Printf("[DEBUG] Start CDN Logs Uploader Config deleting (id=%s)\n", configID)
	config := m.(*Config)
	client := config.CDNClient

	id, err := strconv.ParseInt(configID, 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := client.LogsUploader().ConfigDelete(ctx, id); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	log.Println("[DEBUG] Finish CDN Logs Uploader Config deleting")
	return nil
}
