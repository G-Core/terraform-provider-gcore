// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package fastedge_template

import (
	"context"

	"github.com/G-Core/gcore-go/fastedge"
	"github.com/G-Core/gcore-go/packages/param"
	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type FastedgeTemplateDataSourceModel struct {
	ID         types.Int64                                                         `tfsdk:"id" path:"template_id,computed"`
	TemplateID types.Int64                                                         `tfsdk:"template_id" path:"template_id,optional"`
	APIType    types.String                                                        `tfsdk:"api_type" json:"api_type,computed"`
	BinaryID   types.Int64                                                         `tfsdk:"binary_id" json:"binary_id,computed"`
	LongDescr  types.String                                                        `tfsdk:"long_descr" json:"long_descr,computed"`
	Name       types.String                                                        `tfsdk:"name" json:"name,computed"`
	Owned      types.Bool                                                          `tfsdk:"owned" json:"owned,computed"`
	ShortDescr types.String                                                        `tfsdk:"short_descr" json:"short_descr,computed"`
	Params     customfield.NestedObjectList[FastedgeTemplateParamsDataSourceModel] `tfsdk:"params" json:"params,computed"`
	FindOneBy  *FastedgeTemplateFindOneByDataSourceModel                           `tfsdk:"find_one_by"`
}

func (m *FastedgeTemplateDataSourceModel) toListParams(_ context.Context) (params fastedge.TemplateListParams, diags diag.Diagnostics) {
	params = fastedge.TemplateListParams{}

	if !m.FindOneBy.APIType.IsNull() {
		params.APIType = fastedge.TemplateListParamsAPIType(m.FindOneBy.APIType.ValueString())
	}
	if !m.FindOneBy.OnlyMine.IsNull() {
		params.OnlyMine = param.NewOpt(m.FindOneBy.OnlyMine.ValueBool())
	}

	return
}

type FastedgeTemplateParamsDataSourceModel struct {
	DataType  types.String `tfsdk:"data_type" json:"data_type,computed"`
	Mandatory types.Bool   `tfsdk:"mandatory" json:"mandatory,computed"`
	Name      types.String `tfsdk:"name" json:"name,computed"`
	Descr     types.String `tfsdk:"descr" json:"descr,computed"`
	Metadata  types.String `tfsdk:"metadata" json:"metadata,computed"`
}

type FastedgeTemplateFindOneByDataSourceModel struct {
	APIType  types.String `tfsdk:"api_type" query:"api_type,optional"`
	OnlyMine types.Bool   `tfsdk:"only_mine" query:"only_mine,computed_optional"`
}
