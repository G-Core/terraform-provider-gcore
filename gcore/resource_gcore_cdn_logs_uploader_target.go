package gcore

import (
	"context"
	"fmt"
	"log"
	"reflect"
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
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "HTTP method to use for the request. Supported values are 'POST' and 'PUT'.",
				Default:      "POST",
				ValidateFunc: validation.StringInSlice([]string{"POST", "PUT"}, false),
			},
			"url": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "URL to which logs should be uploaded.",
			},
			"headers": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"timeout_seconds": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Timeout for the HTTP request in seconds. Default value is 30.",
				Default:     30,
			},
			"use_compression": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Default value is false.",
				Default:     false,
			},
			"response_actions": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of actions to perform based on the response from the server.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"match_status_code": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "HTTP status code to match. If not specified, no status code will be matched.",
							Default:     0,
						},
						"match_payload": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Payload to match in the response. If not specified, no payload will be matched.",
							Default:     "",
						},
						"description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Default value is empty string.",
							Default:     "",
						},
						"action": {
							Type:         schema.TypeString,
							Description:  "Action to perform if the response matches the specified criteria. Supported values are 'drop', 'retry', and 'append'.",
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"drop", "retry", "append"}, false),
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
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"access_key_id": {
					Type:        schema.TypeString,
					Description: "Access key ID for the Amazon S3 account.",
					Required:    true,
				},
				"secret_access_key": {
					Type:        schema.TypeString,
					Description: "Secret access key for the Amazon S3 account.",
					Required:    true,
					Sensitive:   true,
				},
				"region": {
					Type:        schema.TypeString,
					Description: "Region of the Amazon S3 bucket.",
					Required:    true,
				},
				"bucket_name": {
					Type:        schema.TypeString,
					Description: "Name of the Amazon S3 bucket.",
					Required:    true,
				},
				"directory": {
					Type:        schema.TypeString,
					Description: "Directory in the Amazon S3 bucket where logs will be uploaded.",
					Optional:    true,
				},
			},
		},
	}
}

func schemaForS3Oss() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"access_key_id": {
					Type:        schema.TypeString,
					Description: "Access key ID for the OSS account.",
					Required:    true,
				},
				"secret_access_key": {
					Type:        schema.TypeString,
					Description: "Secret access key for the OSS account.",
					Required:    true,
					Sensitive:   true,
				},
				"region": {
					Type:        schema.TypeString,
					Description: "Region of the OSS bucket.",
					Optional:    true,
				},
				"bucket_name": {
					Type:        schema.TypeString,
					Description: "Name of the OSS bucket.",
					Required:    true,
				},
				"directory": {
					Type:        schema.TypeString,
					Description: "Directory in the OSS bucket where logs will be uploaded.",
					Optional:    true,
				},
			},
		},
	}
}

func schemaForS3Other() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"access_key_id": {
					Type:        schema.TypeString,
					Description: "Access key ID for the S3-compatible storage account.",
					Required:    true,
				},
				"secret_access_key": {
					Type:        schema.TypeString,
					Description: "Secret access key for the S3-compatible storage account.",
					Required:    true,
					Sensitive:   true,
				},
				"region": {
					Type:        schema.TypeString,
					Description: "Region of the S3-compatible storage bucket.",
					Required:    true,
				},
				"bucket_name": {
					Type:        schema.TypeString,
					Description: "Name of the S3-compatible storage bucket.",
					Required:    true,
				},
				"directory": {
					Type:        schema.TypeString,
					Description: "Directory in the S3-compatible storage bucket where logs will be uploaded.",
					Optional:    true,
				},
				"endpoint": {
					Type:        schema.TypeString,
					Description: "Endpoint of the S3-compatible storage service.",
					Required:    true,
				},
				"use_path_style": {
					Type:        schema.TypeBool,
					Optional:    true,
					Description: "Default value is true.",
					Default:     true,
				},
			},
		},
	}
}

