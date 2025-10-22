// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_reserved_fixed_ip

import (
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/apijson"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

type CloudReservedFixedIPModel struct {
	ID                  types.String                                                               `tfsdk:"id" json:"id,computed"`
	ProjectID           types.Int64                                                                `tfsdk:"project_id" path:"project_id,optional"`
	RegionID            types.Int64                                                                `tfsdk:"region_id" path:"region_id,optional"`
	Type                types.String                                                               `tfsdk:"type" json:"type,required,no_refresh"`
	IPAddress           types.String                                                               `tfsdk:"ip_address" json:"ip_address,optional,no_refresh"`
	IPFamily            types.String                                                               `tfsdk:"ip_family" json:"ip_family,optional,no_refresh"`
	NetworkID           types.String                                                               `tfsdk:"network_id" json:"network_id,optional"`
	SubnetID            types.String                                                               `tfsdk:"subnet_id" json:"subnet_id,optional"`
	PortID              types.String                                                               `tfsdk:"port_id" json:"port_id,computed_optional"`
	IsVip               types.Bool                                                                 `tfsdk:"is_vip" json:"is_vip,optional"`
	CreatedAt           timetypes.RFC3339                                                          `tfsdk:"created_at" json:"created_at,computed" format:"date-time"`
	CreatorTaskID       types.String                                                               `tfsdk:"creator_task_id" json:"creator_task_id,computed"`
	FixedIPAddress      types.String                                                               `tfsdk:"fixed_ip_address" json:"fixed_ip_address,computed"`
	FixedIpv6Address    types.String                                                               `tfsdk:"fixed_ipv6_address" json:"fixed_ipv6_address,computed"`
	IsExternal          types.Bool                                                                 `tfsdk:"is_external" json:"is_external,computed"`
	Name                types.String                                                               `tfsdk:"name" json:"name,computed"`
	Region              types.String                                                               `tfsdk:"region" json:"region,computed"`
	Status              types.String                                                               `tfsdk:"status" json:"status,computed"`
	SubnetV6ID          types.String                                                               `tfsdk:"subnet_v6_id" json:"subnet_v6_id,computed"`
	TaskID              types.String                                                               `tfsdk:"task_id" json:"task_id,computed"`
	UpdatedAt           timetypes.RFC3339                                                          `tfsdk:"updated_at" json:"updated_at,computed" format:"date-time"`
	Tasks               customfield.List[types.String]                                             `tfsdk:"tasks" json:"tasks,computed,no_refresh"`
	AllowedAddressPairs customfield.NestedObjectList[CloudReservedFixedIPAllowedAddressPairsModel] `tfsdk:"allowed_address_pairs" json:"allowed_address_pairs,computed"`
	Attachments         customfield.NestedObjectList[CloudReservedFixedIPAttachmentsModel]         `tfsdk:"attachments" json:"attachments,computed"`
	Network             customfield.NestedObject[CloudReservedFixedIPNetworkModel]                 `tfsdk:"network" json:"network,computed"`
	Reservation         customfield.NestedObject[CloudReservedFixedIPReservationModel]             `tfsdk:"reservation" json:"reservation,computed"`
}

func (m CloudReservedFixedIPModel) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(m)
}

func (m CloudReservedFixedIPModel) MarshalJSONForUpdate(state CloudReservedFixedIPModel) (data []byte, err error) {
	return apijson.MarshalForPatch(m, state)
}

type CloudReservedFixedIPAllowedAddressPairsModel struct {
	IPAddress  types.String `tfsdk:"ip_address" json:"ip_address,computed"`
	MacAddress types.String `tfsdk:"mac_address" json:"mac_address,computed"`
}

type CloudReservedFixedIPAttachmentsModel struct {
	ResourceID   types.String `tfsdk:"resource_id" json:"resource_id,computed"`
	ResourceType types.String `tfsdk:"resource_type" json:"resource_type,computed"`
}

type CloudReservedFixedIPNetworkModel struct {
	ID                  types.String                                                       `tfsdk:"id" json:"id,computed"`
	CreatedAt           timetypes.RFC3339                                                  `tfsdk:"created_at" json:"created_at,computed" format:"date-time"`
	CreatorTaskID       types.String                                                       `tfsdk:"creator_task_id" json:"creator_task_id,computed"`
	Default             types.Bool                                                         `tfsdk:"default" json:"default,computed"`
	External            types.Bool                                                         `tfsdk:"external" json:"external,computed"`
	Mtu                 types.Int64                                                        `tfsdk:"mtu" json:"mtu,computed"`
	Name                types.String                                                       `tfsdk:"name" json:"name,computed"`
	PortSecurityEnabled types.Bool                                                         `tfsdk:"port_security_enabled" json:"port_security_enabled,computed"`
	ProjectID           types.Int64                                                        `tfsdk:"project_id" json:"project_id,computed"`
	Region              types.String                                                       `tfsdk:"region" json:"region,computed"`
	RegionID            types.Int64                                                        `tfsdk:"region_id" json:"region_id,computed"`
	SegmentationID      types.Int64                                                        `tfsdk:"segmentation_id" json:"segmentation_id,computed"`
	Shared              types.Bool                                                         `tfsdk:"shared" json:"shared,computed"`
	Subnets             customfield.List[types.String]                                     `tfsdk:"subnets" json:"subnets,computed"`
	Tags                customfield.NestedObjectList[CloudReservedFixedIPNetworkTagsModel] `tfsdk:"tags" json:"tags,computed"`
	TaskID              types.String                                                       `tfsdk:"task_id" json:"task_id,computed"`
	Type                types.String                                                       `tfsdk:"type" json:"type,computed"`
	UpdatedAt           timetypes.RFC3339                                                  `tfsdk:"updated_at" json:"updated_at,computed" format:"date-time"`
}

type CloudReservedFixedIPNetworkTagsModel struct {
	Key      types.String `tfsdk:"key" json:"key,computed"`
	ReadOnly types.Bool   `tfsdk:"read_only" json:"read_only,computed"`
	Value    types.String `tfsdk:"value" json:"value,computed"`
}

type CloudReservedFixedIPReservationModel struct {
	ResourceID   types.String `tfsdk:"resource_id" json:"resource_id,computed"`
	ResourceType types.String `tfsdk:"resource_type" json:"resource_type,computed"`
	Status       types.String `tfsdk:"status" json:"status,computed"`
}
