package storage_sftp_test

import (
	"context"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/G-Core/terraform-provider-gcore/internal/services/storage_sftp"
)

// Terraform prints the offending config source line for any diagnostic that
// carries an attribute path, and Sensitive does not redact it. For password_wo
// that line is the literal password, so no diagnostic may point there.

// a schema validator reports on its own attribute's path, so password_wo must
// not carry any; the length check lives in passwordPairValidator instead
func TestStorageSftpPasswordAttributeHasNoValidators(t *testing.T) {
	t.Parallel()

	attribute, ok := storage_sftp.ResourceSchema(context.Background()).Attributes["password_wo"].(rschema.StringAttribute)
	if !ok {
		t.Fatal("password_wo is no longer a schema.StringAttribute")
	}
	if len(attribute.Validators) != 0 {
		t.Fatalf("password_wo carries %d attribute validators; their diagnostics point at password_wo, "+
			"whose configuration line Terraform prints verbatim. Validate it in passwordPairValidator instead.",
			len(attribute.Validators))
	}
}

// each case exercises one diagnostic that ConfigValidators or ModifyPlan can
// produce; every case must yield a diagnostic, and none may point at
// password_wo
func TestStorageSftpDiagnosticsNeverPointAtPassword(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	schema := storage_sftp.ResourceSchema(ctx)
	forbidden := path.Root("password_wo")

	unknownPassword := tftypes.NewValue(tftypes.String, tftypes.UnknownValue)
	unknownVersion := tftypes.NewValue(tftypes.Number, tftypes.UnknownValue)
	unknownBool := tftypes.NewValue(tftypes.Bool, tftypes.UnknownValue)

	cases := map[string]struct {
		config map[string]tftypes.Value
		plan   map[string]tftypes.Value
		state  map[string]tftypes.Value
	}{
		"password without a version": {
			config: map[string]tftypes.Value{"name": str("keep"), "password_wo": str("Unit-Test-Pw-8Kd2xQ")},
			plan:   map[string]tftypes.Value{"name": str("keep"), "has_password": unknownBool},
		},
		"version without a password": {
			config: map[string]tftypes.Value{"name": str("keep"), "password_wo_version": num(1)},
			plan:   map[string]tftypes.Value{"name": str("keep"), "password_wo_version": num(1), "has_password": unknownBool},
		},
		"password too short": {
			config: map[string]tftypes.Value{"name": str("keep"), "password_wo": str("short"), "password_wo_version": num(1)},
			plan:   map[string]tftypes.Value{"name": str("keep"), "password_wo_version": num(1), "has_password": unknownBool},
		},
		"password too long": {
			config: map[string]tftypes.Value{
				"name": str("keep"), "password_wo": str("Pw-" + strings.Repeat("x", 61)), "password_wo_version": num(1),
			},
			plan: map[string]tftypes.Value{"name": str("keep"), "password_wo_version": num(1), "has_password": unknownBool},
		},
		"unknown password without a version": {
			// the config validator skips unknowns; the plan check must report
			config: map[string]tftypes.Value{"name": str("keep"), "password_wo": unknownPassword},
			plan:   map[string]tftypes.Value{"name": str("keep"), "has_password": unknownBool},
		},
		"unknown version without a password": {
			config: map[string]tftypes.Value{"name": str("keep"), "password_wo_version": unknownVersion},
			plan:   map[string]tftypes.Value{"name": str("keep"), "password_wo_version": unknownVersion, "has_password": unknownBool},
		},
		"adoption warning": {
			config: map[string]tftypes.Value{
				"name": str("keep"), "password_wo": unknownPassword, "password_wo_version": num(1),
			},
			plan:  map[string]tftypes.Value{"name": str("keep"), "password_wo_version": num(1), "has_password": unknownBool},
			state: map[string]tftypes.Value{"name": str("keep"), "has_password": boolean(true)},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			config := storageSftpValues(ctx, t, tc.config)
			plan := storageSftpValues(ctx, t, tc.plan)
			stateRaw := storageSftpNull(ctx, t)
			if tc.state != nil {
				stateRaw = storageSftpValues(ctx, t, tc.state)
			}

			r := storage_sftp.NewResource()

			validated, ok := r.(resource.ResourceWithConfigValidators)
			if !ok {
				t.Fatal("the storage_sftp resource no longer implements ConfigValidators")
			}
			var diags diag.Diagnostics
			for _, v := range validated.ConfigValidators(ctx) {
				resp := &resource.ValidateConfigResponse{}
				v.ValidateResource(ctx, resource.ValidateConfigRequest{
					Config: tfsdk.Config{Raw: config, Schema: schema},
				}, resp)
				diags.Append(resp.Diagnostics...)
			}

			modifier, ok := r.(resource.ResourceWithModifyPlan)
			if !ok {
				t.Fatal("the storage_sftp resource no longer implements ModifyPlan")
			}
			planResp := &resource.ModifyPlanResponse{Plan: tfsdk.Plan{Raw: plan, Schema: schema}}
			modifier.ModifyPlan(ctx, resource.ModifyPlanRequest{
				Config: tfsdk.Config{Raw: config, Schema: schema},
				Plan:   tfsdk.Plan{Raw: plan, Schema: schema},
				State:  tfsdk.State{Raw: stateRaw, Schema: schema},
			}, planResp)
			diags.Append(planResp.Diagnostics...)

			if len(diags) == 0 {
				t.Fatal("expected at least one diagnostic; this case no longer guards anything")
			}
			for _, d := range diags {
				withPath, ok := d.(diag.DiagnosticWithPath)
				if !ok {
					continue
				}
				if withPath.Path().Equal(forbidden) {
					t.Fatalf("diagnostic %q is attributed to password_wo, whose configuration line Terraform would print: %s", d.Summary(), d.Detail())
				}
			}
		})
	}
}
