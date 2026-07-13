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

type CDNLogsUploaderTargetDataSourceModel struct {
	ID                     types.Int64                                                          `tfsdk:"id" path:"id,computed_optional"`
	ClientID               types.Int64                                                          `tfsdk:"client_id" json:"client_id,computed"`
	Created                timetypes.RFC3339                                                    `tfsdk:"created" json:"created,computed" format:"date-time"`
	Description            types.String                                                         `tfsdk:"description" json:"description,computed"`
	Name                   types.String                                                         `tfsdk:"name" json:"name,computed"`
	StorageType            types.String                                                         `tfsdk:"storage_type" json:"storage_type,computed"`
	Updated                timetypes.RFC3339                                                    `tfsdk:"updated" json:"updated,computed" format:"date-time"`
	RelatedUploaderConfigs customfield.List[types.Int64]                                        `tfsdk:"related_uploader_configs" json:"related_uploader_configs,computed"`
	Config                 customfield.NestedObject[CDNLogsUploaderTargetConfigDataSourceModel] `tfsdk:"config" json:"config,computed"`
	Status                 jsontypes.Normalized                                                 `tfsdk:"status" json:"status,computed"`
	FindOneBy              *CDNLogsUploaderTargetFindOneByDataSourceModel                       `tfsdk:"find_one_by"`
}

func (m *CDNLogsUploaderTargetDataSourceModel) toListParams(_ context.Context) (params cdn.LogsUploaderTargetListParams, diags diag.Diagnostics) {
	mFindOneByConfigIDs := []int64{}
	if m.FindOneBy.ConfigIDs != nil {
		for _, item := range *m.FindOneBy.ConfigIDs {
			mFindOneByConfigIDs = append(mFindOneByConfigIDs, item.ValueInt64())
		}
	}

	params = cdn.LogsUploaderTargetListParams{
		ConfigIDs: mFindOneByConfigIDs,
	}

	if !m.FindOneBy.Search.IsNull() {
		params.Search = param.NewOpt(m.FindOneBy.Search.ValueString())
	}

	return
}

type CDNLogsUploaderTargetConfigDataSourceModel struct {
	AccessKeyID    types.String                                                               `tfsdk:"access_key_id" json:"access_key_id,computed"`
	BucketName     types.String                                                               `tfsdk:"bucket_name" json:"bucket_name,computed"`
	Directory      types.String                                                               `tfsdk:"directory" json:"directory,computed"`
	Endpoint       types.String                                                               `tfsdk:"endpoint" json:"endpoint,computed"`
	Region         types.String                                                               `tfsdk:"region" json:"region,computed"`
	UsePathStyle   types.Bool                                                                 `tfsdk:"use_path_style" json:"use_path_style,computed"`
	Hostname       types.String                                                               `tfsdk:"hostname" json:"hostname,computed"`
	TimeoutSeconds types.Int64                                                                `tfsdk:"timeout_seconds" json:"timeout_seconds,computed"`
	User           types.String                                                               `tfsdk:"user" json:"user,computed"`
	KeyPassphrase  types.String                                                               `tfsdk:"key_passphrase" json:"key_passphrase,computed"`
	Password       types.String                                                               `tfsdk:"password" json:"password,computed"`
	PrivateKey     types.String                                                               `tfsdk:"private_key" json:"private_key,computed"`
	Append         customfield.NestedObject[CDNLogsUploaderTargetConfigAppendDataSourceModel] `tfsdk:"append" json:"append,computed"`
	Auth           customfield.NestedObject[CDNLogsUploaderTargetConfigAuthDataSourceModel]   `tfsdk:"auth" json:"auth,computed"`
	ContentType    types.String                                                               `tfsdk:"content_type" json:"content_type,computed"`
	Retry          customfield.NestedObject[CDNLogsUploaderTargetConfigRetryDataSourceModel]  `tfsdk:"retry" json:"retry,computed"`
	Upload         customfield.NestedObject[CDNLogsUploaderTargetConfigUploadDataSourceModel] `tfsdk:"upload" json:"upload,computed"`
	AccountName    types.String                                                               `tfsdk:"account_name" json:"account_name,computed"`
	ContainerName  types.String                                                               `tfsdk:"container_name" json:"container_name,computed"`
	LogStore       types.String                                                               `tfsdk:"log_store" json:"log_store,computed"`
	Project        types.String                                                               `tfsdk:"project" json:"project,computed"`
	Topic          types.String                                                               `tfsdk:"topic" json:"topic,computed"`
}

