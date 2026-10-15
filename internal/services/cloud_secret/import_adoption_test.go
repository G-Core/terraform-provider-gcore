package cloud_secret

import (
	"context"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/G-Core/terraform-provider-gcore/internal/importadoption"
	"github.com/G-Core/terraform-provider-gcore/internal/importadoption/adoptiontest"
)

// TestImportAdoptionWiringPinned asserts every sealed hook is literally
// present in resource.go, read off disk. resource.go is generated; these
// hooks are custom hunks, the class of change that silently evaporates on a
// bad reseal.
func TestImportAdoptionWiringPinned(t *testing.T) {
	t.Parallel()

	adoptiontest.AssertWired(t, "resource.go",
		adoptiontest.StandardWiring("CloudSecretResource", "secretImportAdoption", "secretImportAdoptionPlan").
			With("ModifyPlan", "requirePayloadOnCreate(ctx,req,resp)"))

	// Advance is one phase per Read; a second call retires the marker in one
	// refresh and the apply after `terraform import` plans a replacement.
	if n := adoptiontest.CountInMethod(t, "resource.go", "CloudSecretResource", "Read", "secretImportAdoption.Advance(ctx,resp.Private)"); n != 1 {
		t.Errorf("Read calls secretImportAdoption.Advance %d times, want exactly 1", n)
	}
}

// adoptingAttributes are the two no_refresh attributes that carry the
// marker-keyed modifier: both independently force replacement on an imported
// plan (confirmed against real infra), so both must be pinned.
var adoptingAttributes = []string{"payload", "payload_wo_version"}

// TestSchemaModifiersPinned pins the schema half: payload and
// payload_wo_version carry the local wrapper modifiers (which fall through to
// the marker-keyed modifier for every transition except an in-place
// removal), by %T. A reseal that reverts to the bare
// planmodifiers.*RequiresReplaceUnlessAdopting call, or to the framework's
// unconditional RequiresReplace(), fails this closed.
func TestSchemaModifiersPinned(t *testing.T) {
	t.Parallel()

	attrs := ResourceSchema(context.Background()).Attributes
	adoptiontest.AssertModifierTypes(t, attrs, "payload", "cloud_secret.objectReplaceUnlessAdoptingOrRemovedModifier")
	adoptiontest.AssertModifierTypes(t, attrs, "payload_wo_version", "cloud_secret.int64ReplaceUnlessAdoptingOrRemovedModifier")

	// %T cannot see the constructor argument: each call must carry
	// secretImportAdoption.Key, or the modifier reads a key nothing writes
	// and every import plans a replacement.
	adoptiontest.AssertAttributeCallArgs(t, "schema.go", "payload",
		"objectReplaceUnlessAdoptingOrRemoved", "secretImportAdoption.Key")
	adoptiontest.AssertAttributeCallArgs(t, "schema.go", "payload_wo_version",
		"int64ReplaceUnlessAdoptingOrRemoved", "secretImportAdoption.Key")
}

// TestPayloadAttributesOptionalPinned pins payload and payload_wo_version as
// Optional, not Required or Computed. A reseal can revert this hand flip
// (the annotation-vs-inference trap that also bit healthmonitor): with either
// attribute Required, importing a secret without its material becomes
// impossible, because the config could never omit it.
func TestPayloadAttributesOptionalPinned(t *testing.T) {
	t.Parallel()

	attrs := ResourceSchema(context.Background()).Attributes
	for _, name := range adoptingAttributes {
		a, ok := attrs[name]
		if !ok {
			t.Fatalf("attribute %q not in schema", name)
		}
		if !a.IsOptional() || a.IsRequired() || a.IsComputed() {
			t.Errorf("%s: want Optional (and not Required or Computed); a reseal reverted the hand flip back to "+
				"Required, which makes importing a secret without its material impossible", name)
		}
	}
}

// TestPayloadHasNoAttributeValidators pins that payload and its children
// carry no attribute-level validators, and that ConfigValidators is empty.
// A validator attached to payload (or a child, or a resourcevalidator that
// targets it) reports its diagnostic on that attribute's path, and Terraform
// then prints the offending config source line - which can hold the private
// key. Every check in this package must stay path-less (see payload_rules.go).
func TestPayloadHasNoAttributeValidators(t *testing.T) {
	t.Parallel()

	payload, ok := ResourceSchema(context.Background()).Attributes["payload"].(schema.SingleNestedAttribute)
	if !ok {
		t.Fatal("payload is no longer a schema.SingleNestedAttribute")
	}
	if len(payload.Validators) != 0 {
		t.Fatalf("payload carries %d attribute validators; its configuration line can hold the private key, "+
			"which Terraform would print verbatim for any diagnostic attributed to this path", len(payload.Validators))
	}
	for name, attribute := range payload.Attributes {
		child, ok := attribute.(schema.StringAttribute)
		if !ok {
			t.Fatalf("payload.%s is no longer a schema.StringAttribute", name)
		}
		if len(child.Validators) != 0 {
			t.Fatalf("payload.%s carries %d attribute validators; its configuration line can hold the private key, "+
				"which Terraform would print verbatim for any diagnostic attributed to this path", name, len(child.Validators))
		}
	}

	if got := (&CloudSecretResource{}).ConfigValidators(context.Background()); len(got) != 0 {
		t.Fatalf("ConfigValidators() = %v, want empty; a resourcevalidator targeting payload prints the "+
			"offending config source line, which can hold the private key", got)
	}
}