func schemaForFTP() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"user": {
					Type:        schema.TypeString,
					Description: "Username for the FTP account.",
					Required:    true,
				},
				"hostname": {
					Type:        schema.TypeString,
					Description: "Hostname or IP address of the FTP server.",
					Required:    true,
				},
				"timeout_seconds": {
					Type:        schema.TypeInt,
					Description: "Timeout for the FTP connection in seconds. Default value is 10.",
					Optional:    true,
					Default:     10,
				},
				"directory": {
					Type:        schema.TypeString,
					Description: "Directory on the FTP server where logs will be uploaded.",
					Optional:    true,
				},
				"password": {
					Type:        schema.TypeString,
					Description: "Password for the FTP account.",
					Required:    true,
					Sensitive:   true,
				},
			},
		},
	}
}

func schemaForSFTP() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"user": {
					Type:        schema.TypeString,
					Description: "Username for the SFTP account.",
					Required:    true,
				},
				"hostname": {
					Type:        schema.TypeString,
					Description: "Hostname or IP address of the SFTP server.",
					Required:    true,
				},
				"timeout_seconds": {
					Type:        schema.TypeInt,
					Description: "Timeout for the SFTP connection in seconds. Default value is 10.",
					Optional:    true,
					Default:     10,
				},
				"directory": {
					Type:        schema.TypeString,
					Description: "Directory on the SFTP server where logs will be uploaded.",
					Optional:    true,
				},
				"password": {
					Type:        schema.TypeString,
					Description: "Password for the SFTP account.",
					Optional:    true,
					Sensitive:   true,
				},
				"private_key": {
					Type:        schema.TypeString,
					Description: "Private key for the SFTP account.",
					Optional:    true,
					Sensitive:   true,
				},
				"key_passphrase": {
					Type:        schema.TypeString,
					Description: "Passphrase for the private key.",
					Optional:    true,
					Sensitive:   true,
				},
			},
		},
	}
}

func schemaForHTTP() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"upload": {
					Type:        schema.TypeList,
					Description: "Upload action configuration.",
					Required:    true,
					MinItems:    1,
					MaxItems:    1,
					Elem:        schemaForHTTPAction(),
				},
				"append": {
					Type:        schema.TypeList,
					Description: "Append action configuration.",
					Optional:    true,
					MaxItems:    1,
					Elem:        schemaForHTTPAction(),
				},
				"retry": {
					Type:        schema.TypeList,
					Description: "Retry action configuration.",
					Optional:    true,
					MaxItems:    1,
					Elem:        schemaForHTTPAction(),
				},
				"auth": {
					Type:        schema.TypeList,
					Description: "Authentication configuration for HTTP target.",
					Optional:    true,
					MinItems:    1,
					MaxItems:    1,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"type": {
								Type:         schema.TypeString,
								Description:  "Type of authentication. Supported values are 'token'.",
								Required:     true,
								ValidateFunc: validation.StringInSlice([]string{"token"}, false),
							},
							"config": {
								Type:        schema.TypeList,
								Description: "Configuration for the authentication type.",
								Required:    true,
								MinItems:    1,
								MaxItems:    1,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"header_name": {
											Type:        schema.TypeString,
											Description: "Name of the header to which the token will be added.",
											Required:    true,
										},
										"token": {
											Type:        schema.TypeString,
											Description: "Token to be used for authentication.",
											Required:    true,
											Sensitive:   true,
										},
									},
								},
							},
						},
					},
				},
				"content_type": {
					Type:         schema.TypeString,
					Optional:     true,
					Description:  "Content type of the logs being uploaded. Supported values are 'json' and 'text'.",
					ValidateFunc: validation.StringInSlice([]string{"json", "text"}, false),
					Default:      "text",
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
				Description: "Name of the target. Default value is \"Target\".",
				Default:     "Target",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the target. Default value is empty string.",
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
	if configList[0] == nil {
		return fmt.Errorf("the 'config' field must contain at least one configuration")
	}

	config := configList[0].(map[string]interface{})
	counter := 0
	for _, value := range config {
		value := value.([]interface{})
		if len(value) != 0 {
			counter += 1
		}
	}

	if counter != 1 {
		return fmt.Errorf("only one storage type should be provided in the 'config' field")
	}

	return nil
}

func resourceCDNLogsUploaderTargetCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start CDN Logs Uploader Target creating")
	config := m.(*Config)
	client := config.CDNClient

	var req logsuploader.TargetCreateRequest
	req.Name = d.Get("name").(string)
	req.Description = d.Get("description").(string)

	configList := d.Get("config").([]interface{})
	configAttr := configList[0].(map[string]interface{})
	for key, value := range configAttr {
		value := value.([]interface{})
		if len(value) != 0 {
			configDict := value[0].(map[string]interface{})
			if key == "http" {
				configDict = convertHttpConfigToDict(configDict)
			}
			req.Config = sanitizeConfig(configDict)
			req.StorageType = logsuploader.StorageType(key)
			break
		}
	}

	result, err := client.LogsUploader().TargetCreate(ctx, &req)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d", result.ID))
	log.Printf("[DEBUG] Finish CDN Logs Uploader Target creating (id=%d)\n", result.ID)
	return resourceCDNLogsUploaderTargetRead(ctx, d, m)
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
		if isNotFoundError(err) {
			d.SetId("") // Resource not found, remove from state
			return diag.Diagnostics{
				{
					Severity: diag.Warning,
					Summary:  fmt.Sprintf("Logs Uploader Target \"%s\" not found, removing from state", d.Get("name")),
				},
			}
		}

		return diag.FromErr(err)
	}

	d.Set("name", result.Name)
	d.Set("description", result.Description)

	mergedConfig := mergeStateConfig(result, d)
	configData := make(map[string]interface{})
	storageTypes := []string{"s3_gcore", "s3_amazon", "s3_oss", "s3_other", "s3_v1", "ftp", "sftp", "http"}
	for _, storageType := range storageTypes {
		if storageType == string(result.StorageType) {
			if storageType == "http" {
				// Convert each dict under "http" to a list (as provider sdk doesn't support dicts)
				for nestedKey, nestedValue := range mergedConfig {
					if nestedKey == "auth" {
						authConfig := nestedValue.(map[string]interface{})
						for authNestedKey, authNestedValue := range authConfig {
							authConfig[authNestedKey] = dictToList(authNestedValue)
						}
					}
					mergedConfig[nestedKey] = dictToList(nestedValue)
				}
			}
			configData[storageType] = []interface{}{mergedConfig}
		} else {
			configData[storageType] = []interface{}{}
		}
	}
	err = d.Set("config", []interface{}{configData})
	if err != nil {
		return diag.Errorf("failed to save 'config' field in state: %s", err)
	}

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
	req.Name = d.Get("name").(string)
	req.Description = d.Get("description").(string)

	configList := d.Get("config").([]interface{})
	configAttr := configList[0].(map[string]interface{})
	for key, value := range configAttr {
		value := value.([]interface{})
		if len(value) != 0 {
			configDict := value[0].(map[string]interface{})
			if key == "http" {
				configDict = convertHttpConfigToDict(configDict)
			}
			req.Config = sanitizeConfig(configDict)
			req.StorageType = logsuploader.StorageType(key)
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
			stateConfigList := configMap[string(result.StorageType)].([]interface{})
			if len(stateConfigList) > 0 {
				stateConfig := stateConfigList[0].(map[string]interface{})
				for k, v := range stateConfig {
					cleanedConfig[k] = v
				}
			}
		}
	}
	for k, v := range result.Config {
		// Sensitive values like passwords are returned as "*****" from the API. We want to use values from
		// the state instead to avoid unnecessary diffs.
		if v != "*****" {
			cleanedConfig[k] = v
		}
	}
	return cleanedConfig
}

func listToDict(nestedValue interface{}) interface{} {
	if reflect.TypeOf(nestedValue).Kind() == reflect.Slice {
		extractedValue := nestedValue.([]interface{})
		if len(extractedValue) > 0 {
			return listToDict(extractedValue[0])
		} else {
			return nil
		}
	}
	return nestedValue
}

func dictToList(value interface{}) interface{} {
	if value != nil && reflect.TypeOf(value).Kind() == reflect.Map {
		return []interface{}{value}
	}
	return value
}

func convertHttpConfigToDict(configDict map[string]interface{}) map[string]interface{} {
	// Convert each list under "http" to a dict (as provider sdk uses lists, but API expects dict structure).
	for nestedKey, nestedValue := range configDict {
		configDict[nestedKey] = listToDict(nestedValue)
		if nestedKey == "auth" {
			authConfig := configDict[nestedKey].(map[string]interface{})
			for authNestedKey, authNestedValue := range authConfig {
				authConfig[authNestedKey] = listToDict(authNestedValue)
			}
			configDict[nestedKey] = authConfig
		}
	}
	return configDict
}
