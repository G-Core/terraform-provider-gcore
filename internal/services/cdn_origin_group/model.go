// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cdn_origin_group

import (
	"github.com/G-Core/terraform-provider-gcore/internal/apijson"
	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type CDNOriginGroupModel struct {
	ID                  types.Int64                                              `tfsdk:"id" json:"id,computed"`
	Name                types.String                                             `tfsdk:"name" json:"name,required"`
	AuthType            types.String                                             `tfsdk:"auth_type" json:"auth_type,optional"`
	UseNext             types.Bool                                               `tfsdk:"use_next" json:"use_next,optional"`
	Auth                *CDNOriginGroupAuthModel                                 `tfsdk:"auth" json:"auth,optional"`
	ProxyNextUpstream   customfield.List[types.String]                           `tfsdk:"proxy_next_upstream" json:"proxy_next_upstream,computed_optional"`
	Sources             customfield.NestedObjectList[CDNOriginGroupSourcesModel] `tfsdk:"sources" json:"sources,computed_optional,no_refresh"`
	HasRelatedResources types.Bool                                               `tfsdk:"has_related_resources" json:"has_related_resources,computed"`
}

func (m CDNOriginGroupModel) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(m)
}

func (m CDNOriginGroupModel) MarshalJSONForUpdate(state CDNOriginGroupModel) (data []byte, err error) {
	return apijson.MarshalForPatch(m, state)
}

type CDNOriginGroupAuthModel struct {
	S3AccessKeyID        types.String `tfsdk:"s3_access_key_id" json:"s3_access_key_id,required"`
	S3BucketName         types.String `tfsdk:"s3_bucket_name" json:"s3_bucket_name,required"`
	S3SecretAccessKey    types.String `tfsdk:"s3_secret_access_key" json:"s3_secret_access_key,required"`
	S3Type               types.String `tfsdk:"s3_type" json:"s3_type,required"`
	S3Region             types.String `tfsdk:"s3_region" json:"s3_region,optional"`
	S3StorageHostname    types.String `tfsdk:"s3_storage_hostname" json:"s3_storage_hostname,optional"`
	S3CredentialsVersion types.Int64  `tfsdk:"s3_credentials_version"` // Trigger for credential updates, not sent to API
}

type CDNOriginGroupSourcesModel struct {
	Source             types.String                                               `tfsdk:"source" json:"source,optional"`
	Backup             types.Bool                                                 `tfsdk:"backup" json:"backup,computed_optional"`
	Enabled            types.Bool                                                 `tfsdk:"enabled" json:"enabled,computed_optional"`
	HostHeaderOverride types.String                                               `tfsdk:"host_header_override" json:"host_header_override,optional"`
	Tag                types.String                                               `tfsdk:"tag" json:"tag,computed_optional"`
	Config             customfield.NestedObject[CDNOriginGroupSourcesConfigModel] `tfsdk:"config" json:"config,computed_optional"`
	OriginType         types.String                                               `tfsdk:"origin_type" json:"origin_type,optional"`
}

type CDNOriginGroupSourcesConfigModel struct {
	S3AccessKeyID     types.String `tfsdk:"s3_access_key_id" json:"s3_access_key_id,required"`
	S3BucketName      types.String `tfsdk:"s3_bucket_name" json:"s3_bucket_name,required"`
	S3SecretAccessKey types.String `tfsdk:"s3_secret_access_key" json:"s3_secret_access_key,required"`
	S3Type            types.String `tfsdk:"s3_type" json:"s3_type,required"`
	S3AuthType        types.String `tfsdk:"s3_auth_type" json:"s3_auth_type,computed_optional"`
	S3Region          types.String `tfsdk:"s3_region" json:"s3_region,optional"`
	S3StorageHostname types.String `tfsdk:"s3_storage_hostname" json:"s3_storage_hostname,optional"`
}
