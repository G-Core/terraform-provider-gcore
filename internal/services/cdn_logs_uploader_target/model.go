// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cdn_logs_uploader_target

import (
	"github.com/G-Core/terraform-provider-gcore/internal/apijson"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type CDNLogsUploaderTargetModel struct {
	ID          types.Int64                       `tfsdk:"id" json:"id,computed"`
	StorageType types.String                      `tfsdk:"storage_type" json:"storage_type,required"`
	Config      *CDNLogsUploaderTargetConfigModel `tfsdk:"config" json:"config,required"`
	Description types.String                      `tfsdk:"description" json:"description,computed_optional"`
	Name        types.String                      `tfsdk:"name" json:"name,computed_optional"`
	ClientID    types.Int64                       `tfsdk:"client_id" json:"client_id,computed"`
	Created     timetypes.RFC3339                 `tfsdk:"created" json:"created,computed" format:"date-time"`
}

func (m CDNLogsUploaderTargetModel) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(m)
}

func (m CDNLogsUploaderTargetModel) MarshalJSONForUpdate(state CDNLogsUploaderTargetModel) (data []byte, err error) {
	return apijson.MarshalForPatch(m, state)
}

type CDNLogsUploaderTargetConfigModel struct {
	AccessKeyID     types.String                            `tfsdk:"access_key_id" json:"access_key_id,optional"`
	BucketName      types.String                            `tfsdk:"bucket_name" json:"bucket_name,optional"`
	Directory       types.String                            `tfsdk:"directory" json:"directory,optional"`
	Endpoint        types.String                            `tfsdk:"endpoint" json:"endpoint,optional"`
	Region          types.String                            `tfsdk:"region" json:"region,optional"`
	SecretAccessKey types.String                            `tfsdk:"secret_access_key" json:"secret_access_key,optional,no_refresh"`
	UsePathStyle    types.Bool                              `tfsdk:"use_path_style" json:"use_path_style,computed_optional"`
	Hostname        types.String                            `tfsdk:"hostname" json:"hostname,optional"`
	Password        types.String                            `tfsdk:"password" json:"password,optional,no_refresh"`
	TimeoutSeconds  types.Int64                             `tfsdk:"timeout_seconds" json:"timeout_seconds,optional"`
	User            types.String                            `tfsdk:"user" json:"user,optional"`
	KeyPassphrase   types.String                            `tfsdk:"key_passphrase" json:"key_passphrase,optional,no_refresh"`
	PrivateKey      types.String                            `tfsdk:"private_key" json:"private_key,optional,no_refresh"`
	Upload          *CDNLogsUploaderTargetConfigUploadModel `tfsdk:"upload" json:"upload,optional"`
	Append          *CDNLogsUploaderTargetConfigAppendModel `tfsdk:"append" json:"append,optional"`
	Auth            *CDNLogsUploaderTargetConfigAuthModel   `tfsdk:"auth" json:"auth,optional"`
	ContentType     types.String                            `tfsdk:"content_type" json:"content_type,computed_optional"`
	Retry           *CDNLogsUploaderTargetConfigRetryModel  `tfsdk:"retry" json:"retry,optional"`
	AccountName     types.String                            `tfsdk:"account_name" json:"account_name,optional"`
	ContainerName   types.String                            `tfsdk:"container_name" json:"container_name,optional"`
	LogStore        types.String                            `tfsdk:"log_store" json:"log_store,optional"`
	Project         types.String                            `tfsdk:"project" json:"project,optional"`
	Topic           types.String                            `tfsdk:"topic" json:"topic,optional"`
}

type CDNLogsUploaderTargetConfigUploadModel struct {
	URL             types.String                                              `tfsdk:"url" json:"url,required"`
	Headers         *map[string]types.String                                  `tfsdk:"headers" json:"headers,optional"`
	Method          types.String                                              `tfsdk:"method" json:"method,computed_optional"`
	ResponseActions *[]*CDNLogsUploaderTargetConfigUploadResponseActionsModel `tfsdk:"response_actions" json:"response_actions,computed_optional"`
	TimeoutSeconds  types.Int64                                               `tfsdk:"timeout_seconds" json:"timeout_seconds,computed_optional"`
	UseCompression  types.Bool                                                `tfsdk:"use_compression" json:"use_compression,computed_optional"`
}

