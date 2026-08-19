package cloud_load_balancer_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/G-Core/terraform-provider-gcore/internal/services/cloud_load_balancer"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

// TestCloudLoadBalancerImportSafePlanModifiers pins the import-safe plan
// modifiers on the create-only attributes the API never returns in GET
// responses (`no_refresh`): vip_network_id and vip_subnet_id.
//
// This is a narrow regression pin on two named attributes, not a schema policy
// test: it asserts nothing about computed attributes in general and adds no
// convention the rest of the repo has to follow. It exists because the schema.go
// change is a modification of generated lines, which is the shape most likely to
// be reverted by a conflicting reseal of the sealed custom code - and the revert
// is silent, restoring the built-in RequiresReplace() and with it a plan that
// destroys and recreates a live load balancer after `terraform import`.
//
// The assertion is deliberately positive-only. An earlier version also scanned
// for the built-in modifier by its upstream unexported type name, which would
// have degraded into a silent no-op if the framework ever renamed it. Requiring
// our own modifier to be present fails closed and covers the same regression:
// a reseal that re-emits the built-in necessarily drops ours.
func TestCloudLoadBalancerImportSafePlanModifiers(t *testing.T) {
	t.Parallel()

	rs := cloud_load_balancer.ResourceSchema(context.TODO())

	for _, name := range []string{"vip_network_id", "vip_subnet_id"} {
		attr, ok := rs.Attributes[name].(schema.StringAttribute)
		if !ok {
			t.Fatalf("%s: expected a schema.StringAttribute, got %T", name, rs.Attributes[name])
		}
		assertModifierPresent(t, name,
			modifierTypeNames(attr.PlanModifiers),
			"planmodifiers.stringRequiresReplaceUnlessAdoptingModifier",
			"planmodifiers.StringRequiresReplaceUnlessAdopting(importAdoptionPrivateKey)",
		)
	}
}

func assertModifierPresent(t *testing.T, attrName string, got []string, wantType, wantConstructor string) {
	t.Helper()

	for _, typeName := range got {
		if typeName == wantType {
			return
		}
	}

	t.Errorf("%s: missing %s - `terraform import` will plan a destroy+recreate; plan modifiers are %v",
		attrName, wantConstructor, got)
}

func modifierTypeNames[T any](modifiers []T) []string {
	names := make([]string, 0, len(modifiers))
	for _, m := range modifiers {
		names = append(names, fmt.Sprintf("%T", m))
	}
	return names
}
