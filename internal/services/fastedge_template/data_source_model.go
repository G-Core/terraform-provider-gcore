// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package fastedge_template

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

type FastedgeTemplateDataSourceModel struct {
	ID         types.Int64                                                         `tfsdk:"id" path:"id,required"`
	APIType    types.String                                                        `tfsdk:"api_type" json:"api_type,computed"`
	BinaryID   types.Int64                                                         `tfsdk:"binary_id" json:"binary_id,computed"`
	LongDescr  types.String                                                        `tfsdk:"long_descr" json:"long_descr,computed"`
	Name       types.String                                                        `tfsdk:"name" json:"name,computed"`
	Owned      types.Bool                                                          `tfsdk:"owned" json:"owned,computed"`
	ShortDescr types.String                                                        `tfsdk:"short_descr" json:"short_descr,computed"`
	Params     customfield.NestedObjectList[FastedgeTemplateParamsDataSourceModel] `tfsdk:"params" json:"params,computed"`
}

type FastedgeTemplateParamsDataSourceModel struct {
	DataType  types.String `tfsdk:"data_type" json:"data_type,computed"`
	Mandatory types.Bool   `tfsdk:"mandatory" json:"mandatory,computed"`
	Name      types.String `tfsdk:"name" json:"name,computed"`
	Descr     types.String `tfsdk:"descr" json:"descr,computed"`
}
