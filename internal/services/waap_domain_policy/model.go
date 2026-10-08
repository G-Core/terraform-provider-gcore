// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package waap_domain_policy

import (
	"github.com/G-Core/terraform-provider-gcore/internal/apijson"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type WaapDomainPolicyModel struct {
	ID       types.String `tfsdk:"id" json:"-,computed"`
	PolicyID types.String `tfsdk:"policy_id" path:"policy_id,required"`
	DomainID types.Int64  `tfsdk:"domain_id" path:"domain_id,required"`
	Mode     types.Bool   `tfsdk:"mode" json:"mode,required"`
}

func (m WaapDomainPolicyModel) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(m)
}

func (m WaapDomainPolicyModel) MarshalJSONForUpdate(state WaapDomainPolicyModel) (data []byte, err error) {
	return apijson.MarshalForPatch(m, state)
}
