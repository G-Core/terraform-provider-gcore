// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cdn_logs_uploader_policy

import (
	"context"

	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/float64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSourceWithConfigValidators = (*CDNLogsUploaderPolicyDataSource)(nil)

func DataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Logs uploader policies define how CDN logs are formatted and delivered, including field selection, field ordering, delimiters, delivery frequency, and file size limits.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Required: true,
			},
			"client_id": schema.Int64Attribute{
				Description: "Client that owns the policy.",
				Computed:    true,
			},
			"created": schema.StringAttribute{
				Description: "Time when logs uploader policy was created.",
				Computed:    true,
				CustomType:  timetypes.RFC3339Type{},
			},
			"date_format": schema.StringAttribute{
				Description: "Date format for logs.",
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "Description of the policy.",
				Computed:    true,
			},
			"escape_special_characters": schema.BoolAttribute{
				Description: "When set to true, the service sanitizes string values by escaping characters that may be unsafe for transport, logging, or downstream processing.\n\nThe following categories of characters are escaped:\n- Control and non-printable characters\n- Quotation marks and escape characters\n- Characters outside the standard ASCII range\n\nThe resulting output contains only printable ASCII characters.",
				Computed:    true,
			},
			"field_delimiter": schema.StringAttribute{
				Description: "Field delimiter for logs.",
				Computed:    true,
			},
			"field_separator": schema.StringAttribute{
				Description: "Field separator for logs.",
				Computed:    true,
			},
			"file_name_template": schema.StringAttribute{
				Description: "Template for log file name.",
				Computed:    true,
			},
			"format_type": schema.StringAttribute{
				Description: "Format type for logs.\n\nPossible values:\n- **\"\"** - empty, it means it will apply the format configurations from the policy.\n- **\"json\"** - output the logs as json lines.\nAvailable values: \"json\", \"\".",
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive("json", ""),
				},
			},
			"include_empty_logs": schema.BoolAttribute{
				Description: "Include empty logs in the upload.",
				Computed:    true,
			},
			"include_shield_logs": schema.BoolAttribute{
				Description: "Include logs from origin shielding in the upload.",
				Computed:    true,
			},
			"log_sample_rate": schema.Float64Attribute{
				Description: "Sampling rate for logs. A value between 0 and 1 that determines the fraction of log entries to collect.\n\n- **1** - collect all logs (default).\n- **0.5** - collect approximately 50% of logs.\n- **0** - collect no logs (effectively disables logging without removing the policy).",
				Computed:    true,
				Validators: []validator.Float64{
					float64validator.Between(0, 1),
				},
			},
			"name": schema.StringAttribute{
				Description: "Name of the policy.",
				Computed:    true,
			},
			"retry_interval_minutes": schema.Int64Attribute{
				Description: "Interval in minutes to retry failed uploads.",
				Computed:    true,
				Validators: []validator.Int64{
					int64validator.Between(5, 60),
				},
			},
			"rotate_interval_minutes": schema.Int64Attribute{
				Description: "Interval in minutes to rotate logs.",
				Computed:    true,
				Validators: []validator.Int64{
					int64validator.Between(5, 60),
				},
			},
			"rotate_threshold_lines": schema.Int64Attribute{
				Description: "Threshold in lines to rotate logs.",
				Computed:    true,
				Validators: []validator.Int64{
					int64validator.Between(0, 10000),
				},
			},
			"rotate_threshold_mb": schema.Int64Attribute{
				Description: "Threshold in MB to rotate logs.",
				Computed:    true,
				Validators: []validator.Int64{
					int64validator.Between(5, 500),
				},
			},
			"updated": schema.StringAttribute{
				Description: "Time when logs uploader policy was updated.",
				Computed:    true,
				CustomType:  timetypes.RFC3339Type{},
			},
			"fields": schema.ListAttribute{
				Description: "List of fields to include in logs.",
				Computed:    true,
				CustomType:  customfield.NewListType[types.String](ctx),
				ElementType: types.StringType,
			},
			"related_uploader_configs": schema.ListAttribute{
				Description: "List of logs uploader configs that use this policy.",
				Computed:    true,
				CustomType:  customfield.NewListType[types.Int64](ctx),
				ElementType: types.Int64Type,
			},
			"tags": schema.MapAttribute{
				Description: "Tags allow for dynamic decoration of logs by adding predefined fields to the log format. These tags serve as customizable key-value pairs that can be included in log entries to enhance context and readability.",
				Computed:    true,
				CustomType:  customfield.NewMapType[types.String](ctx),
				ElementType: types.StringType,
			},
		},
	}
}

func (d *CDNLogsUploaderPolicyDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = DataSourceSchema(ctx)
}

func (d *CDNLogsUploaderPolicyDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{}
}
