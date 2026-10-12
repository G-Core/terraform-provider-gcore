// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package waap_domain_advanced_rule

import (
	"github.com/G-Core/terraform-provider-gcore/internal/apijson"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type WaapDomainAdvancedRuleModel struct {
	ID          types.Int64                        `tfsdk:"id" json:"id,computed"`
	DomainID    types.Int64                        `tfsdk:"domain_id" path:"domain_id,required"`
	Enabled     types.Bool                         `tfsdk:"enabled" json:"enabled,required"`
	Name        types.String                       `tfsdk:"name" json:"name,required"`
	Source      types.String                       `tfsdk:"source" json:"source,required"`
	Action      *WaapDomainAdvancedRuleActionModel `tfsdk:"action" json:"action,required"`
	Description types.String                       `tfsdk:"description" json:"description,optional"`
	Phase       types.String                       `tfsdk:"phase" json:"phase,computed_optional"`
}

func (m WaapDomainAdvancedRuleModel) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(m)
}

func (m WaapDomainAdvancedRuleModel) MarshalJSONForUpdate(state WaapDomainAdvancedRuleModel) (data []byte, err error) {
	return apijson.MarshalForUpdate(m, state)
}

type WaapDomainAdvancedRuleActionModel struct {
	Allow     *map[string]jsontypes.Normalized        `tfsdk:"allow" json:"allow,optional"`
	Block     *WaapDomainAdvancedRuleActionBlockModel `tfsdk:"block" json:"block,optional"`
	Captcha   *map[string]jsontypes.Normalized        `tfsdk:"captcha" json:"captcha,optional"`
	Handshake *map[string]jsontypes.Normalized        `tfsdk:"handshake" json:"handshake,optional"`
	Monitor   *map[string]jsontypes.Normalized        `tfsdk:"monitor" json:"monitor,optional"`
	Tag       *WaapDomainAdvancedRuleActionTagModel   `tfsdk:"tag" json:"tag,optional"`
}

type WaapDomainAdvancedRuleActionBlockModel struct {
	ActionDuration types.String `tfsdk:"action_duration" json:"action_duration,optional"`
	StatusCode     types.Int64  `tfsdk:"status_code" json:"status_code,optional"`
}

type WaapDomainAdvancedRuleActionTagModel struct {
	Tags *[]types.String `tfsdk:"tags" json:"tags,required"`
}
