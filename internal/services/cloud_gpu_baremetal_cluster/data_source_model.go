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

type CloudGPUBaremetalClusterDataSourceModel struct {
	ID              types.String                                                                     `tfsdk:"id" path:"cluster_id,computed"`
	ClusterID       types.String                                                                     `tfsdk:"cluster_id" path:"cluster_id,optional"`
	ProjectID       types.Int64                                                                      `tfsdk:"project_id" path:"project_id,optional"`
	RegionID        types.Int64                                                                      `tfsdk:"region_id" path:"region_id,optional"`
	CreatedAt       timetypes.RFC3339                                                                `tfsdk:"created_at" json:"created_at,computed" format:"date-time"`
	Flavor          types.String                                                                     `tfsdk:"flavor" json:"flavor,computed"`
	ImageID         types.String                                                                     `tfsdk:"image_id" json:"image_id,computed"`
	ManagedBy       types.String                                                                     `tfsdk:"managed_by" json:"managed_by,computed"`
	Name            types.String                                                                     `tfsdk:"name" json:"name,computed"`
	ServersCount    types.Int64                                                                      `tfsdk:"servers_count" json:"servers_count,computed"`
	Status          types.String                                                                     `tfsdk:"status" json:"status,computed"`
	UpdatedAt       timetypes.RFC3339                                                                `tfsdk:"updated_at" json:"updated_at,computed" format:"date-time"`
	ServersIDs      customfield.List[types.String]                                                   `tfsdk:"servers_ids" json:"servers_ids,computed"`
	ServersSettings customfield.NestedObject[CloudGPUBaremetalClusterServersSettingsDataSourceModel] `tfsdk:"servers_settings" json:"servers_settings,computed"`
	Tags            customfield.NestedObjectList[CloudGPUBaremetalClusterTagsDataSourceModel]        `tfsdk:"tags" json:"tags,computed"`
	FindOneBy       *CloudGPUBaremetalClusterFindOneByDataSourceModel                                `tfsdk:"find_one_by"`
}

func (m *CloudGPUBaremetalClusterDataSourceModel) toReadParams(_ context.Context) (params cloud.GPUBaremetalClusterGetParams, diags diag.Diagnostics) {
	params = cloud.GPUBaremetalClusterGetParams{}

	if !m.ProjectID.IsNull() {
		params.ProjectID = param.NewOpt(m.ProjectID.ValueInt64())
	}
	if !m.RegionID.IsNull() {
		params.RegionID = param.NewOpt(m.RegionID.ValueInt64())
	}

	return
}

func (m *CloudGPUBaremetalClusterDataSourceModel) toListParams(ctx context.Context) (params cloud.GPUBaremetalClusterListParams, diags diag.Diagnostics) {
	mFindOneByManagedBy := []string{}
	diags.Append(m.FindOneBy.ManagedBy.ElementsAs(ctx, &mFindOneByManagedBy, true)...)
	if diags.HasError() {
		return
	}

	params = cloud.GPUBaremetalClusterListParams{
		ManagedBy: mFindOneByManagedBy,
	}

	if !m.ProjectID.IsNull() {
		params.ProjectID = param.NewOpt(m.ProjectID.ValueInt64())
	}
	if !m.RegionID.IsNull() {
		params.RegionID = param.NewOpt(m.RegionID.ValueInt64())
	}
	if !m.FindOneBy.Limit.IsNull() {
		params.Limit = param.NewOpt(m.FindOneBy.Limit.ValueInt64())
	}

	return
}

type CloudGPUBaremetalClusterServersSettingsDataSourceModel struct {
	FileShares     customfield.NestedObjectList[CloudGPUBaremetalClusterServersSettingsFileSharesDataSourceModel]     `tfsdk:"file_shares" json:"file_shares,computed"`
	Interfaces     customfield.NestedObjectList[CloudGPUBaremetalClusterServersSettingsInterfacesDataSourceModel]     `tfsdk:"interfaces" json:"interfaces,computed"`
	SecurityGroups customfield.NestedObjectList[CloudGPUBaremetalClusterServersSettingsSecurityGroupsDataSourceModel] `tfsdk:"security_groups" json:"security_groups,computed"`
	SSHKeyName     types.String                                                                                       `tfsdk:"ssh_key_name" json:"ssh_key_name,computed"`
	UserData       types.String                                                                                       `tfsdk:"user_data" json:"user_data,computed"`
}

type CloudGPUBaremetalClusterServersSettingsFileSharesDataSourceModel struct {
	ID        types.String `tfsdk:"id" json:"id,computed"`
	MountPath types.String `tfsdk:"mount_path" json:"mount_path,computed"`
}

type CloudGPUBaremetalClusterServersSettingsInterfacesDataSourceModel struct {
	IPFamily   types.String                                                                                         `tfsdk:"ip_family" json:"ip_family,computed"`
	Name       types.String                                                                                         `tfsdk:"name" json:"name,computed"`
	Type       types.String                                                                                         `tfsdk:"type" json:"type,computed"`
	FloatingIP customfield.NestedObject[CloudGPUBaremetalClusterServersSettingsInterfacesFloatingIPDataSourceModel] `tfsdk:"floating_ip" json:"floating_ip,computed"`
	NetworkID  types.String                                                                                         `tfsdk:"network_id" json:"network_id,computed"`
	SubnetID   types.String                                                                                         `tfsdk:"subnet_id" json:"subnet_id,computed"`
	IPAddress  types.String                                                                                         `tfsdk:"ip_address" json:"ip_address,computed"`
}

type CloudGPUBaremetalClusterServersSettingsInterfacesFloatingIPDataSourceModel struct {
	Source types.String `tfsdk:"source" json:"source,computed"`
}

type CloudGPUBaremetalClusterServersSettingsSecurityGroupsDataSourceModel struct {
	ID   types.String `tfsdk:"id" json:"id,computed"`
	Name types.String `tfsdk:"name" json:"name,computed"`
}

type CloudGPUBaremetalClusterTagsDataSourceModel struct {
	Key      types.String `tfsdk:"key" json:"key,computed"`
	ReadOnly types.Bool   `tfsdk:"read_only" json:"read_only,computed"`
	Value    types.String `tfsdk:"value" json:"value,computed"`
}

type CloudGPUBaremetalClusterFindOneByDataSourceModel struct {
	Limit     types.Int64                    `tfsdk:"limit" query:"limit,computed_optional"`
	ManagedBy customfield.List[types.String] `tfsdk:"managed_by" query:"managed_by,computed_optional"`
}
