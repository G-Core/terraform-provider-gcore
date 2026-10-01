package cloud_volume

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// checkCreatePrerequisites mirrors the API's create-time rules at plan time so
// a replacement does not destroy a volume and then 400 on create:
//   - source = "image"    requires image_id      (CreateVolumeFromImageSerializer)
//   - source = "snapshot" requires snapshot_id   (CreateVolumeFromSnapshotSerializer)
//   - attachment_tag      requires instance_id_to_attach_to
//
// It is a ModifyPlan hook, not a ConfigValidator, because the decision needs
// prior state: a 2.0.0-rc.1 state can hold exactly these shapes (source =
// "image" adopted on import without image_id; attachment_tag recorded without
// instance_id_to_attach_to) and a config that merely restates that state must
// keep planning no changes. A violated rule is therefore
//   - an error when the plan creates the volume (no prior state),
//   - an error when the plan replaces it - re-derived per attribute from the
//     rules its schema plan modifier applies (replacementDerived), because the
//     framework hands ModifyPlan an empty RequiresReplace list - or when one
//     of the two attributes of the rule itself changed (that change is what
//     would trigger the replacement, or adopts an uncreatable configuration
//     after import),
//   - a warning otherwise, so the poisoned shape stays visible on every plan.
//
// An unknown planned value on the rule's own attributes is left to the API
// (same stance as RequireSizeUnlessDerived); on the replace-bearing attributes
// it counts as a change exactly when it is configured, see attrChanged.
// Deliberately NOT checked: "new-volume forbids image_id/snapshot_id" - the API
// ignores them today and existing configs carry them.
//
// `terraform apply -replace=` is invisible to ModifyPlan (the plan still looks
// like an update), so a poisoned volume replaced that way is destroyed and
// then refused; the warning names it.
func checkCreatePrerequisites(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.Plan.Raw.IsNull() { // destroy
		return
	}
	creating := req.State.Raw.IsNull()

	str := func(from interface {
		GetAttribute(context.Context, path.Path, any) diag.Diagnostics
	}, name string) types.String {
		var v types.String
		if from.GetAttribute(ctx, path.Root(name), &v).HasError() {
			return types.StringUnknown() // unreadable: leave it to the API
		}
		return v
	}

	violations := []prerequisite{}
	source := str(req.Plan, "source")
	for _, rule := range []prerequisite{
		{attr: "source", equals: "image", requires: "image_id"},
		{attr: "source", equals: "snapshot", requires: "snapshot_id"},
		{attr: "attachment_tag", requires: "instance_id_to_attach_to"},
	} {
		have := source
		if rule.attr != "source" {
			have = str(req.Plan, rule.attr)
		}
		if have.IsNull() || have.IsUnknown() {
			continue
		}
		if rule.equals != "" && !strings.EqualFold(have.ValueString(), rule.equals) {
			continue
		}
		if need := str(req.Plan, rule.requires); need.IsNull() {
			violations = append(violations, rule)
		}
	}
	if len(violations) == 0 {
		return
	}

	adopting, diags := importAdoption.Pending(ctx, req.Private)
	resp.Diagnostics.Append(diags...)
	replacing := creating || replacementDerived(ctx, req.Config, req.Plan, req.State, adopting)
	for _, rule := range violations {
		changed := attrChanged(ctx, req.Config, req.Plan, req.State, rule.attr) || attrChanged(ctx, req.Config, req.Plan, req.State, rule.requires)
		if replacing || changed {
			resp.Diagnostics.AddAttributeError(path.Root(rule.requires), "Missing required argument", rule.errorDetail())
			continue
		}
		resp.Diagnostics.AddAttributeWarning(path.Root(rule.requires), "Recorded configuration cannot be re-created", rule.warningDetail())
	}
}

// prerequisite is one create-time rule: when attr (equal to equals, if set) is
// configured, requires must be too.
type prerequisite struct {
	attr, equals, requires string
}

func (p prerequisite) errorDetail() string {
	if p.equals != "" {
		return fmt.Sprintf("%q is required when %s is %q: the API refuses to create the volume "+
			"without it, and a replacement would destroy the volume before that refusal.", p.requires, p.attr, p.equals)
	}
	return fmt.Sprintf("%q is only used together with %q; set both or neither.", p.attr, p.requires)
}

func (p prerequisite) warningDetail() string {
	what := fmt.Sprintf("%s = %q without %s", p.attr, p.equals, p.requires)
	if p.equals == "" {
		what = fmt.Sprintf("%s without %s", p.attr, p.requires)
	}
	return fmt.Sprintf("State records %s. The API refuses to create such a volume, so the "+
		"recorded value was accepted by an earlier provider version without being applied. This "+
		"plan does not replace the volume, but any plan that does - including `terraform apply "+
		"-replace=`, which this check cannot see - would destroy it and then fail to create it. "+
		"Set %q, or remove the resource from state and import it again to correct the record.",
		what, p.requires)
}

