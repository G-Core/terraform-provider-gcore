// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_gpu_baremetal_cluster

import (
	"github.com/G-Core/terraform-provider-gcore/internal/apijson"
	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type CloudGPUBaremetalClusterModel struct {
	ID                types.String                                  `tfsdk:"id" json:"id,computed"`
	ProjectID         types.Int64                                   `tfsdk:"project_id" path:"project_id,optional"`
	RegionID          types.Int64                                   `tfsdk:"region_id" path:"region_id,optional"`
	Flavor            types.String                                  `tfsdk:"flavor" json:"flavor,required"`
	ImageID           types.String                                  `tfsdk:"image_id" json:"image_id,required"`
	ServersCount      types.Int64                                   `tfsdk:"servers_count" json:"servers_count,required"`
	ServersSettings   *CloudGPUBaremetalClusterServersSettingsModel `tfsdk:"servers_settings" json:"servers_settings,required"`
	Name              types.String                                  `tfsdk:"name" json:"name,required"`
	Tags              customfield.Map[types.String]                 `tfsdk:"tags" json:"tags,computed_optional,no_refresh"`
	CreatedAt         timetypes.RFC3339                             `tfsdk:"created_at" json:"created_at,computed" format:"date-time"`
	HasPendingChanges types.Bool                                    `tfsdk:"has_pending_changes" json:"has_pending_changes,computed"`
	ManagedBy         types.String                                  `tfsdk:"managed_by" json:"managed_by,computed"`
	Status            types.String                                  `tfsdk:"status" json:"status,computed"`
	UpdatedAt         timetypes.RFC3339                             `tfsdk:"updated_at" json:"updated_at,computed" format:"date-time"`
	ServersIDs        customfield.List[types.String]                `tfsdk:"servers_ids" json:"servers_ids,computed"`
	Tasks             customfield.List[types.String]                `tfsdk:"tasks" json:"tasks,computed,no_refresh"`
}

func (m CloudGPUBaremetalClusterModel) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(m)
}

func (m CloudGPUBaremetalClusterModel) MarshalJSONForUpdate(state CloudGPUBaremetalClusterModel) (data []byte, err error) {
	return apijson.MarshalForPatch(m, state)
}

type CloudGPUBaremetalClusterServersSettingsModel struct {
	Interfaces     *[]*CloudGPUBaremetalClusterServersSettingsInterfacesModel                               `tfsdk:"interfaces" json:"interfaces,required"`
	Credentials    *CloudGPUBaremetalClusterServersSettingsCredentialsModel                                 `tfsdk:"credentials" json:"credentials,optional,no_refresh"`
	FileShares     customfield.NestedObjectList[CloudGPUBaremetalClusterServersSettingsFileSharesModel]     `tfsdk:"file_shares" json:"file_shares,computed_optional"`
	SecurityGroups customfield.NestedObjectList[CloudGPUBaremetalClusterServersSettingsSecurityGroupsModel] `tfsdk:"security_groups" json:"security_groups,computed_optional"`
	UserData       types.String                                                                             `tfsdk:"user_data" json:"user_data,computed_optional,no_refresh"`
}

type CloudGPUBaremetalClusterServersSettingsInterfacesModel struct {
	Type       types.String                                                      `tfsdk:"type" json:"type,required"`
	IPFamily   types.String                                                      `tfsdk:"ip_family" json:"ip_family,computed_optional"`
	Name       types.String                                                      `tfsdk:"name" json:"name,computed_optional"`
	NetworkID  types.String                                                      `tfsdk:"network_id" json:"network_id,optional"`
	SubnetID   types.String                                                      `tfsdk:"subnet_id" json:"subnet_id,optional"`
	FloatingIP *CloudGPUBaremetalClusterServersSettingsInterfacesFloatingIPModel `tfsdk:"floating_ip" json:"floating_ip,optional"`
}

type CloudGPUBaremetalClusterServersSettingsInterfacesFloatingIPModel struct {
	Source types.String `tfsdk:"source" json:"source,required"`
}

type CloudGPUBaremetalClusterServersSettingsCredentialsModel struct {
	Password          types.String `tfsdk:"password_wo" json:"password,optional,no_refresh"`
	SSHKeyName        types.String `tfsdk:"ssh_key_name" json:"ssh_key_name,optional"`
	Username          types.String `tfsdk:"username" json:"username,optional"`
	PasswordWoVersion types.Int64  `tfsdk:"password_wo_version"`
}

type CloudGPUBaremetalClusterServersSettingsFileSharesModel struct {
	ID        types.String `tfsdk:"id" json:"id,required"`
	MountPath types.String `tfsdk:"mount_path" json:"mount_path,required"`
}

type CloudGPUBaremetalClusterServersSettingsSecurityGroupsModel struct {
	ID types.String `tfsdk:"id" json:"id,required"`
}
