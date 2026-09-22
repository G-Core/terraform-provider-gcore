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

type CloudGPUVirtualClusterDataSourceModel struct {
	ID                types.String                                                                   `tfsdk:"id" path:"cluster_id,computed"`
	ClusterID         types.String                                                                   `tfsdk:"cluster_id" path:"cluster_id,optional"`
	ProjectID         types.Int64                                                                    `tfsdk:"project_id" path:"project_id,optional"`
	RegionID          types.Int64                                                                    `tfsdk:"region_id" path:"region_id,optional"`
	CreatedAt         timetypes.RFC3339                                                              `tfsdk:"created_at" json:"created_at,computed" format:"date-time"`
	Flavor            types.String                                                                   `tfsdk:"flavor" json:"flavor,computed"`
	HasPendingChanges types.Bool                                                                     `tfsdk:"has_pending_changes" json:"has_pending_changes,computed"`
	Name              types.String                                                                   `tfsdk:"name" json:"name,computed"`
	ServersCount      types.Int64                                                                    `tfsdk:"servers_count" json:"servers_count,computed"`
	Status            types.String                                                                   `tfsdk:"status" json:"status,computed"`
	UpdatedAt         timetypes.RFC3339                                                              `tfsdk:"updated_at" json:"updated_at,computed" format:"date-time"`
	ServersIDs        customfield.List[types.String]                                                 `tfsdk:"servers_ids" json:"servers_ids,computed"`
	ServersSettings   customfield.NestedObject[CloudGPUVirtualClusterServersSettingsDataSourceModel] `tfsdk:"servers_settings" json:"servers_settings,computed"`
	Tags              customfield.NestedObjectList[CloudGPUVirtualClusterTagsDataSourceModel]        `tfsdk:"tags" json:"tags,computed"`
	FindOneBy         *CloudGPUVirtualClusterFindOneByDataSourceModel                                `tfsdk:"find_one_by"`
}

func (m *CloudGPUVirtualClusterDataSourceModel) toReadParams(_ context.Context) (params cloud.GPUVirtualClusterGetParams, diags diag.Diagnostics) {
	params = cloud.GPUVirtualClusterGetParams{}

	if !m.ProjectID.IsNull() {
		params.ProjectID = param.NewOpt(m.ProjectID.ValueInt64())
	}
	if !m.RegionID.IsNull() {
		params.RegionID = param.NewOpt(m.RegionID.ValueInt64())
	}

	return
}

