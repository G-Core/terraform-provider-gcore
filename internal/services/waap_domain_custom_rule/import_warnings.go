package waap_domain_custom_rule

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/tidwall/gjson"
)

// responseOnlyKinds are condition types the API still returns on rules created
// before it stopped accepting them, but rejects on create and update. The
// resource has no attribute for them, so such a condition imports with every
// type null.
var responseOnlyKinds = []string{"response_header", "response_header_exists"}

// responseOnlyConditionWarnings warns on import about conditions the resource
// cannot represent. Every update sends the full condition list, so the next
// apply that changes the rule removes them. Only import needs this: no rule
// created through the current API can carry these conditions.
func responseOnlyConditionWarnings(raw []byte) diag.Diagnostics {
	var diags diag.Diagnostics
	for i, condition := range gjson.GetBytes(raw, "conditions").Array() {
		for _, kind := range responseOnlyKinds {
			if !condition.Get(kind).Exists() {
				continue
			}
			diags.AddAttributeWarning(
				path.Root("conditions").AtListIndex(i),
				"Condition cannot be managed by Terraform",
				fmt.Sprintf("Condition %d of this rule uses %q, which the API no longer accepts on create or update. "+
					"Terraform cannot represent it, so it is imported as an empty condition. The next apply that "+
					"changes this rule removes it.", i, kind),
			)
		}
	}
	return diags
}
