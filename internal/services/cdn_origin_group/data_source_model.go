// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cdn_origin_group

import (
	"context"

	"github.com/G-Core/gcore-go/cdn"
	"github.com/G-Core/gcore-go/packages/param"
	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type CDNOriginGroupDataSourceModel struct {
	ID                  types.Int64                                                       `tfsdk:"id" path:"origin_group_id,computed"`
	OriginGroupID       types.Int64                                                       `tfsdk:"origin_group_id" path:"origin_group_id,optional"`
	AuthType            types.String                                                      `tfsdk:"auth_type" json:"auth_type,computed"`
	HasRelatedResources types.Bool                                                        `tfsdk:"has_related_resources" json:"has_related_resources,computed"`
	Name                types.String                                                      `tfsdk:"name" json:"name,computed"`
	Path                types.String                                                      `tfsdk:"path" json:"path,computed"`
	UseNext             types.Bool                                                        `tfsdk:"use_next" json:"use_next,computed"`
	ProxyNextUpstream   customfield.List[types.String]                                    `tfsdk:"proxy_next_upstream" json:"proxy_next_upstream,computed"`
	Auth                customfield.NestedObject[CDNOriginGroupAuthDataSourceModel]       `tfsdk:"auth" json:"auth,computed"`
	Sources             customfield.NestedObjectSet[CDNOriginGroupSourcesDataSourceModel] `tfsdk:"sources" json:"sources,computed"`
	FindOneBy           *CDNOriginGroupFindOneByDataSourceModel                           `tfsdk:"find_one_by"`
}

func (m *CDNOriginGroupDataSourceModel) toListParams(_ context.Context) (params cdn.OriginGroupListParams, diags diag.Diagnostics) {
	params = cdn.OriginGroupListParams{}

	if !m.FindOneBy.HasRelatedResources.IsNull() {
		params.HasRelatedResources = param.NewOpt(m.FindOneBy.HasRelatedResources.ValueBool())
	}
	if !m.FindOneBy.Name.IsNull() {
		params.Name = param.NewOpt(m.FindOneBy.Name.ValueString())
	}
	if !m.FindOneBy.Sources.IsNull() {
		params.Sources = param.NewOpt(m.FindOneBy.Sources.ValueString())
	}

	return
}

type CDNOriginGroupAuthDataSourceModel struct {
	S3AccessKeyID     types.String `tfsdk:"s3_access_key_id" json:"s3_access_key_id,computed"`
	S3BucketName      types.String `tfsdk:"s3_bucket_name" json:"s3_bucket_name,computed"`
	S3SecretAccessKey types.String `tfsdk:"s3_secret_access_key" json:"s3_secret_access_key,computed"`
	S3Type            types.String `tfsdk:"s3_type" json:"s3_type,computed"`
	S3Region          types.String `tfsdk:"s3_region" json:"s3_region,computed"`
	S3StorageHostname types.String `tfsdk:"s3_storage_hostname" json:"s3_storage_hostname,computed"`
}

type CDNOriginGroupSourcesDataSourceModel struct {
	Source             types.String                                                         `tfsdk:"source" json:"source,computed"`
	Backup             types.Bool                                                           `tfsdk:"backup" json:"backup,computed"`
	Enabled            types.Bool                                                           `tfsdk:"enabled" json:"enabled,computed"`
	HostHeaderOverride types.String                                                         `tfsdk:"host_header_override" json:"host_header_override,computed"`
	Config             customfield.NestedObject[CDNOriginGroupSourcesConfigDataSourceModel] `tfsdk:"config" json:"config,computed"`
	OriginType         types.String                                                         `tfsdk:"origin_type" json:"origin_type,computed"`
}

type CDNOriginGroupSourcesConfigDataSourceModel struct {
	S3AccessKeyID     types.String `tfsdk:"s3_access_key_id" json:"s3_access_key_id,computed"`
	S3BucketName      types.String `tfsdk:"s3_bucket_name" json:"s3_bucket_name,computed"`
	S3SecretAccessKey types.String `tfsdk:"s3_secret_access_key" json:"s3_secret_access_key,computed"`
	S3Type            types.String `tfsdk:"s3_type" json:"s3_type,computed"`
	S3AuthType        types.String `tfsdk:"s3_auth_type" json:"s3_auth_type,computed"`
	S3Region          types.String `tfsdk:"s3_region" json:"s3_region,computed"`
	S3StorageHostname types.String `tfsdk:"s3_storage_hostname" json:"s3_storage_hostname,computed"`
	AppID             types.String `tfsdk:"app_id" json:"app_id,computed"`
}

type CDNOriginGroupFindOneByDataSourceModel struct {
	HasRelatedResources types.Bool   `tfsdk:"has_related_resources" query:"has_related_resources,optional"`
	Name                types.String `tfsdk:"name" query:"name,optional"`
	Sources             types.String `tfsdk:"sources" query:"sources,optional"`
}
