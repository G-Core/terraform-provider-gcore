// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package waap_domain_firewall_rule

import (
	"context"

	"github.com/G-Core/gcore-go/packages/param"
	"github.com/G-Core/gcore-go/waap"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

type WaapDomainFirewallRulesResultsListDataSourceEnvelope struct {
	Results customfield.NestedObjectList[WaapDomainFirewallRulesItemsDataSourceModel] `json:"results,computed"`
}

type WaapDomainFirewallRulesDataSourceModel struct {
	DomainID    types.Int64                                                               `tfsdk:"domain_id" path:"domain_id,required"`
	Action      types.String                                                              `tfsdk:"action" query:"action,optional"`
	Description types.String                                                              `tfsdk:"description" query:"description,optional"`
	Enabled     types.Bool                                                                `tfsdk:"enabled" query:"enabled,optional"`
	Name        types.String                                                              `tfsdk:"name" query:"name,optional"`
	Ordering    types.String                                                              `tfsdk:"ordering" query:"ordering,optional"`
	MaxItems    types.Int64                                                               `tfsdk:"max_items"`
	Items       customfield.NestedObjectList[WaapDomainFirewallRulesItemsDataSourceModel] `tfsdk:"items"`
}

func (m *WaapDomainFirewallRulesDataSourceModel) toListParams(_ context.Context) (params waap.DomainFirewallRuleListParams, diags diag.Diagnostics) {
	params = waap.DomainFirewallRuleListParams{}

	if !m.Action.IsNull() {
		params.Action = waap.DomainFirewallRuleListParamsAction(m.Action.ValueString())
	}
	if !m.Description.IsNull() {
		params.Description = param.NewOpt(m.Description.ValueString())
	}
	if !m.Enabled.IsNull() {
		params.Enabled = param.NewOpt(m.Enabled.ValueBool())
	}
	if !m.Name.IsNull() {
		params.Name = param.NewOpt(m.Name.ValueString())
	}
	if !m.Ordering.IsNull() {
		params.Ordering = waap.DomainFirewallRuleListParamsOrdering(m.Ordering.ValueString())
	}

	return
}

type WaapDomainFirewallRulesItemsDataSourceModel struct {
	ID          types.Int64                                                                    `tfsdk:"id" json:"id,computed"`
	Action      customfield.NestedObject[WaapDomainFirewallRulesActionDataSourceModel]         `tfsdk:"action" json:"action,computed"`
	Conditions  customfield.NestedObjectList[WaapDomainFirewallRulesConditionsDataSourceModel] `tfsdk:"conditions" json:"conditions,computed"`
	Enabled     types.Bool                                                                     `tfsdk:"enabled" json:"enabled,computed"`
	Name        types.String                                                                   `tfsdk:"name" json:"name,computed"`
	Description types.String                                                                   `tfsdk:"description" json:"description,computed"`
}

type WaapDomainFirewallRulesActionDataSourceModel struct {
	Allow jsontypes.Normalized                                                        `tfsdk:"allow" json:"allow,computed"`
	Block customfield.NestedObject[WaapDomainFirewallRulesActionBlockDataSourceModel] `tfsdk:"block" json:"block,computed"`
}

type WaapDomainFirewallRulesActionBlockDataSourceModel struct {
	ActionDuration types.String `tfsdk:"action_duration" json:"action_duration,computed"`
	StatusCode     types.Int64  `tfsdk:"status_code" json:"status_code,computed"`
}

type WaapDomainFirewallRulesConditionsDataSourceModel struct {
	IP      customfield.NestedObject[WaapDomainFirewallRulesConditionsIPDataSourceModel]      `tfsdk:"ip" json:"ip,computed"`
	IPRange customfield.NestedObject[WaapDomainFirewallRulesConditionsIPRangeDataSourceModel] `tfsdk:"ip_range" json:"ip_range,computed"`
}

type WaapDomainFirewallRulesConditionsIPDataSourceModel struct {
	IPAddress types.String `tfsdk:"ip_address" json:"ip_address,computed"`
	Negation  types.Bool   `tfsdk:"negation" json:"negation,computed"`
}

type WaapDomainFirewallRulesConditionsIPRangeDataSourceModel struct {
	LowerBound types.String `tfsdk:"lower_bound" json:"lower_bound,computed"`
	UpperBound types.String `tfsdk:"upper_bound" json:"upper_bound,computed"`
	Negation   types.Bool   `tfsdk:"negation" json:"negation,computed"`
}
