package cloud_load_balancer_test

import (
	"context"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/G-Core/terraform-provider-gcore/internal/services/cloud_load_balancer"
)

// forbiddenResourceAttributes must never appear in the load balancer *resource*
// schema. Most of them are still exposed by the data sources, which is
// intentional.
//
//   - listeners: listeners are managed by gcore_cloud_load_balancer_listener.
//     Keeping them inline on the resource made every plan report an in-place
//     update, because the API returns listener state the resource cannot own.
//   - task_id, tasks: asynchronous job bookkeeping returned by the API. They
//     change on every write and identify a transient task, not the resource.
//   - updated_at: server metadata that changes on every write, so no plan
//     modifier can honestly stabilise it.
//   - ddos_profile: the resource's Read never asks for the DDoS profile, and the
//     with_ddos query parameter that would return it is being retired by the
//     cloud API, so the attribute is permanently null. A permanently null
//     computed attribute is noise, and pinning it with a plan modifier would
//     break the moment real data ever arrived, because the profile carries
//     volatile fields (status, error_description, options, a re-issuable id).
//     The data sources keep their with_ddos option; only the resource loses it.
//
// This list is a tripwire: the deletions live in custom code that is re-applied
// on every regeneration, and a reseal has silently reverted the listeners
// deletion twice already. ddos_profile is the same story: it was deleted in
// 2.0.0-alpha.3 and a later regeneration put it back.
var forbiddenResourceAttributes = []string{
	"ddos_profile",
	"listeners",
	"task_id",
	"tasks",
	"updated_at",
}

// unstableComputedAttributes are computed attributes that are deliberately
// allowed to re-plan as "(known after apply)" on an in-place update, with the
// reason why. Everything else that is computed must be stabilised.
var unstableComputedAttributes = map[string]string{
	"provisioning_status": "pinning the prior lifecycle status produced " +
		"\"Provider produced inconsistent result after apply\" errors, so it stays " +
		"unpinned until a resize/failover test proves a safe modifier",
}

// stabilisingPlanModifiers are the exact Go type names — package qualifier
// stripped, compared case-insensitively — of the plan modifiers that keep a
// computed attribute's prior value in the plan across an in-place update,
// INCLUDING when that prior value is null.
//
// Matching exactly rather than on a substring is deliberate. The substring
// "stateforunknown" also matches useNonNullStateForUnknown, which leaves a null
// prior value unknown and would therefore satisfy the guard while the attribute
// still drifts. That modifier is handled separately, below.
//
// Add an entry here when a new stabilising modifier lands.
var stabilisingPlanModifiers = []string{
	// {bool,int64,list,map,object,set,string}planmodifier.UseStateForUnknown.
	// As of terraform-plugin-framework v1.19.0 this copies a null prior value
	// too: it bails on req.State.Raw.IsNull() (resource creation), not on
	// req.StateValue.IsNull().
	"useStateForUnknownModifier",
	// internal/planmodifiers.{Bool,Int64,String,Set,Object}UseStateForUnknownInclNull
	"boolUseStateForUnknownInclNull",
	"int64UseStateForUnknownInclNull",
	"stringUseStateForUnknownInclNull",
	"setUseStateForUnknownInclNull",
	"objectUseStateForUnknownInclNull",
	// internal/planmodifiers.UseStateForUnknownIncludingNull{String,Object}
	"useStateForUnknownIncludingNullStringModifier",
	"useStateForUnknownIncludingNullObjectModifier",
	// internal/planmodifiers.ObjectPreserveNullState
	"objectPreserveNullState",
	// internal/planmodifiers.UseStateUnlessCountChanges
	"useStateUnlessCountChangesModifier",
}

// nonNullStabilisingPlanModifiers keep the prior value only when it is non-null;
// a null prior value still re-plans as unknown. They are accepted only for the
// attributes listed in neverNullComputedAttributes.
var nonNullStabilisingPlanModifiers = []string{
	// {bool,int64,list,map,object,set,string}planmodifier.UseNonNullStateForUnknown
	"useNonNullStateForUnknown",
}

// neverNullComputedAttributes may rely on a non-null-only stabiliser, because
// their prior value can never be null once the resource exists.
var neverNullComputedAttributes = map[string]string{
	"id": "the API always returns an id, so the prior state value is never null; " +
		"generated code sets stringplanmodifier.UseNonNullStateForUnknown here",
}

// TestCloudLoadBalancerResourceSchemaOmitsServerOnlyAttributes fails if an
// attribute that was deliberately deleted from the resource schema comes back.
func TestCloudLoadBalancerResourceSchemaOmitsServerOnlyAttributes(t *testing.T) {
	t.Parallel()

	attributes := cloud_load_balancer.ResourceSchema(context.TODO()).Attributes

	for _, name := range forbiddenResourceAttributes {
		if _, found := attributes[name]; found {
			t.Errorf("attribute %q is present in the gcore_cloud_load_balancer resource schema; "+
				"it was removed on purpose and must stay out of the resource (data sources keep it)", name)
		}
	}
}

