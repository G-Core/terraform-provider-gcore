// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_network_router

import (
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/apijson"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

type CloudNetworkRouterModel struct {
	ID                  types.String                                                         `tfsdk:"id" json:"id,computed"`
	ProjectID           types.Int64                                                          `tfsdk:"project_id" path:"project_id,optional"`
	RegionID            types.Int64                                                          `tfsdk:"region_id" path:"region_id,optional"`
	Interfaces          customfield.NestedObjectList[CloudNetworkRouterInterfacesModel]      `tfsdk:"interfaces" json:"interfaces,computed_optional"`
	Name                types.String                                                         `tfsdk:"name" json:"name,required"`
	ExternalGatewayInfo customfield.NestedObject[CloudNetworkRouterExternalGatewayInfoModel] `tfsdk:"external_gateway_info" json:"external_gateway_info,computed_optional"`
	Routes              customfield.NestedObjectList[CloudNetworkRouterRoutesModel]          `tfsdk:"routes" json:"routes,computed_optional"`
	CreatedAt           timetypes.RFC3339                                                    `tfsdk:"created_at" json:"created_at,computed" format:"date-time"`
	CreatorTaskID       types.String                                                         `tfsdk:"creator_task_id" json:"creator_task_id,computed"`
	Distributed         types.Bool                                                           `tfsdk:"distributed" json:"distributed,computed"`
	Region              types.String                                                         `tfsdk:"region" json:"region,computed"`
	Status              types.String                                                         `tfsdk:"status" json:"status,computed"`
	TaskID              types.String                                                         `tfsdk:"task_id" json:"task_id,computed"`
	UpdatedAt           timetypes.RFC3339                                                    `tfsdk:"updated_at" json:"updated_at,computed" format:"date-time"`
	Tasks               customfield.List[types.String]                                       `tfsdk:"tasks" json:"tasks,computed,no_refresh"`
}

func (m CloudNetworkRouterModel) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(m)
}

func (m CloudNetworkRouterModel) MarshalJSONForUpdate(state CloudNetworkRouterModel) (data []byte, err error) {
	// Create a copy of the model to marshal, but force interfaces to equal state
	// so they're not included in the PATCH request (interfaces are managed via attach/detach)
	mCopy := m
	mCopy.Interfaces = state.Interfaces
	return apijson.MarshalForPatch(mCopy, state)
}

type CloudNetworkRouterInterfacesModel struct {
	SubnetID types.String `tfsdk:"subnet_id" json:"subnet_id,required,no_refresh"`
	Type     types.String `tfsdk:"type" json:"type,computed_optional,no_refresh"`
}

type CloudNetworkRouterExternalGatewayInfoModel struct {
	NetworkID  types.String `tfsdk:"network_id" json:"network_id,computed_optional"`
	EnableSnat types.Bool   `tfsdk:"enable_snat" json:"enable_snat,computed_optional"`
	Type       types.String `tfsdk:"type" json:"type,computed_optional,no_refresh"`
}

type CloudNetworkRouterRoutesModel struct {
	Destination types.String `tfsdk:"destination" json:"destination,required"`
	Nexthop     types.String `tfsdk:"nexthop" json:"nexthop,required"`
}
