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

func TestStorageSftpUpdateSendsPasswordDecision(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	schema := storage_sftp.ResourceSchema(ctx)

	const password = "Unit-Test-Pw-8Kd2xQ"
	unknownBool := tftypes.NewValue(tftypes.Bool, tftypes.UnknownValue)

	// one case per PATCH body shape. Which inputs produce which action is
	// TestPasswordActionFor's job; do not re-enumerate the decision here
	cases := map[string]struct {
		config           map[string]tftypes.Value
		plan             map[string]tftypes.Value
		state            map[string]tftypes.Value
		hasPasswordAfter bool
		wantBody         map[string]any
	}{
		// set with a non-empty marshalled body, so the password fields must
		// survive option ordering
		"rotation with a sibling change": {
			config: map[string]tftypes.Value{
				"password_wo": str(password), "password_wo_version": num(2), "is_http_disabled": boolean(true),
			},
			plan: map[string]tftypes.Value{
				"password_wo_version": num(2), "is_http_disabled": boolean(true), "has_password": boolean(true),
			},
			state: map[string]tftypes.Value{
				"password_wo_version": num(1), "is_http_disabled": boolean(false), "has_password": boolean(true),
			},

			hasPasswordAfter: true,
			wantBody:         map[string]any{"is_http_disabled": true, "password_mode": "set", "password": password},
		},
		// none with an empty marshalled body, and no password field may appear
		"clear": {
			config: map[string]tftypes.Value{},
			plan:   map[string]tftypes.Value{"has_password": boolean(false)},
			state:  map[string]tftypes.Value{"password_wo_version": num(1), "has_password": boolean(true)},

			hasPasswordAfter: false,
			wantBody:         map[string]any{"password_mode": "none"},
		},
		// untouched: the body must carry no password field at all
		"sibling-only change with the pair unchanged": {
			config: map[string]tftypes.Value{
				"password_wo": str(password), "password_wo_version": num(1), "is_http_disabled": boolean(true),
			},
			plan: map[string]tftypes.Value{
				"password_wo_version": num(1), "is_http_disabled": boolean(true), "has_password": unknownBool,
			},
			state: map[string]tftypes.Value{
				"password_wo_version": num(1), "is_http_disabled": boolean(false), "has_password": boolean(true),
			},

			hasPasswordAfter: true,
			wantBody:         map[string]any{"is_http_disabled": true},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			var (
				mu       sync.Mutex
				patches  int
				gotPath  string
				gotBody  map[string]any
				echoBody = map[string]any{
					"id": 1, "name": "keep", "location_name": "frn",
					"has_password": tc.hasPasswordAfter, "provisioning_status": "updating",
					"password": "Echoed-Pw-8Kd2xQ",
				}
			)
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				mu.Lock()
				defer mu.Unlock()
				w.Header().Set("Content-Type", "application/json")
				switch r.Method {
				case http.MethodPatch:
					patches++
					gotPath = r.URL.Path
					raw, _ := io.ReadAll(r.Body)
					if err := json.Unmarshal(raw, &gotBody); err != nil {
						t.Errorf("PATCH body is not a JSON object: %s: %q", err, raw)
					}
					_ = json.NewEncoder(w).Encode(echoBody)
				case http.MethodGet:
					_ = json.NewEncoder(w).Encode(map[string]any{
						"id": 1, "name": "keep", "location_name": "frn",
						"has_password": tc.hasPasswordAfter, "provisioning_status": "active",
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

			common := map[string]tftypes.Value{"id": num(1), "name": str("keep"), "location_name": str("frn")}
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

			req := resource.UpdateRequest{
				Config: tfsdk.Config{Raw: storageSftpValues(ctx, t, withCommon(tc.config)), Schema: schema},
				Plan:   tfsdk.Plan{Raw: storageSftpValues(ctx, t, withCommon(tc.plan)), Schema: schema},
				State:  tfsdk.State{Raw: storageSftpValues(ctx, t, withCommon(tc.state)), Schema: schema},
			}
			resp := &resource.UpdateResponse{State: tfsdk.State{Raw: storageSftpNull(ctx, t), Schema: schema}}
			r.Update(ctx, req, resp)
			if resp.Diagnostics.HasError() {
				t.Fatalf("Update failed: %v", resp.Diagnostics.Errors())
			}

			mu.Lock()
			defer mu.Unlock()
			if patches != 1 {
				t.Fatalf("got %d PATCH requests, want exactly 1", patches)
			}
			if gotPath != "/storage/v4/sftp_storages/1" {
				t.Fatalf("PATCH path %s, want /storage/v4/sftp_storages/1", gotPath)
			}
			if !reflect.DeepEqual(gotBody, tc.wantBody) {
				t.Fatalf("PATCH body:\n got %v\nwant %v", gotBody, tc.wantBody)
			}

			var stateVersion, planVersion types.Int64
			var stateHasPassword types.Bool
			var statePassword types.String
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
			if !stateHasPassword.Equal(types.BoolValue(tc.hasPasswordAfter)) {
				t.Fatalf("state has_password = %v, want %t", stateHasPassword, tc.hasPasswordAfter)
			}
			if d := resp.State.GetAttribute(ctx, path.Root("password_wo"), &statePassword); d.HasError() {
				t.Fatalf("reading password_wo from state: %v", d)
			}
			if !statePassword.IsNull() {
				t.Fatalf("state password_wo = %v, want null", statePassword)
			}
		})
	}
}
