package planmodifiers

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// StringRequiresReplaceUnlessMemberOf returns a plan modifier for a
// Computed+Optional parent-id attribute the API reports only through a
// relationship list (cloud_load_balancer_pool listener_id <- listeners[].id,
// load_balancer_id <- loadbalancers[].id). The back-fill after every decode
// can only pick one element of that list; a configuration naming another
// element of the same list still names a true parent and must not plan a
// replacement.
//
// Decision, per plan:
//   - creating or destroying: nothing.
//   - configured value null: nothing (the attribute keeps prior state through
//     the UseNonNullStateForUnknown modifier before it).
//   - configured value unknown (the parent is being replaced): replacement,
//     as RequiresReplaceIfConfigured decides - the pool must not outlive the
//     listener it points at.
//   - planned value equals prior state: nothing.
//   - configured value is the idAttribute of an element of prior state's
//     listAttribute: PlanValue = configured value, no replacement - the apply
//     records the true parent the configuration names (no field is sent for
//     it - the PATCH body is empty and changes nothing; the API cannot
//     re-parent a pool).
//   - otherwise: replacement, the framework's RequiresReplaceIfConfigured
//     decision.
//
// A prior value the list does not contain is left alone when the
// configuration restates it; Terraform ignores a replacement for an unchanged
// value. The resource's Read must replace such a value with a listed one.
func StringRequiresReplaceUnlessMemberOf(listAttribute, idAttribute string) planmodifier.String {
	return stringRequiresReplaceUnlessMemberOfModifier{listAttribute: listAttribute, idAttribute: idAttribute}
}

type stringRequiresReplaceUnlessMemberOfModifier struct {
	listAttribute, idAttribute string
}

func (m stringRequiresReplaceUnlessMemberOfModifier) Description(_ context.Context) string {
	return fmt.Sprintf("Requires replacement when the configured value changes, unless it names an "+
		"element of the resource's %s (by %s).", m.listAttribute, m.idAttribute)
}

func (m stringRequiresReplaceUnlessMemberOfModifier) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (m stringRequiresReplaceUnlessMemberOfModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	if req.State.Raw.IsNull() || req.Plan.Raw.IsNull() {
		return
	}
	if req.ConfigValue.IsNull() {
		return
	}
	// Before the unchanged check: a state-preserving modifier earlier in the
	// chain may already have copied prior state over an unknown configured
	// value, and an unknown parent must replace regardless of that.
	if req.ConfigValue.IsUnknown() {
		resp.PlanValue = req.ConfigValue
		resp.RequiresReplace = true
		return
	}
	if req.PlanValue.Equal(req.StateValue) {
		return
	}
	members, diags := listedIDs(ctx, req.State, m.listAttribute, m.idAttribute)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	for _, id := range members {
		if id == req.ConfigValue.ValueString() {
			resp.PlanValue = req.ConfigValue
			return
		}
	}
	resp.RequiresReplace = true
}

// listedIDs reads state.<listAttribute>[].<idAttribute> through the raw
// terraform value, so it works for any nested-object list type (the provider's
// customfield lists included). A null or unknown list yields no members.
func listedIDs(ctx context.Context, from tfsdk.State, listAttribute, idAttribute string) ([]string, diag.Diagnostics) {
	var list attr.Value
	diags := from.GetAttribute(ctx, path.Root(listAttribute), &list)
	if diags.HasError() || list == nil || list.IsNull() || list.IsUnknown() {
		return nil, diags
	}
	raw, err := list.ToTerraformValue(ctx)
	if err != nil {
		diags.AddAttributeError(
			path.Root(listAttribute),
			"Unable to inspect relationship list",
			fmt.Sprintf("Could not convert %s to a Terraform value: %s", listAttribute, err),
		)
		return nil, diags
	}
	if !raw.IsKnown() || raw.IsNull() {
		return nil, diags
	}
	var elems []tftypes.Value
	if err := raw.As(&elems); err != nil {
		diags.AddAttributeError(
			path.Root(listAttribute),
			"Unable to inspect relationship list",
			fmt.Sprintf("Could not convert %s to a list: %s", listAttribute, err),
		)
		return nil, diags
	}
	var out []string
	for i, e := range elems {
		if !e.IsKnown() || e.IsNull() {
			continue
		}
		var fields map[string]tftypes.Value
		if err := e.As(&fields); err != nil {
			diags.AddAttributeError(
				path.Root(listAttribute),
				"Unable to inspect relationship list",
				fmt.Sprintf("Could not convert %s element %d to an object: %s", listAttribute, i, err),
			)
			return nil, diags
		}
		v, ok := fields[idAttribute]
		if !ok {
			diags.AddAttributeError(
				path.Root(listAttribute),
				"Unable to inspect relationship list",
				fmt.Sprintf("The %s element %d has no %s attribute.", listAttribute, i, idAttribute),
			)
			return nil, diags
		}
		if !v.IsKnown() || v.IsNull() {
			continue
		}
		var id string
		if err := v.As(&id); err != nil {
			diags.AddAttributeError(
				path.Root(listAttribute).AtListIndex(i).AtName(idAttribute),
				"Unable to inspect relationship ID",
				fmt.Sprintf("Could not convert %s to a string: %s", idAttribute, err),
			)
			return nil, diags
		}
		out = append(out, id)
	}
	return out, diags
}
