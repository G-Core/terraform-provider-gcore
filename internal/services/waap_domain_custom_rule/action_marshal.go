package waap_domain_custom_rule

import (
	"github.com/G-Core/terraform-provider-gcore/internal/apijson"
)

// actionBody has the action's fields but not its methods, so encoding it does
// not recurse into MarshalJSONWithState.
type actionBody WaapDomainCustomRuleActionModel

// MarshalJSONWithState encodes the action exactly as on create. The API
// replaces the action as a whole and rejects null for a kind that is not set,
// so the default update encoding (block -> monitor sends "block": null) fails
// with 422.
func (m WaapDomainCustomRuleActionModel) MarshalJSONWithState(_ any, _ any) ([]byte, error) {
	return apijson.MarshalRoot(actionBody(m))
}
