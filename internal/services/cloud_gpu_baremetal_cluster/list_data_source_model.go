// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_gpu_baremetal_cluster

import (
	"context"

	"github.com/G-Core/gcore-go/cloud"
	"github.com/G-Core/gcore-go/packages/param"
	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type CloudGPUBaremetalClustersResultsListDataSourceEnvelope struct {
	Results customfield.NestedObjectList[CloudGPUBaremetalClustersItemsDataSourceModel] `json:"results,computed"`
}

type CloudGPUBaremetalClustersDataSourceModel struct {
	ProjectID    types.Int64                                                                 `tfsdk:"project_id" path:"project_id,optional"`
	RegionID     types.Int64                                                                 `tfsdk:"region_id" path:"region_id,optional"`
	IDs          *[]types.String                                                             `tfsdk:"ids" query:"ids,optional"`
	ImageIDs     *[]types.String                                                             `tfsdk:"image_ids" query:"image_ids,optional"`
	Tags         *map[string]types.String                                                    `tfsdk:"tags" query:"tags,optional"`
	CreatedAt    *CloudGPUBaremetalClustersCreatedAtDataSourceModel                          `tfsdk:"created_at" query:"created_at,optional"`
	Flavor       *CloudGPUBaremetalClustersFlavorDataSourceModel                             `tfsdk:"flavor" query:"flavor,optional"`
	Name         *CloudGPUBaremetalClustersNameDataSourceModel                               `tfsdk:"name" query:"name,optional"`
	ServersCount *CloudGPUBaremetalClustersServersCountDataSourceModel                       `tfsdk:"servers_count" query:"servers_count,optional"`
	TagKey       *CloudGPUBaremetalClustersTagKeyDataSourceModel                             `tfsdk:"tag_key" query:"tag_key,optional"`
	TagValue     *CloudGPUBaremetalClustersTagValueDataSourceModel                           `tfsdk:"tag_value" query:"tag_value,optional"`
	UpdatedAt    *CloudGPUBaremetalClustersUpdatedAtDataSourceModel                          `tfsdk:"updated_at" query:"updated_at,optional"`
	ManagedBy    customfield.List[types.String]                                              `tfsdk:"managed_by" query:"managed_by,computed_optional"`
	MaxItems     types.Int64                                                                 `tfsdk:"max_items"`
	Items        customfield.NestedObjectList[CloudGPUBaremetalClustersItemsDataSourceModel] `tfsdk:"items"`
}

func (m *CloudGPUBaremetalClustersDataSourceModel) toListParams(ctx context.Context) (params cloud.GPUBaremetalClusterListParams, diags diag.Diagnostics) {
	mIDs := []string{}
	if m.IDs != nil {
		for _, item := range *m.IDs {
			mIDs = append(mIDs, item.ValueString())
		}
	}
	mImageIDs := []string{}
	if m.ImageIDs != nil {
		for _, item := range *m.ImageIDs {
			mImageIDs = append(mImageIDs, item.ValueString())
		}
	}
	mManagedBy := []string{}
	diags.Append(m.ManagedBy.ElementsAs(ctx, &mManagedBy, true)...)
	if diags.HasError() {
		return
	}
	mCreatedAtGt, errs := m.CreatedAt.Gt.ValueRFC3339Time()
	diags.Append(errs...)
	mCreatedAtGte, errs := m.CreatedAt.Gte.ValueRFC3339Time()
	diags.Append(errs...)
	mCreatedAtLt, errs := m.CreatedAt.Lt.ValueRFC3339Time()
	diags.Append(errs...)
	mCreatedAtLte, errs := m.CreatedAt.Lte.ValueRFC3339Time()
	diags.Append(errs...)
	mFlavorContains := []string{}
	if m.Flavor.Contains != nil {
		for _, item := range *m.Flavor.Contains {
			if !item.IsNull() {
				mFlavorContains = append(mFlavorContains, item.ValueString())
			}
		}
	}
	mFlavorExact := []string{}
	if m.Flavor.Exact != nil {
		for _, item := range *m.Flavor.Exact {
			if !item.IsNull() {
				mFlavorExact = append(mFlavorExact, item.ValueString())
			}
		}
	}
	mFlavorPrefix := []string{}
	if m.Flavor.Prefix != nil {
		for _, item := range *m.Flavor.Prefix {
			if !item.IsNull() {
				mFlavorPrefix = append(mFlavorPrefix, item.ValueString())
			}
		}
	}
	mFlavorSuffix := []string{}
	if m.Flavor.Suffix != nil {
		for _, item := range *m.Flavor.Suffix {
			if !item.IsNull() {
				mFlavorSuffix = append(mFlavorSuffix, item.ValueString())
			}
		}
	}
	mNameContains := []string{}
	if m.Name.Contains != nil {
		for _, item := range *m.Name.Contains {
			if !item.IsNull() {
				mNameContains = append(mNameContains, item.ValueString())
			}
		}
	}
	mNameExact := []string{}
	if m.Name.Exact != nil {
		for _, item := range *m.Name.Exact {
			if !item.IsNull() {
				mNameExact = append(mNameExact, item.ValueString())
			}
		}
	}
	mNamePrefix := []string{}
	if m.Name.Prefix != nil {
		for _, item := range *m.Name.Prefix {
			if !item.IsNull() {
				mNamePrefix = append(mNamePrefix, item.ValueString())
			}
		}
	}
	mNameSuffix := []string{}
	if m.Name.Suffix != nil {
		for _, item := range *m.Name.Suffix {
			if !item.IsNull() {
				mNameSuffix = append(mNameSuffix, item.ValueString())
			}
		}
	}
	mTagKeyContains := []string{}
	if m.TagKey.Contains != nil {
		for _, item := range *m.TagKey.Contains {
			if !item.IsNull() {
				mTagKeyContains = append(mTagKeyContains, item.ValueString())
			}
		}
	}
	mTagKeyExact := []string{}
	if m.TagKey.Exact != nil {
		for _, item := range *m.TagKey.Exact {
			if !item.IsNull() {
				mTagKeyExact = append(mTagKeyExact, item.ValueString())
			}
		}
	}
	mTagKeyPrefix := []string{}
	if m.TagKey.Prefix != nil {
		for _, item := range *m.TagKey.Prefix {
			if !item.IsNull() {
				mTagKeyPrefix = append(mTagKeyPrefix, item.ValueString())
			}
		}
	}
	mTagKeySuffix := []string{}
	if m.TagKey.Suffix != nil {
		for _, item := range *m.TagKey.Suffix {
			if !item.IsNull() {
				mTagKeySuffix = append(mTagKeySuffix, item.ValueString())
			}
		}
	}
	mTagValueContains := []string{}
	if m.TagValue.Contains != nil {
		for _, item := range *m.TagValue.Contains {
			if !item.IsNull() {
				mTagValueContains = append(mTagValueContains, item.ValueString())
			}
		}
	}
	mTagValueExact := []string{}
	if m.TagValue.Exact != nil {
		for _, item := range *m.TagValue.Exact {
			if !item.IsNull() {
				mTagValueExact = append(mTagValueExact, item.ValueString())
			}
		}
	}
	mTagValuePrefix := []string{}
	if m.TagValue.Prefix != nil {
		for _, item := range *m.TagValue.Prefix {
			if !item.IsNull() {
				mTagValuePrefix = append(mTagValuePrefix, item.ValueString())
			}
		}
	}
	mTagValueSuffix := []string{}
	if m.TagValue.Suffix != nil {
		for _, item := range *m.TagValue.Suffix {
			if !item.IsNull() {
				mTagValueSuffix = append(mTagValueSuffix, item.ValueString())
			}
		}
	}
	mTags := map[string]string{}
	for key, value := range *m.Tags {
		if !value.IsNull() {
			mTags[key] = value.ValueString()
		}
	}
	mUpdatedAtGt, errs := m.UpdatedAt.Gt.ValueRFC3339Time()
	diags.Append(errs...)
	mUpdatedAtGte, errs := m.UpdatedAt.Gte.ValueRFC3339Time()
	diags.Append(errs...)
	mUpdatedAtLt, errs := m.UpdatedAt.Lt.ValueRFC3339Time()
	diags.Append(errs...)
	mUpdatedAtLte, errs := m.UpdatedAt.Lte.ValueRFC3339Time()
	diags.Append(errs...)

	params = cloud.GPUBaremetalClusterListParams{
		IDs:       mIDs,
		ImageIDs:  mImageIDs,
		ManagedBy: mManagedBy,
	}

	if !m.ProjectID.IsNull() {
		params.ProjectID = param.NewOpt(m.ProjectID.ValueInt64())
	}
	if !m.RegionID.IsNull() {
		params.RegionID = param.NewOpt(m.RegionID.ValueInt64())
	}
	if m.CreatedAt != nil {
		paramsCreatedAt := cloud.GPUBaremetalClusterListParamsCreatedAt{}
		if !m.CreatedAt.Gt.IsNull() {
			paramsCreatedAt.Gt = param.NewOpt(mCreatedAtGt)
		}
		if !m.CreatedAt.Gte.IsNull() {
			paramsCreatedAt.Gte = param.NewOpt(mCreatedAtGte)
		}
		if !m.CreatedAt.Lt.IsNull() {
			paramsCreatedAt.Lt = param.NewOpt(mCreatedAtLt)
		}
		if !m.CreatedAt.Lte.IsNull() {
			paramsCreatedAt.Lte = param.NewOpt(mCreatedAtLte)
		}
		params.CreatedAt = paramsCreatedAt
	}
	if m.Flavor != nil {
		params.Flavor = cloud.GPUBaremetalClusterListParamsFlavor{
			Contains: mFlavorContains,
			Exact:    mFlavorExact,
			Prefix:   mFlavorPrefix,
			Suffix:   mFlavorSuffix,
		}
	}
	if m.Name != nil {
		params.Name = cloud.GPUBaremetalClusterListParamsName{
			Contains: mNameContains,
			Exact:    mNameExact,
			Prefix:   mNamePrefix,
			Suffix:   mNameSuffix,
		}
	}
	if m.ServersCount != nil {
		paramsServersCount := cloud.GPUBaremetalClusterListParamsServersCount{}
		if !m.ServersCount.Gt.IsNull() {
			paramsServersCount.Gt = param.NewOpt(m.ServersCount.Gt.ValueInt64())
		}
		if !m.ServersCount.Gte.IsNull() {
			paramsServersCount.Gte = param.NewOpt(m.ServersCount.Gte.ValueInt64())
		}
		if !m.ServersCount.Lt.IsNull() {
			paramsServersCount.Lt = param.NewOpt(m.ServersCount.Lt.ValueInt64())
		}
		if !m.ServersCount.Lte.IsNull() {
			paramsServersCount.Lte = param.NewOpt(m.ServersCount.Lte.ValueInt64())
		}
		params.ServersCount = paramsServersCount
	}
	if m.TagKey != nil {
		params.TagKey = cloud.GPUBaremetalClusterListParamsTagKey{
			Contains: mTagKeyContains,
			Exact:    mTagKeyExact,
			Prefix:   mTagKeyPrefix,
			Suffix:   mTagKeySuffix,
		}
	}
	if m.TagValue != nil {
		params.TagValue = cloud.GPUBaremetalClusterListParamsTagValue{
			Contains: mTagValueContains,
			Exact:    mTagValueExact,
			Prefix:   mTagValuePrefix,
			Suffix:   mTagValueSuffix,
		}
	}
	if m.UpdatedAt != nil {
		paramsUpdatedAt := cloud.GPUBaremetalClusterListParamsUpdatedAt{}
		if !m.UpdatedAt.Gt.IsNull() {
			paramsUpdatedAt.Gt = param.NewOpt(mUpdatedAtGt)
		}
		if !m.UpdatedAt.Gte.IsNull() {
			paramsUpdatedAt.Gte = param.NewOpt(mUpdatedAtGte)
		}
		if !m.UpdatedAt.Lt.IsNull() {
			paramsUpdatedAt.Lt = param.NewOpt(mUpdatedAtLt)
		}
		if !m.UpdatedAt.Lte.IsNull() {
			paramsUpdatedAt.Lte = param.NewOpt(mUpdatedAtLte)
		}
		params.UpdatedAt = paramsUpdatedAt
	}

	return
}

type CloudGPUBaremetalClustersCreatedAtDataSourceModel struct {
	Gt  timetypes.RFC3339 `tfsdk:"gt" json:"gt,optional" format:"date-time"`
	Gte timetypes.RFC3339 `tfsdk:"gte" json:"gte,optional" format:"date-time"`
	Lt  timetypes.RFC3339 `tfsdk:"lt" json:"lt,optional" format:"date-time"`
	Lte timetypes.RFC3339 `tfsdk:"lte" json:"lte,optional" format:"date-time"`
}

type CloudGPUBaremetalClustersFlavorDataSourceModel struct {
	Contains *[]types.String `tfsdk:"contains" json:"contains,optional"`
	Exact    *[]types.String `tfsdk:"exact" json:"exact,optional"`
	Prefix   *[]types.String `tfsdk:"prefix" json:"prefix,optional"`
	Suffix   *[]types.String `tfsdk:"suffix" json:"suffix,optional"`
}

type CloudGPUBaremetalClustersNameDataSourceModel struct {
	Contains *[]types.String `tfsdk:"contains" json:"contains,optional"`
	Exact    *[]types.String `tfsdk:"exact" json:"exact,optional"`
	Prefix   *[]types.String `tfsdk:"prefix" json:"prefix,optional"`
	Suffix   *[]types.String `tfsdk:"suffix" json:"suffix,optional"`
}

type CloudGPUBaremetalClustersServersCountDataSourceModel struct {
	Gt  types.Int64 `tfsdk:"gt" json:"gt,optional"`
	Gte types.Int64 `tfsdk:"gte" json:"gte,optional"`
	Lt  types.Int64 `tfsdk:"lt" json:"lt,optional"`
	Lte types.Int64 `tfsdk:"lte" json:"lte,optional"`
}

type CloudGPUBaremetalClustersTagKeyDataSourceModel struct {
	Contains *[]types.String `tfsdk:"contains" json:"contains,optional"`
	Exact    *[]types.String `tfsdk:"exact" json:"exact,optional"`
	Prefix   *[]types.String `tfsdk:"prefix" json:"prefix,optional"`
	Suffix   *[]types.String `tfsdk:"suffix" json:"suffix,optional"`
}

type CloudGPUBaremetalClustersTagValueDataSourceModel struct {
	Contains *[]types.String `tfsdk:"contains" json:"contains,optional"`
	Exact    *[]types.String `tfsdk:"exact" json:"exact,optional"`
	Prefix   *[]types.String `tfsdk:"prefix" json:"prefix,optional"`
	Suffix   *[]types.String `tfsdk:"suffix" json:"suffix,optional"`
}

type CloudGPUBaremetalClustersUpdatedAtDataSourceModel struct {
	Gt  timetypes.RFC3339 `tfsdk:"gt" json:"gt,optional" format:"date-time"`
	Gte timetypes.RFC3339 `tfsdk:"gte" json:"gte,optional" format:"date-time"`
	Lt  timetypes.RFC3339 `tfsdk:"lt" json:"lt,optional" format:"date-time"`
	Lte timetypes.RFC3339 `tfsdk:"lte" json:"lte,optional" format:"date-time"`
}

type CloudGPUBaremetalClustersItemsDataSourceModel struct {
	ID                types.String                                                                      `tfsdk:"id" json:"id,computed"`
	CreatedAt         timetypes.RFC3339                                                                 `tfsdk:"created_at" json:"created_at,computed" format:"date-time"`
	Flavor            types.String                                                                      `tfsdk:"flavor" json:"flavor,computed"`
	HasPendingChanges types.Bool                                                                        `tfsdk:"has_pending_changes" json:"has_pending_changes,computed"`
	ImageID           types.String                                                                      `tfsdk:"image_id" json:"image_id,computed"`
	ManagedBy         types.String                                                                      `tfsdk:"managed_by" json:"managed_by,computed"`
	Name              types.String                                                                      `tfsdk:"name" json:"name,computed"`
	ServersCount      types.Int64                                                                       `tfsdk:"servers_count" json:"servers_count,computed"`
	ServersIDs        customfield.List[types.String]                                                    `tfsdk:"servers_ids" json:"servers_ids,computed"`
	ServersSettings   customfield.NestedObject[CloudGPUBaremetalClustersServersSettingsDataSourceModel] `tfsdk:"servers_settings" json:"servers_settings,computed"`
	Status            types.String                                                                      `tfsdk:"status" json:"status,computed"`
	Tags              customfield.NestedObjectList[CloudGPUBaremetalClustersTagsDataSourceModel]        `tfsdk:"tags" json:"tags,computed"`
	UpdatedAt         timetypes.RFC3339                                                                 `tfsdk:"updated_at" json:"updated_at,computed" format:"date-time"`
}

type CloudGPUBaremetalClustersServersSettingsDataSourceModel struct {
	FileShares     customfield.NestedObjectList[CloudGPUBaremetalClustersServersSettingsFileSharesDataSourceModel]     `tfsdk:"file_shares" json:"file_shares,computed"`
	Interfaces     customfield.NestedObjectList[CloudGPUBaremetalClustersServersSettingsInterfacesDataSourceModel]     `tfsdk:"interfaces" json:"interfaces,computed"`
	SecurityGroups customfield.NestedObjectList[CloudGPUBaremetalClustersServersSettingsSecurityGroupsDataSourceModel] `tfsdk:"security_groups" json:"security_groups,computed"`
	SSHKeyName     types.String                                                                                        `tfsdk:"ssh_key_name" json:"ssh_key_name,computed"`
	UserData       types.String                                                                                        `tfsdk:"user_data" json:"user_data,computed"`
}

type CloudGPUBaremetalClustersServersSettingsFileSharesDataSourceModel struct {
	ID        types.String `tfsdk:"id" json:"id,computed"`
	MountPath types.String `tfsdk:"mount_path" json:"mount_path,computed"`
}

type CloudGPUBaremetalClustersServersSettingsInterfacesDataSourceModel struct {
	IPFamily       types.String                                                                                                  `tfsdk:"ip_family" json:"ip_family,computed"`
	Name           types.String                                                                                                  `tfsdk:"name" json:"name,computed"`
	SecurityGroups customfield.NestedObjectList[CloudGPUBaremetalClustersServersSettingsInterfacesSecurityGroupsDataSourceModel] `tfsdk:"security_groups" json:"security_groups,computed"`
	Type           types.String                                                                                                  `tfsdk:"type" json:"type,computed"`
	FloatingIP     customfield.NestedObject[CloudGPUBaremetalClustersServersSettingsInterfacesFloatingIPDataSourceModel]         `tfsdk:"floating_ip" json:"floating_ip,computed"`
	NetworkID      types.String                                                                                                  `tfsdk:"network_id" json:"network_id,computed"`
	SubnetID       types.String                                                                                                  `tfsdk:"subnet_id" json:"subnet_id,computed"`
	IPAddress      types.String                                                                                                  `tfsdk:"ip_address" json:"ip_address,computed"`
}

type CloudGPUBaremetalClustersServersSettingsInterfacesSecurityGroupsDataSourceModel struct {
	ID   types.String `tfsdk:"id" json:"id,computed"`
	Name types.String `tfsdk:"name" json:"name,computed"`
}

type CloudGPUBaremetalClustersServersSettingsInterfacesFloatingIPDataSourceModel struct {
	Source types.String `tfsdk:"source" json:"source,computed"`
}

type CloudGPUBaremetalClustersServersSettingsSecurityGroupsDataSourceModel struct {
	ID   types.String `tfsdk:"id" json:"id,computed"`
	Name types.String `tfsdk:"name" json:"name,computed"`
}

type CloudGPUBaremetalClustersTagsDataSourceModel struct {
	Key      types.String `tfsdk:"key" json:"key,computed"`
	ReadOnly types.Bool   `tfsdk:"read_only" json:"read_only,computed"`
	Value    types.String `tfsdk:"value" json:"value,computed"`
}
