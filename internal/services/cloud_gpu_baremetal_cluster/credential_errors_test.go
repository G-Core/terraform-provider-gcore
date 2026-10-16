package cloud_gpu_baremetal_cluster

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/G-Core/terraform-provider-gcore/internal/importadoption/adoptiontest"
)

// credentialsSettingsOrNull builds servers_settings with the given credential
// attributes, or with a null credentials object when creds is nil (a cluster
// created without credentials, or an import whose API response has neither
// username nor ssh_key_name).
func credentialsSettingsOrNull(ctx context.Context, t *testing.T, creds map[string]tftypes.Value) tftypes.Value {
	t.Helper()
	if creds != nil {
		return serversSettingsWithCredentials(ctx, t, creds)
	}
	objType := ResourceSchema(ctx).Type().TerraformType(ctx).(tftypes.Object)
	settingsType := objType.AttributeTypes["servers_settings"].(tftypes.Object)
	vals := map[string]tftypes.Value{}
	for name, attrType := range settingsType.AttributeTypes {
		vals[name] = tftypes.NewValue(attrType, nil)
	}
	return tftypes.NewValue(settingsType, vals)
}

var (
	tfUser    = tftypes.NewValue(tftypes.String, "tfuser")
	verOne    = tftypes.NewValue(tftypes.Number, 1)
	verTwo    = tftypes.NewValue(tftypes.Number, 2)
	pwConfig  = map[string]tftypes.Value{"username": tfUser, "password_wo": tftypes.NewValue(tftypes.String, "secret"), "password_wo_version": verOne}
	pwPlanned = map[string]tftypes.Value{"username": tfUser, "password_wo_version": verOne}
)

type credentialCase struct {
	stateCreds, planCreds, configCreds map[string]tftypes.Value
	imported                           bool
}

func (c credentialCase) modifyPlan(ctx context.Context, t *testing.T) *resource.ModifyPlanResponse {
	t.Helper()
	sch := ResourceSchema(ctx)
	settings := func(creds map[string]tftypes.Value) tftypes.Value {
		return rawObjectWith(ctx, t, map[string]tftypes.Value{"servers_settings": credentialsSettingsOrNull(ctx, t, creds)})
	}
	req := resource.ModifyPlanRequest{
		Plan:   tfsdk.Plan{Raw: settings(c.planCreds), Schema: sch},
		State:  tfsdk.State{Raw: settings(c.stateCreds), Schema: sch},
		Config: tfsdk.Config{Raw: settings(c.configCreds), Schema: sch},
	}
	if c.imported {
		if d := markImported(ctx, adoptiontest.NewPrivate(t, &req)); d.HasError() {
			t.Fatalf("markImported: %v", d)
		}
	}
	resp := &resource.ModifyPlanResponse{}
	(&CloudGPUBaremetalClusterResource{}).ModifyPlan(ctx, req, resp)
	return resp
}

func (c credentialCase) update(ctx context.Context, t *testing.T) *resource.UpdateResponse {
	t.Helper()
	sch := ResourceSchema(ctx)
	settings := func(creds map[string]tftypes.Value) tftypes.Value {
		return rawObjectWith(ctx, t, map[string]tftypes.Value{"servers_settings": credentialsSettingsOrNull(ctx, t, creds)})
	}
	stateRaw := settings(c.stateCreds)
	req := resource.UpdateRequest{
		Plan:  tfsdk.Plan{Raw: settings(c.planCreds), Schema: sch},
		State: tfsdk.State{Raw: stateRaw, Schema: sch},
	}
	if c.imported {
		if d := markImported(ctx, adoptiontest.NewPrivate(t, &req)); d.HasError() {
			t.Fatalf("markImported: %v", d)
		}
	}
	resp := &resource.UpdateResponse{State: tfsdk.State{Raw: stateRaw, Schema: sch}}
	(&CloudGPUBaremetalClusterResource{}).Update(ctx, req, resp)
	return resp
}