func (m *CloudGPUVirtualClusterDataSourceModel) toListParams(_ context.Context) (params cloud.GPUVirtualClusterListParams, diags diag.Diagnostics) {
	mFindOneByIDs := []string{}
	if m.FindOneBy.IDs != nil {
		for _, item := range *m.FindOneBy.IDs {
			mFindOneByIDs = append(mFindOneByIDs, item.ValueString())
		}
	}
	mFindOneByFlavorContains := []string{}
	if m.FindOneBy.Flavor.Contains != nil {
		for _, item := range *m.FindOneBy.Flavor.Contains {
			if !item.IsNull() {
				mFindOneByFlavorContains = append(mFindOneByFlavorContains, item.ValueString())
			}
		}
	}
	mFindOneByFlavorExact := []string{}
	if m.FindOneBy.Flavor.Exact != nil {
		for _, item := range *m.FindOneBy.Flavor.Exact {
			if !item.IsNull() {
				mFindOneByFlavorExact = append(mFindOneByFlavorExact, item.ValueString())
			}
		}
	}
	mFindOneByFlavorPrefix := []string{}
	if m.FindOneBy.Flavor.Prefix != nil {
		for _, item := range *m.FindOneBy.Flavor.Prefix {
			if !item.IsNull() {
				mFindOneByFlavorPrefix = append(mFindOneByFlavorPrefix, item.ValueString())
			}
		}
	}
	mFindOneByFlavorSuffix := []string{}
	if m.FindOneBy.Flavor.Suffix != nil {
		for _, item := range *m.FindOneBy.Flavor.Suffix {
			if !item.IsNull() {
				mFindOneByFlavorSuffix = append(mFindOneByFlavorSuffix, item.ValueString())
			}
		}
	}
	mFindOneByNameContains := []string{}
	if m.FindOneBy.Name.Contains != nil {
		for _, item := range *m.FindOneBy.Name.Contains {
			if !item.IsNull() {
				mFindOneByNameContains = append(mFindOneByNameContains, item.ValueString())
			}
		}
	}
	mFindOneByNameExact := []string{}
	if m.FindOneBy.Name.Exact != nil {
		for _, item := range *m.FindOneBy.Name.Exact {
			if !item.IsNull() {
				mFindOneByNameExact = append(mFindOneByNameExact, item.ValueString())
			}
		}
	}
	mFindOneByNamePrefix := []string{}
	if m.FindOneBy.Name.Prefix != nil {
		for _, item := range *m.FindOneBy.Name.Prefix {
			if !item.IsNull() {
				mFindOneByNamePrefix = append(mFindOneByNamePrefix, item.ValueString())
			}
		}
	}
	mFindOneByNameSuffix := []string{}
	if m.FindOneBy.Name.Suffix != nil {
		for _, item := range *m.FindOneBy.Name.Suffix {
			if !item.IsNull() {
				mFindOneByNameSuffix = append(mFindOneByNameSuffix, item.ValueString())
			}
		}
	}
	mFindOneByTagKeyContains := []string{}
	if m.FindOneBy.TagKey.Contains != nil {
		for _, item := range *m.FindOneBy.TagKey.Contains {
			if !item.IsNull() {
				mFindOneByTagKeyContains = append(mFindOneByTagKeyContains, item.ValueString())
			}
		}
	}
	mFindOneByTagKeyExact := []string{}
	if m.FindOneBy.TagKey.Exact != nil {
		for _, item := range *m.FindOneBy.TagKey.Exact {
			if !item.IsNull() {
				mFindOneByTagKeyExact = append(mFindOneByTagKeyExact, item.ValueString())
			}
		}
	}
	mFindOneByTagKeyPrefix := []string{}
	if m.FindOneBy.TagKey.Prefix != nil {
		for _, item := range *m.FindOneBy.TagKey.Prefix {
			if !item.IsNull() {
				mFindOneByTagKeyPrefix = append(mFindOneByTagKeyPrefix, item.ValueString())
			}
		}
	}
	mFindOneByTagKeySuffix := []string{}
	if m.FindOneBy.TagKey.Suffix != nil {
		for _, item := range *m.FindOneBy.TagKey.Suffix {
			if !item.IsNull() {
				mFindOneByTagKeySuffix = append(mFindOneByTagKeySuffix, item.ValueString())
			}
		}
	}
	mFindOneByTagValueContains := []string{}
	if m.FindOneBy.TagValue.Contains != nil {
		for _, item := range *m.FindOneBy.TagValue.Contains {
			if !item.IsNull() {
				mFindOneByTagValueContains = append(mFindOneByTagValueContains, item.ValueString())
			}
		}
	}
	mFindOneByTagValueExact := []string{}
	if m.FindOneBy.TagValue.Exact != nil {
		for _, item := range *m.FindOneBy.TagValue.Exact {
			if !item.IsNull() {
				mFindOneByTagValueExact = append(mFindOneByTagValueExact, item.ValueString())
			}
		}
	}
	mFindOneByTagValuePrefix := []string{}
	if m.FindOneBy.TagValue.Prefix != nil {
		for _, item := range *m.FindOneBy.TagValue.Prefix {
			if !item.IsNull() {
				mFindOneByTagValuePrefix = append(mFindOneByTagValuePrefix, item.ValueString())
			}
		}
	}
	mFindOneByTagValueSuffix := []string{}
	if m.FindOneBy.TagValue.Suffix != nil {
		for _, item := range *m.FindOneBy.TagValue.Suffix {
			if !item.IsNull() {
				mFindOneByTagValueSuffix = append(mFindOneByTagValueSuffix, item.ValueString())
			}
		}
	}
	mFindOneByTags := map[string][]string{}
	for key, value := range *m.FindOneBy.Tags {
		paramsValue := []string{}
		if value != nil {
			for _, item := range *value {
				if !item.IsNull() {
					paramsValue = append(paramsValue, item.ValueString())
				}
			}
		}
		mFindOneByTags[key] = paramsValue
	}

	params = cloud.GPUVirtualClusterListParams{
		IDs: mFindOneByIDs,
	}

	if !m.ProjectID.IsNull() {
		params.ProjectID = param.NewOpt(m.ProjectID.ValueInt64())
	}
	if !m.RegionID.IsNull() {
		params.RegionID = param.NewOpt(m.RegionID.ValueInt64())
	}
	if m.FindOneBy.CreatedAt != nil {
		paramsCreatedAt := cloud.GPUVirtualClusterListParamsCreatedAt{}
		if !m.FindOneBy.CreatedAt.Gt.IsNull() {
			mFindOneByCreatedAtGt, errs := m.FindOneBy.CreatedAt.Gt.ValueRFC3339Time()
			diags.Append(errs...)
			paramsCreatedAt.Gt = param.NewOpt(mFindOneByCreatedAtGt)
		}
		if !m.FindOneBy.CreatedAt.Gte.IsNull() {
			mFindOneByCreatedAtGte, errs := m.FindOneBy.CreatedAt.Gte.ValueRFC3339Time()
			diags.Append(errs...)
			paramsCreatedAt.Gte = param.NewOpt(mFindOneByCreatedAtGte)
		}
		if !m.FindOneBy.CreatedAt.Lt.IsNull() {
			mFindOneByCreatedAtLt, errs := m.FindOneBy.CreatedAt.Lt.ValueRFC3339Time()
			diags.Append(errs...)
			paramsCreatedAt.Lt = param.NewOpt(mFindOneByCreatedAtLt)
		}
		if !m.FindOneBy.CreatedAt.Lte.IsNull() {
			mFindOneByCreatedAtLte, errs := m.FindOneBy.CreatedAt.Lte.ValueRFC3339Time()
			diags.Append(errs...)
			paramsCreatedAt.Lte = param.NewOpt(mFindOneByCreatedAtLte)
		}
		params.CreatedAt = paramsCreatedAt
	}
	if m.FindOneBy.Flavor != nil {
		params.Flavor = cloud.GPUVirtualClusterListParamsFlavor{
			Contains: mFindOneByFlavorContains,
			Exact:    mFindOneByFlavorExact,
			Prefix:   mFindOneByFlavorPrefix,
			Suffix:   mFindOneByFlavorSuffix,
		}
	}
	if m.FindOneBy.Name != nil {
		params.Name = cloud.GPUVirtualClusterListParamsName{
			Contains: mFindOneByNameContains,
			Exact:    mFindOneByNameExact,
			Prefix:   mFindOneByNamePrefix,
			Suffix:   mFindOneByNameSuffix,
		}
	}
	if m.FindOneBy.ServersCount != nil {
		paramsServersCount := cloud.GPUVirtualClusterListParamsServersCount{}
		if !m.FindOneBy.ServersCount.Gt.IsNull() {
			paramsServersCount.Gt = param.NewOpt(m.FindOneBy.ServersCount.Gt.ValueInt64())
		}
		if !m.FindOneBy.ServersCount.Gte.IsNull() {
			paramsServersCount.Gte = param.NewOpt(m.FindOneBy.ServersCount.Gte.ValueInt64())
		}
		if !m.FindOneBy.ServersCount.Lt.IsNull() {
			paramsServersCount.Lt = param.NewOpt(m.FindOneBy.ServersCount.Lt.ValueInt64())
		}
		if !m.FindOneBy.ServersCount.Lte.IsNull() {
			paramsServersCount.Lte = param.NewOpt(m.FindOneBy.ServersCount.Lte.ValueInt64())
		}
		params.ServersCount = paramsServersCount
	}
	if m.FindOneBy.TagKey != nil {
		params.TagKey = cloud.GPUVirtualClusterListParamsTagKey{
			Contains: mFindOneByTagKeyContains,
			Exact:    mFindOneByTagKeyExact,
			Prefix:   mFindOneByTagKeyPrefix,
			Suffix:   mFindOneByTagKeySuffix,
		}
	}
	if m.FindOneBy.TagValue != nil {
		params.TagValue = cloud.GPUVirtualClusterListParamsTagValue{
			Contains: mFindOneByTagValueContains,
			Exact:    mFindOneByTagValueExact,
			Prefix:   mFindOneByTagValuePrefix,
			Suffix:   mFindOneByTagValueSuffix,
		}
	}
	if m.FindOneBy.UpdatedAt != nil {
		paramsUpdatedAt := cloud.GPUVirtualClusterListParamsUpdatedAt{}
		if !m.FindOneBy.UpdatedAt.Gt.IsNull() {
			mFindOneByUpdatedAtGt, errs := m.FindOneBy.UpdatedAt.Gt.ValueRFC3339Time()
			diags.Append(errs...)
			paramsUpdatedAt.Gt = param.NewOpt(mFindOneByUpdatedAtGt)
		}
		if !m.FindOneBy.UpdatedAt.Gte.IsNull() {
			mFindOneByUpdatedAtGte, errs := m.FindOneBy.UpdatedAt.Gte.ValueRFC3339Time()
			diags.Append(errs...)
			paramsUpdatedAt.Gte = param.NewOpt(mFindOneByUpdatedAtGte)
		}
		if !m.FindOneBy.UpdatedAt.Lt.IsNull() {
			mFindOneByUpdatedAtLt, errs := m.FindOneBy.UpdatedAt.Lt.ValueRFC3339Time()
			diags.Append(errs...)
			paramsUpdatedAt.Lt = param.NewOpt(mFindOneByUpdatedAtLt)
		}
		if !m.FindOneBy.UpdatedAt.Lte.IsNull() {
			mFindOneByUpdatedAtLte, errs := m.FindOneBy.UpdatedAt.Lte.ValueRFC3339Time()
			diags.Append(errs...)
			paramsUpdatedAt.Lte = param.NewOpt(mFindOneByUpdatedAtLte)
		}
		params.UpdatedAt = paramsUpdatedAt
	}

	return
}

