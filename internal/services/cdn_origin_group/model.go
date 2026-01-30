// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cdn_origin_group

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/apijson"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

type CDNOriginGroupModel struct {
	ID                  types.Int64                    `tfsdk:"id" json:"id,computed"`
	Name                types.String                   `tfsdk:"name" json:"name,required"`
	AuthType            types.String                   `tfsdk:"auth_type" json:"auth_type,optional"`
	Path                types.String                   `tfsdk:"path" json:"path,optional"`
	UseNext             types.Bool                     `tfsdk:"use_next" json:"use_next,optional"`
	Auth                *CDNOriginGroupAuthModel       `tfsdk:"auth" json:"auth,optional"`
	Sources             *[]*CDNOriginGroupSourcesModel `tfsdk:"sources" json:"sources,optional"`
	ProxyNextUpstream   customfield.List[types.String] `tfsdk:"proxy_next_upstream" json:"proxy_next_upstream,computed_optional"`
	HasRelatedResources types.Bool                     `tfsdk:"has_related_resources" json:"has_related_resources,computed"`
}

func (m CDNOriginGroupModel) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(m)
}

func (m CDNOriginGroupModel) MarshalJSONForUpdate(state CDNOriginGroupModel) (data []byte, err error) {
	return apijson.MarshalForPatch(m, state)
}

type CDNOriginGroupAuthModel struct {
	S3AccessKeyID     types.String `tfsdk:"s3_access_key_id" json:"s3_access_key_id,required"`
	S3BucketName      types.String `tfsdk:"s3_bucket_name" json:"s3_bucket_name,required"`
	S3SecretAccessKey types.String `tfsdk:"s3_secret_access_key" json:"s3_secret_access_key,required"`
	S3Type            types.String `tfsdk:"s3_type" json:"s3_type,required"`
	S3Region          types.String `tfsdk:"s3_region" json:"s3_region,optional"`
	S3StorageHostname types.String `tfsdk:"s3_storage_hostname" json:"s3_storage_hostname,optional"`
}

type CDNOriginGroupSourcesModel struct {
	Backup  types.Bool   `tfsdk:"backup" json:"backup,optional"`
	Enabled types.Bool   `tfsdk:"enabled" json:"enabled,optional"`
	Source  types.String `tfsdk:"source" json:"source,optional"`
}
