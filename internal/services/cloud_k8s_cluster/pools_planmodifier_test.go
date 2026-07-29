package cloud_k8s_cluster

import (
	"context"
	"math/big"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// poolsObjectType returns the ObjectType of a pools list element from the
// resource schema, so tests use the exact custom types of the real schema.
func poolsObjectType(t *testing.T, ctx context.Context) basetypes.ObjectType {
	t.Helper()
	s := ResourceSchema(ctx)
	listAttr, ok := s.Attributes["pools"].(schema.ListNestedAttribute)
	if !ok {
		t.Fatalf("pools attribute is not a ListNestedAttribute")
	}
	objType, ok := listAttr.NestedObject.Type().(basetypes.ObjectType)
	if !ok {
		t.Fatalf("pools nested object type is not basetypes.ObjectType")
	}
	return objType
}

// tfString / tfNumber / tfStringMap build tftypes values for pool attributes.
func tfString(v string) tftypes.Value { return tftypes.NewValue(tftypes.String, v) }
func tfNumber(v int64) tftypes.Value {
	return tftypes.NewValue(tftypes.Number, new(big.Float).SetInt64(v))
}
func tfStringMap(m map[string]string) tftypes.Value {
	vals := make(map[string]tftypes.Value, len(m))
	for k, v := range m {
		vals[k] = tfString(v)
	}
	return tftypes.NewValue(tftypes.Map{ElementType: tftypes.String}, vals)
}

// buildPool constructs a pool types.Object with all attributes null except
// the provided overrides (keyed by attribute name, tftypes values).
func buildPool(t *testing.T, ctx context.Context, objType basetypes.ObjectType, set map[string]tftypes.Value) types.Object {
	t.Helper()
	tfObjType := objType.TerraformType(ctx).(tftypes.Object)
	attrVals := make(map[string]tftypes.Value, len(tfObjType.AttributeTypes))
	for name, at := range tfObjType.AttributeTypes {
		if v, ok := set[name]; ok {
			attrVals[name] = v
		} else {
			attrVals[name] = tftypes.NewValue(at, nil)
		}
	}
	raw, err := objType.ValueFromTerraform(ctx, tftypes.NewValue(tfObjType, attrVals))
	if err != nil {
		t.Fatalf("building pool object: %v", err)
	}
	obj, ok := raw.(types.Object)
	if !ok {
		t.Fatalf("expected types.Object, got %T", raw)
	}
	return obj
}

func poolList(t *testing.T, ctx context.Context, objType basetypes.ObjectType, pools ...types.Object) types.List {
	t.Helper()
	elems := make([]attr.Value, len(pools))
	for i, p := range pools {
		elems[i] = p
	}
	list, diags := types.ListValue(objType, elems)
	if diags.HasError() {
		t.Fatalf("building pool list: %v", diags)
	}
	return list
}

func runPoolsModifier(t *testing.T, ctx context.Context, state, plan, config types.List) types.List {
	t.Helper()
	req := planmodifier.ListRequest{
		StateValue:  state,
		PlanValue:   plan,
		ConfigValue: config,
	}
	resp := &planmodifier.ListResponse{PlanValue: plan}
	poolsNormalizeOrderPlanModifier().PlanModifyList(ctx, req, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("plan modifier diagnostics: %v", resp.Diagnostics)
	}
	return resp.PlanValue
}

// planPoolNames extracts pool names from the plan list in order.
func planPoolNames(t *testing.T, list types.List) []string {
	t.Helper()
	var names []string
	for _, elem := range list.Elements() {
		names = append(names, getPoolName(elem.(types.Object)))
	}
	return names
}

// planPoolAttr returns the named attribute of the pool with the given name.
func planPoolAttr(t *testing.T, list types.List, poolName, attrName string) attr.Value {
	t.Helper()
	for _, elem := range list.Elements() {
		obj := elem.(types.Object)
		if getPoolName(obj) == poolName {
			return obj.Attributes()[attrName]
		}
	}
	t.Fatalf("pool %q not found in plan", poolName)
	return nil
}

func assertStringEquals(t *testing.T, v attr.Value, want string, msg string) {
	t.Helper()
	s, ok := v.(types.String)
	if !ok || s.IsNull() || s.IsUnknown() || s.ValueString() != want {
		t.Errorf("%s: got %v, want %q", msg, v, want)
	}
}

// mapValueOf converts a (customfield) map attr.Value to map[string]string,
// failing if null/unknown.
func mapValueOf(t *testing.T, ctx context.Context, v attr.Value, msg string) map[string]string {
	t.Helper()
	if v.IsNull() || v.IsUnknown() {
		t.Fatalf("%s: value is null/unknown: %v", msg, v)
	}
	tfVal, err := v.ToTerraformValue(ctx)
	if err != nil {
		t.Fatalf("%s: %v", msg, err)
	}
	var raw map[string]tftypes.Value
	if err := tfVal.As(&raw); err != nil {
		t.Fatalf("%s: %v", msg, err)
	}
	out := make(map[string]string, len(raw))
	for k, mv := range raw {
		var s string
		if err := mv.As(&s); err != nil {
			t.Fatalf("%s: %v", msg, err)
		}
		out[k] = s
	}
	return out
}

// TestPoolsPlanModifier_PureReorder: config reorders the same pools without
// changing any attribute. The plan must preserve CONFIG order (Terraform
// rejects planned values whose non-computed attributes differ from config by
// index), and omitted computed attributes must be filled from state by name.
func TestPoolsPlanModifier_PureReorder(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	objType := poolsObjectType(t, ctx)

	stateA := buildPool(t, ctx, objType, map[string]tftypes.Value{
		"name": tfString("pool-a"), "flavor_id": tfString("f1"),
		"min_node_count": tfNumber(1), "max_node_count": tfNumber(2),
		"labels": tfStringMap(map[string]string{"env": "prod"}),
		"taints": tfStringMap(map[string]string{}),
	})
	stateB := buildPool(t, ctx, objType, map[string]tftypes.Value{
		"name": tfString("pool-b"), "flavor_id": tfString("f1"),
		"min_node_count": tfNumber(1), "max_node_count": tfNumber(3),
		"labels": tfStringMap(map[string]string{"env": "staging"}),
		"taints": tfStringMap(map[string]string{"dedicated": "true:NoSchedule"}),
	})

	// Config order: B then A. Same configured attrs per name; labels/taints omitted.
	configB := buildPool(t, ctx, objType, map[string]tftypes.Value{
		"name": tfString("pool-b"), "flavor_id": tfString("f1"),
		"min_node_count": tfNumber(1), "max_node_count": tfNumber(3),
	})
	configA := buildPool(t, ctx, objType, map[string]tftypes.Value{
		"name": tfString("pool-a"), "flavor_id": tfString("f1"),
		"min_node_count": tfNumber(1), "max_node_count": tfNumber(2),
	})

	state := poolList(t, ctx, objType, stateA, stateB)
	config := poolList(t, ctx, objType, configB, configA)
	plan := config // proposed plan mirrors config here

	result := runPoolsModifier(t, ctx, state, plan, config)

	names := planPoolNames(t, result)
	if len(names) != 2 || names[0] != "pool-b" || names[1] != "pool-a" {
		t.Fatalf("plan must preserve config order [pool-b pool-a], got %v", names)
	}

	// Omitted computed attrs filled from state BY NAME, not index.
	labelsB := mapValueOf(t, ctx, planPoolAttr(t, result, "pool-b", "labels"), "pool-b labels")
	if labelsB["env"] != "staging" {
		t.Errorf("pool-b labels: got %v, want env=staging (from state by name)", labelsB)
	}
	labelsA := mapValueOf(t, ctx, planPoolAttr(t, result, "pool-a", "labels"), "pool-a labels")
	if labelsA["env"] != "prod" {
		t.Errorf("pool-a labels: got %v, want env=prod (from state by name)", labelsA)
	}
	taintsA := mapValueOf(t, ctx, planPoolAttr(t, result, "pool-a", "taints"), "pool-a taints")
	if len(taintsA) != 0 {
		t.Errorf("pool-a taints: got %v, want empty (from state by name)", taintsA)
	}
	taintsB := mapValueOf(t, ctx, planPoolAttr(t, result, "pool-b", "taints"), "pool-b taints")
	if taintsB["dedicated"] != "true:NoSchedule" {
		t.Errorf("pool-b taints: got %v, want dedicated=true:NoSchedule (from state by name)", taintsB)
	}
}

// TestPoolsPlanModifier_NameSwapChangedAttrs: QA issue 1. The config swaps
// pool names while keeping attributes at their positions. The plan must keep
// config order and config-provided attribute values (interpreted as in-place
// updates of both pools by name), never mixing values across names.
func TestPoolsPlanModifier_NameSwapChangedAttrs(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	objType := poolsObjectType(t, ctx)

	stateA := buildPool(t, ctx, objType, map[string]tftypes.Value{
		"name": tfString("pool-a"), "flavor_id": tfString("f1"),
		"min_node_count": tfNumber(1), "max_node_count": tfNumber(2),
		"labels": tfStringMap(map[string]string{"env": "prod"}),
	})
	stateB := buildPool(t, ctx, objType, map[string]tftypes.Value{
		"name": tfString("pool-b"), "flavor_id": tfString("f1"),
		"min_node_count": tfNumber(1), "max_node_count": tfNumber(3),
		"labels": tfStringMap(map[string]string{"env": "staging"}),
	})

	// Config: names swapped, attributes stay with positions.
	configB := buildPool(t, ctx, objType, map[string]tftypes.Value{
		"name": tfString("pool-b"), "flavor_id": tfString("f1"),
		"min_node_count": tfNumber(1), "max_node_count": tfNumber(2),
		"labels": tfStringMap(map[string]string{"env": "prod"}),
	})
	configA := buildPool(t, ctx, objType, map[string]tftypes.Value{
		"name": tfString("pool-a"), "flavor_id": tfString("f1"),
		"min_node_count": tfNumber(1), "max_node_count": tfNumber(3),
		"labels": tfStringMap(map[string]string{"env": "staging"}),
	})

	state := poolList(t, ctx, objType, stateA, stateB)
	config := poolList(t, ctx, objType, configB, configA)
	plan := config

	result := runPoolsModifier(t, ctx, state, plan, config)

	names := planPoolNames(t, result)
	if len(names) != 2 || names[0] != "pool-b" || names[1] != "pool-a" {
		t.Fatalf("plan must preserve config order [pool-b pool-a], got %v", names)
	}

	// Configured values must remain exactly as configured (plan validity).
	labelsB := mapValueOf(t, ctx, planPoolAttr(t, result, "pool-b", "labels"), "pool-b labels")
	if labelsB["env"] != "prod" {
		t.Errorf("pool-b labels: got %v, want env=prod (config value)", labelsB)
	}
	maxB := planPoolAttr(t, result, "pool-b", "max_node_count").(types.Int64)
	if maxB.ValueInt64() != 2 {
		t.Errorf("pool-b max_node_count: got %v, want 2 (config value)", maxB)
	}
	labelsA := mapValueOf(t, ctx, planPoolAttr(t, result, "pool-a", "labels"), "pool-a labels")
	if labelsA["env"] != "staging" {
		t.Errorf("pool-a labels: got %v, want env=staging (config value)", labelsA)
	}
	maxA := planPoolAttr(t, result, "pool-a", "max_node_count").(types.Int64)
	if maxA.ValueInt64() != 3 {
		t.Errorf("pool-a max_node_count: got %v, want 3 (config value)", maxA)
	}
}

// TestPoolsPlanModifier_InsertMiddleWithUpdate: QA issue 2 (plan-level part).
// A pool is inserted mid-list while an existing pool's attribute changes in
// the same plan. Attributes of existing pools must be correlated by name
// (pool-a must NOT get pool-b's taints), and the new pool's omitted computed
// attributes must be unknown.
func TestPoolsPlanModifier_InsertMiddleWithUpdate(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	objType := poolsObjectType(t, ctx)

	stateA := buildPool(t, ctx, objType, map[string]tftypes.Value{
		"name": tfString("pool-a"), "flavor_id": tfString("f1"),
		"min_node_count": tfNumber(1), "max_node_count": tfNumber(2),
		"labels": tfStringMap(map[string]string{}),
		"taints": tfStringMap(map[string]string{}),
	})
	stateB := buildPool(t, ctx, objType, map[string]tftypes.Value{
		"name": tfString("pool-b"), "flavor_id": tfString("f1"),
		"min_node_count": tfNumber(1), "max_node_count": tfNumber(3),
		"labels": tfStringMap(map[string]string{"env": "staging"}),
		"taints": tfStringMap(map[string]string{"dedicated": "true:NoSchedule"}),
	})

	configA := buildPool(t, ctx, objType, map[string]tftypes.Value{
		"name": tfString("pool-a"), "flavor_id": tfString("f1"),
		"min_node_count": tfNumber(1), "max_node_count": tfNumber(3), // changed
		"labels": tfStringMap(map[string]string{}),
	})
	configM := buildPool(t, ctx, objType, map[string]tftypes.Value{
		"name": tfString("middle-pool"), "flavor_id": tfString("f1"),
		"min_node_count": tfNumber(1),
	})
	configB := buildPool(t, ctx, objType, map[string]tftypes.Value{
		"name": tfString("pool-b"), "flavor_id": tfString("f1"),
		"min_node_count": tfNumber(1), "max_node_count": tfNumber(3),
		"labels": tfStringMap(map[string]string{"env": "staging"}),
		"taints": tfStringMap(map[string]string{"dedicated": "true:NoSchedule"}),
	})

	state := poolList(t, ctx, objType, stateA, stateB)
	config := poolList(t, ctx, objType, configA, configM, configB)
	plan := config

	result := runPoolsModifier(t, ctx, state, plan, config)

	names := planPoolNames(t, result)
	if len(names) != 3 || names[0] != "pool-a" || names[1] != "middle-pool" || names[2] != "pool-b" {
		t.Fatalf("plan must preserve config order [pool-a middle-pool pool-b], got %v", names)
	}

	// pool-a: taints omitted in config -> from state pool-a (empty), NOT pool-b's.
	taintsA := mapValueOf(t, ctx, planPoolAttr(t, result, "pool-a", "taints"), "pool-a taints")
	if len(taintsA) != 0 {
		t.Errorf("pool-a taints: got %v, want empty (must not inherit pool-b's taints)", taintsA)
	}
	maxA := planPoolAttr(t, result, "pool-a", "max_node_count").(types.Int64)
	if maxA.ValueInt64() != 3 {
		t.Errorf("pool-a max_node_count: got %v, want 3 (config value)", maxA)
	}

	// middle-pool: omitted computed attrs must be unknown (API decides).
	for _, attrName := range []string{"labels", "taints", "boot_volume_size"} {
		v := planPoolAttr(t, result, "middle-pool", attrName)
		if !v.IsUnknown() {
			t.Errorf("middle-pool %s: got %v, want unknown", attrName, v)
		}
	}

	// pool-b keeps its own configured taints.
	taintsB := mapValueOf(t, ctx, planPoolAttr(t, result, "pool-b", "taints"), "pool-b taints")
	if taintsB["dedicated"] != "true:NoSchedule" {
		t.Errorf("pool-b taints: got %v, want dedicated=true:NoSchedule", taintsB)
	}
	assertStringEquals(t, planPoolAttr(t, result, "pool-b", "name"), "pool-b", "pool-b name")
}
