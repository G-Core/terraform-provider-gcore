// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cdn_origin_group

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

type CDNOriginGroupDataSourceModel struct {
	ID                  types.Int64                                                       `tfsdk:"id" path:"origin_group_id,computed"`
	OriginGroupID       types.Int64                                                       `tfsdk:"origin_group_id" path:"origin_group_id,required"`
	AuthType            types.String                                                      `tfsdk:"auth_type" json:"auth_type,computed"`
	HasRelatedResources types.Bool                                                        `tfsdk:"has_related_resources" json:"has_related_resources,computed"`
	Name                types.String                                                      `tfsdk:"name" json:"name,computed"`
	Path                types.String                                                      `tfsdk:"path" json:"path,computed"`
	UseNext             types.Bool                                                        `tfsdk:"use_next" json:"use_next,computed"`
	ProxyNextUpstream   customfield.List[types.String]                                    `tfsdk:"proxy_next_upstream" json:"proxy_next_upstream,computed"`
	Auth                customfield.NestedObject[CDNOriginGroupAuthDataSourceModel]       `tfsdk:"auth" json:"auth,computed"`
	Sources             customfield.NestedObjectSet[CDNOriginGroupSourcesDataSourceModel] `tfsdk:"sources" json:"sources,computed"`
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
	Backup  types.Bool   `tfsdk:"backup" json:"backup,computed"`
	Enabled types.Bool   `tfsdk:"enabled" json:"enabled,computed"`
	Source  types.String `tfsdk:"source" json:"source,computed"`
}
