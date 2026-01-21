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

type WaapDomainFirewallRuleDataSourceModel struct {
	ID          types.Int64                                                                   `tfsdk:"id" path:"rule_id,computed"`
	RuleID      types.Int64                                                                   `tfsdk:"rule_id" path:"rule_id,optional"`
	DomainID    types.Int64                                                                   `tfsdk:"domain_id" path:"domain_id,required"`
	Description types.String                                                                  `tfsdk:"description" json:"description,computed"`
	Enabled     types.Bool                                                                    `tfsdk:"enabled" json:"enabled,computed"`
	Name        types.String                                                                  `tfsdk:"name" json:"name,computed"`
	Action      customfield.NestedObject[WaapDomainFirewallRuleActionDataSourceModel]         `tfsdk:"action" json:"action,computed"`
	Conditions  customfield.NestedObjectList[WaapDomainFirewallRuleConditionsDataSourceModel] `tfsdk:"conditions" json:"conditions,computed"`
	FindOneBy   *WaapDomainFirewallRuleFindOneByDataSourceModel                               `tfsdk:"find_one_by"`
}

func (m *WaapDomainFirewallRuleDataSourceModel) toReadParams(_ context.Context) (params waap.DomainFirewallRuleGetParams, diags diag.Diagnostics) {
	params = waap.DomainFirewallRuleGetParams{
		DomainID: m.DomainID.ValueInt64(),
	}

	return
}

func (m *WaapDomainFirewallRuleDataSourceModel) toListParams(_ context.Context) (params waap.DomainFirewallRuleListParams, diags diag.Diagnostics) {
	params = waap.DomainFirewallRuleListParams{
		DomainID: m.DomainID.ValueInt64(),
	}

	if !m.FindOneBy.Action.IsNull() {
		params.Action = waap.DomainFirewallRuleListParamsAction(m.FindOneBy.Action.ValueString())
	}
	if !m.FindOneBy.Description.IsNull() {
		params.Description = param.NewOpt(m.FindOneBy.Description.ValueString())
	}
	if !m.FindOneBy.Enabled.IsNull() {
		params.Enabled = param.NewOpt(m.FindOneBy.Enabled.ValueBool())
	}
	if !m.FindOneBy.Name.IsNull() {
		params.Name = param.NewOpt(m.FindOneBy.Name.ValueString())
	}
	if !m.FindOneBy.Ordering.IsNull() {
		params.Ordering = waap.DomainFirewallRuleListParamsOrdering(m.FindOneBy.Ordering.ValueString())
	}

	return
}

type WaapDomainFirewallRuleActionDataSourceModel struct {
	Allow jsontypes.Normalized                                                       `tfsdk:"allow" json:"allow,computed"`
	Block customfield.NestedObject[WaapDomainFirewallRuleActionBlockDataSourceModel] `tfsdk:"block" json:"block,computed"`
}

type WaapDomainFirewallRuleActionBlockDataSourceModel struct {
	ActionDuration types.String `tfsdk:"action_duration" json:"action_duration,computed"`
	StatusCode     types.Int64  `tfsdk:"status_code" json:"status_code,computed"`
}

type WaapDomainFirewallRuleConditionsDataSourceModel struct {
	IP      customfield.NestedObject[WaapDomainFirewallRuleConditionsIPDataSourceModel]      `tfsdk:"ip" json:"ip,computed"`
	IPRange customfield.NestedObject[WaapDomainFirewallRuleConditionsIPRangeDataSourceModel] `tfsdk:"ip_range" json:"ip_range,computed"`
}

type WaapDomainFirewallRuleConditionsIPDataSourceModel struct {
	IPAddress types.String `tfsdk:"ip_address" json:"ip_address,computed"`
	Negation  types.Bool   `tfsdk:"negation" json:"negation,computed"`
}

type WaapDomainFirewallRuleConditionsIPRangeDataSourceModel struct {
	LowerBound types.String `tfsdk:"lower_bound" json:"lower_bound,computed"`
	UpperBound types.String `tfsdk:"upper_bound" json:"upper_bound,computed"`
	Negation   types.Bool   `tfsdk:"negation" json:"negation,computed"`
}

type WaapDomainFirewallRuleFindOneByDataSourceModel struct {
	Action      types.String `tfsdk:"action" query:"action,optional"`
	Description types.String `tfsdk:"description" query:"description,optional"`
	Enabled     types.Bool   `tfsdk:"enabled" query:"enabled,optional"`
	Name        types.String `tfsdk:"name" query:"name,optional"`
	Ordering    types.String `tfsdk:"ordering" query:"ordering,optional"`
}
