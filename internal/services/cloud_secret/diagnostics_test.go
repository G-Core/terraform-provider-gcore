package cloud_secret_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/G-Core/terraform-provider-gcore/internal/services/cloud_secret"
)

// Literal copies of the payload_rules.go summary constants: this test lives
// in cloud_secret_test and cannot reference the unexported originals.
const (
	wantMissingPayloadSummary  = "Missing secret payload"
	wantUnpairedPayloadSummary = "\"payload\" and \"payload_wo_version\" must be set together"
)

// Terraform prints the offending config source line for any diagnostic that
// carries an attribute path, and it does not redact the value. For "payload"
// that line can hold the private key, so no diagnostic from ValidateConfig or
// ModifyPlan may point there. Comparing only the first path step (not a
// string prefix) matters: "payload_wo_version" starts with "payload" but is a
// different, harmless attribute.
func TestCloudSecretDiagnosticsNeverPointAtPayload(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	schema := cloud_secret.ResourceSchema(ctx)
	forbidden := path.Root("payload").Steps()[0]

	payload := payloadValue(t)

	cases := map[string]struct {
		config      map[string]tftypes.Value
		plan        map[string]tftypes.Value
		wantSummary string
	}{
		"payload without version": {
			config:      map[string]tftypes.Value{"name": tfStr("keep"), "payload": payload},
			plan:        map[string]tftypes.Value{"name": tfStr("keep"), "payload": payload},
			wantSummary: wantUnpairedPayloadSummary,
		},
		"version without payload": {
			config:      map[string]tftypes.Value{"name": tfStr("keep"), "payload_wo_version": tfNum(1)},
			plan:        map[string]tftypes.Value{"name": tfStr("keep"), "payload_wo_version": tfNum(1)},
			wantSummary: wantUnpairedPayloadSummary,
		},
		"create with neither": {
			config:      map[string]tftypes.Value{"name": tfStr("keep")},
			plan:        map[string]tftypes.Value{"name": tfStr("keep")},
			wantSummary: wantMissingPayloadSummary,
		},
		"replacement re-plan": {
			// Terraform core re-plans the create half of every replacement
			// with a NULL prior state, so a null state here stands in for
			// that re-plan even though conceptually this is a replacement,
			// not a create.
			config:      map[string]tftypes.Value{"name": tfStr("renamed")},
			plan:        map[string]tftypes.Value{"name": tfStr("renamed")},
			wantSummary: wantMissingPayloadSummary,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			config := secretValues(ctx, t, tc.config)
			plan := secretValues(ctx, t, tc.plan)
			stateRaw := secretNull(ctx, t) // every case above plans against a null prior state

			diags := runValidateAndModifyPlan(ctx, t, schema, config, plan, stateRaw)

			if len(diags) == 0 {
				t.Fatal("expected at least one diagnostic; this case no longer guards anything")
			}

			foundWant := false
			for _, d := range diags {
				if d.Severity() == diag.SeverityError && d.Summary() == tc.wantSummary {
					foundWant = true
				}
				withPath, ok := d.(diag.DiagnosticWithPath)
				if !ok {
					continue
				}
				steps := withPath.Path().Steps()
				if len(steps) > 0 && steps[0].Equal(forbidden) {
					t.Fatalf("diagnostic %q is attributed to payload, whose configuration line Terraform would print: %s", d.Summary(), d.Detail())
				}
			}
			if !foundWant {
				t.Errorf("no error diagnostic with summary %q; got %v", tc.wantSummary, diags)
			}
		})
	}
}

// TestCloudSecretDiagnosticsNoErrorCases covers the two shapes that must NOT
// error: an existing resource that never configures payload, and a create
// that configures both.
func TestCloudSecretDiagnosticsNoErrorCases(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	schema := cloud_secret.ResourceSchema(ctx)
	payload := payloadValue(t)

	t.Run("existing state, both omitted", func(t *testing.T) {
		t.Parallel()

		config := secretValues(ctx, t, map[string]tftypes.Value{"name": tfStr("keep")})
		plan := secretValues(ctx, t, map[string]tftypes.Value{"name": tfStr("keep")})
		state := secretValues(ctx, t, map[string]tftypes.Value{"name": tfStr("keep")})

		diags := runValidateAndModifyPlan(ctx, t, schema, config, plan, state)
		if diags.HasError() {
			t.Errorf("unexpected error diagnostics: %v", diags)
		}
	})

	t.Run("create with both set", func(t *testing.T) {
		t.Parallel()

		config := secretValues(ctx, t, map[string]tftypes.Value{
			"name": tfStr("keep"), "payload": payload, "payload_wo_version": tfNum(1),
		})
		plan := secretValues(ctx, t, map[string]tftypes.Value{
			"name": tfStr("keep"), "payload": payload, "payload_wo_version": tfNum(1),
		})
		state := secretNull(ctx, t)

		diags := runValidateAndModifyPlan(ctx, t, schema, config, plan, state)
		if diags.HasError() {
			t.Errorf("unexpected error diagnostics: %v", diags)
		}
	})
}

// runValidateAndModifyPlan calls ValidateConfig and ModifyPlan the way
// Terraform core would for one plan, with Private left nil (no import
// adoption pending), and returns every diagnostic from both.
func runValidateAndModifyPlan(ctx context.Context, t *testing.T, schema rschema.Schema, config, plan, state tftypes.Value) diag.Diagnostics {
	t.Helper()

	r := cloud_secret.NewResource()

	var diags diag.Diagnostics

	validated, ok := r.(resource.ResourceWithValidateConfig)
	if !ok {
		t.Fatal("the cloud_secret resource no longer implements ValidateConfig")
	}
	validateResp := &resource.ValidateConfigResponse{}
	validated.ValidateConfig(ctx, resource.ValidateConfigRequest{
		Config: tfsdk.Config{Raw: config, Schema: schema},
	}, validateResp)
	diags.Append(validateResp.Diagnostics...)

	modifier, ok := r.(resource.ResourceWithModifyPlan)
	if !ok {
		t.Fatal("the cloud_secret resource no longer implements ModifyPlan")
	}
	planResp := &resource.ModifyPlanResponse{Plan: tfsdk.Plan{Raw: plan, Schema: schema}}
	modifier.ModifyPlan(ctx, resource.ModifyPlanRequest{
		Config: tfsdk.Config{Raw: config, Schema: schema},
		Plan:   tfsdk.Plan{Raw: plan, Schema: schema},
		State:  tfsdk.State{Raw: state, Schema: schema},
	}, planResp)
	diags.Append(planResp.Diagnostics...)

	return diags
}
