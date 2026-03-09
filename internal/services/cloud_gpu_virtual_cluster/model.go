// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_gpu_virtual_cluster

import (
	"github.com/G-Core/terraform-provider-gcore/internal/apijson"
	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type CloudGPUVirtualClusterModel struct {
	ID              types.String                                `tfsdk:"id" json:"id,computed"`
	ProjectID       types.Int64                                 `tfsdk:"project_id" path:"project_id,optional"`
	RegionID        types.Int64                                 `tfsdk:"region_id" path:"region_id,optional"`
	Flavor          types.String                                `tfsdk:"flavor" json:"flavor,required"`
	ServersCount    types.Int64                                 `tfsdk:"servers_count" json:"servers_count,required"`
	ServersSettings *CloudGPUVirtualClusterServersSettingsModel `tfsdk:"servers_settings" json:"servers_settings,required"`
	Tags            customfield.Map[types.String]               `tfsdk:"tags" json:"tags,computed_optional,no_refresh"`
	Name            types.String                                `tfsdk:"name" json:"name,required"`
	CreatedAt       timetypes.RFC3339                           `tfsdk:"created_at" json:"created_at,computed" format:"date-time"`
	Status          types.String                                `tfsdk:"status" json:"status,computed"`
	UpdatedAt       timetypes.RFC3339                           `tfsdk:"updated_at" json:"updated_at,computed" format:"date-time"`
	ServersIDs      customfield.List[types.String]              `tfsdk:"servers_ids" json:"servers_ids,computed"`
}

func (m CloudGPUVirtualClusterModel) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(m)
}

func (m CloudGPUVirtualClusterModel) MarshalJSONForUpdate(state CloudGPUVirtualClusterModel) (data []byte, err error) {
	return apijson.MarshalForPatch(m, state)
}

type CloudGPUVirtualClusterServersSettingsModel struct {
	Interfaces     *[]*CloudGPUVirtualClusterServersSettingsInterfacesModel                               `tfsdk:"interfaces" json:"interfaces,required"`
	Volumes        *[]*CloudGPUVirtualClusterServersSettingsVolumesModel                                  `tfsdk:"volumes" json:"volumes,required"`
	Credentials    *CloudGPUVirtualClusterServersSettingsCredentialsModel                                 `tfsdk:"credentials" json:"credentials,optional,no_refresh"`
	FileShares     customfield.NestedObjectList[CloudGPUVirtualClusterServersSettingsFileSharesModel]     `tfsdk:"file_shares" json:"file_shares,computed_optional"`
	SecurityGroups customfield.NestedObjectList[CloudGPUVirtualClusterServersSettingsSecurityGroupsModel] `tfsdk:"security_groups" json:"security_groups,computed_optional"`
	UserData       types.String                                                                           `tfsdk:"user_data" json:"user_data,computed_optional,no_refresh"`
}

type CloudGPUVirtualClusterServersSettingsInterfacesModel struct {
	Type       types.String                                                    `tfsdk:"type" json:"type,required"`
	IPFamily   types.String                                                    `tfsdk:"ip_family" json:"ip_family,computed_optional"`
	Name       types.String                                                    `tfsdk:"name" json:"name,computed_optional"`
	NetworkID  types.String                                                    `tfsdk:"network_id" json:"network_id,optional"`
	SubnetID   types.String                                                    `tfsdk:"subnet_id" json:"subnet_id,optional"`
	FloatingIP *CloudGPUVirtualClusterServersSettingsInterfacesFloatingIPModel `tfsdk:"floating_ip" json:"floating_ip,optional"`
}

type CloudGPUVirtualClusterServersSettingsInterfacesFloatingIPModel struct {
	Source types.String `tfsdk:"source" json:"source,required"`
}

type CloudGPUVirtualClusterServersSettingsVolumesModel struct {
	BootIndex           types.Int64              `tfsdk:"boot_index" json:"boot_index,required"`
	Name                types.String             `tfsdk:"name" json:"name,required"`
	Size                types.Int64              `tfsdk:"size" json:"size,required"`
	Source              types.String             `tfsdk:"source" json:"source,required,no_refresh"`
	Type                types.String             `tfsdk:"type" json:"type,required"`
	DeleteOnTermination types.Bool               `tfsdk:"delete_on_termination" json:"delete_on_termination,computed_optional"`
	Tags                *map[string]types.String `tfsdk:"tags" json:"tags,optional,no_refresh"`
	ImageID             types.String             `tfsdk:"image_id" json:"image_id,optional"`
}

type CloudGPUVirtualClusterServersSettingsCredentialsModel struct {
	Password          types.String `tfsdk:"password_wo" json:"password,optional,no_refresh"`
	SSHKeyName        types.String `tfsdk:"ssh_key_name" json:"ssh_key_name,optional"`
	Username          types.String `tfsdk:"username" json:"username,optional"`
	PasswordWoVersion types.Int64  `tfsdk:"password_wo_version"`
}

type CloudGPUVirtualClusterServersSettingsFileSharesModel struct {
	ID        types.String `tfsdk:"id" json:"id,required"`
	MountPath types.String `tfsdk:"mount_path" json:"mount_path,required"`
}

type CloudGPUVirtualClusterServersSettingsSecurityGroupsModel struct {
	ID types.String `tfsdk:"id" json:"id,required"`
}
