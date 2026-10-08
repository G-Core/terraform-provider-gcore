// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package waap_domain_policy

import (
	"context"

	"github.com/G-Core/gcore-go/packages/param"
	"github.com/G-Core/gcore-go/waap"
	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type WaapDomainPoliciesResultsListDataSourceEnvelope struct {
	Results customfield.NestedObjectList[WaapDomainPoliciesItemsDataSourceModel] `json:"results,computed"`
}

type WaapDomainPoliciesDataSourceModel struct {
	DomainID  types.Int64                                                          `tfsdk:"domain_id" path:"domain_id,required"`
	RuleSetID types.Int64                                                          `tfsdk:"rule_set_id" query:"rule_set_id,optional"`
	MaxItems  types.Int64                                                          `tfsdk:"max_items"`
	Items     customfield.NestedObjectList[WaapDomainPoliciesItemsDataSourceModel] `tfsdk:"items"`
}

func (m *WaapDomainPoliciesDataSourceModel) toListParams(_ context.Context) (params waap.DomainPolicyListParams, diags diag.Diagnostics) {
	params = waap.DomainPolicyListParams{}

	if !m.RuleSetID.IsNull() {
		params.RuleSetID = param.NewOpt(m.RuleSetID.ValueInt64())
	}

	return
}

type WaapDomainPoliciesItemsDataSourceModel struct {
	ID          types.String `tfsdk:"id" json:"id,computed"`
	Action      types.String `tfsdk:"action" json:"action,computed"`
	Description types.String `tfsdk:"description" json:"description,computed"`
	Group       types.String `tfsdk:"group" json:"group,computed"`
	Mode        types.Bool   `tfsdk:"mode" json:"mode,computed"`
	Name        types.String `tfsdk:"name" json:"name,computed"`
	RuleSetID   types.Int64  `tfsdk:"rule_set_id" json:"rule_set_id,computed"`
}
