// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_network_subnet

import (
	"context"

	"github.com/G-Core/gcore-go/cloud"
	"github.com/G-Core/gcore-go/packages/param"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

type CloudNetworkSubnetDataSourceModel struct {
	ProjectID      types.Int64                                                               `tfsdk:"project_id" path:"project_id,required"`
	RegionID       types.Int64                                                               `tfsdk:"region_id" path:"region_id,required"`
	SubnetID       types.String                                                              `tfsdk:"subnet_id" path:"subnet_id,required"`
	AvailableIPs   types.Int64                                                               `tfsdk:"available_ips" json:"available_ips,computed"`
	Cidr           types.String                                                              `tfsdk:"cidr" json:"cidr,computed"`
	CreatedAt      timetypes.RFC3339                                                         `tfsdk:"created_at" json:"created_at,computed" format:"date-time"`
	CreatorTaskID  types.String                                                              `tfsdk:"creator_task_id" json:"creator_task_id,computed"`
	EnableDhcp     types.Bool                                                                `tfsdk:"enable_dhcp" json:"enable_dhcp,computed"`
	GatewayIP      types.String                                                              `tfsdk:"gateway_ip" json:"gateway_ip,computed"`
	HasRouter      types.Bool                                                                `tfsdk:"has_router" json:"has_router,computed"`
	ID             types.String                                                              `tfsdk:"id" json:"id,computed"`
	IPVersion      types.Int64                                                               `tfsdk:"ip_version" json:"ip_version,computed"`
	Name           types.String                                                              `tfsdk:"name" json:"name,computed"`
	NetworkID      types.String                                                              `tfsdk:"network_id" json:"network_id,computed"`
	Region         types.String                                                              `tfsdk:"region" json:"region,computed"`
	TaskID         types.String                                                              `tfsdk:"task_id" json:"task_id,computed"`
	TotalIPs       types.Int64                                                               `tfsdk:"total_ips" json:"total_ips,computed"`
	UpdatedAt      timetypes.RFC3339                                                         `tfsdk:"updated_at" json:"updated_at,computed" format:"date-time"`
	DNSNameservers customfield.List[types.String]                                            `tfsdk:"dns_nameservers" json:"dns_nameservers,computed"`
	HostRoutes     customfield.NestedObjectList[CloudNetworkSubnetHostRoutesDataSourceModel] `tfsdk:"host_routes" json:"host_routes,computed"`
	Tags           customfield.NestedObjectList[CloudNetworkSubnetTagsDataSourceModel]       `tfsdk:"tags" json:"tags,computed"`
}

func (m *CloudNetworkSubnetDataSourceModel) toReadParams(_ context.Context) (params cloud.NetworkSubnetGetParams, diags diag.Diagnostics) {
	params = cloud.NetworkSubnetGetParams{}

	if !m.ProjectID.IsNull() {
		params.ProjectID = param.NewOpt(m.ProjectID.ValueInt64())
	}
	if !m.RegionID.IsNull() {
		params.RegionID = param.NewOpt(m.RegionID.ValueInt64())
	}

	return
}

type CloudNetworkSubnetHostRoutesDataSourceModel struct {
	Destination types.String `tfsdk:"destination" json:"destination,computed"`
	Nexthop     types.String `tfsdk:"nexthop" json:"nexthop,computed"`
}

type CloudNetworkSubnetTagsDataSourceModel struct {
	Key      types.String `tfsdk:"key" json:"key,computed"`
	ReadOnly types.Bool   `tfsdk:"read_only" json:"read_only,computed"`
	Value    types.String `tfsdk:"value" json:"value,computed"`
}
