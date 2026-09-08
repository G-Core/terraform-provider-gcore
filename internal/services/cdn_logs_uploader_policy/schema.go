// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cdn_logs_uploader_policy

import (
	"context"

	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/float64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/float64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.ResourceWithConfigValidators = (*CDNLogsUploaderPolicyResource)(nil)

func ResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Logs uploader policies define how CDN logs are formatted and delivered, including field selection, field ordering, delimiters, delivery frequency, and file size limits.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:      true,
				PlanModifiers: []planmodifier.Int64{int64planmodifier.UseNonNullStateForUnknown()},
			},
			"date_format": schema.StringAttribute{
				Description: "Date format for logs.",
				Computed:    true,
				Optional:    true,
				Default:     stringdefault.StaticString(""),
			},
			"description": schema.StringAttribute{
				Description: "Description of the policy.",
				Computed:    true,
				Optional:    true,
				Default:     stringdefault.StaticString(""),
			},
			"format_type": schema.StringAttribute{
				Description: "Format type for logs.\n\nPossible values:\n- **\"\"** - empty, it means it will apply the format configurations from the policy.\n- **\"json\"** - output the logs as json lines.\nAvailable values: \"json\", \"\".",
				Computed:    true,
				Optional:    true,
				Default:     stringdefault.StaticString(""),
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive("json", ""),
				},
			},
			"rotate_threshold_mb": schema.Int64Attribute{
				Description: "Threshold in MB to rotate logs.",
				Optional:    true,
				Validators: []validator.Int64{
					int64validator.Between(100, 2147483647),
				},
			},
			"field_remap": schema.MapAttribute{
				Description: "Per-field output-name remap for exported logs. Maps a canonical Gcore field name (from `/cdn/logs_uploader/policies/fields`, and must be present in `fields`) to the field name it should have in the exported logs. Unmapped fields keep their canonical name. Output names (after remapping) must be unique.",
				Optional:    true,
				ElementType: types.StringType,
			},
			"tags": schema.MapAttribute{
				Description: "Tags allow for dynamic decoration of logs by adding predefined fields to the log format. These tags serve as customizable key-value pairs that can be included in log entries to enhance context and readability.",
				Computed:    true,
				Optional:    true,
				ElementType: types.StringType,
				Default:     mapdefault.StaticValue(types.MapValueMust(types.StringType, map[string]attr.Value{})),
			},
			"field_conversions": schema.MapNestedAttribute{
				Description: "Per-field value conversions for exported logs. Maps a canonical Gcore field name to the pipeline applied to its values. Field names are limited to 255 characters and must not be empty. Each key must be present in `fields`, and each conversion type must be listed in that field's `allowed_conversions` from `/cdn/v2/logs_uploader/policies/fields`. Conversions in a pipeline are applied in array order. Values are converted independently of `field_remap`, which renames the exported field: both are keyed on the canonical field name.",
				Optional:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"conversions": schema.ListNestedAttribute{
							Required: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"config": schema.SingleNestedAttribute{
										Required: true,
										Attributes: map[string]schema.Attribute{
											"factor": schema.Float64Attribute{
												Description: "Multiplier applied to the field value.",
												Optional:    true,
												Validators: []validator.Float64{
													float64validator.Between(1e-15, 1000000000000000),
												},
											},
											"precision": schema.Int64Attribute{
												Description: "Optional number of decimal places in the converted value. Must be specified together with `rounding`; returned as `null` when unspecified.",
												Optional:    true,
												Validators: []validator.Int64{
													int64validator.Between(0, 9),
												},
											},
											"rounding": schema.StringAttribute{
												Description: "Optional rounding mode. Must be specified together with `precision`; returned as `null` when unspecified.\nAvailable values: \"nearest\", \"down\", \"up\".",
												Optional:    true,
												Validators: []validator.String{
													stringvalidator.OneOfCaseInsensitive(
														"nearest",
														"down",
														"up",
													),
												},
											},
											"values": schema.MapAttribute{
												Description: "Exact, case-sensitive replacements, matched and written without trimming. Keys must not be empty and are limited to 255 characters; values are limited to 100 characters and may be empty.",
												Optional:    true,
												ElementType: types.StringType,
											},
											"default": schema.StringAttribute{
												Description: "Value used when the input does not match any key in `values`.",
												Optional:    true,
											},
										},
									},
									"type": schema.StringAttribute{
										Description: `Available values: "scale", "replace".`,
										Required:    true,
										Validators: []validator.String{
											stringvalidator.OneOfCaseInsensitive("scale", "replace"),
										},
									},
								},
							},
						},
					},
				},
			},
			"escape_special_characters": schema.BoolAttribute{
				Description: "When set to true, the service sanitizes string values by escaping characters that may be unsafe for transport, logging, or downstream processing.\n\nThe following categories of characters are escaped:\n- Control and non-printable characters\n- Quotation marks and escape characters\n- Characters outside the standard ASCII range\n\nThe resulting output contains only printable ASCII characters.",
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"field_delimiter": schema.StringAttribute{
				Description: "Field delimiter for logs.",
				Computed:    true,
				Optional:    true,
				Default:     stringdefault.StaticString("\""),
			},
			"field_separator": schema.StringAttribute{
				Description: "Field separator for logs.",
				Computed:    true,
				Optional:    true,
				Default:     stringdefault.StaticString(" "),
			},
			"file_name_template": schema.StringAttribute{
				Description: "Template for log file name.",
				Computed:    true,
				Optional:    true,
				Default:     stringdefault.StaticString("{{YYYY}}/{{MM}}/{{DD}}/{{HH}}/{{mm}}/{{ss}}/{{HOST}}_{{CNAME}}_access.log.gz"),
			},
			"include_empty_logs": schema.BoolAttribute{
				Description: "Include empty logs in the upload.",
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"include_shield_logs": schema.BoolAttribute{
				Description: "Include logs from origin shielding in the upload.",
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"log_sample_rate": schema.Float64Attribute{
				Description: "Sampling rate for logs. A value between 0 and 1 that determines the fraction of log entries to collect.\n\n- **1** - collect all logs (default).\n- **0.5** - collect approximately 50% of logs.\n- **0** - collect no logs (effectively disables logging without removing the policy).",
				Computed:    true,
				Optional:    true,
				Validators: []validator.Float64{
					float64validator.Between(0, 1),
				},
				Default: float64default.StaticFloat64(1),
			},
			"name": schema.StringAttribute{
				Description: "Name of the policy.",
				Computed:    true,
				Optional:    true,
				Default:     stringdefault.StaticString("Policy"),
			},
			"retry_interval_minutes": schema.Int64Attribute{
				Description: "Interval in minutes to retry failed uploads.",
				Computed:    true,
				Optional:    true,
				Validators: []validator.Int64{
					int64validator.Between(5, 60),
				},
				Default: int64default.StaticInt64(60),
			},
			"rotate_interval_minutes": schema.Int64Attribute{
				Description: "Interval in minutes to rotate logs.",
				Computed:    true,
				Optional:    true,
				Validators: []validator.Int64{
					int64validator.Between(5, 60),
				},
				Default: int64default.StaticInt64(5),
			},
			"rotate_threshold_lines": schema.Int64Attribute{
				Description: "Threshold in lines to rotate logs.",
				Computed:    true,
				Optional:    true,
				Validators: []validator.Int64{
					int64validator.Between(0, 10000),
				},
				Default: int64default.StaticInt64(0),
			},
			"fields": schema.ListAttribute{
				Description:   "List of fields to include in logs. Duplicate names are allowed for plain text output, but rejected when `format_type` is `json` or a `field_remap` is set (each field becomes a distinct output key).",
				Computed:      true,
				Optional:      true,
				CustomType:    customfield.NewListType[types.String](ctx),
				ElementType:   types.StringType,
				PlanModifiers: []planmodifier.List{listplanmodifier.UseStateForUnknown()},
			},
			"client_id": schema.Int64Attribute{
				Description:   "Client that owns the policy.",
				Computed:      true,
				PlanModifiers: []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
			},
			"created": schema.StringAttribute{
				Description:   "Time when logs uploader policy was created.",
				Computed:      true,
				CustomType:    timetypes.RFC3339Type{},
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"updated": schema.StringAttribute{
				Description: "Time when logs uploader policy was updated.",
				Computed:    true,
				CustomType:  timetypes.RFC3339Type{},
			},
			"related_uploader_configs": schema.ListAttribute{
				Description:   "List of logs uploader configs that use this policy.",
				Computed:      true,
				CustomType:    customfield.NewListType[types.Int64](ctx),
				ElementType:   types.Int64Type,
				PlanModifiers: []planmodifier.List{listplanmodifier.UseStateForUnknown()},
			},
		},
	}
}

func (r *CDNLogsUploaderPolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ResourceSchema(ctx)
}

func (r *CDNLogsUploaderPolicyResource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{}
}
