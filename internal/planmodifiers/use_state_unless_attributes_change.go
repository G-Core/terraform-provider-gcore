package planmodifiers

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// StringUseStateUnlessAttributesChange returns a plan modifier that preserves the
// state value of a computed string attribute unless any of the watched top-level
// attributes change between config and state.
//
// Use this for volatile server-side attributes (e.g. updated_at, status) that
// change whenever the resource is mutated in place: preserving them on a no-op
// plan keeps the plan empty, while leaving them unknown on a real update avoids
// "Provider produced inconsistent result" errors.
//
// During create and replace the prior state is null, so the value is left
// unknown (Terraform core re-plans a replacement create with null prior state).
func StringUseStateUnlessAttributesChange(watchedAttrs ...string) planmodifier.String {
	return useStateUnlessAttributesChangeModifier{watchedAttrs: watchedAttrs}
}

// BoolUseStateUnlessAttributesChange is the bool variant of
// StringUseStateUnlessAttributesChange. See that function for details.
func BoolUseStateUnlessAttributesChange(watchedAttrs ...string) planmodifier.Bool {
	return useStateUnlessAttributesChangeModifier{watchedAttrs: watchedAttrs}
}

type useStateUnlessAttributesChangeModifier struct {
	watchedAttrs []string
}

func (m useStateUnlessAttributesChangeModifier) Description(_ context.Context) string {
	return "Preserves state value unless watched attributes change"
}

func (m useStateUnlessAttributesChangeModifier) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

// watchedAttrsUnchanged reports whether every watched attribute is unchanged.
// An attribute is considered unchanged when its config value is null (the
// framework keeps the prior state value) or when the config value equals the
// state value. Unknown config values are treated as changed (conservative).
func (m useStateUnlessAttributesChangeModifier) watchedAttrsUnchanged(config tfsdk.Config, state tfsdk.State) bool {
	for _, name := range m.watchedAttrs {
		attrPath := tftypes.NewAttributePath().WithAttributeName(name)

		configVal, _, err := tftypes.WalkAttributePath(config.Raw, attrPath)
		if err != nil {
			return false
		}
		cv, ok := configVal.(tftypes.Value)
		if !ok {
			return false
		}
		if cv.IsNull() {
			continue
		}
		if !cv.IsKnown() {
			return false
		}

		stateVal, _, err := tftypes.WalkAttributePath(state.Raw, attrPath)
		if err != nil {
			return false
		}
		sv, ok := stateVal.(tftypes.Value)
		if !ok {
			return false
		}
		if !cv.Equal(sv) {
			return false
		}
	}

	return true
}

func (m useStateUnlessAttributesChangeModifier) PlanModifyString(_ context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	// No state to preserve on create (and on replace, where Terraform core
	// re-plans the replacement create with null prior state).
	if req.State.Raw.IsNull() {
		return
	}

	// Don't override known planned values.
	if !resp.PlanValue.IsUnknown() {
		return
	}

	// Nothing to preserve if the state value was never set.
	if req.StateValue.IsNull() {
		return
	}

	if !m.watchedAttrsUnchanged(req.Config, req.State) {
		return
	}

	resp.PlanValue = req.StateValue
}

func (m useStateUnlessAttributesChangeModifier) PlanModifyBool(_ context.Context, req planmodifier.BoolRequest, resp *planmodifier.BoolResponse) {
	if req.State.Raw.IsNull() {
		return
	}

	if !resp.PlanValue.IsUnknown() {
		return
	}

	if req.StateValue.IsNull() {
		return
	}

	if !m.watchedAttrsUnchanged(req.Config, req.State) {
		return
	}

	resp.PlanValue = req.StateValue
}
