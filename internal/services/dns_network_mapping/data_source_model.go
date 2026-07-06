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

type DNSNetworkMappingDataSourceModel struct {
	ID        types.Int64                                                           `tfsdk:"id" path:"id,computed_optional"`
	Name      types.String                                                          `tfsdk:"name" json:"name,computed"`
	Mapping   customfield.NestedObjectList[DNSNetworkMappingMappingDataSourceModel] `tfsdk:"mapping" json:"mapping,computed"`
	FindOneBy *DNSNetworkMappingFindOneByDataSourceModel                            `tfsdk:"find_one_by"`
}

func (m *DNSNetworkMappingDataSourceModel) toListParams(_ context.Context) (params dns.NetworkMappingListParams, diags diag.Diagnostics) {
	params = dns.NetworkMappingListParams{}

	if !m.FindOneBy.OrderBy.IsNull() {
		params.OrderBy = param.NewOpt(m.FindOneBy.OrderBy.ValueString())
	}
	if !m.FindOneBy.OrderDirection.IsNull() {
		params.OrderDirection = dns.NetworkMappingListParamsOrderDirection(m.FindOneBy.OrderDirection.ValueString())
	}

	return
}

type DNSNetworkMappingMappingDataSourceModel struct {
	Cidr4 customfield.List[types.String] `tfsdk:"cidr4" json:"cidr4,computed"`
	Cidr6 customfield.List[types.String] `tfsdk:"cidr6" json:"cidr6,computed"`
	Tags  customfield.List[types.String] `tfsdk:"tags" json:"tags,computed"`
}

type DNSNetworkMappingFindOneByDataSourceModel struct {
	OrderBy        types.String `tfsdk:"order_by" query:"order_by,optional"`
	OrderDirection types.String `tfsdk:"order_direction" query:"order_direction,optional"`
}
