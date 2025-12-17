// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package dns_network_mapping

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/apijson"
)

type DNSNetworkMappingModel struct {
	ID      types.Int64                       `tfsdk:"id" json:"id,computed"`
	Name    types.String                      `tfsdk:"name" json:"name,optional"`
	Mapping *[]*DNSNetworkMappingMappingModel `tfsdk:"mapping" json:"mapping,optional"`
}

func (m DNSNetworkMappingModel) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(m)
}

func (m DNSNetworkMappingModel) MarshalJSONForUpdate(state DNSNetworkMappingModel) (data []byte, err error) {
	return apijson.MarshalForUpdate(m, state)
}

type DNSNetworkMappingMappingModel struct {
	Cidr4 *[]types.String `tfsdk:"cidr4" json:"cidr4,optional"`
	Cidr6 *[]types.String `tfsdk:"cidr6" json:"cidr6,optional"`
	Tags  *[]types.String `tfsdk:"tags" json:"tags,optional"`
}
