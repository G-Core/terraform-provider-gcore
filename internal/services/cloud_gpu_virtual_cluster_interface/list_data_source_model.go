// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_gpu_virtual_cluster_interface

import (
	"context"

	"github.com/G-Core/gcore-go/cloud"
	"github.com/G-Core/gcore-go/packages/param"
	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type CloudGPUVirtualClusterInterfacesResultsListDataSourceEnvelope struct {
	Results customfield.NestedObjectList[CloudGPUVirtualClusterInterfacesItemsDataSourceModel] `json:"results,computed"`
}

type CloudGPUVirtualClusterInterfacesDataSourceModel struct {
	ClusterID types.String                                                                       `tfsdk:"cluster_id" path:"cluster_id,required"`
	ProjectID types.Int64                                                                        `tfsdk:"project_id" path:"project_id,optional"`
	RegionID  types.Int64                                                                        `tfsdk:"region_id" path:"region_id,optional"`
	MaxItems  types.Int64                                                                        `tfsdk:"max_items"`
	Items     customfield.NestedObjectList[CloudGPUVirtualClusterInterfacesItemsDataSourceModel] `tfsdk:"items"`
}

func (m *CloudGPUVirtualClusterInterfacesDataSourceModel) toListParams(_ context.Context) (params cloud.GPUVirtualClusterInterfaceListParams, diags diag.Diagnostics) {
	params = cloud.GPUVirtualClusterInterfaceListParams{}

	if !m.ProjectID.IsNull() {
		params.ProjectID = param.NewOpt(m.ProjectID.ValueInt64())
	}
	if !m.RegionID.IsNull() {
		params.RegionID = param.NewOpt(m.RegionID.ValueInt64())
	}

	return
}

type CloudGPUVirtualClusterInterfacesItemsDataSourceModel struct {
	FloatingIPs         customfield.NestedObjectList[CloudGPUVirtualClusterInterfacesFloatingIPsDataSourceModel]   `tfsdk:"floating_ips" json:"floating_ips,computed"`
	IPAssignments       customfield.NestedObjectList[CloudGPUVirtualClusterInterfacesIPAssignmentsDataSourceModel] `tfsdk:"ip_assignments" json:"ip_assignments,computed"`
	MacAddress          types.String                                                                               `tfsdk:"mac_address" json:"mac_address,computed"`
	Network             customfield.NestedObject[CloudGPUVirtualClusterInterfacesNetworkDataSourceModel]           `tfsdk:"network" json:"network,computed"`
	NetworkID           types.String                                                                               `tfsdk:"network_id" json:"network_id,computed"`
	PortID              types.String                                                                               `tfsdk:"port_id" json:"port_id,computed"`
	PortSecurityEnabled types.Bool                                                                                 `tfsdk:"port_security_enabled" json:"port_security_enabled,computed"`
}

type CloudGPUVirtualClusterInterfacesFloatingIPsDataSourceModel struct {
	ID                types.String                                                                                 `tfsdk:"id" json:"id,computed"`
	CreatedAt         timetypes.RFC3339                                                                            `tfsdk:"created_at" json:"created_at,computed" format:"date-time"`
	FixedIPAddress    types.String                                                                                 `tfsdk:"fixed_ip_address" json:"fixed_ip_address,computed"`
	FloatingIPAddress types.String                                                                                 `tfsdk:"floating_ip_address" json:"floating_ip_address,computed"`
	PortID            types.String                                                                                 `tfsdk:"port_id" json:"port_id,computed"`
	RouterID          types.String                                                                                 `tfsdk:"router_id" json:"router_id,computed"`
	Status            types.String                                                                                 `tfsdk:"status" json:"status,computed"`
	Tags              customfield.NestedObjectList[CloudGPUVirtualClusterInterfacesFloatingIPsTagsDataSourceModel] `tfsdk:"tags" json:"tags,computed"`
	UpdatedAt         timetypes.RFC3339                                                                            `tfsdk:"updated_at" json:"updated_at,computed" format:"date-time"`
}

type CloudGPUVirtualClusterInterfacesFloatingIPsTagsDataSourceModel struct {
	Key      types.String `tfsdk:"key" json:"key,computed"`
	ReadOnly types.Bool   `tfsdk:"read_only" json:"read_only,computed"`
	Value    types.String `tfsdk:"value" json:"value,computed"`
}

type CloudGPUVirtualClusterInterfacesIPAssignmentsDataSourceModel struct {
	IPAddress types.String `tfsdk:"ip_address" json:"ip_address,computed"`
	SubnetID  types.String `tfsdk:"subnet_id" json:"subnet_id,computed"`
}

type CloudGPUVirtualClusterInterfacesNetworkDataSourceModel struct {
	ID                  types.String                                                                                `tfsdk:"id" json:"id,computed"`
	CreatedAt           timetypes.RFC3339                                                                           `tfsdk:"created_at" json:"created_at,computed" format:"date-time"`
	External            types.Bool                                                                                  `tfsdk:"external" json:"external,computed"`
	Mtu                 types.Int64                                                                                 `tfsdk:"mtu" json:"mtu,computed"`
	Name                types.String                                                                                `tfsdk:"name" json:"name,computed"`
	PortSecurityEnabled types.Bool                                                                                  `tfsdk:"port_security_enabled" json:"port_security_enabled,computed"`
	SegmentationID      types.Int64                                                                                 `tfsdk:"segmentation_id" json:"segmentation_id,computed"`
	Shared              types.Bool                                                                                  `tfsdk:"shared" json:"shared,computed"`
	Subnets             customfield.NestedObjectList[CloudGPUVirtualClusterInterfacesNetworkSubnetsDataSourceModel] `tfsdk:"subnets" json:"subnets,computed"`
	Tags                customfield.NestedObjectList[CloudGPUVirtualClusterInterfacesNetworkTagsDataSourceModel]    `tfsdk:"tags" json:"tags,computed"`
	Type                types.String                                                                                `tfsdk:"type" json:"type,computed"`
	UpdatedAt           timetypes.RFC3339                                                                           `tfsdk:"updated_at" json:"updated_at,computed" format:"date-time"`
}

type CloudGPUVirtualClusterInterfacesNetworkSubnetsDataSourceModel struct {
	ID             types.String                                                                                          `tfsdk:"id" json:"id,computed"`
	AvailableIPs   types.Int64                                                                                           `tfsdk:"available_ips" json:"available_ips,computed"`
	Cidr           types.String                                                                                          `tfsdk:"cidr" json:"cidr,computed"`
	CreatedAt      timetypes.RFC3339                                                                                     `tfsdk:"created_at" json:"created_at,computed" format:"date-time"`
	DNSNameservers customfield.List[types.String]                                                                        `tfsdk:"dns_nameservers" json:"dns_nameservers,computed"`
	EnableDhcp     types.Bool                                                                                            `tfsdk:"enable_dhcp" json:"enable_dhcp,computed"`
	GatewayIP      types.String                                                                                          `tfsdk:"gateway_ip" json:"gateway_ip,computed"`
	HasRouter      types.Bool                                                                                            `tfsdk:"has_router" json:"has_router,computed"`
	HostRoutes     customfield.NestedObjectList[CloudGPUVirtualClusterInterfacesNetworkSubnetsHostRoutesDataSourceModel] `tfsdk:"host_routes" json:"host_routes,computed"`
	IPVersion      types.Int64                                                                                           `tfsdk:"ip_version" json:"ip_version,computed"`
	Name           types.String                                                                                          `tfsdk:"name" json:"name,computed"`
	NetworkID      types.String                                                                                          `tfsdk:"network_id" json:"network_id,computed"`
	Tags           customfield.NestedObjectList[CloudGPUVirtualClusterInterfacesNetworkSubnetsTagsDataSourceModel]       `tfsdk:"tags" json:"tags,computed"`
	TotalIPs       types.Int64                                                                                           `tfsdk:"total_ips" json:"total_ips,computed"`
	UpdatedAt      timetypes.RFC3339                                                                                     `tfsdk:"updated_at" json:"updated_at,computed" format:"date-time"`
}

type CloudGPUVirtualClusterInterfacesNetworkSubnetsHostRoutesDataSourceModel struct {
	Destination types.String `tfsdk:"destination" json:"destination,computed"`
	Nexthop     types.String `tfsdk:"nexthop" json:"nexthop,computed"`
}

type CloudGPUVirtualClusterInterfacesNetworkSubnetsTagsDataSourceModel struct {
	Key      types.String `tfsdk:"key" json:"key,computed"`
	ReadOnly types.Bool   `tfsdk:"read_only" json:"read_only,computed"`
	Value    types.String `tfsdk:"value" json:"value,computed"`
}

type CloudGPUVirtualClusterInterfacesNetworkTagsDataSourceModel struct {
	Key      types.String `tfsdk:"key" json:"key,computed"`
	ReadOnly types.Bool   `tfsdk:"read_only" json:"read_only,computed"`
	Value    types.String `tfsdk:"value" json:"value,computed"`
}
