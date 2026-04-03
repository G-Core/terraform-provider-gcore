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

type CloudNetworkRoutersResultsListDataSourceEnvelope struct {
	Results customfield.NestedObjectList[CloudNetworkRoutersItemsDataSourceModel] `json:"results,computed"`
}

type CloudNetworkRoutersDataSourceModel struct {
	ProjectID types.Int64                                                           `tfsdk:"project_id" path:"project_id,optional"`
	RegionID  types.Int64                                                           `tfsdk:"region_id" path:"region_id,optional"`
	Limit     types.Int64                                                           `tfsdk:"limit" query:"limit,optional"`
	MaxItems  types.Int64                                                           `tfsdk:"max_items"`
	Items     customfield.NestedObjectList[CloudNetworkRoutersItemsDataSourceModel] `tfsdk:"items"`
}

func (m *CloudNetworkRoutersDataSourceModel) toListParams(_ context.Context) (params cloud.NetworkRouterListParams, diags diag.Diagnostics) {
	params = cloud.NetworkRouterListParams{}

	if !m.ProjectID.IsNull() {
		params.ProjectID = param.NewOpt(m.ProjectID.ValueInt64())
	}
	if !m.RegionID.IsNull() {
		params.RegionID = param.NewOpt(m.RegionID.ValueInt64())
	}
	if !m.Limit.IsNull() {
		params.Limit = param.NewOpt(m.Limit.ValueInt64())
	}

	return
}

type CloudNetworkRoutersItemsDataSourceModel struct {
	ID                  types.String                                                                    `tfsdk:"id" json:"id,computed"`
	CreatedAt           timetypes.RFC3339                                                               `tfsdk:"created_at" json:"created_at,computed" format:"date-time"`
	Distributed         types.Bool                                                                      `tfsdk:"distributed" json:"distributed,computed"`
	Interfaces          customfield.NestedObjectList[CloudNetworkRoutersInterfacesDataSourceModel]      `tfsdk:"interfaces" json:"interfaces,computed"`
	Name                types.String                                                                    `tfsdk:"name" json:"name,computed"`
	ProjectID           types.Int64                                                                     `tfsdk:"project_id" json:"project_id,computed"`
	Region              types.String                                                                    `tfsdk:"region" json:"region,computed"`
	RegionID            types.Int64                                                                     `tfsdk:"region_id" json:"region_id,computed"`
	Routes              customfield.NestedObjectList[CloudNetworkRoutersRoutesDataSourceModel]          `tfsdk:"routes" json:"routes,computed"`
	Status              types.String                                                                    `tfsdk:"status" json:"status,computed"`
	UpdatedAt           timetypes.RFC3339                                                               `tfsdk:"updated_at" json:"updated_at,computed" format:"date-time"`
	CreatorTaskID       types.String                                                                    `tfsdk:"creator_task_id" json:"creator_task_id,computed"`
	ExternalGatewayInfo customfield.NestedObject[CloudNetworkRoutersExternalGatewayInfoDataSourceModel] `tfsdk:"external_gateway_info" json:"external_gateway_info,computed"`
}

type CloudNetworkRoutersInterfacesDataSourceModel struct {
	IPAssignments customfield.NestedObjectList[CloudNetworkRoutersInterfacesIPAssignmentsDataSourceModel] `tfsdk:"ip_assignments" json:"ip_assignments,computed"`
	NetworkID     types.String                                                                            `tfsdk:"network_id" json:"network_id,computed"`
	PortID        types.String                                                                            `tfsdk:"port_id" json:"port_id,computed"`
	MacAddress    types.String                                                                            `tfsdk:"mac_address" json:"mac_address,computed"`
}

type CloudNetworkRoutersInterfacesIPAssignmentsDataSourceModel struct {
	IPAddress types.String `tfsdk:"ip_address" json:"ip_address,computed"`
	SubnetID  types.String `tfsdk:"subnet_id" json:"subnet_id,computed"`
}

type CloudNetworkRoutersRoutesDataSourceModel struct {
	Destination types.String `tfsdk:"destination" json:"destination,computed"`
	Nexthop     types.String `tfsdk:"nexthop" json:"nexthop,computed"`
}

type CloudNetworkRoutersExternalGatewayInfoDataSourceModel struct {
	EnableSnat       types.Bool                                                                                          `tfsdk:"enable_snat" json:"enable_snat,computed"`
	ExternalFixedIPs customfield.NestedObjectList[CloudNetworkRoutersExternalGatewayInfoExternalFixedIPsDataSourceModel] `tfsdk:"external_fixed_ips" json:"external_fixed_ips,computed"`
	NetworkID        types.String                                                                                        `tfsdk:"network_id" json:"network_id,computed"`
}

type CloudNetworkRoutersExternalGatewayInfoExternalFixedIPsDataSourceModel struct {
	IPAddress types.String `tfsdk:"ip_address" json:"ip_address,computed"`
	SubnetID  types.String `tfsdk:"subnet_id" json:"subnet_id,computed"`
}