type CDNLogsUploaderTargetConfigAppendDataSourceModel struct {
	Headers         customfield.Map[types.String]                                                                 `tfsdk:"headers" json:"headers,computed"`
	Method          types.String                                                                                  `tfsdk:"method" json:"method,computed"`
	ResponseActions customfield.NestedObjectList[CDNLogsUploaderTargetConfigAppendResponseActionsDataSourceModel] `tfsdk:"response_actions" json:"response_actions,computed"`
	TimeoutSeconds  types.Int64                                                                                   `tfsdk:"timeout_seconds" json:"timeout_seconds,computed"`
	URL             types.String                                                                                  `tfsdk:"url" json:"url,computed"`
	UseCompression  types.Bool                                                                                    `tfsdk:"use_compression" json:"use_compression,computed"`
}

type CDNLogsUploaderTargetConfigAppendResponseActionsDataSourceModel struct {
	Action          types.String `tfsdk:"action" json:"action,computed"`
	Description     types.String `tfsdk:"description" json:"description,computed"`
	MatchPayload    types.String `tfsdk:"match_payload" json:"match_payload,computed"`
	MatchStatusCode types.Int64  `tfsdk:"match_status_code" json:"match_status_code,computed"`
}

type CDNLogsUploaderTargetConfigAuthDataSourceModel struct {
	Config customfield.NestedObject[CDNLogsUploaderTargetConfigAuthConfigDataSourceModel] `tfsdk:"config" json:"config,computed"`
	Type   types.String                                                                   `tfsdk:"type" json:"type,computed"`
}

type CDNLogsUploaderTargetConfigAuthConfigDataSourceModel struct {
	Token           types.String `tfsdk:"token" json:"token,computed"`
	HeaderName      types.String `tfsdk:"header_name" json:"header_name,computed"`
	AccountKey      types.String `tfsdk:"account_key" json:"account_key,computed"`
	AccessKeyID     types.String `tfsdk:"access_key_id" json:"access_key_id,computed"`
	SecretAccessKey types.String `tfsdk:"secret_access_key" json:"secret_access_key,computed"`
}

type CDNLogsUploaderTargetConfigRetryDataSourceModel struct {
	Headers         customfield.Map[types.String]                                                                `tfsdk:"headers" json:"headers,computed"`
	Method          types.String                                                                                 `tfsdk:"method" json:"method,computed"`
	ResponseActions customfield.NestedObjectList[CDNLogsUploaderTargetConfigRetryResponseActionsDataSourceModel] `tfsdk:"response_actions" json:"response_actions,computed"`
	TimeoutSeconds  types.Int64                                                                                  `tfsdk:"timeout_seconds" json:"timeout_seconds,computed"`
	URL             types.String                                                                                 `tfsdk:"url" json:"url,computed"`
	UseCompression  types.Bool                                                                                   `tfsdk:"use_compression" json:"use_compression,computed"`
}

type CDNLogsUploaderTargetConfigRetryResponseActionsDataSourceModel struct {
	Action          types.String `tfsdk:"action" json:"action,computed"`
	Description     types.String `tfsdk:"description" json:"description,computed"`
	MatchPayload    types.String `tfsdk:"match_payload" json:"match_payload,computed"`
	MatchStatusCode types.Int64  `tfsdk:"match_status_code" json:"match_status_code,computed"`
}

type CDNLogsUploaderTargetConfigUploadDataSourceModel struct {
	Headers         customfield.Map[types.String]                                                                 `tfsdk:"headers" json:"headers,computed"`
	Method          types.String                                                                                  `tfsdk:"method" json:"method,computed"`
	ResponseActions customfield.NestedObjectList[CDNLogsUploaderTargetConfigUploadResponseActionsDataSourceModel] `tfsdk:"response_actions" json:"response_actions,computed"`
	TimeoutSeconds  types.Int64                                                                                   `tfsdk:"timeout_seconds" json:"timeout_seconds,computed"`
	URL             types.String                                                                                  `tfsdk:"url" json:"url,computed"`
	UseCompression  types.Bool                                                                                    `tfsdk:"use_compression" json:"use_compression,computed"`
}

type CDNLogsUploaderTargetConfigUploadResponseActionsDataSourceModel struct {
	Action          types.String `tfsdk:"action" json:"action,computed"`
	Description     types.String `tfsdk:"description" json:"description,computed"`
	MatchPayload    types.String `tfsdk:"match_payload" json:"match_payload,computed"`
	MatchStatusCode types.Int64  `tfsdk:"match_status_code" json:"match_status_code,computed"`
}

type CDNLogsUploaderTargetFindOneByDataSourceModel struct {
	ConfigIDs *[]types.Int64 `tfsdk:"config_ids" query:"config_ids,optional"`
	Search    types.String   `tfsdk:"search" query:"search,optional"`
}
