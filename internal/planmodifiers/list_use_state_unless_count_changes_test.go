package planmodifiers_test

import (
	"context"
	"testing"

	"github.com/G-Core/terraform-provider-gcore/internal/planmodifiers"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestUseStateUnlessCountChanges(t *testing.T) {
	t.Parallel()

	tfSchema := schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"servers_count": schema.Int64Attribute{
				Required: true,
			},
			"servers_ids": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
			},
		},
	}

	schemaType := tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"id":            tftypes.String,
			"servers_count": tftypes.Number,
			"servers_ids":   tftypes.List{ElementType: tftypes.String},
		},
	}

	stateValue, _ := types.ListValueFrom(context.Background(), types.StringType, []string{"server-1", "server-2"})

	makeState := func(id string, count int, serverIDs []string) tfsdk.State {
		idVal := tftypes.NewValue(tftypes.String, id)
		countVal := tftypes.NewValue(tftypes.Number, count)
		var idsVal tftypes.Value
		if serverIDs != nil {
			elems := make([]tftypes.Value, len(serverIDs))
			for i, s := range serverIDs {
				elems[i] = tftypes.NewValue(tftypes.String, s)
			}
			idsVal = tftypes.NewValue(tftypes.List{ElementType: tftypes.String}, elems)
		} else {
			idsVal = tftypes.NewValue(tftypes.List{ElementType: tftypes.String}, tftypes.UnknownValue)
		}
		return tfsdk.State{
			Raw: tftypes.NewValue(schemaType, map[string]tftypes.Value{
				"id":            idVal,
				"servers_count": countVal,
				"servers_ids":   idsVal,
			}),
			Schema: tfSchema,
		}
	}

	makePlan := func(id interface{}, count int) tfsdk.Plan {
		var idVal tftypes.Value
		if id == nil {
			idVal = tftypes.NewValue(tftypes.String, tftypes.UnknownValue)
		} else {
			idVal = tftypes.NewValue(tftypes.String, id.(string))
		}
		return tfsdk.Plan{
			Raw: tftypes.NewValue(schemaType, map[string]tftypes.Value{
				"id":            idVal,
				"servers_count": tftypes.NewValue(tftypes.Number, count),
				"servers_ids":   tftypes.NewValue(tftypes.List{ElementType: tftypes.String}, tftypes.UnknownValue),
			}),
			Schema: tfSchema,
		}
	}

	t.Run("new resource - no state - don't preserve", func(t *testing.T) {
		t.Parallel()
		resp := &planmodifier.ListResponse{PlanValue: types.ListUnknown(types.StringType)}
		planmodifiers.UseStateUnlessCountChanges("servers_count").PlanModifyList(
			context.Background(),
			planmodifier.ListRequest{
				StateValue: types.ListNull(types.StringType),
				PlanValue:  types.ListUnknown(types.StringType),
				State:      tfsdk.State{Raw: tftypes.NewValue(schemaType, nil), Schema: tfSchema},
				Plan:       makePlan(nil, 2),
				Path:       path.Root("servers_ids"),
			},
			resp,
		)
		if resp.Diagnostics.HasError() {
			t.Fatalf("unexpected error: %s", resp.Diagnostics.Errors())
		}
		if !resp.PlanValue.IsUnknown() {
			t.Errorf("expected unknown, got %v", resp.PlanValue)
		}
	})

	t.Run("no change in count - preserve state", func(t *testing.T) {
		t.Parallel()
		resp := &planmodifier.ListResponse{PlanValue: types.ListUnknown(types.StringType)}
		planmodifiers.UseStateUnlessCountChanges("servers_count").PlanModifyList(
			context.Background(),
			planmodifier.ListRequest{
				StateValue: stateValue,
				PlanValue:  types.ListUnknown(types.StringType),
				State:      makeState("cluster-1", 2, []string{"server-1", "server-2"}),
				Plan:       makePlan("cluster-1", 2),
				Path:       path.Root("servers_ids"),
			},
			resp,
		)
		if resp.Diagnostics.HasError() {
			t.Fatalf("unexpected error: %s", resp.Diagnostics.Errors())
		}
		if !resp.PlanValue.Equal(stateValue) {
			t.Errorf("expected state to be preserved, got %v", resp.PlanValue)
		}
	})

	t.Run("count changed - don't preserve state", func(t *testing.T) {
		t.Parallel()
		resp := &planmodifier.ListResponse{PlanValue: types.ListUnknown(types.StringType)}
		planmodifiers.UseStateUnlessCountChanges("servers_count").PlanModifyList(
			context.Background(),
			planmodifier.ListRequest{
				StateValue: stateValue,
				PlanValue:  types.ListUnknown(types.StringType),
				State:      makeState("cluster-1", 2, []string{"server-1", "server-2"}),
				Plan:       makePlan("cluster-1", 3),
				Path:       path.Root("servers_ids"),
			},
			resp,
		)
		if resp.Diagnostics.HasError() {
			t.Fatalf("unexpected error: %s", resp.Diagnostics.Errors())
		}
		if !resp.PlanValue.IsUnknown() {
			t.Errorf("expected unknown, got %v", resp.PlanValue)
		}
	})

	t.Run("plan id unknown but count unchanged - preserve state", func(t *testing.T) {
		// Attribute plan modifiers run in nondeterministic order: on a no-op
		// re-plan, id may still be unknown (marked by the framework, not yet
		// restored by its own UseStateForUnknown-style modifier) when this
		// modifier runs. It must NOT treat that as a replacement. Actual
		// replacement is covered by the null-prior-state case: Terraform core
		// re-plans replacement creates with null prior state.
		t.Parallel()
		resp := &planmodifier.ListResponse{PlanValue: types.ListUnknown(types.StringType)}
		planmodifiers.UseStateUnlessCountChanges("servers_count").PlanModifyList(
			context.Background(),
			planmodifier.ListRequest{
				StateValue: stateValue,
				PlanValue:  types.ListUnknown(types.StringType),
				State:      makeState("cluster-1", 2, []string{"server-1", "server-2"}),
				Plan:       makePlan(nil, 2), // id unknown, count unchanged
				Path:       path.Root("servers_ids"),
			},
			resp,
		)
		if resp.Diagnostics.HasError() {
			t.Fatalf("unexpected error: %s", resp.Diagnostics.Errors())
		}
		if !resp.PlanValue.Equal(stateValue) {
			t.Errorf("expected state to be preserved, got %v", resp.PlanValue)
		}
	})

	t.Run("plan value already known - don't override", func(t *testing.T) {
		t.Parallel()
		knownValue, _ := types.ListValueFrom(context.Background(), types.StringType, []string{"server-3"})
		resp := &planmodifier.ListResponse{PlanValue: knownValue}
		planmodifiers.UseStateUnlessCountChanges("servers_count").PlanModifyList(
			context.Background(),
			planmodifier.ListRequest{
				StateValue: stateValue,
				PlanValue:  knownValue,
				State:      makeState("cluster-1", 2, []string{"server-1", "server-2"}),
				Plan:       makePlan("cluster-1", 2),
				Path:       path.Root("servers_ids"),
			},
			resp,
		)
		if resp.Diagnostics.HasError() {
			t.Fatalf("unexpected error: %s", resp.Diagnostics.Errors())
		}
		if !resp.PlanValue.Equal(knownValue) {
			t.Errorf("expected known value preserved, got %v", resp.PlanValue)
		}
	})
}
