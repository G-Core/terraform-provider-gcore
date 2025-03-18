package gcore

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/G-Core/gcorelabscdn-go/logsuploader"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func schemaForHTTPAction() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"method": {
				Type:     schema.TypeString,
				Required: true,
			},
			"url": {
				Type:     schema.TypeString,
				Required: true,
			},
			"headers": {
				Type:     schema.TypeMap,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"timeout_seconds": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"use_compression": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"response_actions": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"match_status_code": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"match_payload": {
							Type:     schema.TypeString,
							Required: true,
						},
						"description": {
							Type:     schema.TypeString,
							Required: true,
						},
						"action": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func schemaForS3Amazon() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"access_key_id": {
					Type:     schema.TypeString,
					Required: true,
				},
				"secret_access_key": {
					Type:     schema.TypeString,
					Required: true,
				},
				"region": {
					Type:     schema.TypeString,
					Required: true,
				},
				"bucket_name": {
					Type:     schema.TypeString,
					Required: true,
				},
				"directory": {
					Type:     schema.TypeString,
					Optional: true,
				},
			},
		},
	}
}

func schemaForS3Oss() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"access_key_id": {
					Type:     schema.TypeString,
					Required: true,
				},
				"secret_access_key": {
					Type:     schema.TypeString,
					Required: true,
				},
				"region": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"bucket_name": {
					Type:     schema.TypeString,
					Required: true,
				},
				"directory": {
					Type:     schema.TypeString,
					Optional: true,
				},
			},
		},
	}
}

func schemaForS3Other() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"access_key_id": {
					Type:     schema.TypeString,
					Required: true,
				},
				"region": {
					Type:     schema.TypeString,
					Required: true,
				},
				"bucket_name": {
					Type:     schema.TypeString,
					Required: true,
				},
				"directory": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"endpoint": {
					Type:     schema.TypeString,
					Required: true,
				},
				"use_path_style": {
					Type:     schema.TypeBool,
					Required: true,
				},
			},
		},
	}
}

func schemaForFTP() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"user": {
					Type:     schema.TypeString,
					Required: true,
				},
				"hostname": {
					Type:     schema.TypeString,
					Required: true,
				},
				"timeout_seconds": {
					Type:     schema.TypeInt,
					Optional: true,
					Computed: true,
				},
				"directory": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
				"password": {
					Type:     schema.TypeString,
					Required: true,
				},
			},
		},
	}
}

func schemaForSFTP() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"user": {
					Type:     schema.TypeString,
					Required: true,
				},
				"hostname": {
					Type:     schema.TypeString,
					Required: true,
				},
				"timeout_seconds": {
					Type:     schema.TypeInt,
					Optional: true,
					Computed: true,
				},
				"directory": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"password": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"private_key": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"key_passphrase": {
					Type:     schema.TypeString,
					Optional: true,
				},
			},
		},
	}
}

func schemaForHTTP() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"upload": {
					Type:     schema.TypeList,
					Required: true,
					Elem:     schemaForHTTPAction(),
				},
				"append": {
					Type:     schema.TypeList,
					Optional: true,
					Elem:     schemaForHTTPAction(),
				},
				"retry": {
					Type:     schema.TypeList,
					Optional: true,
					Elem:     schemaForHTTPAction(),
				},
				"auth": {
					Type:     schema.TypeList,
					Required: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"type": {
								Type:     schema.TypeString,
								Required: true,
							},
							"config": {
								Type:     schema.TypeList,
								Required: true,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"header_name": {
											Type:     schema.TypeString,
											Required: true,
										},
										"token": {
											Type:     schema.TypeString,
											Required: true,
										},
									},
								},
							},
						},
					},
				},
				"payload_type": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
			},
		},
	}
}

func resourceCDNLogsUploaderTarget() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Name of the target.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the target.",
			},
			"storage_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"s3_gcore", "s3_amazon", "s3_oss", "s3_other", "s3_v1", "ftp", "sftp", "http"}, false),
				Description:  "Type of storage for logs.",
			},
			"config": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Config for specific storage type.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"s3_gcore":  schemaForS3Other(),
						"s3_amazon": schemaForS3Amazon(),
						"s3_oss":    schemaForS3Oss(),
						"s3_other":  schemaForS3Other(),
						"s3_v1":     schemaForS3Other(),
						"ftp":       schemaForFTP(),
						"sftp":      schemaForSFTP(),
						"http":      schemaForHTTP(),
					},
				},
			},
		},
		CreateContext: resourceCDNLogsUploaderTargetCreate,
		ReadContext:   resourceCDNLogsUploaderTargetRead,
		UpdateContext: resourceCDNLogsUploaderTargetUpdate,
		DeleteContext: resourceCDNLogsUploaderTargetDelete,
		CustomizeDiff: customizeDiffStorageTypeConfig,
		Description:   "Represent CDN logs uploader target",
	}
}