type CloudGPUVirtualClusterServersSettingsDataSourceModel struct {
	FileShares     customfield.NestedObjectList[CloudGPUVirtualClusterServersSettingsFileSharesDataSourceModel]     `tfsdk:"file_shares" json:"file_shares,computed"`
	Interfaces     customfield.NestedObjectList[CloudGPUVirtualClusterServersSettingsInterfacesDataSourceModel]     `tfsdk:"interfaces" json:"interfaces,computed"`
	SecurityGroups customfield.NestedObjectList[CloudGPUVirtualClusterServersSettingsSecurityGroupsDataSourceModel] `tfsdk:"security_groups" json:"security_groups,computed"`
	SSHKeyName     types.String                                                                                     `tfsdk:"ssh_key_name" json:"ssh_key_name,computed"`
	UserData       types.String                                                                                     `tfsdk:"user_data" json:"user_data,computed"`
	Volumes        customfield.NestedObjectList[CloudGPUVirtualClusterServersSettingsVolumesDataSourceModel]        `tfsdk:"volumes" json:"volumes,computed"`
}

type CloudGPUVirtualClusterServersSettingsFileSharesDataSourceModel struct {
	ID        types.String `tfsdk:"id" json:"id,computed"`
	MountPath types.String `tfsdk:"mount_path" json:"mount_path,computed"`
}

