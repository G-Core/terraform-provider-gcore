// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package fastedge_app

import (
	"context"

	"github.com/G-Core/gcore-go/fastedge"
	"github.com/G-Core/gcore-go/packages/param"
	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type FastedgeAppDataSourceModel struct {
	ID           types.Int64                                                    `tfsdk:"id" path:"id,optional"`
	APIType      types.String                                                   `tfsdk:"api_type" json:"api_type,computed"`
	Binary       types.Int64                                                    `tfsdk:"binary" json:"binary,computed"`
	Comment      types.String                                                   `tfsdk:"comment" json:"comment,computed"`
	Debug        types.Bool                                                     `tfsdk:"debug" json:"debug,computed"`
	DebugUntil   timetypes.RFC3339                                              `tfsdk:"debug_until" json:"debug_until,computed" format:"date-time"`
	Log          types.String                                                   `tfsdk:"log" json:"log,computed"`
	Name         types.String                                                   `tfsdk:"name" json:"name,computed"`
	Plan         types.String                                                   `tfsdk:"plan" json:"plan,computed"`
	PlanID       types.Int64                                                    `tfsdk:"plan_id" json:"plan_id,computed"`
	Status       types.Int64                                                    `tfsdk:"status" json:"status,computed"`
	Template     types.Int64                                                    `tfsdk:"template" json:"template,computed"`
	TemplateName types.String                                                   `tfsdk:"template_name" json:"template_name,computed"`
	URL          types.String                                                   `tfsdk:"url" json:"url,computed"`
	Env          customfield.Map[types.String]                                  `tfsdk:"env" json:"env,computed"`
	Networks     customfield.List[types.String]                                 `tfsdk:"networks" json:"networks,computed"`
	RspHeaders   customfield.Map[types.String]                                  `tfsdk:"rsp_headers" json:"rsp_headers,computed"`
	Secrets      customfield.NestedObjectMap[FastedgeAppSecretsDataSourceModel] `tfsdk:"secrets" json:"secrets,computed"`
	Stores       customfield.NestedObjectMap[FastedgeAppStoresDataSourceModel]  `tfsdk:"stores" json:"stores,computed"`
	FindOneBy    *FastedgeAppFindOneByDataSourceModel                           `tfsdk:"find_one_by"`
}

func (m *FastedgeAppDataSourceModel) toListParams(_ context.Context) (params fastedge.AppListParams, diags diag.Diagnostics) {
	params = fastedge.AppListParams{}

	if !m.FindOneBy.APIType.IsNull() {
		params.APIType = fastedge.AppListParamsAPIType(m.FindOneBy.APIType.ValueString())
	}
	if !m.FindOneBy.Binary.IsNull() {
		params.Binary = param.NewOpt(m.FindOneBy.Binary.ValueInt64())
	}
	if !m.FindOneBy.Name.IsNull() {
		params.Name = param.NewOpt(m.FindOneBy.Name.ValueString())
	}
	if !m.FindOneBy.Ordering.IsNull() {
		params.Ordering = fastedge.AppListParamsOrdering(m.FindOneBy.Ordering.ValueString())
	}
	if !m.FindOneBy.Plan.IsNull() {
		params.Plan = param.NewOpt(m.FindOneBy.Plan.ValueInt64())
	}
	if !m.FindOneBy.Status.IsNull() {
		params.Status = param.NewOpt(m.FindOneBy.Status.ValueInt64())
	}
	if !m.FindOneBy.Template.IsNull() {
		params.Template = param.NewOpt(m.FindOneBy.Template.ValueInt64())
	}

	return
}

type FastedgeAppSecretsDataSourceModel struct {
	ID      types.Int64  `tfsdk:"id" json:"id,computed"`
	Comment types.String `tfsdk:"comment" json:"comment,computed"`
	Name    types.String `tfsdk:"name" json:"name,computed"`
}

type FastedgeAppStoresDataSourceModel struct {
	ID      types.Int64  `tfsdk:"id" json:"id,computed"`
	Comment types.String `tfsdk:"comment" json:"comment,computed"`
	Name    types.String `tfsdk:"name" json:"name,computed"`
}

type FastedgeAppFindOneByDataSourceModel struct {
	APIType  types.String `tfsdk:"api_type" query:"api_type,optional"`
	Binary   types.Int64  `tfsdk:"binary" query:"binary,optional"`
	Name     types.String `tfsdk:"name" query:"name,optional"`
	Ordering types.String `tfsdk:"ordering" query:"ordering,optional"`
	Plan     types.Int64  `tfsdk:"plan" query:"plan,optional"`
	Status   types.Int64  `tfsdk:"status" query:"status,optional"`
	Template types.Int64  `tfsdk:"template" query:"template,optional"`
}
