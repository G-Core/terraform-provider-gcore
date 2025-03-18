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

func schemaForHTTPAction() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"method": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "POST",
			},
			"url": {
				Type:     schema.TypeString,
				Required: true,
			},
			"headers": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"timeout_seconds": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  30,
			},
			"use_compression": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"response_actions": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"match_status_code": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  0,
						},
						"match_payload": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "",
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "",
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
					Type:      schema.TypeString,
					Required:  true,
					Sensitive: true,
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
					Type:      schema.TypeString,
					Required:  true,
					Sensitive: true,
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
				"secret_access_key": {
					Type:      schema.TypeString,
					Required:  true,
					Sensitive: true,
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
					Optional: true,
					Default:  true,
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
					Default:  10,
				},
				"directory": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"password": {
					Type:      schema.TypeString,
					Required:  true,
					Sensitive: true,
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
					Default:  10,
				},
				"directory": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"password": {
					Type:      schema.TypeString,
					Optional:  true,
					Sensitive: true,
				},
				"private_key": {
					Type:      schema.TypeString,
					Optional:  true,
					Sensitive: true,
				},
				"key_passphrase": {
					Type:      schema.TypeString,
					Optional:  true,
					Sensitive: true,
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
											Type:      schema.TypeString,
											Required:  true,
											Sensitive: true,
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
					Default:  "text",
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
				Description: "Name of the target.",
				Default:     "Target",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the target.",
				Default:     "",
			},
			"config": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Config for specific storage type.",
				MinItems:    1,
				MaxItems:    1,
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
	configList := diff.Get("config").([]interface{})
	config := configList[0].(map[string]interface{})
	counter := 0

	for _, value := range config {
		value := value.([]interface{})
		if len(value) != 0 {
			counter += 1
		}
	}

	if counter != 1 {
		return fmt.Errorf("Only one storage type should be provided in the 'config' field.")
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

	configList := d.Get("config").([]interface{})
	configAttr := configList[0].(map[string]interface{})
	for key, value := range configAttr {
		value := value.([]interface{})
		if len(value) != 0 {
			req.StorageType = logsuploader.StorageType(key)
			req.Config = sanitizeConfig(value[0].(map[string]interface{}))
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

	mergedConfig := mergeStateConfig(result, d)
	configData := make(map[string]interface{})
	storageTypes := []string{"s3_gcore", "s3_amazon", "s3_oss", "s3_other", "s3_v1", "ftp", "sftp", "http"}
	for _, storageType := range storageTypes {
		if storageType == string(result.StorageType) {
			configData[storageType] = []interface{}{mergedConfig}
		} else {
			configData[storageType] = []interface{}{}
		}
	}
	configList := make([]interface{}, 1)
	configList[0] = configData
	d.Set("config", configList)

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

	configList := d.Get("config").([]interface{})
	configAttr := configList[0].(map[string]interface{})
	for key, value := range configAttr {
		value := value.([]interface{})
		if len(value) != 0 {
			req.StorageType = logsuploader.StorageType(key)
			req.Config = sanitizeConfig(value[0].(map[string]interface{}))
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

func sanitizeConfig(rawConfig map[string]interface{}) map[string]interface{} {
	// Iterate over the raw config and ensure empty values are set to nil
	// because tf provider SDK v2 doesn't allow to set nil values for types like string
	// and API expects nil values for optional fields instead of empty strings.
	sanitizedConfig := make(map[string]interface{})
	for k, v := range rawConfig {
		if v == "" || v == nil {
			sanitizedConfig[k] = nil
		} else {
			sanitizedConfig[k] = v
		}
	}

	return sanitizedConfig
}

func mergeStateConfig(result *logsuploader.Target, d *schema.ResourceData) map[string]interface{} {
	// For each field that's in the state but not in the result, set it to value from the state.
	// This is necessary because the state may contain fields that are not returned by the API like password.
	// Otherwise terraform would always show the password as a diff, because it would be set as empty in the result.
	cleanedConfig := make(map[string]interface{})
	if configList, ok := d.Get("config").([]interface{}); ok && len(configList) > 0 {
		if configMap, ok := configList[0].(map[string]interface{}); ok {
			stateConfig := configMap[string(result.StorageType)].([]interface{})[0].(map[string]interface{})
			for k, v := range stateConfig {
				cleanedConfig[k] = v
			}
		}
	}
	for k, v := range result.Config {
		cleanedConfig[k] = v
	}
	return cleanedConfig
}
