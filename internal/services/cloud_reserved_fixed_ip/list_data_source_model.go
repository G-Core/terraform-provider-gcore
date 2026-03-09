// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_reserved_fixed_ip

import (
	"context"

	"github.com/G-Core/gcore-go/cloud"
	"github.com/G-Core/gcore-go/packages/param"
	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type CloudReservedFixedIPsResultsListDataSourceEnvelope struct {
	Results customfield.NestedObjectList[CloudReservedFixedIPsItemsDataSourceModel] `json:"results,computed"`
}

type CloudReservedFixedIPsDataSourceModel struct {
	ProjectID     types.Int64                                                             `tfsdk:"project_id" path:"project_id,optional"`
	RegionID      types.Int64                                                             `tfsdk:"region_id" path:"region_id,optional"`
	AvailableOnly types.Bool                                                              `tfsdk:"available_only" query:"available_only,optional"`
	DeviceID      types.String                                                            `tfsdk:"device_id" query:"device_id,optional"`
	ExternalOnly  types.Bool                                                              `tfsdk:"external_only" query:"external_only,optional"`
	InternalOnly  types.Bool                                                              `tfsdk:"internal_only" query:"internal_only,optional"`
	IPAddress     types.String                                                            `tfsdk:"ip_address" query:"ip_address,optional"`
	OrderBy       types.String                                                            `tfsdk:"order_by" query:"order_by,optional"`
	VipOnly       types.Bool                                                              `tfsdk:"vip_only" query:"vip_only,optional"`
	MaxItems      types.Int64                                                             `tfsdk:"max_items"`
	Items         customfield.NestedObjectList[CloudReservedFixedIPsItemsDataSourceModel] `tfsdk:"items"`
}

func (m *CloudReservedFixedIPsDataSourceModel) toListParams(_ context.Context) (params cloud.ReservedFixedIPListParams, diags diag.Diagnostics) {
	params = cloud.ReservedFixedIPListParams{}

	if !m.ProjectID.IsNull() {
		params.ProjectID = param.NewOpt(m.ProjectID.ValueInt64())
	}
	if !m.RegionID.IsNull() {
		params.RegionID = param.NewOpt(m.RegionID.ValueInt64())
	}
	if !m.AvailableOnly.IsNull() {
		params.AvailableOnly = param.NewOpt(m.AvailableOnly.ValueBool())
	}
	if !m.DeviceID.IsNull() {
		params.DeviceID = param.NewOpt(m.DeviceID.ValueString())
	}
	if !m.ExternalOnly.IsNull() {
		params.ExternalOnly = param.NewOpt(m.ExternalOnly.ValueBool())
	}
	if !m.InternalOnly.IsNull() {
		params.InternalOnly = param.NewOpt(m.InternalOnly.ValueBool())
	}
	if !m.IPAddress.IsNull() {
		params.IPAddress = param.NewOpt(m.IPAddress.ValueString())
	}
	if !m.OrderBy.IsNull() {
		params.OrderBy = param.NewOpt(m.OrderBy.ValueString())
	}
	if !m.VipOnly.IsNull() {
		params.VipOnly = param.NewOpt(m.VipOnly.ValueBool())
	}

	return
}

type CloudReservedFixedIPsItemsDataSourceModel struct {
	AllowedAddressPairs customfield.NestedObjectList[CloudReservedFixedIPsAllowedAddressPairsDataSourceModel] `tfsdk:"allowed_address_pairs" json:"allowed_address_pairs,computed"`
	Attachments         customfield.NestedObjectList[CloudReservedFixedIPsAttachmentsDataSourceModel]         `tfsdk:"attachments" json:"attachments,computed"`
	CreatedAt           timetypes.RFC3339                                                                     `tfsdk:"created_at" json:"created_at,computed" format:"date-time"`
	IsExternal          types.Bool                                                                            `tfsdk:"is_external" json:"is_external,computed"`
	IsVip               types.Bool                                                                            `tfsdk:"is_vip" json:"is_vip,computed"`
	Name                types.String                                                                          `tfsdk:"name" json:"name,computed"`
	Network             customfield.NestedObject[CloudReservedFixedIPsNetworkDataSourceModel]                 `tfsdk:"network" json:"network,computed"`
	NetworkID           types.String                                                                          `tfsdk:"network_id" json:"network_id,computed"`
	PortID              types.String                                                                          `tfsdk:"port_id" json:"port_id,computed"`
	Region              types.String                                                                          `tfsdk:"region" json:"region,computed"`
	RegionID            types.Int64                                                                           `tfsdk:"region_id" json:"region_id,computed"`
	Reservation         customfield.NestedObject[CloudReservedFixedIPsReservationDataSourceModel]             `tfsdk:"reservation" json:"reservation,computed"`
	Status              types.String                                                                          `tfsdk:"status" json:"status,computed"`
	UpdatedAt           timetypes.RFC3339                                                                     `tfsdk:"updated_at" json:"updated_at,computed" format:"date-time"`
	CreatorTaskID       types.String                                                                          `tfsdk:"creator_task_id" json:"creator_task_id,computed"`
	FixedIPAddress      types.String                                                                          `tfsdk:"fixed_ip_address" json:"fixed_ip_address,computed"`
	FixedIpv6Address    types.String                                                                          `tfsdk:"fixed_ipv6_address" json:"fixed_ipv6_address,computed"`
	ProjectID           types.Int64                                                                           `tfsdk:"project_id" json:"project_id,computed"`
	SubnetID            types.String                                                                          `tfsdk:"subnet_id" json:"subnet_id,computed"`
	SubnetV6ID          types.String                                                                          `tfsdk:"subnet_v6_id" json:"subnet_v6_id,computed"`
}

type CloudReservedFixedIPsAllowedAddressPairsDataSourceModel struct {
	IPAddress  types.String `tfsdk:"ip_address" json:"ip_address,computed"`
	MacAddress types.String `tfsdk:"mac_address" json:"mac_address,computed"`
}

type CloudReservedFixedIPsAttachmentsDataSourceModel struct {
	ResourceID   types.String `tfsdk:"resource_id" json:"resource_id,computed"`
	ResourceType types.String `tfsdk:"resource_type" json:"resource_type,computed"`
}

type CloudReservedFixedIPsNetworkDataSourceModel struct {
	ID                  types.String                                                                  `tfsdk:"id" json:"id,computed"`
	CreatedAt           timetypes.RFC3339                                                             `tfsdk:"created_at" json:"created_at,computed" format:"date-time"`
	CreatorTaskID       types.String                                                                  `tfsdk:"creator_task_id" json:"creator_task_id,computed"`
	Default             types.Bool                                                                    `tfsdk:"default" json:"default,computed"`
	External            types.Bool                                                                    `tfsdk:"external" json:"external,computed"`
	Mtu                 types.Int64                                                                   `tfsdk:"mtu" json:"mtu,computed"`
	Name                types.String                                                                  `tfsdk:"name" json:"name,computed"`
	PortSecurityEnabled types.Bool                                                                    `tfsdk:"port_security_enabled" json:"port_security_enabled,computed"`
	ProjectID           types.Int64                                                                   `tfsdk:"project_id" json:"project_id,computed"`
	Region              types.String                                                                  `tfsdk:"region" json:"region,computed"`
	RegionID            types.Int64                                                                   `tfsdk:"region_id" json:"region_id,computed"`
	SegmentationID      types.Int64                                                                   `tfsdk:"segmentation_id" json:"segmentation_id,computed"`
	Shared              types.Bool                                                                    `tfsdk:"shared" json:"shared,computed"`
	Subnets             customfield.List[types.String]                                                `tfsdk:"subnets" json:"subnets,computed"`
	Tags                customfield.NestedObjectList[CloudReservedFixedIPsNetworkTagsDataSourceModel] `tfsdk:"tags" json:"tags,computed"`
	Type                types.String                                                                  `tfsdk:"type" json:"type,computed"`
	UpdatedAt           timetypes.RFC3339                                                             `tfsdk:"updated_at" json:"updated_at,computed" format:"date-time"`
}

type CloudReservedFixedIPsNetworkTagsDataSourceModel struct {
	Key      types.String `tfsdk:"key" json:"key,computed"`
	ReadOnly types.Bool   `tfsdk:"read_only" json:"read_only,computed"`
	Value    types.String `tfsdk:"value" json:"value,computed"`
}

type CloudReservedFixedIPsReservationDataSourceModel struct {
	ResourceID   types.String `tfsdk:"resource_id" json:"resource_id,computed"`
	ResourceType types.String `tfsdk:"resource_type" json:"resource_type,computed"`
	Status       types.String `tfsdk:"status" json:"status,computed"`
}
