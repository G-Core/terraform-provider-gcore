// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package waap_domain_advanced_rule

import (
	"context"

	"github.com/G-Core/gcore-go/waap"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

type WaapDomainAdvancedRuleDataSourceModel struct {
	DomainID    types.Int64                                                           `tfsdk:"domain_id" path:"domain_id,required"`
	RuleID      types.Int64                                                           `tfsdk:"rule_id" path:"rule_id,required"`
	Description types.String                                                          `tfsdk:"description" json:"description,computed"`
	Enabled     types.Bool                                                            `tfsdk:"enabled" json:"enabled,computed"`
	ID          types.Int64                                                           `tfsdk:"id" json:"id,computed"`
	Name        types.String                                                          `tfsdk:"name" json:"name,computed"`
	Phase       types.String                                                          `tfsdk:"phase" json:"phase,computed"`
	Source      types.String                                                          `tfsdk:"source" json:"source,computed"`
	Action      customfield.NestedObject[WaapDomainAdvancedRuleActionDataSourceModel] `tfsdk:"action" json:"action,computed"`
}

func (m *WaapDomainAdvancedRuleDataSourceModel) toReadParams(_ context.Context) (params waap.DomainAdvancedRuleGetParams, diags diag.Diagnostics) {
	params = waap.DomainAdvancedRuleGetParams{
		DomainID: m.DomainID.ValueInt64(),
	}

	return
}

type WaapDomainAdvancedRuleActionDataSourceModel struct {
	Allow     jsontypes.Normalized                                                       `tfsdk:"allow" json:"allow,computed"`
	Block     customfield.NestedObject[WaapDomainAdvancedRuleActionBlockDataSourceModel] `tfsdk:"block" json:"block,computed"`
	Captcha   jsontypes.Normalized                                                       `tfsdk:"captcha" json:"captcha,computed"`
	Handshake jsontypes.Normalized                                                       `tfsdk:"handshake" json:"handshake,computed"`
	Monitor   jsontypes.Normalized                                                       `tfsdk:"monitor" json:"monitor,computed"`
	Tag       customfield.NestedObject[WaapDomainAdvancedRuleActionTagDataSourceModel]   `tfsdk:"tag" json:"tag,computed"`
}

type WaapDomainAdvancedRuleActionBlockDataSourceModel struct {
	ActionDuration types.String `tfsdk:"action_duration" json:"action_duration,computed"`
	StatusCode     types.Int64  `tfsdk:"status_code" json:"status_code,computed"`
}

type WaapDomainAdvancedRuleActionTagDataSourceModel struct {
	Tags customfield.List[types.String] `tfsdk:"tags" json:"tags,computed"`
}
