// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package dns_network_mapping

import (
	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type DNSNetworkMappingDataSourceModel struct {
	ID      types.Int64                                                           `tfsdk:"id" path:"id,required"`
	Name    types.String                                                          `tfsdk:"name" json:"name,computed"`
	Mapping customfield.NestedObjectList[DNSNetworkMappingMappingDataSourceModel] `tfsdk:"mapping" json:"mapping,computed"`
}

type DNSNetworkMappingMappingDataSourceModel struct {
	Cidr4 customfield.List[types.String] `tfsdk:"cidr4" json:"cidr4,computed"`
	Cidr6 customfield.List[types.String] `tfsdk:"cidr6" json:"cidr6,computed"`
	Tags  customfield.List[types.String] `tfsdk:"tags" json:"tags,computed"`
}
