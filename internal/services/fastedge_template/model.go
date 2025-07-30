// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package fastedge_template

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/apijson"
)

type FastedgeTemplateModel struct {
	ID         types.Int64                     `tfsdk:"id" json:"id,computed"`
	BinaryID   types.Int64                     `tfsdk:"binary_id" json:"binary_id,required"`
	Name       types.String                    `tfsdk:"name" json:"name,required"`
	Owned      types.Bool                      `tfsdk:"owned" json:"owned,required"`
	Params     *[]*FastedgeTemplateParamsModel `tfsdk:"params" json:"params,required"`
	LongDescr  types.String                    `tfsdk:"long_descr" json:"long_descr,optional"`
	ShortDescr types.String                    `tfsdk:"short_descr" json:"short_descr,optional"`
	APIType    types.String                    `tfsdk:"api_type" json:"api_type,computed"`
}

func (m FastedgeTemplateModel) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(m)
}

func (m FastedgeTemplateModel) MarshalJSONForUpdate(state FastedgeTemplateModel) (data []byte, err error) {
	return apijson.MarshalForUpdate(m, state)
}

type FastedgeTemplateParamsModel struct {
	DataType  types.String `tfsdk:"data_type" json:"data_type,required"`
	Mandatory types.Bool   `tfsdk:"mandatory" json:"mandatory,computed_optional"`
	Name      types.String `tfsdk:"name" json:"name,required"`
	Descr     types.String `tfsdk:"descr" json:"descr,optional"`
}
