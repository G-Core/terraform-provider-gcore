// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package waap_domain_policy

import (
	"context"

	"github.com/G-Core/gcore-go/waap"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type WaapDomainPolicyDataSourceModel struct {
	ID       types.String `tfsdk:"id" path:"policy_id,computed"`
	PolicyID types.String `tfsdk:"policy_id" path:"policy_id,required"`
	DomainID types.Int64  `tfsdk:"domain_id" path:"domain_id,required"`
	Mode     types.Bool   `tfsdk:"mode" json:"mode,computed"`
}

func (m *WaapDomainPolicyDataSourceModel) toReadParams(_ context.Context) (params waap.DomainPolicyGetParams, diags diag.Diagnostics) {
	params = waap.DomainPolicyGetParams{
		DomainID: m.DomainID.ValueInt64(),
	}

	return
}
