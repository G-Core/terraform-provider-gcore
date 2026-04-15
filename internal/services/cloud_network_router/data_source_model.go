// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_network_router

import (
	"context"

	"github.com/G-Core/gcore-go/cloud"
	"github.com/G-Core/gcore-go/packages/param"
	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type CloudNetworkRouterDataSourceModel struct {
	ID                  types.String                                                                   `tfsdk:"id" path:"router_id,computed"`
	RouterID            types.String                                                                   `tfsdk:"router_id" path:"router_id,optional"`
	ProjectID           types.Int64                                                                    `tfsdk:"project_id" path:"project_id,optional"`
	RegionID            types.Int64                                                                    `tfsdk:"region_id" path:"region_id,optional"`
	CreatedAt           timetypes.RFC3339                                                              `tfsdk:"created_at" json:"created_at,computed" format:"date-time"`
	CreatorTaskID       types.String                                                                   `tfsdk:"creator_task_id" json:"creator_task_id,computed"`
	Distributed         types.Bool                                                                     `tfsdk:"distributed" json:"distributed,computed"`
	Name                types.String                                                                   `tfsdk:"name" json:"name,computed"`
	Region              types.String                                                                   `tfsdk:"region" json:"region,computed"`
	Status              types.String                                                                   `tfsdk:"status" json:"status,computed"`
	TaskID              types.String                                                                   `tfsdk:"task_id" json:"task_id,computed"`
	UpdatedAt           timetypes.RFC3339                                                              `tfsdk:"updated_at" json:"updated_at,computed" format:"date-time"`
	ExternalGatewayInfo customfield.NestedObject[CloudNetworkRouterExternalGatewayInfoDataSourceModel] `tfsdk:"external_gateway_info" json:"external_gateway_info,computed"`
	Interfaces          customfield.NestedObjectList[CloudNetworkRouterInterfacesDataSourceModel]      `tfsdk:"interfaces" json:"interfaces,computed"`
	Routes              customfield.NestedObjectList[CloudNetworkRouterRoutesDataSourceModel]          `tfsdk:"routes" json:"routes,computed"`
	FindOneBy           *CloudNetworkRouterFindOneByDataSourceModel                                    `tfsdk:"find_one_by"`
}

func (m *CloudNetworkRouterDataSourceModel) toReadParams(_ context.Context) (params cloud.NetworkRouterGetParams, diags diag.Diagnostics) {
	params = cloud.NetworkRouterGetParams{}

	if !m.ProjectID.IsNull() {
		params.ProjectID = param.NewOpt(m.ProjectID.ValueInt64())
	}
	if !m.RegionID.IsNull() {
		params.RegionID = param.NewOpt(m.RegionID.ValueInt64())
	}

	return
}

func (m *CloudNetworkRouterDataSourceModel) toListParams(_ context.Context) (params cloud.NetworkRouterListParams, diags diag.Diagnostics) {
	params = cloud.NetworkRouterListParams{}

	if !m.ProjectID.IsNull() {
		params.ProjectID = param.NewOpt(m.ProjectID.ValueInt64())
	}
	if !m.RegionID.IsNull() {
		params.RegionID = param.NewOpt(m.RegionID.ValueInt64())
	}
	if !m.FindOneBy.Name.IsNull() {
		params.Name = param.NewOpt(m.FindOneBy.Name.ValueString())
	}

	return
}

type CloudNetworkRouterExternalGatewayInfoDataSourceModel struct {
	EnableSnat       types.Bool                                                                                         `tfsdk:"enable_snat" json:"enable_snat,computed"`
	ExternalFixedIPs customfield.NestedObjectList[CloudNetworkRouterExternalGatewayInfoExternalFixedIPsDataSourceModel] `tfsdk:"external_fixed_ips" json:"external_fixed_ips,computed"`
	NetworkID        types.String                                                                                       `tfsdk:"network_id" json:"network_id,computed"`
}

type CloudNetworkRouterExternalGatewayInfoExternalFixedIPsDataSourceModel struct {
	IPAddress types.String `tfsdk:"ip_address" json:"ip_address,computed"`
	SubnetID  types.String `tfsdk:"subnet_id" json:"subnet_id,computed"`
}

type CloudNetworkRouterInterfacesDataSourceModel struct {
	IPAssignments customfield.NestedObjectList[CloudNetworkRouterInterfacesIPAssignmentsDataSourceModel] `tfsdk:"ip_assignments" json:"ip_assignments,computed"`
	NetworkID     types.String                                                                           `tfsdk:"network_id" json:"network_id,computed"`
	PortID        types.String                                                                           `tfsdk:"port_id" json:"port_id,computed"`
	MacAddress    types.String                                                                           `tfsdk:"mac_address" json:"mac_address,computed"`
}

type CloudNetworkRouterInterfacesIPAssignmentsDataSourceModel struct {
	IPAddress types.String `tfsdk:"ip_address" json:"ip_address,computed"`
	SubnetID  types.String `tfsdk:"subnet_id" json:"subnet_id,computed"`
}

type CloudNetworkRouterRoutesDataSourceModel struct {
	Destination types.String `tfsdk:"destination" json:"destination,computed"`
	Nexthop     types.String `tfsdk:"nexthop" json:"nexthop,computed"`
}

type CloudNetworkRouterFindOneByDataSourceModel struct {
	Name types.String `tfsdk:"name" query:"name,optional"`
}
