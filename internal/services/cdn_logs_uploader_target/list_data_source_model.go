// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cdn_logs_uploader_target

import (
	"context"

	"github.com/G-Core/gcore-go/cdn"
	"github.com/G-Core/gcore-go/packages/param"
	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type CDNLogsUploaderTargetsResultsListDataSourceEnvelope struct {
	Results customfield.NestedObjectList[CDNLogsUploaderTargetsItemsDataSourceModel] `json:"results,computed"`
}

type CDNLogsUploaderTargetsDataSourceModel struct {
	Search    types.String                                                             `tfsdk:"search" query:"search,optional"`
	ConfigIDs *[]types.Int64                                                           `tfsdk:"config_ids" query:"config_ids,optional"`
	MaxItems  types.Int64                                                              `tfsdk:"max_items"`
	Items     customfield.NestedObjectList[CDNLogsUploaderTargetsItemsDataSourceModel] `tfsdk:"items"`
}

func (m *CDNLogsUploaderTargetsDataSourceModel) toListParams(_ context.Context) (params cdn.LogsUploaderTargetListParams, diags diag.Diagnostics) {
	mConfigIDs := []int64{}
	if m.ConfigIDs != nil {
		for _, item := range *m.ConfigIDs {
			mConfigIDs = append(mConfigIDs, item.ValueInt64())
		}
	}

	params = cdn.LogsUploaderTargetListParams{
		ConfigIDs: mConfigIDs,
	}

	if !m.Search.IsNull() {
		params.Search = param.NewOpt(m.Search.ValueString())
	}

	return
}

type CDNLogsUploaderTargetsItemsDataSourceModel struct {
	ID                     types.Int64                                                           `tfsdk:"id" json:"id,computed"`
	ClientID               types.Int64                                                           `tfsdk:"client_id" json:"client_id,computed"`
	Config                 customfield.NestedObject[CDNLogsUploaderTargetsConfigDataSourceModel] `tfsdk:"config" json:"config,computed"`
	Created                timetypes.RFC3339                                                     `tfsdk:"created" json:"created,computed" format:"date-time"`
	Description            types.String                                                          `tfsdk:"description" json:"description,computed"`
	Name                   types.String                                                          `tfsdk:"name" json:"name,computed"`
	RelatedUploaderConfigs customfield.List[types.Int64]                                         `tfsdk:"related_uploader_configs" json:"related_uploader_configs,computed"`
	Status                 jsontypes.Normalized                                                  `tfsdk:"status" json:"status,computed"`
	StorageType            types.String                                                          `tfsdk:"storage_type" json:"storage_type,computed"`
	Updated                timetypes.RFC3339                                                     `tfsdk:"updated" json:"updated,computed" format:"date-time"`
}

type CDNLogsUploaderTargetsConfigDataSourceModel struct {
	AccessKeyID    types.String                                                                `tfsdk:"access_key_id" json:"access_key_id,computed"`
	BucketName     types.String                                                                `tfsdk:"bucket_name" json:"bucket_name,computed"`
	Directory      types.String                                                                `tfsdk:"directory" json:"directory,computed"`
	Endpoint       types.String                                                                `tfsdk:"endpoint" json:"endpoint,computed"`
	Region         types.String                                                                `tfsdk:"region" json:"region,computed"`
	UsePathStyle   types.Bool                                                                  `tfsdk:"use_path_style" json:"use_path_style,computed"`
	Hostname       types.String                                                                `tfsdk:"hostname" json:"hostname,computed"`
	TimeoutSeconds types.Int64                                                                 `tfsdk:"timeout_seconds" json:"timeout_seconds,computed"`
	User           types.String                                                                `tfsdk:"user" json:"user,computed"`
	KeyPassphrase  types.String                                                                `tfsdk:"key_passphrase" json:"key_passphrase,computed"`
	Password       types.String                                                                `tfsdk:"password" json:"password,computed"`
	PrivateKey     types.String                                                                `tfsdk:"private_key" json:"private_key,computed"`
	Append         customfield.NestedObject[CDNLogsUploaderTargetsConfigAppendDataSourceModel] `tfsdk:"append" json:"append,computed"`
	Auth           customfield.NestedObject[CDNLogsUploaderTargetsConfigAuthDataSourceModel]   `tfsdk:"auth" json:"auth,computed"`
	ContentType    types.String                                                                `tfsdk:"content_type" json:"content_type,computed"`
	Retry          customfield.NestedObject[CDNLogsUploaderTargetsConfigRetryDataSourceModel]  `tfsdk:"retry" json:"retry,computed"`
	Upload         customfield.NestedObject[CDNLogsUploaderTargetsConfigUploadDataSourceModel] `tfsdk:"upload" json:"upload,computed"`
	AccountName    types.String                                                                `tfsdk:"account_name" json:"account_name,computed"`
	ContainerName  types.String                                                                `tfsdk:"container_name" json:"container_name,computed"`
	LogStore       types.String                                                                `tfsdk:"log_store" json:"log_store,computed"`
	Project        types.String                                                                `tfsdk:"project" json:"project,computed"`
	Topic          types.String                                                                `tfsdk:"topic" json:"topic,computed"`
}

type CDNLogsUploaderTargetsConfigAppendDataSourceModel struct {
	Headers         customfield.Map[types.String]                                                                  `tfsdk:"headers" json:"headers,computed"`
	Method          types.String                                                                                   `tfsdk:"method" json:"method,computed"`
	ResponseActions customfield.NestedObjectList[CDNLogsUploaderTargetsConfigAppendResponseActionsDataSourceModel] `tfsdk:"response_actions" json:"response_actions,computed"`
	TimeoutSeconds  types.Int64                                                                                    `tfsdk:"timeout_seconds" json:"timeout_seconds,computed"`
	URL             types.String                                                                                   `tfsdk:"url" json:"url,computed"`
	UseCompression  types.Bool                                                                                     `tfsdk:"use_compression" json:"use_compression,computed"`
}

type CDNLogsUploaderTargetsConfigAppendResponseActionsDataSourceModel struct {
	Action          types.String `tfsdk:"action" json:"action,computed"`
	Description     types.String `tfsdk:"description" json:"description,computed"`
	MatchPayload    types.String `tfsdk:"match_payload" json:"match_payload,computed"`
	MatchStatusCode types.Int64  `tfsdk:"match_status_code" json:"match_status_code,computed"`
}

type CDNLogsUploaderTargetsConfigAuthDataSourceModel struct {
	Config customfield.NestedObject[CDNLogsUploaderTargetsConfigAuthConfigDataSourceModel] `tfsdk:"config" json:"config,computed"`
	Type   types.String                                                                    `tfsdk:"type" json:"type,computed"`
}

type CDNLogsUploaderTargetsConfigAuthConfigDataSourceModel struct {
	Token           types.String `tfsdk:"token" json:"token,computed"`
	HeaderName      types.String `tfsdk:"header_name" json:"header_name,computed"`
	AccountKey      types.String `tfsdk:"account_key" json:"account_key,computed"`
	AccessKeyID     types.String `tfsdk:"access_key_id" json:"access_key_id,computed"`
	SecretAccessKey types.String `tfsdk:"secret_access_key" json:"secret_access_key,computed"`
}

type CDNLogsUploaderTargetsConfigRetryDataSourceModel struct {
	Headers         customfield.Map[types.String]                                                                 `tfsdk:"headers" json:"headers,computed"`
	Method          types.String                                                                                  `tfsdk:"method" json:"method,computed"`
	ResponseActions customfield.NestedObjectList[CDNLogsUploaderTargetsConfigRetryResponseActionsDataSourceModel] `tfsdk:"response_actions" json:"response_actions,computed"`
	TimeoutSeconds  types.Int64                                                                                   `tfsdk:"timeout_seconds" json:"timeout_seconds,computed"`
	URL             types.String                                                                                  `tfsdk:"url" json:"url,computed"`
	UseCompression  types.Bool                                                                                    `tfsdk:"use_compression" json:"use_compression,computed"`
}

type CDNLogsUploaderTargetsConfigRetryResponseActionsDataSourceModel struct {
	Action          types.String `tfsdk:"action" json:"action,computed"`
	Description     types.String `tfsdk:"description" json:"description,computed"`
	MatchPayload    types.String `tfsdk:"match_payload" json:"match_payload,computed"`
	MatchStatusCode types.Int64  `tfsdk:"match_status_code" json:"match_status_code,computed"`
}

type CDNLogsUploaderTargetsConfigUploadDataSourceModel struct {
	Headers         customfield.Map[types.String]                                                                  `tfsdk:"headers" json:"headers,computed"`
	Method          types.String                                                                                   `tfsdk:"method" json:"method,computed"`
	ResponseActions customfield.NestedObjectList[CDNLogsUploaderTargetsConfigUploadResponseActionsDataSourceModel] `tfsdk:"response_actions" json:"response_actions,computed"`
	TimeoutSeconds  types.Int64                                                                                    `tfsdk:"timeout_seconds" json:"timeout_seconds,computed"`
	URL             types.String                                                                                   `tfsdk:"url" json:"url,computed"`
	UseCompression  types.Bool                                                                                     `tfsdk:"use_compression" json:"use_compression,computed"`
}

type CDNLogsUploaderTargetsConfigUploadResponseActionsDataSourceModel struct {
	Action          types.String `tfsdk:"action" json:"action,computed"`
	Description     types.String `tfsdk:"description" json:"description,computed"`
	MatchPayload    types.String `tfsdk:"match_payload" json:"match_payload,computed"`
	MatchStatusCode types.Int64  `tfsdk:"match_status_code" json:"match_status_code,computed"`
}
