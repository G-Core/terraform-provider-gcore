// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_file_share

import (
	"github.com/G-Core/terraform-provider-gcore/internal/apijson"
	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type CloudFileShareModel struct {
	ID               types.String                                               `tfsdk:"id" json:"id,computed"`
	ProjectID        types.Int64                                                `tfsdk:"project_id" path:"project_id,optional"`
	RegionID         types.Int64                                                `tfsdk:"region_id" path:"region_id,optional"`
	Protocol         types.String                                               `tfsdk:"protocol" json:"protocol,required"`
	Size             types.Int64                                                `tfsdk:"size" json:"size,required"`
	Network          *CloudFileShareNetworkModel                                `tfsdk:"network" json:"network,optional,no_refresh"`
	TypeName         types.String                                               `tfsdk:"type_name" json:"type_name,computed_optional"`
	VolumeType       types.String                                               `tfsdk:"volume_type" json:"volume_type,computed_optional"`
	Access           customfield.NestedObjectList[CloudFileShareAccessModel]    `tfsdk:"access" json:"access,computed_optional,no_refresh"`
	Name             types.String                                               `tfsdk:"name" json:"name,required"`
	Tags             customfield.Map[types.String]                              `tfsdk:"tags" json:"tags,computed_optional,no_refresh"`
	ShareSettings    customfield.NestedObject[CloudFileShareShareSettingsModel] `tfsdk:"share_settings" json:"share_settings,computed_optional"`
	ConnectionPoint  types.String                                               `tfsdk:"connection_point" json:"connection_point,computed"`
	CreatedAt        types.String                                               `tfsdk:"created_at" json:"created_at,computed"`
	CreatorTaskID    types.String                                               `tfsdk:"creator_task_id" json:"creator_task_id,computed"`
	NetworkID        types.String                                               `tfsdk:"network_id" json:"network_id,computed"`
	NetworkName      types.String                                               `tfsdk:"network_name" json:"network_name,computed"`
	Region           types.String                                               `tfsdk:"region" json:"region,computed"`
	ShareNetworkName types.String                                               `tfsdk:"share_network_name" json:"share_network_name,computed"`
	Status           types.String                                               `tfsdk:"status" json:"status,computed"`
	SubnetID         types.String                                               `tfsdk:"subnet_id" json:"subnet_id,computed"`
	SubnetName       types.String                                               `tfsdk:"subnet_name" json:"subnet_name,computed"`
	TaskID           types.String                                               `tfsdk:"task_id" json:"task_id,computed"`
	Tasks            customfield.List[types.String]                             `tfsdk:"tasks" json:"tasks,computed,no_refresh"`
}

func (m CloudFileShareModel) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(m)
}

func (m CloudFileShareModel) MarshalJSONForUpdate(state CloudFileShareModel) (data []byte, err error) {
	return apijson.MarshalForPatch(m, state)
}

type CloudFileShareNetworkModel struct {
	NetworkID types.String `tfsdk:"network_id" json:"network_id,required"`
	SubnetID  types.String `tfsdk:"subnet_id" json:"subnet_id,optional"`
}

type CloudFileShareAccessModel struct {
	AccessMode types.String `tfsdk:"access_mode" json:"access_mode,required"`
	IPAddress  types.String `tfsdk:"ip_address" json:"ip_address,required"`
}

type CloudFileShareShareSettingsModel struct {
	AllowedCharacters types.String `tfsdk:"allowed_characters" json:"allowed_characters,computed_optional"`
	PathLength        types.String `tfsdk:"path_length" json:"path_length,computed_optional"`
	RootSquash        types.Bool   `tfsdk:"root_squash" json:"root_squash,computed_optional"`
}
