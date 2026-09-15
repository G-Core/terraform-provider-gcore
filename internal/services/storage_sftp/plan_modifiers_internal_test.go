package storage_sftp

import (
	"context"
	"reflect"
	"sort"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestReplacementAttributesMirrorSchema(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	resourceSchema := ResourceSchema(ctx)

	// RequiresReplace needs non-null resource values.
	objectType := resourceSchema.Type().TerraformType(ctx).(tftypes.Object)
	attributeValues := map[string]tftypes.Value{}
	for name, attributeType := range objectType.AttributeTypes {
		attributeValues[name] = tftypes.NewValue(attributeType, nil)
	}
	raw := tftypes.NewValue(objectType, attributeValues)
	config := tfsdk.Config{Raw: raw, Schema: resourceSchema}
	state := tfsdk.State{Raw: raw, Schema: resourceSchema}
	plan := tfsdk.Plan{Raw: raw, Schema: resourceSchema}

	var replacing []string
	for name, attribute := range resourceSchema.Attributes {
		requiresReplace := false

		for index, modifier := range attributePlanModifiers(t, name, attribute) {
			attributePath := path.Root(name)

			switch modifier := modifier.(type) {
			case planmodifier.String:
				prior, planned := types.StringValue("prior"), types.StringValue("planned")
				resp := &planmodifier.StringResponse{PlanValue: planned}
				modifier.PlanModifyString(ctx, planmodifier.StringRequest{
					Path:        attributePath,
					Config:      config,
					State:       state,
					Plan:        plan,
					ConfigValue: types.StringNull(),
					StateValue:  prior,
					PlanValue:   planned,
				}, resp)
				requiresReplace = requiresReplace || resp.RequiresReplace
				assertNoErrors(t, name, index, resp.Diagnostics)
			case planmodifier.Int64:
				prior, planned := types.Int64Value(1), types.Int64Value(2)
				resp := &planmodifier.Int64Response{PlanValue: planned}
				modifier.PlanModifyInt64(ctx, planmodifier.Int64Request{
					Path:        attributePath,
					Config:      config,
					State:       state,
					Plan:        plan,
					ConfigValue: types.Int64Null(),
					StateValue:  prior,
					PlanValue:   planned,
				}, resp)
				requiresReplace = requiresReplace || resp.RequiresReplace
				assertNoErrors(t, name, index, resp.Diagnostics)
			case planmodifier.Bool:
				prior, planned := types.BoolValue(false), types.BoolValue(true)
				resp := &planmodifier.BoolResponse{PlanValue: planned}
				modifier.PlanModifyBool(ctx, planmodifier.BoolRequest{
					Path:        attributePath,
					Config:      config,
					State:       state,
					Plan:        plan,
					ConfigValue: types.BoolNull(),
					StateValue:  prior,
					PlanValue:   planned,
				}, resp)
				requiresReplace = requiresReplace || resp.RequiresReplace
				assertNoErrors(t, name, index, resp.Diagnostics)
			default:
				t.Fatalf("attribute %q plan modifier %d has an unhandled type %T; extend this test with a case for it", name, index, modifier)
			}
		}

		if !requiresReplace {
			continue
		}
		replacing = append(replacing, name)

		if _, ok := attribute.(schema.StringAttribute); !ok {
			t.Errorf("attribute %q sets RequiresReplace but is a %T, not a schema.StringAttribute; "+
				"replacementPlanned reads it as types.String, which would add a type-mismatch "+
				"diagnostic and report no replacement", name, attribute)
		}
	}

	want := append([]string(nil), replacementAttributes...)
	sort.Strings(want)
	sort.Strings(replacing)

	if !reflect.DeepEqual(replacing, want) {
		t.Errorf("the schema's RequiresReplace attributes are %v, replacementAttributes is %v; "+
			"update replacementAttributes in plan_modifiers.go", replacing, want)
	}
}

// attributePlanModifiers reads typed PlanModifiers fields in one form.
func attributePlanModifiers(t *testing.T, name string, attribute schema.Attribute) []any {
	t.Helper()

	field := reflect.ValueOf(attribute).FieldByName("PlanModifiers")
	if !field.IsValid() {
		t.Fatalf("attribute %q (%T) has no PlanModifiers field", name, attribute)
	}

	modifiers := make([]any, 0, field.Len())
	for index := 0; index < field.Len(); index++ {
		modifiers = append(modifiers, field.Index(index).Interface())
	}
	return modifiers
}

func assertNoErrors(t *testing.T, name string, index int, diags diag.Diagnostics) {
	t.Helper()

	if diags.HasError() {
		t.Fatalf("attribute %q plan modifier %d errored: %v", name, index, diags.Errors())
	}
}
