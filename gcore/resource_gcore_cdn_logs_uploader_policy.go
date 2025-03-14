package gcore

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/AlekSi/pointer"
	"github.com/G-Core/gcorelabscdn-go/logsuploader"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCDNLogsUploaderPolicy() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"include_empty_logs": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Include empty logs in the upload.",
			},
			"include_shield_logs": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Include logs from origin shielding in the upload.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Name of the policy.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the policy.",
			},
			"retry_interval_minutes": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Interval in minutes to retry failed uploads.",
			},
			"rotate_interval_minutes": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Interval in minutes to rotate logs.",
			},
			"rotate_threshold_mb": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Threshold in MB to rotate logs.",
			},
			"rotate_threshold_lines": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Threshold in lines to rotate logs.",
			},
			"date_format": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Date format for logs.",
			},
			"field_delimiter": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Field delimiter for logs.",
			},
			"field_separator": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Field separator for logs.",
			},
			"fields": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "List of fields to include in logs.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"file_name_template": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Template for log file name.",
			},
			"format_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Format type for logs.",
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: `Tags allow for dynamic decoration of logs by adding predefined fields to the log format. These tags serve as customizable key-value pairs that can be included in log entries to enhance context and readability.`,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
		CreateContext: resourceCDNLogsUploaderPolicyCreate,
		ReadContext:   resourceCDNLogsUploaderPolicyRead,
		UpdateContext: resourceCDNLogsUploaderPolicyUpdate,
		DeleteContext: resourceCDNLogsUploaderPolicyDelete,
		Description:   "Represent CDN logs uploader policy",
	}
}

func resourceCDNLogsUploaderPolicyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start CDN Logs Uploader Policy creating")
	config := m.(*Config)
	client := config.CDNClient

	var req logsuploader.PolicyCreateRequest

	if v, ok := d.GetOk("include_empty_logs"); ok {
		req.IncludeEmptyLogs = v.(bool)
	}

	if v, ok := d.GetOk("include_shield_logs"); ok {
		req.IncludeShieldLogs = v.(bool)
	}

	if v, ok := d.GetOk("name"); ok {
		req.Name = v.(string)
	}

	if v, ok := d.GetOk("description"); ok {
		req.Description = v.(string)
	}

	if v, ok := d.GetOk("retry_interval_minutes"); ok {
		req.RetryIntervalMinutes = v.(int)
	}

	if v, ok := d.GetOk("rotate_interval_minutes"); ok {
		req.RotateIntervalMinutes = v.(int)
	}

	if v, ok := d.GetOk("rotate_threshold_lines"); ok {
		req.RotateThresholdLines = v.(int)
	}

	if v, ok := d.GetOk("date_format"); ok {
		req.DateFormat = v.(string)
	}

	if v, ok := d.GetOk("field_delimiter"); ok {
		req.FieldDelimiter = v.(string)
	}

	if v, ok := d.GetOk("field_separator"); ok {
		req.FieldSeparator = v.(string)
	}

	if v, ok := d.GetOk("file_name_template"); ok {
		req.FileNameTemplate = v.(string)
	}

	if v, ok := d.GetOk("format_type"); ok {
		req.FormatType = v.(string)
	}

	if v, ok := d.GetOk("rotate_threshold_mb"); ok {
		req.RotateThresholdMB = pointer.ToInt(v.(int))
	} else if d.Get("rotate_threshold_mb") == nil {
		req.RotateThresholdMB = nil
	}

	if v, ok := d.GetOk("fields"); ok {
		fields := make([]string, len(v.([]interface{})))
		for i, v := range v.([]interface{}) {
			fields[i] = v.(string)
		}
		req.Fields = fields
	}
	if v, ok := d.GetOk("tags"); ok {
		tagsInterface := v.(map[string]interface{})
		tags := make(map[string]string, len(tagsInterface))
		for k, value := range tagsInterface {
			tags[k] = value.(string)
		}
		req.Tags = tags
	}

	result, err := client.LogsUploader().PolicyCreate(ctx, &req)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d", result.ID))
	resourceCDNLogsUploaderPolicyRead(ctx, d, m)

	log.Printf("[DEBUG] Finish CDN Logs Uploader Policy creating (id=%d)\n", result.ID)
	return nil
}

func resourceCDNLogsUploaderPolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	policyID := d.Id()
	log.Printf("[DEBUG] Start CDN Logs Uploader Policy reading (id=%s)\n", policyID)
	config := m.(*Config)
	client := config.CDNClient

	id, err := strconv.ParseInt(policyID, 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}

	result, err := client.LogsUploader().PolicyGet(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("include_empty_logs", result.IncludeEmptyLogs)
	d.Set("include_shield_logs", result.IncludeShieldLogs)
	d.Set("name", result.Name)
	d.Set("description", result.Description)
	d.Set("retry_interval_minutes", result.RetryIntervalMinutes)
	d.Set("rotate_interval_minutes", result.RotateIntervalMinutes)
	d.Set("rotate_threshold_mb", result.RotateThresholdMB)
	d.Set("rotate_threshold_lines", result.RotateThresholdLines)
	d.Set("date_format", result.DateFormat)
	d.Set("field_delimiter", result.FieldDelimiter)
	d.Set("field_separator", result.FieldSeparator)
	d.Set("fields", result.Fields)
	d.Set("file_name_template", result.FileNameTemplate)
	d.Set("format_type", result.FormatType)
	d.Set("tags", result.Tags)

	log.Println("[DEBUG] Finish CDN Logs Uploader Policy reading")
	return nil
}

func resourceCDNLogsUploaderPolicyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	policyID := d.Id()
	log.Printf("[DEBUG] Start CDN Logs Uploader Policy updating (id=%s)\n", policyID)
	config := m.(*Config)
	client := config.CDNClient

	id, err := strconv.ParseInt(policyID, 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}

	var req logsuploader.PolicyUpdateRequest
	if v, ok := d.GetOk("include_empty_logs"); ok {
		req.IncludeEmptyLogs = v.(bool)
	}

	if v, ok := d.GetOk("include_shield_logs"); ok {
		req.IncludeShieldLogs = v.(bool)
	}

	if v, ok := d.GetOk("name"); ok {
		req.Name = v.(string)
	}

	if v, ok := d.GetOk("description"); ok {
		req.Description = v.(string)
	}

	if v, ok := d.GetOk("retry_interval_minutes"); ok {
		req.RetryIntervalMinutes = v.(int)
	}

	if v, ok := d.GetOk("rotate_interval_minutes"); ok {
		req.RotateIntervalMinutes = v.(int)
	}

	if v, ok := d.GetOk("rotate_threshold_lines"); ok {
		req.RotateThresholdLines = v.(int)
	}

	if v, ok := d.GetOk("date_format"); ok {
		req.DateFormat = v.(string)
	}

	if v, ok := d.GetOk("field_delimiter"); ok {
		req.FieldDelimiter = v.(string)
	}

	if v, ok := d.GetOk("field_separator"); ok {
		req.FieldSeparator = v.(string)
	}

	if v, ok := d.GetOk("file_name_template"); ok {
		req.FileNameTemplate = v.(string)
	}

	if v, ok := d.GetOk("format_type"); ok {
		req.FormatType = v.(string)
	}

	if v, ok := d.GetOk("rotate_threshold_mb"); ok {
		req.RotateThresholdMB = pointer.ToInt(v.(int))
	} else if d.Get("rotate_threshold_mb") == nil {
		req.RotateThresholdMB = nil
	}

	if v, ok := d.GetOk("fields"); ok {
		fields := make([]string, len(v.([]interface{})))
		for i, v := range v.([]interface{}) {
			fields[i] = v.(string)
		}
		req.Fields = fields
	}
	if v, ok := d.GetOk("tags"); ok {
		tagsInterface := v.(map[string]interface{})
		tags := make(map[string]string, len(tagsInterface))
		for k, value := range tagsInterface {
			tags[k] = value.(string)
		}
		req.Tags = tags
	}

	if _, err := client.LogsUploader().PolicyUpdate(ctx, id, &req); err != nil {
		return diag.FromErr(err)
	}

	log.Println("[DEBUG] Finish CDN Logs Uploader Policy updating")
	return resourceCDNLogsUploaderPolicyRead(ctx, d, m)
}

func resourceCDNLogsUploaderPolicyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	policyID := d.Id()
	log.Printf("[DEBUG] Start CDN Logs Uploader Policy deleting (id=%s)\n", policyID)
	config := m.(*Config)
	client := config.CDNClient

	id, err := strconv.ParseInt(policyID, 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := client.LogsUploader().PolicyDelete(ctx, id); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	log.Println("[DEBUG] Finish CDN Logs Uploader Policy deleting")
	return nil
}
