// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package fastedge_app

import (
	"context"

	"github.com/G-Core/gcore-go/fastedge"
	"github.com/G-Core/gcore-go/packages/param"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

type FastedgeAppsAppsListDataSourceEnvelope struct {
	Apps customfield.NestedObjectList[FastedgeAppsItemsDataSourceModel] `json:"apps,computed"`
}

type FastedgeAppsDataSourceModel struct {
	APIType  types.String                                                   `tfsdk:"api_type" query:"api_type,optional"`
	Binary   types.Int64                                                    `tfsdk:"binary" query:"binary,optional"`
	Limit    types.Int64                                                    `tfsdk:"limit" query:"limit,optional"`
	Name     types.String                                                   `tfsdk:"name" query:"name,optional"`
	Ordering types.String                                                   `tfsdk:"ordering" query:"ordering,optional"`
	Plan     types.Int64                                                    `tfsdk:"plan" query:"plan,optional"`
	Status   types.Int64                                                    `tfsdk:"status" query:"status,optional"`
	Template types.Int64                                                    `tfsdk:"template" query:"template,optional"`
	MaxItems types.Int64                                                    `tfsdk:"max_items"`
	Items    customfield.NestedObjectList[FastedgeAppsItemsDataSourceModel] `tfsdk:"items"`
}

func (m *FastedgeAppsDataSourceModel) toListParams(_ context.Context) (params fastedge.AppListParams, diags diag.Diagnostics) {
	params = fastedge.AppListParams{}

	if !m.APIType.IsNull() {
		params.APIType = fastedge.AppListParamsAPIType(m.APIType.ValueString())
	}
	if !m.Binary.IsNull() {
		params.Binary = param.NewOpt(m.Binary.ValueInt64())
	}
	if !m.Limit.IsNull() {
		params.Limit = param.NewOpt(m.Limit.ValueInt64())
	}
	if !m.Name.IsNull() {
		params.Name = param.NewOpt(m.Name.ValueString())
	}
	if !m.Ordering.IsNull() {
		params.Ordering = fastedge.AppListParamsOrdering(m.Ordering.ValueString())
	}
	if !m.Plan.IsNull() {
		params.Plan = param.NewOpt(m.Plan.ValueInt64())
	}
	if !m.Status.IsNull() {
		params.Status = param.NewOpt(m.Status.ValueInt64())
	}
	if !m.Template.IsNull() {
		params.Template = param.NewOpt(m.Template.ValueInt64())
	}

	return
}

type FastedgeAppsItemsDataSourceModel struct {
	ID            types.Int64                    `tfsdk:"id" json:"id,computed"`
	APIType       types.String                   `tfsdk:"api_type" json:"api_type,computed"`
	Binary        types.Int64                    `tfsdk:"binary" json:"binary,computed"`
	Name          types.String                   `tfsdk:"name" json:"name,computed"`
	PlanID        types.Int64                    `tfsdk:"plan_id" json:"plan_id,computed"`
	Status        types.Int64                    `tfsdk:"status" json:"status,computed"`
	Comment       types.String                   `tfsdk:"comment" json:"comment,computed"`
	Debug         types.Bool                     `tfsdk:"debug" json:"debug,computed"`
	DebugUntil    timetypes.RFC3339              `tfsdk:"debug_until" json:"debug_until,computed" format:"date-time"`
	Networks      customfield.List[types.String] `tfsdk:"networks" json:"networks,computed"`
	Plan          types.String                   `tfsdk:"plan" json:"plan,computed"`
	Template      types.Int64                    `tfsdk:"template" json:"template,computed"`
	TemplateName  types.String                   `tfsdk:"template_name" json:"template_name,computed"`
	UpgradeableTo types.Int64                    `tfsdk:"upgradeable_to" json:"upgradeable_to,computed"`
	URL           types.String                   `tfsdk:"url" json:"url,computed"`
}