type CDNLogsUploaderTargetConfigUploadResponseActionsModel struct {
	Action          types.String `tfsdk:"action" json:"action,required"`
	Description     types.String `tfsdk:"description" json:"description,computed_optional"`
	MatchPayload    types.String `tfsdk:"match_payload" json:"match_payload,computed_optional"`
	MatchStatusCode types.Int64  `tfsdk:"match_status_code" json:"match_status_code,computed_optional"`
}

type CDNLogsUploaderTargetConfigAppendModel struct {
	URL             types.String                                              `tfsdk:"url" json:"url,required"`
	Headers         *map[string]types.String                                  `tfsdk:"headers" json:"headers,optional"`
	Method          types.String                                              `tfsdk:"method" json:"method,computed_optional"`
	ResponseActions *[]*CDNLogsUploaderTargetConfigAppendResponseActionsModel `tfsdk:"response_actions" json:"response_actions,computed_optional"`
	TimeoutSeconds  types.Int64                                               `tfsdk:"timeout_seconds" json:"timeout_seconds,computed_optional"`
	UseCompression  types.Bool                                                `tfsdk:"use_compression" json:"use_compression,computed_optional"`
}

type CDNLogsUploaderTargetConfigAppendResponseActionsModel struct {
	Action          types.String `tfsdk:"action" json:"action,required"`
	Description     types.String `tfsdk:"description" json:"description,computed_optional"`
	MatchPayload    types.String `tfsdk:"match_payload" json:"match_payload,computed_optional"`
	MatchStatusCode types.Int64  `tfsdk:"match_status_code" json:"match_status_code,computed_optional"`
}

type CDNLogsUploaderTargetConfigAuthModel struct {
	Config *CDNLogsUploaderTargetConfigAuthConfigModel `tfsdk:"config" json:"config,required"`
	Type   types.String                                `tfsdk:"type" json:"type,required"`
}

type CDNLogsUploaderTargetConfigAuthConfigModel struct {
	Token           types.String `tfsdk:"token" json:"token,optional,no_refresh"`
	HeaderName      types.String `tfsdk:"header_name" json:"header_name,optional"`
	AccountKey      types.String `tfsdk:"account_key" json:"account_key,optional,no_refresh"`
	AccessKeyID     types.String `tfsdk:"access_key_id" json:"access_key_id,optional"`
	SecretAccessKey types.String `tfsdk:"secret_access_key" json:"secret_access_key,optional,no_refresh"`
}

type CDNLogsUploaderTargetConfigRetryModel struct {
	URL             types.String                                             `tfsdk:"url" json:"url,required"`
	Headers         *map[string]types.String                                 `tfsdk:"headers" json:"headers,optional"`
	Method          types.String                                             `tfsdk:"method" json:"method,computed_optional"`
	ResponseActions *[]*CDNLogsUploaderTargetConfigRetryResponseActionsModel `tfsdk:"response_actions" json:"response_actions,computed_optional"`
	TimeoutSeconds  types.Int64                                              `tfsdk:"timeout_seconds" json:"timeout_seconds,computed_optional"`
	UseCompression  types.Bool                                               `tfsdk:"use_compression" json:"use_compression,computed_optional"`
}

type CDNLogsUploaderTargetConfigRetryResponseActionsModel struct {
	Action          types.String `tfsdk:"action" json:"action,required"`
	Description     types.String `tfsdk:"description" json:"description,computed_optional"`
	MatchPayload    types.String `tfsdk:"match_payload" json:"match_payload,computed_optional"`
	MatchStatusCode types.Int64  `tfsdk:"match_status_code" json:"match_status_code,computed_optional"`
}
