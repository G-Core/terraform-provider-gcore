// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package fastedge_template

import (
	"context"

	"github.com/G-Core/gcore-go/fastedge"
	"github.com/G-Core/gcore-go/packages/param"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

type FastedgeTemplatesTemplatesListDataSourceEnvelope struct {
	Templates customfield.NestedObjectList[FastedgeTemplatesItemsDataSourceModel] `json:"templates,computed"`
}

type FastedgeTemplatesDataSourceModel struct {
	APIType  types.String                                                        `tfsdk:"api_type" query:"api_type,optional"`
	Limit    types.Int64                                                         `tfsdk:"limit" query:"limit,optional"`
	OnlyMine types.Bool                                                          `tfsdk:"only_mine" query:"only_mine,computed_optional"`
	MaxItems types.Int64                                                         `tfsdk:"max_items"`
	Items    customfield.NestedObjectList[FastedgeTemplatesItemsDataSourceModel] `tfsdk:"items"`
}

func (m *FastedgeTemplatesDataSourceModel) toListParams(_ context.Context) (params fastedge.TemplateListParams, diags diag.Diagnostics) {
	params = fastedge.TemplateListParams{}

	if !m.APIType.IsNull() {
		params.APIType = fastedge.TemplateListParamsAPIType(m.APIType.ValueString())
	}
	if !m.Limit.IsNull() {
		params.Limit = param.NewOpt(m.Limit.ValueInt64())
	}
	if !m.OnlyMine.IsNull() {
		params.OnlyMine = param.NewOpt(m.OnlyMine.ValueBool())
	}

	return
}

type FastedgeTemplatesItemsDataSourceModel struct {
	ID         types.Int64  `tfsdk:"id" json:"id,computed"`
	APIType    types.String `tfsdk:"api_type" json:"api_type,computed"`
	Name       types.String `tfsdk:"name" json:"name,computed"`
	Owned      types.Bool   `tfsdk:"owned" json:"owned,computed"`
	LongDescr  types.String `tfsdk:"long_descr" json:"long_descr,computed"`
	ShortDescr types.String `tfsdk:"short_descr" json:"short_descr,computed"`
}
