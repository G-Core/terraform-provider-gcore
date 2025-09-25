// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_network_subnet

import (
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/apijson"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

type CloudNetworkSubnetModel struct {
	ID                     types.String                                                    `tfsdk:"id" json:"id,computed"`
	ProjectID              types.Int64                                                     `tfsdk:"project_id" path:"project_id,optional"`
	RegionID               types.Int64                                                     `tfsdk:"region_id" path:"region_id,optional"`
	Cidr                   types.String                                                    `tfsdk:"cidr" json:"cidr,required"`
	NetworkID              types.String                                                    `tfsdk:"network_id" json:"network_id,required"`
	IPVersion              types.Int64                                                     `tfsdk:"ip_version" json:"ip_version,computed_optional"`
	RouterIDToConnect      types.String                                                    `tfsdk:"router_id_to_connect" json:"router_id_to_connect,optional,no_refresh"`
	ConnectToNetworkRouter types.Bool                                                      `tfsdk:"connect_to_network_router" json:"connect_to_network_router,computed_optional,no_refresh"`
	Name                   types.String                                                    `tfsdk:"name" json:"name,required"`
	Tags                   *map[string]types.String                                        `tfsdk:"tags" json:"tags,optional,no_refresh"`
	EnableDhcp             types.Bool                                                      `tfsdk:"enable_dhcp" json:"enable_dhcp,computed_optional"`
	GatewayIP              types.String                                                    `tfsdk:"gateway_ip" json:"gateway_ip,computed_optional"`
	DNSNameservers         customfield.List[types.String]                                  `tfsdk:"dns_nameservers" json:"dns_nameservers,computed_optional"`
	HostRoutes             customfield.NestedObjectList[CloudNetworkSubnetHostRoutesModel] `tfsdk:"host_routes" json:"host_routes,computed_optional"`
	AvailableIPs           types.Int64                                                     `tfsdk:"available_ips" json:"available_ips,computed"`
	CreatedAt              timetypes.RFC3339                                               `tfsdk:"created_at" json:"created_at,computed" format:"date-time"`
	CreatorTaskID          types.String                                                    `tfsdk:"creator_task_id" json:"creator_task_id,computed"`
	HasRouter              types.Bool                                                      `tfsdk:"has_router" json:"has_router,computed"`
	Region                 types.String                                                    `tfsdk:"region" json:"region,computed"`
	TaskID                 types.String                                                    `tfsdk:"task_id" json:"task_id,computed"`
	TotalIPs               types.Int64                                                     `tfsdk:"total_ips" json:"total_ips,computed"`
	UpdatedAt              timetypes.RFC3339                                               `tfsdk:"updated_at" json:"updated_at,computed" format:"date-time"`
	Tasks                  customfield.List[types.String]                                  `tfsdk:"tasks" json:"tasks,computed,no_refresh"`
}

func (m CloudNetworkSubnetModel) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(m)
}

func (m CloudNetworkSubnetModel) MarshalJSONForUpdate(state CloudNetworkSubnetModel) (data []byte, err error) {
	return apijson.MarshalForPatch(m, state)
}

type CloudNetworkSubnetHostRoutesModel struct {
	Destination types.String `tfsdk:"destination" json:"destination,required"`
	Nexthop     types.String `tfsdk:"nexthop" json:"nexthop,required"`
}
