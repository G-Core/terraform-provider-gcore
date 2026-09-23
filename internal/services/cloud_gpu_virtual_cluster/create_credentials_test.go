package cloud_gpu_virtual_cluster

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/G-Core/gcore-go"
	"github.com/G-Core/gcore-go/option"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// clusterWithCredentials builds a raw cluster whose attributes are null except
// the required scalars and the given credentials.
func clusterWithCredentials(ctx context.Context, creds map[string]tftypes.Value) tftypes.Value {
	objType := ResourceSchema(ctx).Type().TerraformType(ctx).(tftypes.Object)
	settingsType := objType.AttributeTypes["servers_settings"].(tftypes.Object)
	credType := settingsType.AttributeTypes["credentials"].(tftypes.Object)

	nulls := func(o tftypes.Object) map[string]tftypes.Value {
		vals := map[string]tftypes.Value{}
		for name, attrType := range o.AttributeTypes {
			vals[name] = tftypes.NewValue(attrType, nil)
		}
		return vals
	}

	credVals := nulls(credType)
	for name, val := range creds {
		credVals[name] = val
	}
	settingsVals := nulls(settingsType)
	settingsVals["credentials"] = tftypes.NewValue(credType, credVals)

	vals := nulls(objType)
	vals["name"] = tftypes.NewValue(tftypes.String, "c")
	vals["flavor"] = tftypes.NewValue(tftypes.String, "f")
	vals["project_id"] = tftypes.NewValue(tftypes.Number, 1)
	vals["region_id"] = tftypes.NewValue(tftypes.Number, 1)
	vals["servers_count"] = tftypes.NewValue(tftypes.Number, 1)
	vals["servers_settings"] = tftypes.NewValue(settingsType, settingsVals)

	return tftypes.NewValue(objType, vals)
}

// TestCreateAcceptsPasswordCredentials: Terraform sends the write-only
// password_wo as null in the plan, even at apply. The re-check must read the
// config, where it is set, or it rejects every username/password cluster.
func TestCreateAcceptsPasswordCredentials(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	var seen []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		seen = append(seen, r.Method+" "+r.URL.Path)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"message": "stub"}`))
	}))
	t.Cleanup(server.Close)
	client := gcore.NewClient(option.WithBaseURL(server.URL), option.WithAPIKey("test"))

	cluster := func(password tftypes.Value) tftypes.Value {
		return clusterWithCredentials(ctx, map[string]tftypes.Value{
			"username":            tftypes.NewValue(tftypes.String, "ubuntu"),
			"password_wo":         password,
			"password_wo_version": tftypes.NewValue(tftypes.Number, 1),
		})
	}
	sch := ResourceSchema(ctx)
	resp := &resource.CreateResponse{State: tfsdk.State{Schema: sch}}
	(&CloudGPUVirtualClusterResource{client: &client}).Create(ctx, resource.CreateRequest{
		Plan:   tfsdk.Plan{Raw: cluster(tftypes.NewValue(tftypes.String, nil)), Schema: sch},
		Config: tfsdk.Config{Raw: cluster(tftypes.NewValue(tftypes.String, "secret")), Schema: sch},
	}, resp)

	for _, d := range resp.Diagnostics.Errors() {
		if strings.Contains(d.Summary(), "credentials") || strings.Contains(d.Summary(), "password_wo") {
			t.Fatalf("valid password credentials rejected: %s: %s", d.Summary(), d.Detail())
		}
	}
	if len(seen) == 0 {
		t.Fatalf("expected Create to reach the API, got: %v", resp.Diagnostics)
	}
}
