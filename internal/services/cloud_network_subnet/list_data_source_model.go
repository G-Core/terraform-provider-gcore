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

type CloudNetworkSubnetsResultsListDataSourceEnvelope struct {
	Results customfield.NestedObjectList[CloudNetworkSubnetsItemsDataSourceModel] `json:"results,computed"`
}

type CloudNetworkSubnetsDataSourceModel struct {
	ProjectID   types.Int64                                                           `tfsdk:"project_id" path:"project_id,optional"`
	RegionID    types.Int64                                                           `tfsdk:"region_id" path:"region_id,optional"`
	NetworkID   types.String                                                          `tfsdk:"network_id" query:"network_id,optional"`
	TagKeyValue types.String                                                          `tfsdk:"tag_key_value" query:"tag_key_value,optional"`
	TagKey      *[]types.String                                                       `tfsdk:"tag_key" query:"tag_key,optional"`
	Limit       types.Int64                                                           `tfsdk:"limit" query:"limit,computed_optional"`
	OrderBy     types.String                                                          `tfsdk:"order_by" query:"order_by,computed_optional"`
	MaxItems    types.Int64                                                           `tfsdk:"max_items"`
	Items       customfield.NestedObjectList[CloudNetworkSubnetsItemsDataSourceModel] `tfsdk:"items"`
}

func (m *CloudNetworkSubnetsDataSourceModel) toListParams(_ context.Context) (params cloud.NetworkSubnetListParams, diags diag.Diagnostics) {
	mTagKey := []string{}
	if m.TagKey != nil {
		for _, item := range *m.TagKey {
			mTagKey = append(mTagKey, item.ValueString())
		}
	}

	params = cloud.NetworkSubnetListParams{
		TagKey: mTagKey,
	}

	if !m.Limit.IsNull() {
		params.Limit = param.NewOpt(m.Limit.ValueInt64())
	}
	if !m.NetworkID.IsNull() {
		params.NetworkID = param.NewOpt(m.NetworkID.ValueString())
	}
	if !m.OrderBy.IsNull() {
		params.OrderBy = cloud.NetworkSubnetListParamsOrderBy(m.OrderBy.ValueString())
	}
	if !m.TagKeyValue.IsNull() {
		params.TagKeyValue = param.NewOpt(m.TagKeyValue.ValueString())
	}
	if !m.ProjectID.IsNull() {
		params.ProjectID = param.NewOpt(m.ProjectID.ValueInt64())
	}
	if !m.RegionID.IsNull() {
		params.RegionID = param.NewOpt(m.RegionID.ValueInt64())
	}

	return
}

type CloudNetworkSubnetsItemsDataSourceModel struct {
	Cidr           types.String                                                               `tfsdk:"cidr" json:"cidr,computed"`
	CreatedAt      timetypes.RFC3339                                                          `tfsdk:"created_at" json:"created_at,computed" format:"date-time"`
	EnableDhcp     types.Bool                                                                 `tfsdk:"enable_dhcp" json:"enable_dhcp,computed"`
	IPVersion      types.Int64                                                                `tfsdk:"ip_version" json:"ip_version,computed"`
	Name           types.String                                                               `tfsdk:"name" json:"name,computed"`
	NetworkID      types.String                                                               `tfsdk:"network_id" json:"network_id,computed"`
	ProjectID      types.Int64                                                                `tfsdk:"project_id" json:"project_id,computed"`
	Region         types.String                                                               `tfsdk:"region" json:"region,computed"`
	RegionID       types.Int64                                                                `tfsdk:"region_id" json:"region_id,computed"`
	Tags           customfield.NestedObjectList[CloudNetworkSubnetsTagsDataSourceModel]       `tfsdk:"tags" json:"tags,computed"`
	UpdatedAt      timetypes.RFC3339                                                          `tfsdk:"updated_at" json:"updated_at,computed" format:"date-time"`
	ID             types.String                                                               `tfsdk:"id" json:"id,computed"`
	AvailableIPs   types.Int64                                                                `tfsdk:"available_ips" json:"available_ips,computed"`
	CreatorTaskID  types.String                                                               `tfsdk:"creator_task_id" json:"creator_task_id,computed"`
	DNSNameservers customfield.List[types.String]                                             `tfsdk:"dns_nameservers" json:"dns_nameservers,computed"`
	GatewayIP      types.String                                                               `tfsdk:"gateway_ip" json:"gateway_ip,computed"`
	HasRouter      types.Bool                                                                 `tfsdk:"has_router" json:"has_router,computed"`
	HostRoutes     customfield.NestedObjectList[CloudNetworkSubnetsHostRoutesDataSourceModel] `tfsdk:"host_routes" json:"host_routes,computed"`
	TaskID         types.String                                                               `tfsdk:"task_id" json:"task_id,computed"`
	TotalIPs       types.Int64                                                                `tfsdk:"total_ips" json:"total_ips,computed"`
}

type CloudNetworkSubnetsTagsDataSourceModel struct {
	Key      types.String `tfsdk:"key" json:"key,computed"`
	ReadOnly types.Bool   `tfsdk:"read_only" json:"read_only,computed"`
	Value    types.String `tfsdk:"value" json:"value,computed"`
}

type CloudNetworkSubnetsHostRoutesDataSourceModel struct {
	Destination types.String `tfsdk:"destination" json:"destination,computed"`
	Nexthop     types.String `tfsdk:"nexthop" json:"nexthop,computed"`
}
