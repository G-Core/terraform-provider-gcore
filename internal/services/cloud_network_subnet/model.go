// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_network_subnet

import (
	"github.com/G-Core/terraform-provider-gcore/internal/apijson"
	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type CloudNetworkSubnetModel struct {
	ID             types.String                                                    `tfsdk:"id" json:"id,computed"`
	ProjectID      types.Int64                                                     `tfsdk:"project_id" path:"project_id,optional"`
	RegionID       types.Int64                                                     `tfsdk:"region_id" path:"region_id,optional"`
	Cidr           types.String                                                    `tfsdk:"cidr" json:"cidr,required"`
	NetworkID      types.String                                                    `tfsdk:"network_id" json:"network_id,required"`
	IPVersion      types.Int64                                                     `tfsdk:"ip_version" json:"ip_version,computed_optional"`
	Name           types.String                                                    `tfsdk:"name" json:"name,required"`
	Tags           *map[string]types.String                                        `tfsdk:"tags" json:"tags,optional,no_refresh"`
	EnableDhcp     types.Bool                                                      `tfsdk:"enable_dhcp" json:"enable_dhcp,computed_optional"`
	GatewayIP      types.String                                                    `tfsdk:"gateway_ip" json:"gateway_ip,computed_optional"`
	DNSNameservers customfield.List[types.String]                                  `tfsdk:"dns_nameservers" json:"dns_nameservers,computed_optional"`
	HostRoutes     customfield.NestedObjectList[CloudNetworkSubnetHostRoutesModel] `tfsdk:"host_routes" json:"host_routes,computed_optional"`
	AvailableIPs   types.Int64                                                     `tfsdk:"available_ips" json:"available_ips,computed"`
	CreatedAt      timetypes.RFC3339                                               `tfsdk:"created_at" json:"created_at,computed" format:"date-time"`
	CreatorTaskID  types.String                                                    `tfsdk:"creator_task_id" json:"creator_task_id,computed"`
	HasRouter      types.Bool                                                      `tfsdk:"has_router" json:"has_router,computed"`
	Region         types.String                                                    `tfsdk:"region" json:"region,computed"`
	TotalIPs       types.Int64                                                     `tfsdk:"total_ips" json:"total_ips,computed"`
	UpdatedAt      timetypes.RFC3339                                               `tfsdk:"updated_at" json:"updated_at,computed" format:"date-time"`
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