// TestAdoptionWarningExplainsWriteOnlyPayload pins the secret-specific
// transport paragraph: the shared default says the configured values are
// recorded, which reads as if editing the write-only material were tracked.
func TestAdoptionWarningExplainsWriteOnlyPayload(t *testing.T) {
	t.Parallel()

	detail := secretImportAdoptionPlan.Warning.Diagnostic([]string{"payload", "payload_wo_version"}).Detail()
	for _, want := range []string{"write-only", "Bump \"payload_wo_version\" whenever the material changes"} {
		if !strings.Contains(detail, want) {
			t.Errorf("adoption warning lacks %q:\n%s", want, detail)
		}
	}
	if strings.Contains(detail, "trusted as-is") {
		t.Errorf("adoption warning still uses the default transport paragraph:\n%s", detail)
	}
}

func TestCloudSecretImportAdoptionWarning(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	resourceSchema := ResourceSchema(ctx)
	objectType := resourceSchema.Type().TerraformType(ctx).(tftypes.Object)
	payloadType := objectType.AttributeTypes["payload"].(tftypes.Object)

	configuredPayload := tftypes.NewValue(payloadType, map[string]tftypes.Value{
		"certificate_wo":       tftypes.NewValue(tftypes.String, "certificate-material"),
		"certificate_chain_wo": tftypes.NewValue(tftypes.String, "chain-material"),
		"private_key_wo":       tftypes.NewValue(tftypes.String, "private-key-material"),
	})
	plannedPayload := tftypes.NewValue(payloadType, map[string]tftypes.Value{
		"certificate_wo":       tftypes.NewValue(tftypes.String, nil),
		"certificate_chain_wo": tftypes.NewValue(tftypes.String, nil),
		"private_key_wo":       tftypes.NewValue(tftypes.String, nil),
	})

	configRaw := cloudSecretTestValue(objectType, map[string]tftypes.Value{
		"name":               tftypes.NewValue(tftypes.String, "imported-secret"),
		"payload":            configuredPayload,
		"payload_wo_version": tftypes.NewValue(tftypes.Number, int64(1)),
	})
	planRaw := cloudSecretTestValue(objectType, map[string]tftypes.Value{
		"name":               tftypes.NewValue(tftypes.String, "imported-secret"),
		"payload":            plannedPayload,
		"payload_wo_version": tftypes.NewValue(tftypes.Number, int64(1)),
	})
	liveStateRaw := cloudSecretTestValue(objectType, map[string]tftypes.Value{
		"name": tftypes.NewValue(tftypes.String, "imported-secret"),
	})
	nullStateRaw := tftypes.NewValue(objectType, nil)

	tests := []struct {
		name        string
		state       tftypes.Value
		markImport  bool
		wantWarning bool
	}{
		{name: "import adoption warns", state: liveStateRaw, markImport: true, wantWarning: true},
		{name: "create does not warn", state: nullStateRaw},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			req := resource.ModifyPlanRequest{
				Config: tfsdk.Config{Raw: configRaw, Schema: resourceSchema},
				Plan:   tfsdk.Plan{Raw: planRaw, Schema: resourceSchema},
				State:  tfsdk.State{Raw: tc.state, Schema: resourceSchema},
			}
			if tc.markImport {
				private := adoptiontest.NewPrivate(t, &req)
				if diags := secretImportAdoption.Mark(ctx, private); diags.HasError() {
					t.Fatalf("marking import adoption: %v", diags)
				}
			}

			resp := &resource.ModifyPlanResponse{Plan: req.Plan}
			(&CloudSecretResource{}).ModifyPlan(ctx, req, resp)
			if resp.Diagnostics.HasError() {
				t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
			}

			warnings := resp.Diagnostics.Warnings()
			if !tc.wantWarning {
				if len(warnings) != 0 {
					t.Fatalf("unexpected warnings: %v", warnings)
				}
				return
			}
			if len(warnings) != 1 || warnings[0].Summary() != importadoption.Summary {
				t.Fatalf("expected one import adoption warning, got %v", warnings)
			}
			for _, secret := range []string{"certificate-material", "chain-material", "private-key-material"} {
				if strings.Contains(warnings[0].Detail(), secret) {
					t.Fatalf("warning contains configured secret material %q", secret)
				}
			}
		})
	}
}

func cloudSecretTestValue(objectType tftypes.Object, overrides map[string]tftypes.Value) tftypes.Value {
	values := make(map[string]tftypes.Value, len(objectType.AttributeTypes))
	for name, attributeType := range objectType.AttributeTypes {
		values[name] = tftypes.NewValue(attributeType, nil)
	}
	for name, value := range overrides {
		values[name] = value
	}
	return tftypes.NewValue(objectType, values)
}