func customizeDiffStorageTypeConfig(ctx context.Context, diff *schema.ResourceDiff, v interface{}) error {
	storageType := diff.Get("storage_type").(string)
	configList := diff.Get("config").([]interface{})
	config := configList[0].(map[string]interface{})

	for key, value := range config {
		value := value.([]interface{})
		if len(value) != 0 && key != storageType {
			return fmt.Errorf("invalid key %s in config for storage_type %s", key, storageType)
		}
	}

	return nil
}

func resourceCDNLogsUploaderTargetCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start CDN Logs Uploader Target creating")
	config := m.(*Config)
	client := config.CDNClient

	var req logsuploader.TargetCreateRequest

	if v, ok := d.GetOk("name"); ok {
		req.Name = v.(string)
	}

	if v, ok := d.GetOk("description"); ok {
		req.Description = v.(string)
	}

	req.StorageType = logsuploader.StorageType(d.Get("storage_type").(string))

	configList := d.Get("config").([]interface{})
	configAttr := configList[0].(map[string]interface{})
	for _, value := range configAttr {
		value := value.([]interface{})
		if len(value) != 0 {
			req.Config = value[0].(map[string]interface{})
			break
		}
	}

	result, err := client.LogsUploader().TargetCreate(ctx, &req)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d", result.ID))
	resourceCDNLogsUploaderTargetRead(ctx, d, m)

	log.Printf("[DEBUG] Finish CDN Logs Uploader Target creating (id=%d)\n", result.ID)
	return nil
}

func resourceCDNLogsUploaderTargetRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	targetID := d.Id()
	log.Printf("[DEBUG] Start CDN Logs Uploader Target reading (id=%s)\n", targetID)
	config := m.(*Config)
	client := config.CDNClient

	id, err := strconv.ParseInt(targetID, 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}

	result, err := client.LogsUploader().TargetGet(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("name", result.Name)
	d.Set("description", result.Description)
	d.Set("name", result.Name)
	d.Set("storage_type", result.StorageType)
	d.Set("config", result.Config)

	log.Println("[DEBUG] Finish CDN Logs Uploader Target reading value s%", d)
	log.Println("[DEBUG] Finish CDN Logs Uploader Target reading")
	return nil
}

func resourceCDNLogsUploaderTargetUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	targetID := d.Id()
	log.Printf("[DEBUG] Start CDN Logs Uploader Target updating (id=%s)\n", targetID)
	config := m.(*Config)
	client := config.CDNClient

	id, err := strconv.ParseInt(targetID, 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}

	var req logsuploader.TargetUpdateRequest

	if v, ok := d.GetOk("name"); ok {
		req.Name = v.(string)
	}

	if v, ok := d.GetOk("description"); ok {
		req.Description = v.(string)
	}

	req.StorageType = logsuploader.StorageType(d.Get("storage_type").(string))

	configList := d.Get("config").([]interface{})
	configAttr := configList[0].(map[string]interface{})
	for _, value := range configAttr {
		value := value.([]interface{})
		if len(value) != 0 {
			req.Config = value[0].(map[string]interface{})
			break
		}
	}

	if _, err := client.LogsUploader().TargetUpdate(ctx, id, &req); err != nil {
		return diag.FromErr(err)
	}

	log.Println("[DEBUG] Finish CDN Logs Uploader Target updating")
	return resourceCDNLogsUploaderTargetRead(ctx, d, m)
}

func resourceCDNLogsUploaderTargetDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	targetID := d.Id()
	log.Printf("[DEBUG] Start CDN Logs Uploader Target deleting (id=%s)\n", targetID)
	config := m.(*Config)
	client := config.CDNClient

	id, err := strconv.ParseInt(targetID, 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := client.LogsUploader().TargetDelete(ctx, id); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	log.Println("[DEBUG] Finish CDN Logs Uploader Target deleting")
	return nil
}