// TestCloudLoadBalancerResourceSchemaStabilisesComputedAttributes fails if a
// computed attribute has no plan modifier holding its prior state, which is what
// makes it flip to "(known after apply)" on an unrelated in-place update.
func TestCloudLoadBalancerResourceSchemaStabilisesComputedAttributes(t *testing.T) {
	t.Parallel()

	attributes := cloud_load_balancer.ResourceSchema(context.TODO()).Attributes

	for name, attribute := range attributes {
		if !attribute.IsComputed() {
			continue
		}

		modifiers := planModifierTypeNames(t, name, attribute)
		stabilised := hasPlanModifier(modifiers, stabilisingPlanModifiers)
		nonNullStabilised := hasPlanModifier(modifiers, nonNullStabilisingPlanModifiers)

		if reason, exempt := unstableComputedAttributes[name]; exempt {
			if stabilised || nonNullStabilised {
				t.Errorf("computed attribute %q now carries a stabilising plan modifier %v, but it is "+
					"listed as deliberately unstable (%s); if the modifier is intended, drop %q from "+
					"unstableComputedAttributes", name, modifiers, reason, name)
			}
			continue
		}

		if reason, neverNull := neverNullComputedAttributes[name]; neverNull {
			if !stabilised && !nonNullStabilised {
				t.Errorf("computed attribute %q has no stabilising plan modifier (%s); plan modifiers "+
					"present: %v, expected one of %v or %v",
					name, reason, modifiers, stabilisingPlanModifiers, nonNullStabilisingPlanModifiers)
			}
			continue
		}

		if !stabilised {
			extra := ""
			if nonNullStabilised {
				extra = "; the modifiers present only preserve a NON-NULL prior value, so a null " +
					"prior value still re-plans as unknown"
			}
			t.Errorf("computed attribute %q has no stabilising plan modifier, so it will re-plan as "+
				"\"(known after apply)\" on every in-place update%s; plan modifiers present: %v, expected "+
				"one named any of %v", name, extra, modifiers, stabilisingPlanModifiers)
		}
	}
}

// stableComputedAttributeNames derives, from the live resource schema, the names
// of the computed attributes whose planned value must stay known across an
// in-place update: every computed attribute except the deliberately unstable
// ones. Deriving it means a newly generated computed attribute is covered
// automatically, instead of rotting out of a hand-maintained list.
//
// Configurable-only attributes (Required or Optional without Computed) are not
// included: Terraform never plans them as unknown from the provider side.
func stableComputedAttributeNames() []string {
	attributes := cloud_load_balancer.ResourceSchema(context.TODO()).Attributes

	names := make([]string, 0, len(attributes))
	for name, attribute := range attributes {
		if !attribute.IsComputed() {
			continue
		}
		if _, exempt := unstableComputedAttributes[name]; exempt {
			continue
		}
		names = append(names, name)
	}
	sort.Strings(names)

	return names
}

// planModifierTypeNames returns the Go type names of the plan modifiers declared
// on a schema attribute. Every attribute type in the resource schema package
// carries a PlanModifiers slice; the test fails loudly if a new one does not, so
// the guard can never silently pass.
func planModifierTypeNames(t *testing.T, name string, attribute schema.Attribute) []string {
	t.Helper()

	value := reflect.ValueOf(attribute)
	if value.Kind() != reflect.Struct {
		t.Fatalf("attribute %q has unexpected kind %s (%T)", name, value.Kind(), attribute)
	}

	field := value.FieldByName("PlanModifiers")
	if !field.IsValid() || field.Kind() != reflect.Slice {
		t.Fatalf("attribute %q (%T) exposes no PlanModifiers slice; update this guard for that attribute type", name, attribute)
	}

	names := make([]string, 0, field.Len())
	for i := 0; i < field.Len(); i++ {
		names = append(names, reflect.TypeOf(field.Index(i).Interface()).String())
	}
	return names
}

// hasPlanModifier reports whether any of the declared plan modifier type names
// matches one of the wanted names exactly, ignoring the package qualifier and
// letter case. An exact match is what keeps useStateForUnknownModifier and
// useNonNullStateForUnknown apart.
func hasPlanModifier(modifierTypeNames, wanted []string) bool {
	for _, modifier := range modifierTypeNames {
		base := modifier
		if index := strings.LastIndex(base, "."); index >= 0 {
			base = base[index+1:]
		}
		for _, want := range wanted {
			if strings.EqualFold(base, want) {
				return true
			}
		}
	}
	return false
}
