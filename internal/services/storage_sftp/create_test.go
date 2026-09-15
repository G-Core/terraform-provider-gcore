package storage_sftp_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sync"
	"testing"

	"github.com/G-Core/gcore-go"
	"github.com/G-Core/gcore-go/option"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/G-Core/terraform-provider-gcore/internal/services/storage_sftp"
)

// Create must read the write-only password from config (it is null in the
// plan), derive password_mode for the POST body, and keep the echoed password
// out of state.
func TestStorageSftpCreateSendsPasswordDecision(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	schema := storage_sftp.ResourceSchema(ctx)

	const password = "Unit-Test-Pw-8Kd2xQ"
	unknownBool := tftypes.NewValue(tftypes.Bool, tftypes.UnknownValue)

	cases := map[string]struct {
		config      map[string]tftypes.Value
		plan        map[string]tftypes.Value
		hasPassword bool
		wantBody    map[string]any
	}{
		"with the pair": {
			config: map[string]tftypes.Value{"password_wo": str(password), "password_wo_version": num(1)},
			plan:   map[string]tftypes.Value{"password_wo_version": num(1), "has_password": unknownBool},

			hasPassword: true,
			wantBody: map[string]any{
				"name": "keep", "location_name": "frn",
				"password_mode": "set", "password": password,
			},
		},
		"without the pair": {
			config: map[string]tftypes.Value{},
			plan:   map[string]tftypes.Value{"has_password": unknownBool},

			hasPassword: false,
			wantBody: map[string]any{
				"name": "keep", "location_name": "frn",
				"password_mode": "none",
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			var (
				mu      sync.Mutex
				posts   int
				gotPath string
				gotBody map[string]any
			)
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				mu.Lock()
				defer mu.Unlock()
				w.Header().Set("Content-Type", "application/json")
				switch r.Method {
				case http.MethodPost:
					posts++
					gotPath = r.URL.Path
					raw, _ := io.ReadAll(r.Body)
					if err := json.Unmarshal(raw, &gotBody); err != nil {
						t.Errorf("POST body is not a JSON object: %s: %q", err, raw)
					}
					_ = json.NewEncoder(w).Encode(map[string]any{
						"id": 1, "name": "keep", "location_name": "frn",
						"has_password": tc.hasPassword, "provisioning_status": "creating",
						"password": "Echoed-Pw-8Kd2xQ",
					})
				case http.MethodGet:
					_ = json.NewEncoder(w).Encode(map[string]any{
						"id": 1, "name": "keep", "location_name": "frn",
						"has_password": tc.hasPassword, "provisioning_status": "active",
					})
				default:
					t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
					w.WriteHeader(http.StatusMethodNotAllowed)
				}
			}))
			defer srv.Close()

			client := gcore.NewClient(
				option.WithBaseURL(srv.URL),
				option.WithAPIKey("unit-test"),
				option.WithMaxRetries(0),
			)

			r := storage_sftp.NewResource()
			configurable, ok := r.(resource.ResourceWithConfigure)
			if !ok {
				t.Fatal("the storage_sftp resource no longer implements Configure")
			}
			configureResp := &resource.ConfigureResponse{}
			configurable.Configure(ctx, resource.ConfigureRequest{ProviderData: &client}, configureResp)
			if configureResp.Diagnostics.HasError() {
				t.Fatalf("configuring the resource: %v", configureResp.Diagnostics)
			}

			common := map[string]tftypes.Value{"name": str("keep"), "location_name": str("frn")}
			withCommon := func(values map[string]tftypes.Value) map[string]tftypes.Value {
				merged := map[string]tftypes.Value{}
				for k, v := range common {
					merged[k] = v
				}
				for k, v := range values {
					merged[k] = v
				}
				return merged
			}

			req := resource.CreateRequest{
				Config: tfsdk.Config{Raw: storageSftpValues(ctx, t, withCommon(tc.config)), Schema: schema},
				Plan:   tfsdk.Plan{Raw: storageSftpValues(ctx, t, withCommon(tc.plan)), Schema: schema},
			}
			resp := &resource.CreateResponse{State: tfsdk.State{Raw: storageSftpNull(ctx, t), Schema: schema}}
			r.Create(ctx, req, resp)
			if resp.Diagnostics.HasError() {
				t.Fatalf("Create failed: %v", resp.Diagnostics.Errors())
			}

			mu.Lock()
			defer mu.Unlock()
			if posts != 1 {
				t.Fatalf("got %d POST requests, want exactly 1", posts)
			}
			if gotPath != "/storage/v4/sftp_storages" {
				t.Fatalf("POST path %s, want /storage/v4/sftp_storages", gotPath)
			}
			if !reflect.DeepEqual(gotBody, tc.wantBody) {
				t.Fatalf("POST body:\n got %v\nwant %v", gotBody, tc.wantBody)
			}

			var statePassword types.String
			var stateVersion, planVersion types.Int64
			var stateHasPassword types.Bool
			if d := resp.State.GetAttribute(ctx, path.Root("password_wo"), &statePassword); d.HasError() {
				t.Fatalf("reading password_wo from state: %v", d)
			}
			if !statePassword.IsNull() {
				t.Fatalf("state password_wo = %v, want null", statePassword)
			}
			if d := resp.State.GetAttribute(ctx, path.Root("password_wo_version"), &stateVersion); d.HasError() {
				t.Fatalf("reading password_wo_version from state: %v", d)
			}
			if d := req.Plan.GetAttribute(ctx, path.Root("password_wo_version"), &planVersion); d.HasError() {
				t.Fatalf("reading password_wo_version from plan: %v", d)
			}
			if !stateVersion.Equal(planVersion) {
				t.Fatalf("state password_wo_version = %v, want the planned %v", stateVersion, planVersion)
			}
			if d := resp.State.GetAttribute(ctx, path.Root("has_password"), &stateHasPassword); d.HasError() {
				t.Fatalf("reading has_password from state: %v", d)
			}
			if !stateHasPassword.Equal(types.BoolValue(tc.hasPassword)) {
				t.Fatalf("state has_password = %v, want %t", stateHasPassword, tc.hasPassword)
			}
		})
	}
}
