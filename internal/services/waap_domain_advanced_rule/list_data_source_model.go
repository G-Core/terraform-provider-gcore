// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package waap_domain_advanced_rule

import (
	"context"

	"github.com/G-Core/gcore-go/packages/param"
	"github.com/G-Core/gcore-go/waap"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

type WaapDomainAdvancedRulesResultsListDataSourceEnvelope struct {
	Results customfield.NestedObjectList[WaapDomainAdvancedRulesItemsDataSourceModel] `json:"results,computed"`
}

type WaapDomainAdvancedRulesDataSourceModel struct {
	DomainID    types.Int64                                                               `tfsdk:"domain_id" path:"domain_id,required"`
	Action      types.String                                                              `tfsdk:"action" query:"action,optional"`
	Description types.String                                                              `tfsdk:"description" query:"description,optional"`
	Enabled     types.Bool                                                                `tfsdk:"enabled" query:"enabled,optional"`
	Name        types.String                                                              `tfsdk:"name" query:"name,optional"`
	Ordering    types.String                                                              `tfsdk:"ordering" query:"ordering,optional"`
	Phase       types.String                                                              `tfsdk:"phase" query:"phase,optional"`
	Limit       types.Int64                                                               `tfsdk:"limit" query:"limit,computed_optional"`
	MaxItems    types.Int64                                                               `tfsdk:"max_items"`
	Items       customfield.NestedObjectList[WaapDomainAdvancedRulesItemsDataSourceModel] `tfsdk:"items"`
}

func (m *WaapDomainAdvancedRulesDataSourceModel) toListParams(_ context.Context) (params waap.DomainAdvancedRuleListParams, diags diag.Diagnostics) {
	params = waap.DomainAdvancedRuleListParams{}

	if !m.Action.IsNull() {
		params.Action = waap.DomainAdvancedRuleListParamsAction(m.Action.ValueString())
	}
	if !m.Description.IsNull() {
		params.Description = param.NewOpt(m.Description.ValueString())
	}
	if !m.Enabled.IsNull() {
		params.Enabled = param.NewOpt(m.Enabled.ValueBool())
	}
	if !m.Limit.IsNull() {
		params.Limit = param.NewOpt(m.Limit.ValueInt64())
	}
	if !m.Name.IsNull() {
		params.Name = param.NewOpt(m.Name.ValueString())
	}
	if !m.Ordering.IsNull() {
		params.Ordering = waap.DomainAdvancedRuleListParamsOrdering(m.Ordering.ValueString())
	}
	if !m.Phase.IsNull() {
		params.Phase = waap.DomainAdvancedRuleListParamsPhase(m.Phase.ValueString())
	}

	return
}

type WaapDomainAdvancedRulesItemsDataSourceModel struct {
	ID          types.Int64                                                            `tfsdk:"id" json:"id,computed"`
	Action      customfield.NestedObject[WaapDomainAdvancedRulesActionDataSourceModel] `tfsdk:"action" json:"action,computed"`
	Enabled     types.Bool                                                             `tfsdk:"enabled" json:"enabled,computed"`
	Name        types.String                                                           `tfsdk:"name" json:"name,computed"`
	Source      types.String                                                           `tfsdk:"source" json:"source,computed"`
	Description types.String                                                           `tfsdk:"description" json:"description,computed"`
	Phase       types.String                                                           `tfsdk:"phase" json:"phase,computed"`
}

type WaapDomainAdvancedRulesActionDataSourceModel struct {
	Allow     jsontypes.Normalized                                                        `tfsdk:"allow" json:"allow,computed"`
	Block     customfield.NestedObject[WaapDomainAdvancedRulesActionBlockDataSourceModel] `tfsdk:"block" json:"block,computed"`
	Captcha   jsontypes.Normalized                                                        `tfsdk:"captcha" json:"captcha,computed"`
	Handshake jsontypes.Normalized                                                        `tfsdk:"handshake" json:"handshake,computed"`
	Monitor   jsontypes.Normalized                                                        `tfsdk:"monitor" json:"monitor,computed"`
	Tag       customfield.NestedObject[WaapDomainAdvancedRulesActionTagDataSourceModel]   `tfsdk:"tag" json:"tag,computed"`
}

type WaapDomainAdvancedRulesActionBlockDataSourceModel struct {
	ActionDuration types.String `tfsdk:"action_duration" json:"action_duration,computed"`
	StatusCode     types.Int64  `tfsdk:"status_code" json:"status_code,computed"`
}

type WaapDomainAdvancedRulesActionTagDataSourceModel struct {
	Tags customfield.List[types.String] `tfsdk:"tags" json:"tags,computed"`
}
