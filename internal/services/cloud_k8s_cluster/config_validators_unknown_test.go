package cloud_k8s_cluster

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// nullConfigWith builds a resource config where every top-level attribute is
// null except the given overrides.
func nullConfigWith(ctx context.Context, t *testing.T, overrides map[string]tftypes.Value) tfsdk.Config {
	t.Helper()
	schema := ResourceSchema(ctx)
	objType := schema.Type().TerraformType(ctx).(tftypes.Object)
	vals := map[string]tftypes.Value{}
	for name, attrType := range objType.AttributeTypes {
		vals[name] = tftypes.NewValue(attrType, nil)
	}
	for name, val := range overrides {
		vals[name] = val
	}
	return tfsdk.Config{Raw: tftypes.NewValue(objType, vals), Schema: schema}
}

// The whole pools list is unknown, e.g. built with a for expression over a
// for_each resource that does not exist in state yet. None of the resource
// config validators may fail with a Value Conversion Error.
func TestConfigValidatorsUnknownPools(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	schema := ResourceSchema(ctx)
	objType := schema.Type().TerraformType(ctx).(tftypes.Object)
	config := nullConfigWith(ctx, t, map[string]tftypes.Value{
		"pools": tftypes.NewValue(objType.AttributeTypes["pools"], tftypes.UnknownValue),
	})

	req := resource.ValidateConfigRequest{Config: config}
	for _, v := range (&CloudK8SClusterResource{}).ConfigValidators(ctx) {
		resp := &resource.ValidateConfigResponse{}
		v.ValidateResource(ctx, req, resp)
		if resp.Diagnostics.HasError() {
			t.Fatalf("validator %T: expected no error for unknown pools, got: %v", v, resp.Diagnostics)
		}
	}
}

// A known pools list containing a wholly-unknown element must be skipped by
// poolFlavorValidator, not crash value conversion.
func TestPoolFlavorValidatorUnknownPoolElement(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	schema := ResourceSchema(ctx)
	objType := schema.Type().TerraformType(ctx).(tftypes.Object)
	listType := objType.AttributeTypes["pools"].(tftypes.List)
	unknownPool := tftypes.NewValue(listType.ElementType, tftypes.UnknownValue)
	config := nullConfigWith(ctx, t, map[string]tftypes.Value{
		"pools": tftypes.NewValue(listType, []tftypes.Value{unknownPool}),
	})

	resp := &resource.ValidateConfigResponse{}
	(&poolFlavorValidator{}).ValidateResource(ctx, resource.ValidateConfigRequest{Config: config}, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("expected no error for a pools list with an unknown element, got: %v", resp.Diagnostics)
	}
}

// A VM-based pool flavor with null boot volume attributes must still be
// rejected by poolFlavorValidator.
func TestPoolFlavorValidatorVMFlavorMissingAttributes(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	schema := ResourceSchema(ctx)
	objType := schema.Type().TerraformType(ctx).(tftypes.Object)
	listType := objType.AttributeTypes["pools"].(tftypes.List)
	elemType := listType.ElementType.(tftypes.Object)
	vals := map[string]tftypes.Value{}
	for name, attrType := range elemType.AttributeTypes {
		vals[name] = tftypes.NewValue(attrType, nil)
	}
	vals["flavor_id"] = tftypes.NewValue(tftypes.String, "g1-standard-1-2")
	pool := tftypes.NewValue(elemType, vals)

	config := nullConfigWith(ctx, t, map[string]tftypes.Value{
		"pools": tftypes.NewValue(listType, []tftypes.Value{pool}),
	})

	resp := &resource.ValidateConfigResponse{}
	(&poolFlavorValidator{}).ValidateResource(ctx, resource.ValidateConfigRequest{Config: config}, resp)
	if got := resp.Diagnostics.ErrorsCount(); got != 3 {
		t.Fatalf("expected 3 errors (servergroup_policy, boot_volume_size, boot_volume_type), got %d: %v", got, resp.Diagnostics)
	}
}
