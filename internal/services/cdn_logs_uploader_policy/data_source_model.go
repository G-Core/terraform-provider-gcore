// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cdn_logs_uploader_policy

import (
	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type CDNLogsUploaderPolicyDataSourceModel struct {
	ID                      types.Int64                    `tfsdk:"id" path:"id,required"`
	ClientID                types.Int64                    `tfsdk:"client_id" json:"client_id,computed"`
	Created                 timetypes.RFC3339              `tfsdk:"created" json:"created,computed" format:"date-time"`
	DateFormat              types.String                   `tfsdk:"date_format" json:"date_format,computed"`
	Description             types.String                   `tfsdk:"description" json:"description,computed"`
	EscapeSpecialCharacters types.Bool                     `tfsdk:"escape_special_characters" json:"escape_special_characters,computed"`
	FieldDelimiter          types.String                   `tfsdk:"field_delimiter" json:"field_delimiter,computed"`
	FieldSeparator          types.String                   `tfsdk:"field_separator" json:"field_separator,computed"`
	FileNameTemplate        types.String                   `tfsdk:"file_name_template" json:"file_name_template,computed"`
	FormatType              types.String                   `tfsdk:"format_type" json:"format_type,computed"`
	IncludeEmptyLogs        types.Bool                     `tfsdk:"include_empty_logs" json:"include_empty_logs,computed"`
	IncludeShieldLogs       types.Bool                     `tfsdk:"include_shield_logs" json:"include_shield_logs,computed"`
	LogSampleRate           types.Float64                  `tfsdk:"log_sample_rate" json:"log_sample_rate,computed"`
	Name                    types.String                   `tfsdk:"name" json:"name,computed"`
	RetryIntervalMinutes    types.Int64                    `tfsdk:"retry_interval_minutes" json:"retry_interval_minutes,computed"`
	RotateIntervalMinutes   types.Int64                    `tfsdk:"rotate_interval_minutes" json:"rotate_interval_minutes,computed"`
	RotateThresholdLines    types.Int64                    `tfsdk:"rotate_threshold_lines" json:"rotate_threshold_lines,computed"`
	RotateThresholdMB       types.Int64                    `tfsdk:"rotate_threshold_mb" json:"rotate_threshold_mb,computed"`
	Updated                 timetypes.RFC3339              `tfsdk:"updated" json:"updated,computed" format:"date-time"`
	Fields                  customfield.List[types.String] `tfsdk:"fields" json:"fields,computed"`
	RelatedUploaderConfigs  customfield.List[types.Int64]  `tfsdk:"related_uploader_configs" json:"related_uploader_configs,computed"`
	Tags                    customfield.Map[types.String]  `tfsdk:"tags" json:"tags,computed"`
}
