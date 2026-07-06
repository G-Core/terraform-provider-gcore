// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package dns_network_mapping

import (
	"context"

	"github.com/G-Core/gcore-go/dns"
	"github.com/G-Core/gcore-go/packages/param"
	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type DNSNetworkMappingsNetworkMappingsListDataSourceEnvelope struct {
	NetworkMappings customfield.NestedObjectList[DNSNetworkMappingsItemsDataSourceModel] `json:"network_mappings,computed"`
}

type DNSNetworkMappingsDataSourceModel struct {
	OrderBy        types.String                                                         `tfsdk:"order_by" query:"order_by,optional"`
	OrderDirection types.String                                                         `tfsdk:"order_direction" query:"order_direction,optional"`
	MaxItems       types.Int64                                                          `tfsdk:"max_items"`
	Items          customfield.NestedObjectList[DNSNetworkMappingsItemsDataSourceModel] `tfsdk:"items"`
}

func (m *DNSNetworkMappingsDataSourceModel) toListParams(_ context.Context) (params dns.NetworkMappingListParams, diags diag.Diagnostics) {
	params = dns.NetworkMappingListParams{}

	if !m.OrderBy.IsNull() {
		params.OrderBy = param.NewOpt(m.OrderBy.ValueString())
	}
	if !m.OrderDirection.IsNull() {
		params.OrderDirection = dns.NetworkMappingListParamsOrderDirection(m.OrderDirection.ValueString())
	}

	return
}

type DNSNetworkMappingsItemsDataSourceModel struct {
	ID      types.Int64                                                            `tfsdk:"id" json:"id,computed"`
	Mapping customfield.NestedObjectList[DNSNetworkMappingsMappingDataSourceModel] `tfsdk:"mapping" json:"mapping,computed"`
	Name    types.String                                                           `tfsdk:"name" json:"name,computed"`
}

type DNSNetworkMappingsMappingDataSourceModel struct {
	Cidr4 customfield.List[types.String] `tfsdk:"cidr4" json:"cidr4,computed"`
	Cidr6 customfield.List[types.String] `tfsdk:"cidr6" json:"cidr6,computed"`
	Tags  customfield.List[types.String] `tfsdk:"tags" json:"tags,computed"`
}
