// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package waap_domain_firewall_rule

import (
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/apijson"
)

type WaapDomainFirewallRuleModel struct {
	ID          types.Int64                               `tfsdk:"id" json:"id,computed"`
	DomainID    types.Int64                               `tfsdk:"domain_id" path:"domain_id,required"`
	Enabled     types.Bool                                `tfsdk:"enabled" json:"enabled,required"`
	Name        types.String                              `tfsdk:"name" json:"name,required"`
	Action      *WaapDomainFirewallRuleActionModel        `tfsdk:"action" json:"action,required"`
	Conditions  *[]*WaapDomainFirewallRuleConditionsModel `tfsdk:"conditions" json:"conditions,required"`
	Description types.String                              `tfsdk:"description" json:"description,optional"`
}

func (m WaapDomainFirewallRuleModel) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(m)
}

func (m WaapDomainFirewallRuleModel) MarshalJSONForUpdate(state WaapDomainFirewallRuleModel) (data []byte, err error) {
	return apijson.MarshalForPatch(m, state)
}

type WaapDomainFirewallRuleActionModel struct {
	Allow jsontypes.Normalized                    `tfsdk:"allow" json:"allow,optional"`
	Block *WaapDomainFirewallRuleActionBlockModel `tfsdk:"block" json:"block,optional"`
}

type WaapDomainFirewallRuleActionBlockModel struct {
	ActionDuration types.String `tfsdk:"action_duration" json:"action_duration,optional"`
	StatusCode     types.Int64  `tfsdk:"status_code" json:"status_code,optional"`
}

type WaapDomainFirewallRuleConditionsModel struct {
	IP      *WaapDomainFirewallRuleConditionsIPModel      `tfsdk:"ip" json:"ip,optional"`
	IPRange *WaapDomainFirewallRuleConditionsIPRangeModel `tfsdk:"ip_range" json:"ip_range,optional"`
}

type WaapDomainFirewallRuleConditionsIPModel struct {
	IPAddress types.String `tfsdk:"ip_address" json:"ip_address,required"`
	Negation  types.Bool   `tfsdk:"negation" json:"negation,computed_optional"`
}

type WaapDomainFirewallRuleConditionsIPRangeModel struct {
	LowerBound types.String `tfsdk:"lower_bound" json:"lower_bound,required"`
	UpperBound types.String `tfsdk:"upper_bound" json:"upper_bound,required"`
	Negation   types.Bool   `tfsdk:"negation" json:"negation,computed_optional"`
}
