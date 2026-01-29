// This file provides custom Terraform plan modifiers specific to this resource.
package cloud_floating_ip

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// UnknownOnPortChange returns a plan modifier for computed attributes that
// depend on port_id. It marks the value as unknown when port_id is changing,
// allowing Terraform to accept whatever value the API returns.
func UnknownOnPortChange() planmodifier.String {
	return unknownOnPortChangeModifier{}
}

type unknownOnPortChangeModifier struct{}

func (m unknownOnPortChangeModifier) Description(_ context.Context) string {
	return "Marks value as unknown when port_id is changing."
}

func (m unknownOnPortChangeModifier) MarkdownDescription(_ context.Context) string {
	return "Marks value as unknown when port_id is changing."
}

func (m unknownOnPortChangeModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	// During Create (no prior state), computed fields should be unknown
	if req.State.Raw.IsNull() {
		resp.PlanValue = types.StringUnknown()
		return
	}

	// Get port_id from config (to avoid plan modifier ordering issues) and state
	var portIDConfig types.String
	var portIDState types.String

	diags := req.Config.GetAttribute(ctx, path.Root("port_id"), &portIDConfig)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	_ = req.State.GetAttribute(ctx, path.Root("port_id"), &portIDState)

	// If port_id is changing, mark as unknown (API will determine value)
	if !portIDConfig.Equal(portIDState) {
		resp.PlanValue = types.StringUnknown()
		return
	}

	// port_id is NOT changing - keep the state value
	resp.PlanValue = req.StateValue
}

// ComputedIfPortSet returns a plan modifier for fixed_ip_address that:
// - Uses the config value when explicitly set
// - Uses unknown when port_id is CHANGING (API will compute new value)
// - Uses state value when port_id is not changing and state has a value
// - Uses null when port_id is null (detaching)
func ComputedIfPortSet() planmodifier.String {
	return computedIfPortSetModifier{}
}

type computedIfPortSetModifier struct{}

func (m computedIfPortSetModifier) Description(_ context.Context) string {
	return "Uses unknown when port_id is changing, preserving state when stable."
}

func (m computedIfPortSetModifier) MarkdownDescription(_ context.Context) string {
	return "Uses unknown when port_id is changing, preserving state when stable."
}

func (m computedIfPortSetModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	// If config has an explicit value, use it
	if !req.ConfigValue.IsNull() && !req.ConfigValue.IsUnknown() {
		resp.PlanValue = req.ConfigValue
		return
	}

	// Get port_id from config (to avoid plan modifier ordering issues) and state
	var portIDConfig types.String
	var portIDState types.String

	diags := req.Config.GetAttribute(ctx, path.Root("port_id"), &portIDConfig)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	diags = req.State.GetAttribute(ctx, path.Root("port_id"), &portIDState)
	if diags.HasError() {
		// State might not exist (new resource) - that's ok
		portIDState = types.StringNull()
	}

	// If port_id is null in config, fixed_ip_address should also be null (detaching)
	if portIDConfig.IsNull() {
		resp.PlanValue = types.StringNull()
		return
	}

	// If port_id is changing (different from state), mark as unknown (API will compute)
	if !portIDConfig.Equal(portIDState) {
		resp.PlanValue = types.StringUnknown()
		return
	}

	// port_id is NOT changing - keep the state value if we have one
	if !req.StateValue.IsNull() {
		resp.PlanValue = req.StateValue
		return
	}

	// No state value and port_id is set - API will compute
	resp.PlanValue = types.StringUnknown()
}
