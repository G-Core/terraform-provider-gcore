package planmodifiers_test

import (
	"context"
	"testing"

	"github.com/G-Core/terraform-provider-gcore/internal/planmodifiers"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestRequiresReplaceOnConfigChange(t *testing.T) {
	t.Parallel()

	objAttrTypes := map[string]attr.Type{
		"name":     types.StringType,
		"computed": types.StringType,
		"optional": types.StringType,
	}

	tests := []struct {
		name           string
		stateValue     types.Object
		configValue    types.Object
		planValue      types.Object
		expectsReplace bool
	}{
		{
			name:       "null state - no replacement",
			stateValue: types.ObjectNull(objAttrTypes),
			configValue: types.ObjectValueMust(objAttrTypes, map[string]attr.Value{
				"name":     types.StringValue("test"),
				"computed": types.StringNull(),
				"optional": types.StringNull(),
			}),
			planValue: types.ObjectValueMust(objAttrTypes, map[string]attr.Value{
				"name":     types.StringValue("test"),
				"computed": types.StringUnknown(),
				"optional": types.StringNull(),
			}),
			expectsReplace: false,
		},
		{
			name: "no config change - no replacement",
			stateValue: types.ObjectValueMust(objAttrTypes, map[string]attr.Value{
				"name":     types.StringValue("test"),
				"computed": types.StringValue("computed-val"),
				"optional": types.StringNull(),
			}),
			configValue: types.ObjectValueMust(objAttrTypes, map[string]attr.Value{
				"name":     types.StringValue("test"),
				"computed": types.StringNull(),
				"optional": types.StringNull(),
			}),
			planValue: types.ObjectValueMust(objAttrTypes, map[string]attr.Value{
				"name":     types.StringValue("test"),
				"computed": types.StringUnknown(),
				"optional": types.StringNull(),
			}),
			expectsReplace: false,
		},
		{
			name: "config field changed - requires replacement",
			stateValue: types.ObjectValueMust(objAttrTypes, map[string]attr.Value{
				"name":     types.StringValue("old"),
				"computed": types.StringValue("computed-val"),
				"optional": types.StringNull(),
			}),
			configValue: types.ObjectValueMust(objAttrTypes, map[string]attr.Value{
				"name":     types.StringValue("new"),
				"computed": types.StringNull(),
				"optional": types.StringNull(),
			}),
			planValue: types.ObjectValueMust(objAttrTypes, map[string]attr.Value{
				"name":     types.StringValue("new"),
				"computed": types.StringUnknown(),
				"optional": types.StringNull(),
			}),
			expectsReplace: true,
		},
		{
			name: "only computed field differs - no replacement",
			stateValue: types.ObjectValueMust(objAttrTypes, map[string]attr.Value{
				"name":     types.StringValue("test"),
				"computed": types.StringValue("old-computed"),
				"optional": types.StringNull(),
			}),
			configValue: types.ObjectValueMust(objAttrTypes, map[string]attr.Value{
				"name":     types.StringValue("test"),
				"computed": types.StringNull(),
				"optional": types.StringNull(),
			}),
			planValue: types.ObjectValueMust(objAttrTypes, map[string]attr.Value{
				"name":     types.StringValue("test"),
				"computed": types.StringUnknown(),
				"optional": types.StringNull(),
			}),
			expectsReplace: false,
		},
		{
			name: "null config - no replacement",
			stateValue: types.ObjectValueMust(objAttrTypes, map[string]attr.Value{
				"name":     types.StringValue("test"),
				"computed": types.StringValue("computed-val"),
				"optional": types.StringNull(),
			}),
			configValue:    types.ObjectNull(objAttrTypes),
			planValue:      types.ObjectNull(objAttrTypes),
			expectsReplace: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			req := planmodifier.ObjectRequest{
				StateValue:  tc.stateValue,
				ConfigValue: tc.configValue,
				PlanValue:   tc.planValue,
				Path:        path.Root("servers_settings"),
			}

			resp := &planmodifier.ObjectResponse{
				PlanValue: tc.planValue,
			}

			modifier := planmodifiers.RequiresReplaceOnConfigChange() // no ignored fields
			modifier.PlanModifyObject(context.Background(), req, resp)

			if resp.Diagnostics.HasError() {
				t.Fatalf("unexpected error: %s", resp.Diagnostics.Errors())
			}

			if resp.RequiresReplace != tc.expectsReplace {
				t.Errorf("expected RequiresReplace=%v, got %v", tc.expectsReplace, resp.RequiresReplace)
			}
		})
	}
}

