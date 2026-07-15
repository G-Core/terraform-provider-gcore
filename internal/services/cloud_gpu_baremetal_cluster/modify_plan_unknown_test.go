package cloud_gpu_baremetal_cluster

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// rawObjectWith builds a raw resource object value where every top-level
// attribute is null except the given overrides.
func rawObjectWith(ctx context.Context, t *testing.T, overrides map[string]tftypes.Value) tftypes.Value {
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
	return tftypes.NewValue(objType, vals)
}

// serversSettingsWithCredentials builds a servers_settings object whose
// credentials attributes are null except the given overrides.
func serversSettingsWithCredentials(ctx context.Context, t *testing.T, credOverrides map[string]tftypes.Value) tftypes.Value {
	t.Helper()
	schema := ResourceSchema(ctx)
	objType := schema.Type().TerraformType(ctx).(tftypes.Object)
	settingsType := objType.AttributeTypes["servers_settings"].(tftypes.Object)
	vals := map[string]tftypes.Value{}
	for name, attrType := range settingsType.AttributeTypes {
		vals[name] = tftypes.NewValue(attrType, nil)
	}
	credType := settingsType.AttributeTypes["credentials"].(tftypes.Object)
	credVals := map[string]tftypes.Value{}
	for name, attrType := range credType.AttributeTypes {
		credVals[name] = tftypes.NewValue(attrType, nil)
	}
	for name, val := range credOverrides {
		credVals[name] = val
	}
	vals["credentials"] = tftypes.NewValue(credType, credVals)
	return tftypes.NewValue(settingsType, vals)
}

func runModifyPlan(ctx context.Context, t *testing.T, planRaw, stateRaw tftypes.Value) *resource.ModifyPlanResponse {
	t.Helper()
	schema := ResourceSchema(ctx)
	resp := &resource.ModifyPlanResponse{}
	(&CloudGPUBaremetalClusterResource{}).ModifyPlan(ctx, resource.ModifyPlanRequest{
		Plan:  tfsdk.Plan{Raw: planRaw, Schema: schema},
		State: tfsdk.State{Raw: stateRaw, Schema: schema},
	}, resp)
	return resp
}

// The whole servers_settings object is unknown in the plan (e.g. derived from
// resources that do not exist in state yet). ModifyPlan must not fail with a
// Value Conversion Error.
func TestModifyPlanUnknownServersSettings(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	schema := ResourceSchema(ctx)
	objType := schema.Type().TerraformType(ctx).(tftypes.Object)
	planRaw := rawObjectWith(ctx, t, map[string]tftypes.Value{
		"servers_settings": tftypes.NewValue(objType.AttributeTypes["servers_settings"], tftypes.UnknownValue),
	})
	stateRaw := rawObjectWith(ctx, t, map[string]tftypes.Value{
		"servers_settings": serversSettingsWithCredentials(ctx, t, map[string]tftypes.Value{
			"username": tftypes.NewValue(tftypes.String, "admin"),
		}),
	})

	resp := runModifyPlan(ctx, t, planRaw, stateRaw)
	if resp.Diagnostics.HasError() {
		t.Fatalf("expected no error for unknown servers_settings, got: %v", resp.Diagnostics)
	}
}

// Changing username between state and plan must still be rejected.
func TestModifyPlanUsernameChangeRejected(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	planRaw := rawObjectWith(ctx, t, map[string]tftypes.Value{
		"servers_settings": serversSettingsWithCredentials(ctx, t, map[string]tftypes.Value{
			"username": tftypes.NewValue(tftypes.String, "new-user"),
		}),
	})
	stateRaw := rawObjectWith(ctx, t, map[string]tftypes.Value{
		"servers_settings": serversSettingsWithCredentials(ctx, t, map[string]tftypes.Value{
			"username": tftypes.NewValue(tftypes.String, "old-user"),
		}),
	})

	resp := runModifyPlan(ctx, t, planRaw, stateRaw)
	if !resp.Diagnostics.HasError() {
		t.Fatal("expected an error when username changes between state and plan")
	}
}

// When credentials were unknown at plan time the ModifyPlan guard is
// skipped; Update must re-check with the resolved values and reject an
// unsupported username change instead of silently dropping it.
func TestUpdateUsernameChangeRejected(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	schema := ResourceSchema(ctx)
	planRaw := rawObjectWith(ctx, t, map[string]tftypes.Value{
		"servers_settings": serversSettingsWithCredentials(ctx, t, map[string]tftypes.Value{
			"username": tftypes.NewValue(tftypes.String, "new-user"),
		}),
	})
	stateRaw := rawObjectWith(ctx, t, map[string]tftypes.Value{
		"servers_settings": serversSettingsWithCredentials(ctx, t, map[string]tftypes.Value{
			"username": tftypes.NewValue(tftypes.String, "old-user"),
		}),
	})

	resp := &resource.UpdateResponse{State: tfsdk.State{Raw: stateRaw, Schema: schema}}
	(&CloudGPUBaremetalClusterResource{}).Update(ctx, resource.UpdateRequest{
		Plan:  tfsdk.Plan{Raw: planRaw, Schema: schema},
		State: tfsdk.State{Raw: stateRaw, Schema: schema},
	}, resp)
	if !resp.Diagnostics.HasError() {
		t.Fatal("expected Update to reject a resolved username change")
	}
}

// Changing only ssh_key_name is supported and must not be rejected.
func TestModifyPlanSSHKeyChangeAllowed(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	planRaw := rawObjectWith(ctx, t, map[string]tftypes.Value{
		"servers_settings": serversSettingsWithCredentials(ctx, t, map[string]tftypes.Value{
			"ssh_key_name": tftypes.NewValue(tftypes.String, "new-key"),
		}),
	})
	stateRaw := rawObjectWith(ctx, t, map[string]tftypes.Value{
		"servers_settings": serversSettingsWithCredentials(ctx, t, map[string]tftypes.Value{
			"ssh_key_name": tftypes.NewValue(tftypes.String, "old-key"),
		}),
	})

	resp := runModifyPlan(ctx, t, planRaw, stateRaw)
	if resp.Diagnostics.HasError() {
		t.Fatalf("expected no error for ssh_key_name-only change, got: %v", resp.Diagnostics)
	}
}
