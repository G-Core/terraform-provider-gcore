// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package waap_organization

import (
	"context"

	"github.com/G-Core/gcore-go/packages/param"
	"github.com/G-Core/gcore-go/waap"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

type WaapOrganizationsResultsListDataSourceEnvelope struct {
	Results customfield.NestedObjectList[WaapOrganizationsItemsDataSourceModel] `json:"results,computed"`
}

type WaapOrganizationsDataSourceModel struct {
	Name     types.String                                                        `tfsdk:"name" query:"name,optional"`
	Ordering types.String                                                        `tfsdk:"ordering" query:"ordering,optional"`
	Limit    types.Int64                                                         `tfsdk:"limit" query:"limit,computed_optional"`
	MaxItems types.Int64                                                         `tfsdk:"max_items"`
	Items    customfield.NestedObjectList[WaapOrganizationsItemsDataSourceModel] `tfsdk:"items"`
}

func (m *WaapOrganizationsDataSourceModel) toListParams(_ context.Context) (params waap.OrganizationListParams, diags diag.Diagnostics) {
	params = waap.OrganizationListParams{}

	if !m.Limit.IsNull() {
		params.Limit = param.NewOpt(m.Limit.ValueInt64())
	}
	if !m.Name.IsNull() {
		params.Name = param.NewOpt(m.Name.ValueString())
	}
	if !m.Ordering.IsNull() {
		params.Ordering = waap.OrganizationListParamsOrdering(m.Ordering.ValueString())
	}

	return
}

type WaapOrganizationsItemsDataSourceModel struct {
	ID   types.Int64  `tfsdk:"id" json:"id,computed"`
	Name types.String `tfsdk:"name" json:"name,computed"`
}
