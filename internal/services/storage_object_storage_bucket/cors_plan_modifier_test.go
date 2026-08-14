package storage_object_storage_bucket

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// allowedOriginsAttr returns the cors.allowed_origins schema attribute, failing
// the test if it stopped being Computed+Optional — the shape that makes the
// framework mark the attribute unknown when the cors block is present but the
// attribute is omitted. If that changes, the guards below are no longer
// describing reality and should be re-evaluated rather than silently kept.
func allowedOriginsAttr(t *testing.T, ctx context.Context) schema.ListAttribute {
	t.Helper()

	corsRaw, ok := ResourceSchema(ctx).Attributes["cors"]
	if !ok {
		t.Fatal("cors attribute missing from resource schema")
	}
	cors, ok := corsRaw.(schema.SingleNestedAttribute)
	if !ok {
		t.Fatalf("cors is %T, expected schema.SingleNestedAttribute", corsRaw)
	}
	originsRaw, ok := cors.Attributes["allowed_origins"]
	if !ok {
		t.Fatal("allowed_origins attribute missing from the cors schema")
	}
	origins, ok := originsRaw.(schema.ListAttribute)
	if !ok {
		t.Fatalf("allowed_origins is %T, expected schema.ListAttribute", originsRaw)
	}
	if !origins.Computed || !origins.Optional {
		t.Fatalf("allowed_origins configurability changed (Computed=%v Optional=%v); re-evaluate the unknown-value guards", origins.Computed, origins.Optional)
	}
	return origins
}

// rawBucketWith builds a raw resource object where every top-level attribute is
// null except the given overrides.
func rawBucketWith(ctx context.Context, t *testing.T, overrides map[string]tftypes.Value) tftypes.Value {
	t.Helper()
	objType := ResourceSchema(ctx).Type().TerraformType(ctx).(tftypes.Object)
	vals := map[string]tftypes.Value{}
	for name, attrType := range objType.AttributeTypes {
		vals[name] = tftypes.NewValue(attrType, nil)
	}
	for name, val := range overrides {
		vals[name] = val
	}
	return tftypes.NewValue(objType, vals)
}

// rawCors builds a raw cors object carrying the given allowed_origins value.
func rawCors(ctx context.Context, t *testing.T, allowedOrigins tftypes.Value) tftypes.Value {
	t.Helper()
	objType := ResourceSchema(ctx).Type().TerraformType(ctx).(tftypes.Object)
	corsType := objType.AttributeTypes["cors"].(tftypes.Object)
	return tftypes.NewValue(corsType, map[string]tftypes.Value{
		"allowed_origins": allowedOrigins,
	})
}

// originsTerraformType returns the terraform type of cors.allowed_origins.
func originsTerraformType(ctx context.Context, t *testing.T) tftypes.Type {
	t.Helper()
	objType := ResourceSchema(ctx).Type().TerraformType(ctx).(tftypes.Object)
	return objType.AttributeTypes["cors"].(tftypes.Object).AttributeTypes["allowed_origins"]
}

// TestCorsBlockWithoutAllowedOriginsDecodes is the regression guard for the
// "Value Conversion Error ... Path: cors.allowed_origins" apply crash. Writing a
// cors block without allowed_origins makes the framework plan that Computed
// attribute as unknown; if the model field is retyped back to a plain
// *[]types.String — a Go type with no representation for unknown — this Get is
// exactly the call that fails inside Create, before any request is sent.
func TestCorsBlockWithoutAllowedOriginsDecodes(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	planRaw := rawBucketWith(ctx, t, map[string]tftypes.Value{
		"storage_id": tftypes.NewValue(tftypes.Number, 1),
		"name":       tftypes.NewValue(tftypes.String, "tf-test-bucket"),
		"cors":       rawCors(ctx, t, tftypes.NewValue(originsTerraformType(ctx, t), tftypes.UnknownValue)),
	})

	var data *StorageObjectStorageBucketModel
	diags := tfsdk.Plan{Raw: planRaw, Schema: ResourceSchema(ctx)}.Get(ctx, &data)
	if diags.HasError() {
		t.Fatalf("planning a cors block without allowed_origins must decode into the model, got: %v", diags)
	}
	if data.Cors == nil {
		t.Fatal("expected a non-nil cors after decoding a cors block")
	}
	if !data.Cors.AllowedOrigins.IsUnknown() {
		t.Fatalf("expected allowed_origins to carry the unknown through to the model, got %v", data.Cors.AllowedOrigins)
	}
}

