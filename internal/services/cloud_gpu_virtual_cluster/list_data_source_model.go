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
	ProjectID types.Int64                                                               `tfsdk:"project_id" path:"project_id,optional"`
	RegionID  types.Int64                                                               `tfsdk:"region_id" path:"region_id,optional"`
	MaxItems  types.Int64                                                               `tfsdk:"max_items"`
	Items     customfield.NestedObjectList[CloudGPUVirtualClustersItemsDataSourceModel] `tfsdk:"items"`
}

func (m *CloudGPUVirtualClustersDataSourceModel) toListParams(_ context.Context) (params cloud.GPUVirtualClusterListParams, diags diag.Diagnostics) {
	params = cloud.GPUVirtualClusterListParams{}

	if !m.ProjectID.IsNull() {
		params.ProjectID = param.NewOpt(m.ProjectID.ValueInt64())
	}
	if !m.RegionID.IsNull() {
		params.RegionID = param.NewOpt(m.RegionID.ValueInt64())
	}

	return
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
	IPFamily   types.String                                                                                        `tfsdk:"ip_family" json:"ip_family,computed"`
	Name       types.String                                                                                        `tfsdk:"name" json:"name,computed"`
	Type       types.String                                                                                        `tfsdk:"type" json:"type,computed"`
	FloatingIP customfield.NestedObject[CloudGPUVirtualClustersServersSettingsInterfacesFloatingIPDataSourceModel] `tfsdk:"floating_ip" json:"floating_ip,computed"`
	NetworkID  types.String                                                                                        `tfsdk:"network_id" json:"network_id,computed"`
	SubnetID   types.String                                                                                        `tfsdk:"subnet_id" json:"subnet_id,computed"`
	IPAddress  types.String                                                                                        `tfsdk:"ip_address" json:"ip_address,computed"`
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