// The credential error points to the documented ignore_changes block only
// for an imported cluster. A cluster created without credentials has the
// same null credentials in state, so state alone must not trigger the hint.
func TestUnsupportedCredentialUpdateImportHint(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	for name, tc := range map[string]struct {
		credentialCase
		wantHint bool
	}{
		"imported, API without username": {credentialCase{planCreds: pwPlanned, configCreds: pwConfig, imported: true}, true},
		"imported, API with username": {credentialCase{stateCreds: map[string]tftypes.Value{"username": tfUser},
			planCreds: pwPlanned, configCreds: pwConfig, imported: true}, true},
		// QA: created by Terraform without credentials, then credentials added.
		"created without credentials, password added": {credentialCase{planCreds: pwPlanned, configCreds: pwConfig}, false},
		"created with a password, version bumped": {credentialCase{stateCreds: pwPlanned,
			planCreds:   map[string]tftypes.Value{"username": tfUser, "password_wo_version": verTwo},
			configCreds: map[string]tftypes.Value{"username": tfUser, "password_wo": tftypes.NewValue(tftypes.String, "s2"), "password_wo_version": verTwo}}, false},
		"created with a key, switched to a password": {credentialCase{
			stateCreds: map[string]tftypes.Value{"ssh_key_name": tftypes.NewValue(tftypes.String, "k1")},
			planCreds:  pwPlanned, configCreds: pwConfig}, false},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			for step, diags := range map[string]diag.Diagnostics{
				"ModifyPlan": tc.modifyPlan(ctx, t).Diagnostics,
				"Update":     tc.update(ctx, t).Diagnostics,
			} {
				errs := diags.Errors()
				if len(errs) != 1 || errs[0].Summary() != unsupportedCredentialUpdateSummary {
					t.Fatalf("%s: want one credential error, got %v", step, diags)
				}
				hasHint := strings.Contains(errs[0].Detail(), "ignore_changes = [servers_settings.credentials]")
				if hasHint != tc.wantHint {
					t.Fatalf("%s: import hint = %v, want %v: %q", step, hasHint, tc.wantHint, errs[0].Detail())
				}
			}
		})
	}
}

// ImportState is the only writer of the mark; Create starts from empty
// private state, so a replaced cluster is no longer treated as imported.
func TestImportMark(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	var req resource.ModifyPlanRequest
	if wasImported(ctx, req.Private) {
		t.Fatal("nil private state must read as not imported")
	}
	p := adoptiontest.NewPrivate(t, &req)
	if wasImported(ctx, req.Private) {
		t.Fatal("empty private state must read as not imported")
	}
	if d := markImported(ctx, p); d.HasError() {
		t.Fatalf("markImported: %v", d)
	}
	if !wasImported(ctx, req.Private) {
		t.Fatal("mark not read back")
	}
}

// ImportState writes the mark and the private state it returns is what later
// plans read: a credential change after a real import gets the hint.
func TestImportStateMarksImported(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	res, _, _ := newStubbedResource(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, clusterResponse("imported", ""))
	})
	sch := ResourceSchema(ctx)
	importResp := &resource.ImportStateResponse{State: tfsdk.State{
		Raw: tftypes.NewValue(sch.Type().TerraformType(ctx), nil), Schema: sch,
	}}
	adoptiontest.NewPrivate(t, importResp)
	res.ImportState(ctx, resource.ImportStateRequest{ID: "1/2/cluster-1"}, importResp)
	if importResp.Diagnostics.HasError() {
		t.Fatalf("ImportState: %v", importResp.Diagnostics)
	}
	if !wasImported(ctx, importResp.Private) {
		t.Fatal("ImportState did not mark the cluster as imported")
	}

	settings := func(creds map[string]tftypes.Value) tftypes.Value {
		return rawObjectWith(ctx, t, map[string]tftypes.Value{"servers_settings": credentialsSettingsOrNull(ctx, t, creds)})
	}
	planResp := &resource.ModifyPlanResponse{}
	res.ModifyPlan(ctx, resource.ModifyPlanRequest{
		Plan:    tfsdk.Plan{Raw: settings(pwPlanned), Schema: sch},
		State:   tfsdk.State{Raw: settings(nil), Schema: sch},
		Config:  tfsdk.Config{Raw: settings(pwConfig), Schema: sch},
		Private: importResp.Private,
	}, planResp)
	errs := planResp.Diagnostics.Errors()
	if len(errs) != 1 || !strings.Contains(errs[0].Detail(), "This cluster was imported") {
		t.Fatalf("want the credential error with the import hint, got %v", planResp.Diagnostics)
	}
}
