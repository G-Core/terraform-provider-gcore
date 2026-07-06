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

type CDNOriginGroupsResultsListDataSourceEnvelope struct {
	Results customfield.NestedObjectList[CDNOriginGroupsItemsDataSourceModel] `json:"results,computed"`
}

type CDNOriginGroupsDataSourceModel struct {
	HasRelatedResources types.Bool                                                        `tfsdk:"has_related_resources" query:"has_related_resources,optional"`
	Name                types.String                                                      `tfsdk:"name" query:"name,optional"`
	Sources             types.String                                                      `tfsdk:"sources" query:"sources,optional"`
	MaxItems            types.Int64                                                       `tfsdk:"max_items"`
	Items               customfield.NestedObjectList[CDNOriginGroupsItemsDataSourceModel] `tfsdk:"items"`
}

func (m *CDNOriginGroupsDataSourceModel) toListParams(_ context.Context) (params cdn.OriginGroupListParams, diags diag.Diagnostics) {
	params = cdn.OriginGroupListParams{}

	if !m.HasRelatedResources.IsNull() {
		params.HasRelatedResources = param.NewOpt(m.HasRelatedResources.ValueBool())
	}
	if !m.Name.IsNull() {
		params.Name = param.NewOpt(m.Name.ValueString())
	}
	if !m.Sources.IsNull() {
		params.Sources = param.NewOpt(m.Sources.ValueString())
	}

	return
}

type CDNOriginGroupsItemsDataSourceModel struct {
	ID                  types.Int64                                                        `tfsdk:"id" json:"id,computed"`
	Name                types.String                                                       `tfsdk:"name" json:"name,computed"`
	Sources             customfield.NestedObjectSet[CDNOriginGroupsSourcesDataSourceModel] `tfsdk:"sources" json:"sources,computed"`
	AuthType            types.String                                                       `tfsdk:"auth_type" json:"auth_type,computed"`
	HasRelatedResources types.Bool                                                         `tfsdk:"has_related_resources" json:"has_related_resources,computed"`
	Path                types.String                                                       `tfsdk:"path" json:"path,computed"`
	ProxyNextUpstream   customfield.List[types.String]                                     `tfsdk:"proxy_next_upstream" json:"proxy_next_upstream,computed"`
	UseNext             types.Bool                                                         `tfsdk:"use_next" json:"use_next,computed"`
	Auth                customfield.NestedObject[CDNOriginGroupsAuthDataSourceModel]       `tfsdk:"auth" json:"auth,computed"`
}

type CDNOriginGroupsSourcesDataSourceModel struct {
	Source             types.String                                                          `tfsdk:"source" json:"source,computed"`
	Backup             types.Bool                                                            `tfsdk:"backup" json:"backup,computed"`
	Enabled            types.Bool                                                            `tfsdk:"enabled" json:"enabled,computed"`
	HostHeaderOverride types.String                                                          `tfsdk:"host_header_override" json:"host_header_override,computed"`
	Config             customfield.NestedObject[CDNOriginGroupsSourcesConfigDataSourceModel] `tfsdk:"config" json:"config,computed"`
	OriginType         types.String                                                          `tfsdk:"origin_type" json:"origin_type,computed"`
}

type CDNOriginGroupsSourcesConfigDataSourceModel struct {
	S3AccessKeyID     types.String `tfsdk:"s3_access_key_id" json:"s3_access_key_id,computed"`
	S3BucketName      types.String `tfsdk:"s3_bucket_name" json:"s3_bucket_name,computed"`
	S3SecretAccessKey types.String `tfsdk:"s3_secret_access_key" json:"s3_secret_access_key,computed"`
	S3Type            types.String `tfsdk:"s3_type" json:"s3_type,computed"`
	S3AuthType        types.String `tfsdk:"s3_auth_type" json:"s3_auth_type,computed"`
	S3Region          types.String `tfsdk:"s3_region" json:"s3_region,computed"`
	S3StorageHostname types.String `tfsdk:"s3_storage_hostname" json:"s3_storage_hostname,computed"`
	AppID             types.String `tfsdk:"app_id" json:"app_id,computed"`
}

type CDNOriginGroupsAuthDataSourceModel struct {
	S3AccessKeyID     types.String `tfsdk:"s3_access_key_id" json:"s3_access_key_id,computed"`
	S3BucketName      types.String `tfsdk:"s3_bucket_name" json:"s3_bucket_name,computed"`
	S3SecretAccessKey types.String `tfsdk:"s3_secret_access_key" json:"s3_secret_access_key,computed"`
	S3Type            types.String `tfsdk:"s3_type" json:"s3_type,computed"`
	S3Region          types.String `tfsdk:"s3_region" json:"s3_region,computed"`
	S3StorageHostname types.String `tfsdk:"s3_storage_hostname" json:"s3_storage_hostname,computed"`
}
