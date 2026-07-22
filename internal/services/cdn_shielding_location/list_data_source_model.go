// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cdn_shielding_location

import (
	"context"

	"github.com/G-Core/gcore-go/cdn"
	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type CDNShieldingLocationsResultsListDataSourceEnvelope struct {
	Results customfield.NestedObjectList[CDNShieldingLocationsItemsDataSourceModel] `json:"results,computed"`
}

type CDNShieldingLocationsDataSourceModel struct {
	MaxItems types.Int64                                                             `tfsdk:"max_items"`
	Items    customfield.NestedObjectList[CDNShieldingLocationsItemsDataSourceModel] `tfsdk:"items"`
}

func (m *CDNShieldingLocationsDataSourceModel) toListParams(_ context.Context) (params cdn.ShieldingLocationListParams, diags diag.Diagnostics) {
	params = cdn.ShieldingLocationListParams{}

	return
}

type CDNShieldingLocationsItemsDataSourceModel struct {
	ID         types.Int64  `tfsdk:"id" json:"id,computed"`
	City       types.String `tfsdk:"city" json:"city,computed"`
	Country    types.String `tfsdk:"country" json:"country,computed"`
	Datacenter types.String `tfsdk:"datacenter" json:"datacenter,computed"`
}