func TestRequiresReplaceOnConfigChange_IgnoreFields(t *testing.T) {
	t.Parallel()

	objAttrTypes := map[string]attr.Type{
		"interfaces": types.StringType,
		"user_data":  types.StringType,
		"credentials": types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"ssh_key_name": types.StringType,
			},
		},
	}

	credentialsType := types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"ssh_key_name": types.StringType,
		},
	}

	tests := []struct {
		name           string
		ignoreFields   []string
		stateValue     types.Object
		configValue    types.Object
		planValue      types.Object
		expectsReplace bool
	}{
		{
			name:         "ignored field changed - no replacement",
			ignoreFields: []string{"user_data", "credentials"},
			stateValue: types.ObjectValueMust(objAttrTypes, map[string]attr.Value{
				"interfaces":  types.StringValue("external"),
				"user_data":   types.StringValue("old-data"),
				"credentials": types.ObjectNull(credentialsType.AttrTypes),
			}),
			configValue: types.ObjectValueMust(objAttrTypes, map[string]attr.Value{
				"interfaces":  types.StringValue("external"),
				"user_data":   types.StringValue("new-data"),
				"credentials": types.ObjectNull(credentialsType.AttrTypes),
			}),
			planValue: types.ObjectValueMust(objAttrTypes, map[string]attr.Value{
				"interfaces":  types.StringValue("external"),
				"user_data":   types.StringValue("new-data"),
				"credentials": types.ObjectNull(credentialsType.AttrTypes),
			}),
			expectsReplace: false,
		},
		{
			name:         "non-ignored field changed - requires replacement",
			ignoreFields: []string{"user_data", "credentials"},
			stateValue: types.ObjectValueMust(objAttrTypes, map[string]attr.Value{
				"interfaces":  types.StringValue("external"),
				"user_data":   types.StringValue("data"),
				"credentials": types.ObjectNull(credentialsType.AttrTypes),
			}),
			configValue: types.ObjectValueMust(objAttrTypes, map[string]attr.Value{
				"interfaces":  types.StringValue("subnet"),
				"user_data":   types.StringValue("data"),
				"credentials": types.ObjectNull(credentialsType.AttrTypes),
			}),
			planValue: types.ObjectValueMust(objAttrTypes, map[string]attr.Value{
				"interfaces":  types.StringValue("subnet"),
				"user_data":   types.StringValue("data"),
				"credentials": types.ObjectNull(credentialsType.AttrTypes),
			}),
			expectsReplace: true,
		},
		{
			name:         "ignored credentials changed - no replacement",
			ignoreFields: []string{"user_data", "credentials"},
			stateValue: types.ObjectValueMust(objAttrTypes, map[string]attr.Value{
				"interfaces": types.StringValue("external"),
				"user_data":  types.StringNull(),
				"credentials": types.ObjectValueMust(credentialsType.AttrTypes, map[string]attr.Value{
					"ssh_key_name": types.StringValue("old-key"),
				}),
			}),
			configValue: types.ObjectValueMust(objAttrTypes, map[string]attr.Value{
				"interfaces": types.StringValue("external"),
				"user_data":  types.StringNull(),
				"credentials": types.ObjectValueMust(credentialsType.AttrTypes, map[string]attr.Value{
					"ssh_key_name": types.StringValue("new-key"),
				}),
			}),
			planValue: types.ObjectValueMust(objAttrTypes, map[string]attr.Value{
				"interfaces": types.StringValue("external"),
				"user_data":  types.StringNull(),
				"credentials": types.ObjectValueMust(credentialsType.AttrTypes, map[string]attr.Value{
					"ssh_key_name": types.StringValue("new-key"),
				}),
			}),
			expectsReplace: false,
		},
		{
			name:         "both ignored and non-ignored changed - requires replacement",
			ignoreFields: []string{"user_data", "credentials"},
			stateValue: types.ObjectValueMust(objAttrTypes, map[string]attr.Value{
				"interfaces":  types.StringValue("external"),
				"user_data":   types.StringValue("old-data"),
				"credentials": types.ObjectNull(credentialsType.AttrTypes),
			}),
			configValue: types.ObjectValueMust(objAttrTypes, map[string]attr.Value{
				"interfaces":  types.StringValue("subnet"),
				"user_data":   types.StringValue("new-data"),
				"credentials": types.ObjectNull(credentialsType.AttrTypes),
			}),
			planValue: types.ObjectValueMust(objAttrTypes, map[string]attr.Value{
				"interfaces":  types.StringValue("subnet"),
				"user_data":   types.StringValue("new-data"),
				"credentials": types.ObjectNull(credentialsType.AttrTypes),
			}),
			expectsReplace: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			req := planmodifier.ObjectRequest{
				StateValue:  tc.stateValue,
				ConfigValue: tc.configValue,
				PlanValue:   tc.planValue,
				Path:        path.Root("servers_settings"),
			}

			resp := &planmodifier.ObjectResponse{
				PlanValue: tc.planValue,
			}

			modifier := planmodifiers.RequiresReplaceOnConfigChange(tc.ignoreFields...)
			modifier.PlanModifyObject(context.Background(), req, resp)

			if resp.Diagnostics.HasError() {
				t.Fatalf("unexpected error: %s", resp.Diagnostics.Errors())
			}

			if resp.RequiresReplace != tc.expectsReplace {
				t.Errorf("expected RequiresReplace=%v, got %v", tc.expectsReplace, resp.RequiresReplace)
			}
		})
	}
}