type CloudGPUVirtualClusterServersSettingsInterfacesDataSourceModel struct {
	IPFamily       types.String                                                                                               `tfsdk:"ip_family" json:"ip_family,computed"`
	Name           types.String                                                                                               `tfsdk:"name" json:"name,computed"`
	SecurityGroups customfield.NestedObjectList[CloudGPUVirtualClusterServersSettingsInterfacesSecurityGroupsDataSourceModel] `tfsdk:"security_groups" json:"security_groups,computed"`
	Type           types.String                                                                                               `tfsdk:"type" json:"type,computed"`
	FloatingIP     customfield.NestedObject[CloudGPUVirtualClusterServersSettingsInterfacesFloatingIPDataSourceModel]         `tfsdk:"floating_ip" json:"floating_ip,computed"`
	NetworkID      types.String                                                                                               `tfsdk:"network_id" json:"network_id,computed"`
	SubnetID       types.String                                                                                               `tfsdk:"subnet_id" json:"subnet_id,computed"`
	IPAddress      types.String                                                                                               `tfsdk:"ip_address" json:"ip_address,computed"`
}

type CloudGPUVirtualClusterServersSettingsInterfacesSecurityGroupsDataSourceModel struct {
	ID   types.String `tfsdk:"id" json:"id,computed"`
	Name types.String `tfsdk:"name" json:"name,computed"`
}