// replaceRule names the schema plan modifier whose decision replacementDerived
// re-derives for one attribute. The rules must match schema.go; the wiring test
// pins the modifier types there, TestReplaceBearingMatchesSchema pins this map
// against them.
type replaceRule int

const (
	// replaceUnlessAdopting: planmodifiers.*RequiresReplaceUnlessAdopting -
	// any change replaces, except a null prior value while the import-adoption
	// marker is live.
	replaceUnlessAdopting replaceRule = iota
	// replaceIfConfigured: RequiresReplaceIfConfigured - a change replaces
	// only when the attribute is configured (config non-null); an unconfigured
	// value going to null (project_id/region_id written by ImportState, omitted
	// from the configuration) is not a replacement.
	replaceIfConfigured
	// replaceIfConfiguredCaseInsensitive: replaceIfConfigured behind
	// StringUseStateForCaseInsensitiveEqual - a configured value equal to prior
	// state ignoring case is not a change.
	replaceIfConfiguredCaseInsensitive
	// replaceIfStateNotNull: ListRequiresReplaceIfNotNull - a change replaces
	// only when prior state is not null (a null prior value is adopted in place).
	replaceIfStateNotNull
)

// replaceBearing lists every attribute whose plan modifier can plan a
// replacement, with the rule that modifier applies.
var replaceBearing = map[string]replaceRule{
	"source":                   replaceUnlessAdopting,
	"attachment_tag":           replaceUnlessAdopting,
	"image_id":                 replaceUnlessAdopting,
	"instance_id_to_attach_to": replaceUnlessAdopting,
	"snapshot_id":              replaceUnlessAdopting,
	"project_id":               replaceIfConfigured,
	"region_id":                replaceIfConfigured,
	"type_name":                replaceIfConfiguredCaseInsensitive,
	"lifecycle_policy_ids":     replaceIfStateNotNull,
}

// replacementDerived reports whether any replace-bearing attribute's plan
// modifier would have planned a replacement, re-derived from config, plan and
// prior state with the rules in replaceBearing. adopting is the import-adoption
// marker's presence (importAdoption.Pending).
func replacementDerived(ctx context.Context, config tfsdk.Config, plan tfsdk.Plan, state tfsdk.State, adopting bool) bool {
	for name, rule := range replaceBearing {
		if attrReplaces(ctx, config, plan, state, name, rule, adopting) {
			return true
		}
	}
	return false
}

// attrReplaces is replacementDerived for one attribute.
func attrReplaces(ctx context.Context, config tfsdk.Config, plan tfsdk.Plan, state tfsdk.State, name string, rule replaceRule, adopting bool) bool {
	c, p, s, ok := attrValues(ctx, config, plan, state, name)
	if !ok || !attrChanged(ctx, config, plan, state, name) {
		return false
	}
	switch rule {
	case replaceIfConfigured:
		return !c.IsNull()
	case replaceIfConfiguredCaseInsensitive:
		if c.IsNull() {
			return false
		}
		ps, pok := p.(types.String)
		ss, sok := s.(types.String)
		if pok && sok && !ps.IsNull() && !ps.IsUnknown() && !ss.IsNull() && !ss.IsUnknown() {
			return !strings.EqualFold(ps.ValueString(), ss.ValueString())
		}
		return true
	case replaceIfStateNotNull:
		return !s.IsNull()
	default: // replaceUnlessAdopting
		return !(s.IsNull() && adopting)
	}
}

// attrChanged reports whether the planned value differs from prior state. An
// unknown planned value counts as a change exactly when the attribute is
// configured (config non-null): that is the shape every replace modifier in
// schema.go replaces on (a configured unknown - a value from a resource created
// in the same apply - is a change; an unconfigured Computed attribute left
// unknown by the framework, type_name under -refresh=false on an rc.1 state, is
// not).
func attrChanged(ctx context.Context, config tfsdk.Config, plan tfsdk.Plan, state tfsdk.State, name string) bool {
	c, p, s, ok := attrValues(ctx, config, plan, state, name)
	if !ok {
		return false
	}
	if p.IsUnknown() {
		return !c.IsNull()
	}
	return !p.Equal(s)
}

func attrValues(ctx context.Context, config tfsdk.Config, plan tfsdk.Plan, state tfsdk.State, name string) (c, p, s attr.Value, ok bool) {
	if config.GetAttribute(ctx, path.Root(name), &c).HasError() ||
		plan.GetAttribute(ctx, path.Root(name), &p).HasError() ||
		state.GetAttribute(ctx, path.Root(name), &s).HasError() {
		return nil, nil, nil, false
	}
	if c == nil || p == nil || s == nil {
		return nil, nil, nil, false
	}
	return c, p, s, true
}
