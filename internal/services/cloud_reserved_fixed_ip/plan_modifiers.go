// This file provides custom Terraform plan modifiers specific to this resource.
package cloud_reserved_fixed_ip

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

// reservedFixedIPTypeModifier returns a plan modifier for the write-only "type"
// field on reserved fixed IPs. It combines two behaviors:
//
//  1. Import-safe replacement: when the state value is null (e.g., after import
//     when the API didn't return the field), the config value is adopted into
//     the plan without triggering replacement.
//
//  2. Subnet/any_subnet equivalence: the API cannot distinguish resources
//     created with type "subnet" from those created with type "any_subnet"
//     (both return subnet_id and network_id). During import, type is inferred
//     as "subnet" by default. If the user's config says "any_subnet", the
//     modifier silently adopts the config value instead of triggering
//     replacement, because both types produce identical infrastructure.
//
// For all other cases (e.g., changing from "external" to "subnet"), the
// modifier requires resource replacement as expected.
func reservedFixedIPTypeModifier() planmodifier.String {
	return reservedFixedIPTypePlanModifier{}
}

type reservedFixedIPTypePlanModifier struct{}

func (m reservedFixedIPTypePlanModifier) Description(_ context.Context) string {
	return "Requires replacement when the type changes, except for subnet/any_subnet " +
		"equivalence after import. Both types produce identical infrastructure and " +
		"cannot be distinguished from the API response."
}

func (m reservedFixedIPTypePlanModifier) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (m reservedFixedIPTypePlanModifier) PlanModifyString(_ context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	// If config is null or unknown, no action needed.
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	// If state is null or unknown (new resource or post-import with no value),
	// adopt the config value into the plan without triggering replacement.
	if req.StateValue.IsNull() || req.StateValue.IsUnknown() {
		resp.PlanValue = req.ConfigValue
		return
	}

	// Both config and state have known values.
	if req.ConfigValue.Equal(req.StateValue) {
		return
	}

	// Special case: "subnet" and "any_subnet" are functionally equivalent.
	// The API always returns both subnet_id and network_id for non-external
	// IPs, making the two types indistinguishable. During import, we infer
	// "subnet" as the default. If the user's config says "any_subnet" (or
	// vice versa), adopt the config value without requiring replacement.
	configVal := req.ConfigValue.ValueString()
	stateVal := req.StateValue.ValueString()
	if isSubnetType(configVal) && isSubnetType(stateVal) {
		resp.PlanValue = req.ConfigValue
		return
	}

	// For any other type change (e.g., external → subnet), require replacement.
	resp.RequiresReplace = true
}

// isSubnetType returns true if the type is "subnet" or "any_subnet".
func isSubnetType(t string) bool {
	return t == "subnet" || t == "any_subnet"
}