type CloudGPUVirtualClusterServersSettingsInterfacesFloatingIPDataSourceModel struct {
	Source types.String `tfsdk:"source" json:"source,computed"`
}

type CloudGPUVirtualClusterServersSettingsSecurityGroupsDataSourceModel struct {
	ID   types.String `tfsdk:"id" json:"id,computed"`
	Name types.String `tfsdk:"name" json:"name,computed"`
}

type CloudGPUVirtualClusterServersSettingsVolumesDataSourceModel struct {
	BootIndex           types.Int64                                                                                   `tfsdk:"boot_index" json:"boot_index,computed"`
	DeleteOnTermination types.Bool                                                                                    `tfsdk:"delete_on_termination" json:"delete_on_termination,computed"`
	ImageID             types.String                                                                                  `tfsdk:"image_id" json:"image_id,computed"`
	Name                types.String                                                                                  `tfsdk:"name" json:"name,computed"`
	Size                types.Int64                                                                                   `tfsdk:"size" json:"size,computed"`
	Tags                customfield.NestedObjectList[CloudGPUVirtualClusterServersSettingsVolumesTagsDataSourceModel] `tfsdk:"tags" json:"tags,computed"`
	Type                types.String                                                                                  `tfsdk:"type" json:"type,computed"`
}

type CloudGPUVirtualClusterServersSettingsVolumesTagsDataSourceModel struct {
	Key      types.String `tfsdk:"key" json:"key,computed"`
	ReadOnly types.Bool   `tfsdk:"read_only" json:"read_only,computed"`
	Value    types.String `tfsdk:"value" json:"value,computed"`
}

type CloudGPUVirtualClusterTagsDataSourceModel struct {
	Key      types.String `tfsdk:"key" json:"key,computed"`
	ReadOnly types.Bool   `tfsdk:"read_only" json:"read_only,computed"`
	Value    types.String `tfsdk:"value" json:"value,computed"`
}

type CloudGPUVirtualClusterFindOneByDataSourceModel struct {
	CreatedAt    *CloudGPUVirtualClustersCreatedAtDataSourceModel    `tfsdk:"created_at" query:"created_at,optional"`
	Flavor       *CloudGPUVirtualClustersFlavorDataSourceModel       `tfsdk:"flavor" query:"flavor,optional"`
	IDs          *[]types.String                                     `tfsdk:"ids" query:"ids,optional"`
	Name         *CloudGPUVirtualClustersNameDataSourceModel         `tfsdk:"name" query:"name,optional"`
	ServersCount *CloudGPUVirtualClustersServersCountDataSourceModel `tfsdk:"servers_count" query:"servers_count,optional"`
	TagKey       *CloudGPUVirtualClustersTagKeyDataSourceModel       `tfsdk:"tag_key" query:"tag_key,optional"`
	TagValue     *CloudGPUVirtualClustersTagValueDataSourceModel     `tfsdk:"tag_value" query:"tag_value,optional"`
	Tags         *map[string]*[]types.String                         `tfsdk:"tags" query:"tags,optional"`
	UpdatedAt    *CloudGPUVirtualClustersUpdatedAtDataSourceModel    `tfsdk:"updated_at" query:"updated_at,optional"`
}