// planAllowedOrigins runs the attribute's configured plan modifiers in order and
// returns the resulting planned value.
func planAllowedOrigins(ctx context.Context, t *testing.T, plan, config, state types.List, stateRaw tftypes.Value) types.List {
	t.Helper()
	origins := allowedOriginsAttr(t, ctx)
	current := plan
	for _, pm := range origins.PlanModifiers {
		resp := &planmodifier.ListResponse{PlanValue: current}
		pm.PlanModifyList(ctx, planmodifier.ListRequest{
			PlanValue:   current,
			ConfigValue: config,
			StateValue:  state,
			State:       tfsdk.State{Raw: stateRaw, Schema: ResourceSchema(ctx)},
			Plan:        tfsdk.Plan{Raw: tftypes.Value{}, Schema: ResourceSchema(ctx)},
		}, resp)
		current = resp.PlanValue
	}
	return current
}

// TestAllowedOriginsOmittedOnCreateStaysUnknown documents the create path: there
// is no prior state to reuse, so the attribute is left unknown, omitted from the
// request body, and resolved from the response when it is decoded. Planning a
// concrete value here instead would be a behaviour change — null is undefined for
// this endpoint, and an empty list is its documented "remove CORS" sentinel.
func TestAllowedOriginsOmittedOnCreateStaysUnknown(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	got := planAllowedOrigins(ctx, t,
		types.ListUnknown(types.StringType),                                  // framework marked the computed attribute unknown
		types.ListNull(types.StringType),                                     // allowed_origins omitted from config
		types.ListNull(types.StringType),                                     // no prior value
		tftypes.NewValue(ResourceSchema(ctx).Type().TerraformType(ctx), nil), // no prior state: create
	)

	if !got.IsUnknown() {
		t.Fatalf("expected allowed_origins to stay unknown on create, got %v", got)
	}
}

// TestAllowedOriginsOmittedOnUpdateReusesState covers the update path: omitting
// allowed_origins from a cors block that already has origins must re-plan the
// prior value, so the PATCH body omits the field and the configured origins are
// preserved. Without this the attribute would replan as unknown on every apply
// and never converge to an empty plan.
func TestAllowedOriginsOmittedOnUpdateReusesState(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	priorOrigins := types.ListValueMust(types.StringType, []attr.Value{types.StringValue("https://example.com")})
	stateRaw := rawBucketWith(ctx, t, map[string]tftypes.Value{
		"storage_id": tftypes.NewValue(tftypes.Number, 1),
		"name":       tftypes.NewValue(tftypes.String, "tf-test-bucket"),
		"cors": rawCors(ctx, t, tftypes.NewValue(originsTerraformType(ctx, t), []tftypes.Value{
			tftypes.NewValue(tftypes.String, "https://example.com"),
		})),
	})

	got := planAllowedOrigins(ctx, t,
		types.ListUnknown(types.StringType), // computed attribute marked unknown
		types.ListNull(types.StringType),    // allowed_origins omitted from config
		priorOrigins,                        // prior state has origins
		stateRaw,
	)

	if got.IsUnknown() {
		t.Fatal("allowed_origins must not stay unknown on update; the plan would never converge")
	}
	if !got.Equal(priorOrigins) {
		t.Fatalf("expected the prior origins to be re-planned, got %v", got)
	}
}
