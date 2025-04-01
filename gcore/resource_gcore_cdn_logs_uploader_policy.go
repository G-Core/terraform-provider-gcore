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
				Description: "Include empty logs in the upload. Default value is false.",
				Default:     false,
			},
			"include_shield_logs": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Include logs from origin shielding in the upload. Default value is false.",
				Default:     false,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the policy. Default value is \"Policy\".",
				Default:     "Policy",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the policy. Default value is empty string.",
				Default:     "",
			},
			"retry_interval_minutes": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Interval in minutes to retry failed uploads. Default value is 60.",
				Default:     60,
			},
			"rotate_interval_minutes": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Interval in minutes to rotate logs. Default value is 5.",
				Default:     5,
			},
			"rotate_threshold_mb": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Threshold in MB to rotate logs. Default value is nil.",
				Default:     nil,
			},
			"rotate_threshold_lines": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Threshold in lines to rotate logs. Default value is 0.",
				Default:     0,
			},
			"date_format": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Date format for logs. Default value is empty string.",
				Default:     "",
			},
			"field_delimiter": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Field delimiter for logs. Default value is \\\".",
				Default:     "\"",
			},
			"field_separator": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Field separator for logs. Default value is a space character.",
				Default:     " ",
			},
			"fields": {
				Type:        schema.TypeList,
				Required:    true, // Due to known limitation of default in list type in SDK v2 - https://github.com/hashicorp/terraform-plugin-sdk/issues/142
				Description: "List of fields to include in logs.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"file_name_template": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Template for log file name. Default value is \"{{YYYY}}/{{MM}}/{{DD}}/{{HH}}/{{mm}}/{{ss}}/{{HOST}}_{{CNAME}}_access.log.gz\".",
				Default:     "{{YYYY}}/{{MM}}/{{DD}}/{{HH}}/{{mm}}/{{ss}}/{{HOST}}_{{CNAME}}_access.log.gz",
			},
			"format_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Format type for logs. Default value is empty string.",
				Default:     "",
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: `Tags allow for dynamic decoration of logs by adding predefined fields to the log format. These tags serve as customizable key-value pairs that can be included in log entries to enhance context and readability.`,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Default:     map[string]interface{}{},
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
	req.IncludeEmptyLogs = d.Get("include_empty_logs").(bool)
	req.IncludeShieldLogs = d.Get("include_shield_logs").(bool)
	req.Name = d.Get("name").(string)
	req.Description = d.Get("description").(string)
	req.RetryIntervalMinutes = d.Get("retry_interval_minutes").(int)
	req.RotateIntervalMinutes = d.Get("rotate_interval_minutes").(int)
	req.RotateThresholdLines = d.Get("rotate_threshold_lines").(int)
	req.DateFormat = d.Get("date_format").(string)
	req.FieldDelimiter = d.Get("field_delimiter").(string)
	req.FieldSeparator = d.Get("field_separator").(string)
	req.FileNameTemplate = d.Get("file_name_template").(string)
	req.FormatType = d.Get("format_type").(string)

	if v, ok := d.GetOk("rotate_threshold_mb"); ok {
		req.RotateThresholdMB = pointer.ToInt(v.(int))
	} else {
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
	req.IncludeEmptyLogs = d.Get("include_empty_logs").(bool)
	req.IncludeShieldLogs = d.Get("include_shield_logs").(bool)
	req.Name = d.Get("name").(string)
	req.Description = d.Get("description").(string)
	req.RetryIntervalMinutes = d.Get("retry_interval_minutes").(int)
	req.RotateIntervalMinutes = d.Get("rotate_interval_minutes").(int)
	req.RotateThresholdLines = d.Get("rotate_threshold_lines").(int)
	req.DateFormat = d.Get("date_format").(string)
	req.FieldDelimiter = d.Get("field_delimiter").(string)
	req.FieldSeparator = d.Get("field_separator").(string)
	req.FileNameTemplate = d.Get("file_name_template").(string)
	req.FormatType = d.Get("format_type").(string)

	if v, ok := d.GetOk("rotate_threshold_mb"); ok {
		req.RotateThresholdMB = pointer.ToInt(v.(int))
	} else {
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
	} else {
		req.Tags = make(map[string]string)
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
