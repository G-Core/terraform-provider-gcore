// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cdn_logs_uploader_policy

import (
	"context"

	"github.com/G-Core/gcore-go/cdn"
	"github.com/G-Core/gcore-go/packages/param"
	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type CDNLogsUploaderPoliciesResultsListDataSourceEnvelope struct {
	Results customfield.NestedObjectList[CDNLogsUploaderPoliciesItemsDataSourceModel] `json:"results,computed"`
}

type CDNLogsUploaderPoliciesDataSourceModel struct {
	Search    types.String                                                              `tfsdk:"search" query:"search,optional"`
	ConfigIDs *[]types.Int64                                                            `tfsdk:"config_ids" query:"config_ids,optional"`
	MaxItems  types.Int64                                                               `tfsdk:"max_items"`
	Items     customfield.NestedObjectList[CDNLogsUploaderPoliciesItemsDataSourceModel] `tfsdk:"items"`
}

func (m *CDNLogsUploaderPoliciesDataSourceModel) toListParams(_ context.Context) (params cdn.LogsUploaderPolicyListParams, diags diag.Diagnostics) {
	mConfigIDs := []int64{}
	if m.ConfigIDs != nil {
		for _, item := range *m.ConfigIDs {
			mConfigIDs = append(mConfigIDs, item.ValueInt64())
		}
	}

	params = cdn.LogsUploaderPolicyListParams{
		ConfigIDs: mConfigIDs,
	}

	if !m.Search.IsNull() {
		params.Search = param.NewOpt(m.Search.ValueString())
	}

	return
}

type CDNLogsUploaderPoliciesItemsDataSourceModel struct {
	ID                      types.Int64                    `tfsdk:"id" json:"id,computed"`
	ClientID                types.Int64                    `tfsdk:"client_id" json:"client_id,computed"`
	Created                 timetypes.RFC3339              `tfsdk:"created" json:"created,computed" format:"date-time"`
	DateFormat              types.String                   `tfsdk:"date_format" json:"date_format,computed"`
	Description             types.String                   `tfsdk:"description" json:"description,computed"`
	EscapeSpecialCharacters types.Bool                     `tfsdk:"escape_special_characters" json:"escape_special_characters,computed"`
	FieldDelimiter          types.String                   `tfsdk:"field_delimiter" json:"field_delimiter,computed"`
	FieldRemap              customfield.Map[types.String]  `tfsdk:"field_remap" json:"field_remap,computed"`
	FieldSeparator          types.String                   `tfsdk:"field_separator" json:"field_separator,computed"`
	Fields                  customfield.List[types.String] `tfsdk:"fields" json:"fields,computed"`
	FileNameTemplate        types.String                   `tfsdk:"file_name_template" json:"file_name_template,computed"`
	FormatType              types.String                   `tfsdk:"format_type" json:"format_type,computed"`
	IncludeEmptyLogs        types.Bool                     `tfsdk:"include_empty_logs" json:"include_empty_logs,computed"`
	IncludeShieldLogs       types.Bool                     `tfsdk:"include_shield_logs" json:"include_shield_logs,computed"`
	LogSampleRate           types.Float64                  `tfsdk:"log_sample_rate" json:"log_sample_rate,computed"`
	Name                    types.String                   `tfsdk:"name" json:"name,computed"`
	RelatedUploaderConfigs  customfield.List[types.Int64]  `tfsdk:"related_uploader_configs" json:"related_uploader_configs,computed"`
	RetryIntervalMinutes    types.Int64                    `tfsdk:"retry_interval_minutes" json:"retry_interval_minutes,computed"`
	RotateIntervalMinutes   types.Int64                    `tfsdk:"rotate_interval_minutes" json:"rotate_interval_minutes,computed"`
	RotateThresholdLines    types.Int64                    `tfsdk:"rotate_threshold_lines" json:"rotate_threshold_lines,computed"`
	RotateThresholdMB       types.Int64                    `tfsdk:"rotate_threshold_mb" json:"rotate_threshold_mb,computed"`
	Tags                    customfield.Map[types.String]  `tfsdk:"tags" json:"tags,computed"`
	Updated                 timetypes.RFC3339              `tfsdk:"updated" json:"updated,computed" format:"date-time"`
}
