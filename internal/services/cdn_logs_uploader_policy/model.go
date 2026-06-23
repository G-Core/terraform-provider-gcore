// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cdn_logs_uploader_policy

import (
	"encoding/json"
	"reflect"

	"github.com/G-Core/terraform-provider-gcore/internal/apijson"
	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/tidwall/sjson"
)

type CDNLogsUploaderPolicyModel struct {
	ID                      types.Int64                    `tfsdk:"id" json:"id,computed"`
	DateFormat              types.String                   `tfsdk:"date_format" json:"date_format,computed_optional"`
	Description             types.String                   `tfsdk:"description" json:"description,computed_optional"`
	FormatType              types.String                   `tfsdk:"format_type" json:"format_type,computed_optional"`
	RotateThresholdMB       types.Int64                    `tfsdk:"rotate_threshold_mb" json:"rotate_threshold_mb,optional"`
	Tags                    *map[string]types.String       `tfsdk:"tags" json:"tags,computed_optional"`
	EscapeSpecialCharacters types.Bool                     `tfsdk:"escape_special_characters" json:"escape_special_characters,computed_optional"`
	FieldDelimiter          types.String                   `tfsdk:"field_delimiter" json:"field_delimiter,computed_optional"`
	FieldSeparator          types.String                   `tfsdk:"field_separator" json:"field_separator,computed_optional"`
	FileNameTemplate        types.String                   `tfsdk:"file_name_template" json:"file_name_template,computed_optional"`
	IncludeEmptyLogs        types.Bool                     `tfsdk:"include_empty_logs" json:"include_empty_logs,computed_optional"`
	IncludeShieldLogs       types.Bool                     `tfsdk:"include_shield_logs" json:"include_shield_logs,computed_optional"`
	LogSampleRate           types.Float64                  `tfsdk:"log_sample_rate" json:"log_sample_rate,computed_optional"`
	Name                    types.String                   `tfsdk:"name" json:"name,computed_optional"`
	RetryIntervalMinutes    types.Int64                    `tfsdk:"retry_interval_minutes" json:"retry_interval_minutes,computed_optional"`
	RotateIntervalMinutes   types.Int64                    `tfsdk:"rotate_interval_minutes" json:"rotate_interval_minutes,computed_optional"`
	RotateThresholdLines    types.Int64                    `tfsdk:"rotate_threshold_lines" json:"rotate_threshold_lines,computed_optional"`
	Fields                  customfield.List[types.String] `tfsdk:"fields" json:"fields,computed_optional"`
	ClientID                types.Int64                    `tfsdk:"client_id" json:"client_id,computed"`
	Created                 timetypes.RFC3339              `tfsdk:"created" json:"created,computed" format:"date-time"`
	Updated                 timetypes.RFC3339              `tfsdk:"updated" json:"updated,computed" format:"date-time"`
	RelatedUploaderConfigs  customfield.List[types.Int64]  `tfsdk:"related_uploader_configs" json:"related_uploader_configs,computed"`
}

func (m CDNLogsUploaderPolicyModel) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(m)
}

func (m CDNLogsUploaderPolicyModel) MarshalJSONForUpdate(state CDNLogsUploaderPolicyModel) (data []byte, err error) {
	planTags := m.Tags
	// Prevent MarshalForPatch from patching tags key-by-key (which sends null
	// for removed keys). We'll inject the full tags value below instead.
	m.Tags = state.Tags
	result, err := apijson.MarshalForPatch(m, state)
	if err != nil {
		return nil, err
	}
	if !reflect.DeepEqual(planTags, state.Tags) {
		tags := map[string]string{}
		if planTags != nil {
			for k, v := range *planTags {
				tags[k] = v.ValueString()
			}
		}
		tagsJSON, _ := json.Marshal(tags)
		result, _ = sjson.SetRawBytes(result, "tags", tagsJSON)
	}
	return result, nil
}
