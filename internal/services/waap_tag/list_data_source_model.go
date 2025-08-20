// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package waap_tag

import (
	"context"

	"github.com/G-Core/gcore-go/packages/param"
	"github.com/G-Core/gcore-go/waap"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

type WaapTagsResultsListDataSourceEnvelope struct {
	Results customfield.NestedObjectList[WaapTagsItemsDataSourceModel] `json:"results,computed"`
}

type WaapTagsDataSourceModel struct {
	Name         types.String                                               `tfsdk:"name" query:"name,optional"`
	Ordering     types.String                                               `tfsdk:"ordering" query:"ordering,optional"`
	ReadableName types.String                                               `tfsdk:"readable_name" query:"readable_name,optional"`
	Reserved     types.Bool                                                 `tfsdk:"reserved" query:"reserved,optional"`
	Limit        types.Int64                                                `tfsdk:"limit" query:"limit,computed_optional"`
	MaxItems     types.Int64                                                `tfsdk:"max_items"`
	Items        customfield.NestedObjectList[WaapTagsItemsDataSourceModel] `tfsdk:"items"`
}

func (m *WaapTagsDataSourceModel) toListParams(_ context.Context) (params waap.TagListParams, diags diag.Diagnostics) {
	params = waap.TagListParams{}

	if !m.Limit.IsNull() {
		params.Limit = param.NewOpt(m.Limit.ValueInt64())
	}
	if !m.Name.IsNull() {
		params.Name = param.NewOpt(m.Name.ValueString())
	}
	if !m.Ordering.IsNull() {
		params.Ordering = waap.TagListParamsOrdering(m.Ordering.ValueString())
	}
	if !m.ReadableName.IsNull() {
		params.ReadableName = param.NewOpt(m.ReadableName.ValueString())
	}
	if !m.Reserved.IsNull() {
		params.Reserved = param.NewOpt(m.Reserved.ValueBool())
	}

	return
}

type WaapTagsItemsDataSourceModel struct {
	Description  types.String `tfsdk:"description" json:"description,computed"`
	Name         types.String `tfsdk:"name" json:"name,computed"`
	ReadableName types.String `tfsdk:"readable_name" json:"readable_name,computed"`
}
