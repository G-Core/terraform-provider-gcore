package cloud_instance

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

// interfaceElement builds a single interfaces list element with all attributes
// null except the given overrides.
func interfaceElement(ctx context.Context, t *testing.T, overrides map[string]tftypes.Value) (tftypes.List, tftypes.Value) {
	t.Helper()
	schema := ResourceSchema(ctx)
	objType := schema.Type().TerraformType(ctx).(tftypes.Object)
	listType := objType.AttributeTypes["interfaces"].(tftypes.List)
	elemType := listType.ElementType.(tftypes.Object)
	vals := map[string]tftypes.Value{}
	for name, attrType := range elemType.AttributeTypes {
		vals[name] = tftypes.NewValue(attrType, nil)
	}
	for name, val := range overrides {
		vals[name] = val
	}
	return listType, tftypes.NewValue(elemType, vals)
}

func runInterfaceTypeValidator(ctx context.Context, config tfsdk.Config) *resource.ValidateConfigResponse {
	resp := &resource.ValidateConfigResponse{}
	(&interfaceTypeValidator{}).ValidateResource(ctx, resource.ValidateConfigRequest{Config: config}, resp)
	return resp
}

// The whole interfaces list is unknown, e.g. built with a for expression over
// a for_each resource that does not exist in state yet. The validator must
// not fail with a Value Conversion Error.
func TestInterfaceTypeValidatorUnknownInterfaces(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	schema := ResourceSchema(ctx)
	objType := schema.Type().TerraformType(ctx).(tftypes.Object)
	config := nullConfigWith(ctx, t, map[string]tftypes.Value{
		"interfaces": tftypes.NewValue(objType.AttributeTypes["interfaces"], tftypes.UnknownValue),
	})

	resp := runInterfaceTypeValidator(ctx, config)
	if resp.Diagnostics.HasError() {
		t.Fatalf("expected no error for unknown interfaces, got: %v", resp.Diagnostics)
	}
}

func TestInterfaceTypeValidatorNullInterfaces(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	resp := runInterfaceTypeValidator(ctx, nullConfigWith(ctx, t, nil))
	if resp.Diagnostics.HasError() {
		t.Fatalf("expected no error for null interfaces, got: %v", resp.Diagnostics)
	}
}

// A reserved_fixed_ip interface with a null port_id must still be rejected.
func TestInterfaceTypeValidatorMissingPortID(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	listType, elem := interfaceElement(ctx, t, map[string]tftypes.Value{
		"type": tftypes.NewValue(tftypes.String, "reserved_fixed_ip"),
	})
	config := nullConfigWith(ctx, t, map[string]tftypes.Value{
		"interfaces": tftypes.NewValue(listType, []tftypes.Value{elem}),
	})

	resp := runInterfaceTypeValidator(ctx, config)
	if !resp.Diagnostics.HasError() {
		t.Fatal("expected an error for reserved_fixed_ip interface without port_id")
	}
}

// An unknown port_id (e.g. referencing a resource not created yet) is valid
// during planning.
func TestInterfaceTypeValidatorUnknownPortID(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	listType, elem := interfaceElement(ctx, t, map[string]tftypes.Value{
		"type":    tftypes.NewValue(tftypes.String, "reserved_fixed_ip"),
		"port_id": tftypes.NewValue(tftypes.String, tftypes.UnknownValue),
	})
	config := nullConfigWith(ctx, t, map[string]tftypes.Value{
		"interfaces": tftypes.NewValue(listType, []tftypes.Value{elem}),
	})

	resp := runInterfaceTypeValidator(ctx, config)
	if resp.Diagnostics.HasError() {
		t.Fatalf("expected no error for unknown port_id, got: %v", resp.Diagnostics)
	}
}

// A known list containing a wholly-unknown element (e.g. a single element
// derived from a not-yet-created resource object via a conditional
// expression) must be skipped, not crash value conversion.
func TestInterfaceTypeValidatorUnknownElement(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	listType, knownElem := interfaceElement(ctx, t, map[string]tftypes.Value{
		"type":    tftypes.NewValue(tftypes.String, "reserved_fixed_ip"),
		"port_id": tftypes.NewValue(tftypes.String, "port-id"),
	})
	unknownElem := tftypes.NewValue(listType.ElementType, tftypes.UnknownValue)
	config := nullConfigWith(ctx, t, map[string]tftypes.Value{
		"interfaces": tftypes.NewValue(listType, []tftypes.Value{unknownElem, knownElem}),
	})

	resp := runInterfaceTypeValidator(ctx, config)
	if resp.Diagnostics.HasError() {
		t.Fatalf("expected no error for a list with an unknown element, got: %v", resp.Diagnostics)
	}
}

// Validation errors for subnet/any_subnet types are preserved.
func TestInterfaceTypeValidatorMissingSubnetAndNetwork(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	listType, subnetElem := interfaceElement(ctx, t, map[string]tftypes.Value{
		"type": tftypes.NewValue(tftypes.String, "subnet"),
	})
	_, anySubnetElem := interfaceElement(ctx, t, map[string]tftypes.Value{
		"type": tftypes.NewValue(tftypes.String, "any_subnet"),
	})
	config := nullConfigWith(ctx, t, map[string]tftypes.Value{
		"interfaces": tftypes.NewValue(listType, []tftypes.Value{subnetElem, anySubnetElem}),
	})

	resp := runInterfaceTypeValidator(ctx, config)
	if got := resp.Diagnostics.ErrorsCount(); got != 2 {
		t.Fatalf("expected 2 errors (missing subnet_id and network_id), got %d: %v", got, resp.Diagnostics)
	}
}
