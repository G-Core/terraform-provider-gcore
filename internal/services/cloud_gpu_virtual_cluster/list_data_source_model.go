// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_gpu_virtual_cluster

import (
	"context"

	"github.com/G-Core/gcore-go/cloud"
	"github.com/G-Core/gcore-go/packages/param"
	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type CloudGPUVirtualClustersResultsListDataSourceEnvelope struct {
	Results customfield.NestedObjectList[CloudGPUVirtualClustersItemsDataSourceModel] `json:"results,computed"`
}

type CloudGPUVirtualClustersDataSourceModel struct {
	ProjectID    types.Int64                                                               `tfsdk:"project_id" path:"project_id,optional"`
	RegionID     types.Int64                                                               `tfsdk:"region_id" path:"region_id,optional"`
	IDs          *[]types.String                                                           `tfsdk:"ids" query:"ids,optional"`
	Tags         *map[string]types.String                                                  `tfsdk:"tags" query:"tags,optional"`
	CreatedAt    *CloudGPUVirtualClustersCreatedAtDataSourceModel                          `tfsdk:"created_at" query:"created_at,optional"`
	Flavor       *CloudGPUVirtualClustersFlavorDataSourceModel                             `tfsdk:"flavor" query:"flavor,optional"`
	Name         *CloudGPUVirtualClustersNameDataSourceModel                               `tfsdk:"name" query:"name,optional"`
	ServersCount *CloudGPUVirtualClustersServersCountDataSourceModel                       `tfsdk:"servers_count" query:"servers_count,optional"`
	TagKey       *CloudGPUVirtualClustersTagKeyDataSourceModel                             `tfsdk:"tag_key" query:"tag_key,optional"`
	TagValue     *CloudGPUVirtualClustersTagValueDataSourceModel                           `tfsdk:"tag_value" query:"tag_value,optional"`
	UpdatedAt    *CloudGPUVirtualClustersUpdatedAtDataSourceModel                          `tfsdk:"updated_at" query:"updated_at,optional"`
	MaxItems     types.Int64                                                               `tfsdk:"max_items"`
	Items        customfield.NestedObjectList[CloudGPUVirtualClustersItemsDataSourceModel] `tfsdk:"items"`
}

func (m *CloudGPUVirtualClustersDataSourceModel) toListParams(_ context.Context) (params cloud.GPUVirtualClusterListParams, diags diag.Diagnostics) {
	mIDs := []string{}
	if m.IDs != nil {
		for _, item := range *m.IDs {
			mIDs = append(mIDs, item.ValueString())
		}
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

	params = cloud.GPUVirtualClusterListParams{
		IDs: mIDs,
	}

	if !m.ProjectID.IsNull() {
		params.ProjectID = param.NewOpt(m.ProjectID.ValueInt64())
	}
	if !m.RegionID.IsNull() {
		params.RegionID = param.NewOpt(m.RegionID.ValueInt64())
	}
	if m.CreatedAt != nil {
		paramsCreatedAt := cloud.GPUVirtualClusterListParamsCreatedAt{}
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
		params.Flavor = cloud.GPUVirtualClusterListParamsFlavor{
			Contains: mFlavorContains,
			Exact:    mFlavorExact,
			Prefix:   mFlavorPrefix,
			Suffix:   mFlavorSuffix,
		}
	}
	if m.Name != nil {
		params.Name = cloud.GPUVirtualClusterListParamsName{
			Contains: mNameContains,
			Exact:    mNameExact,
			Prefix:   mNamePrefix,
			Suffix:   mNameSuffix,
		}
	}
	if m.ServersCount != nil {
		paramsServersCount := cloud.GPUVirtualClusterListParamsServersCount{}
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
		params.TagKey = cloud.GPUVirtualClusterListParamsTagKey{
			Contains: mTagKeyContains,
			Exact:    mTagKeyExact,
			Prefix:   mTagKeyPrefix,
			Suffix:   mTagKeySuffix,
		}
	}
	if m.TagValue != nil {
		params.TagValue = cloud.GPUVirtualClusterListParamsTagValue{
			Contains: mTagValueContains,
			Exact:    mTagValueExact,
			Prefix:   mTagValuePrefix,
			Suffix:   mTagValueSuffix,
		}
	}
	if m.UpdatedAt != nil {
		paramsUpdatedAt := cloud.GPUVirtualClusterListParamsUpdatedAt{}
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

type CloudGPUVirtualClustersCreatedAtDataSourceModel struct {
	Gt  timetypes.RFC3339 `tfsdk:"gt" json:"gt,optional" format:"date-time"`
	Gte timetypes.RFC3339 `tfsdk:"gte" json:"gte,optional" format:"date-time"`
	Lt  timetypes.RFC3339 `tfsdk:"lt" json:"lt,optional" format:"date-time"`
	Lte timetypes.RFC3339 `tfsdk:"lte" json:"lte,optional" format:"date-time"`
}

type CloudGPUVirtualClustersFlavorDataSourceModel struct {
	Contains *[]types.String `tfsdk:"contains" json:"contains,optional"`
	Exact    *[]types.String `tfsdk:"exact" json:"exact,optional"`
	Prefix   *[]types.String `tfsdk:"prefix" json:"prefix,optional"`
	Suffix   *[]types.String `tfsdk:"suffix" json:"suffix,optional"`
}

type CloudGPUVirtualClustersNameDataSourceModel struct {
	Contains *[]types.String `tfsdk:"contains" json:"contains,optional"`
	Exact    *[]types.String `tfsdk:"exact" json:"exact,optional"`
	Prefix   *[]types.String `tfsdk:"prefix" json:"prefix,optional"`
	Suffix   *[]types.String `tfsdk:"suffix" json:"suffix,optional"`
}

type CloudGPUVirtualClustersServersCountDataSourceModel struct {
	Gt  types.Int64 `tfsdk:"gt" json:"gt,optional"`
	Gte types.Int64 `tfsdk:"gte" json:"gte,optional"`
	Lt  types.Int64 `tfsdk:"lt" json:"lt,optional"`
	Lte types.Int64 `tfsdk:"lte" json:"lte,optional"`
}

type CloudGPUVirtualClustersTagKeyDataSourceModel struct {
	Contains *[]types.String `tfsdk:"contains" json:"contains,optional"`
	Exact    *[]types.String `tfsdk:"exact" json:"exact,optional"`
	Prefix   *[]types.String `tfsdk:"prefix" json:"prefix,optional"`
	Suffix   *[]types.String `tfsdk:"suffix" json:"suffix,optional"`
}

type CloudGPUVirtualClustersTagValueDataSourceModel struct {
	Contains *[]types.String `tfsdk:"contains" json:"contains,optional"`
	Exact    *[]types.String `tfsdk:"exact" json:"exact,optional"`
	Prefix   *[]types.String `tfsdk:"prefix" json:"prefix,optional"`
	Suffix   *[]types.String `tfsdk:"suffix" json:"suffix,optional"`
}

type CloudGPUVirtualClustersUpdatedAtDataSourceModel struct {
	Gt  timetypes.RFC3339 `tfsdk:"gt" json:"gt,optional" format:"date-time"`
	Gte timetypes.RFC3339 `tfsdk:"gte" json:"gte,optional" format:"date-time"`
	Lt  timetypes.RFC3339 `tfsdk:"lt" json:"lt,optional" format:"date-time"`
	Lte timetypes.RFC3339 `tfsdk:"lte" json:"lte,optional" format:"date-time"`
}

type CloudGPUVirtualClustersItemsDataSourceModel struct {
	ID                types.String                                                                    `tfsdk:"id" json:"id,computed"`
	CreatedAt         timetypes.RFC3339                                                               `tfsdk:"created_at" json:"created_at,computed" format:"date-time"`
	Flavor            types.String                                                                    `tfsdk:"flavor" json:"flavor,computed"`
	HasPendingChanges types.Bool                                                                      `tfsdk:"has_pending_changes" json:"has_pending_changes,computed"`
	Name              types.String                                                                    `tfsdk:"name" json:"name,computed"`
	ServersCount      types.Int64                                                                     `tfsdk:"servers_count" json:"servers_count,computed"`
	ServersIDs        customfield.List[types.String]                                                  `tfsdk:"servers_ids" json:"servers_ids,computed"`
	ServersSettings   customfield.NestedObject[CloudGPUVirtualClustersServersSettingsDataSourceModel] `tfsdk:"servers_settings" json:"servers_settings,computed"`
	Status            types.String                                                                    `tfsdk:"status" json:"status,computed"`
	Tags              customfield.NestedObjectList[CloudGPUVirtualClustersTagsDataSourceModel]        `tfsdk:"tags" json:"tags,computed"`
	UpdatedAt         timetypes.RFC3339                                                               `tfsdk:"updated_at" json:"updated_at,computed" format:"date-time"`
}

type CloudGPUVirtualClustersServersSettingsDataSourceModel struct {
	FileShares     customfield.NestedObjectList[CloudGPUVirtualClustersServersSettingsFileSharesDataSourceModel]     `tfsdk:"file_shares" json:"file_shares,computed"`
	Interfaces     customfield.NestedObjectList[CloudGPUVirtualClustersServersSettingsInterfacesDataSourceModel]     `tfsdk:"interfaces" json:"interfaces,computed"`
	SecurityGroups customfield.NestedObjectList[CloudGPUVirtualClustersServersSettingsSecurityGroupsDataSourceModel] `tfsdk:"security_groups" json:"security_groups,computed"`
	SSHKeyName     types.String                                                                                      `tfsdk:"ssh_key_name" json:"ssh_key_name,computed"`
	UserData       types.String                                                                                      `tfsdk:"user_data" json:"user_data,computed"`
	Volumes        customfield.NestedObjectList[CloudGPUVirtualClustersServersSettingsVolumesDataSourceModel]        `tfsdk:"volumes" json:"volumes,computed"`
}

type CloudGPUVirtualClustersServersSettingsFileSharesDataSourceModel struct {
	ID        types.String `tfsdk:"id" json:"id,computed"`
	MountPath types.String `tfsdk:"mount_path" json:"mount_path,computed"`
}

type CloudGPUVirtualClustersServersSettingsInterfacesDataSourceModel struct {
	IPFamily       types.String                                                                                                `tfsdk:"ip_family" json:"ip_family,computed"`
	Name           types.String                                                                                                `tfsdk:"name" json:"name,computed"`
	SecurityGroups customfield.NestedObjectList[CloudGPUVirtualClustersServersSettingsInterfacesSecurityGroupsDataSourceModel] `tfsdk:"security_groups" json:"security_groups,computed"`
	Type           types.String                                                                                                `tfsdk:"type" json:"type,computed"`
	FloatingIP     customfield.NestedObject[CloudGPUVirtualClustersServersSettingsInterfacesFloatingIPDataSourceModel]         `tfsdk:"floating_ip" json:"floating_ip,computed"`
	NetworkID      types.String                                                                                                `tfsdk:"network_id" json:"network_id,computed"`
	SubnetID       types.String                                                                                                `tfsdk:"subnet_id" json:"subnet_id,computed"`
	IPAddress      types.String                                                                                                `tfsdk:"ip_address" json:"ip_address,computed"`
}

type CloudGPUVirtualClustersServersSettingsInterfacesSecurityGroupsDataSourceModel struct {
	ID   types.String `tfsdk:"id" json:"id,computed"`
	Name types.String `tfsdk:"name" json:"name,computed"`
}

type CloudGPUVirtualClustersServersSettingsInterfacesFloatingIPDataSourceModel struct {
	Source types.String `tfsdk:"source" json:"source,computed"`
}

type CloudGPUVirtualClustersServersSettingsSecurityGroupsDataSourceModel struct {
	ID   types.String `tfsdk:"id" json:"id,computed"`
	Name types.String `tfsdk:"name" json:"name,computed"`
}

type CloudGPUVirtualClustersServersSettingsVolumesDataSourceModel struct {
	BootIndex           types.Int64                                                                                    `tfsdk:"boot_index" json:"boot_index,computed"`
	DeleteOnTermination types.Bool                                                                                     `tfsdk:"delete_on_termination" json:"delete_on_termination,computed"`
	ImageID             types.String                                                                                   `tfsdk:"image_id" json:"image_id,computed"`
	Name                types.String                                                                                   `tfsdk:"name" json:"name,computed"`
	Size                types.Int64                                                                                    `tfsdk:"size" json:"size,computed"`
	Tags                customfield.NestedObjectList[CloudGPUVirtualClustersServersSettingsVolumesTagsDataSourceModel] `tfsdk:"tags" json:"tags,computed"`
	Type                types.String                                                                                   `tfsdk:"type" json:"type,computed"`
}

type CloudGPUVirtualClustersServersSettingsVolumesTagsDataSourceModel struct {
	Key      types.String `tfsdk:"key" json:"key,computed"`
	ReadOnly types.Bool   `tfsdk:"read_only" json:"read_only,computed"`
	Value    types.String `tfsdk:"value" json:"value,computed"`
}

type CloudGPUVirtualClustersTagsDataSourceModel struct {
	Key      types.String `tfsdk:"key" json:"key,computed"`
	ReadOnly types.Bool   `tfsdk:"read_only" json:"read_only,computed"`
	Value    types.String `tfsdk:"value" json:"value,computed"`
}
