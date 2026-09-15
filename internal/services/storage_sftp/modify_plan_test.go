package storage_sftp_test

import (
	"context"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/G-Core/terraform-provider-gcore/internal/services/storage_sftp"
)

func TestStorageSftpModifyPlanPasswordDecisions(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	schema := storage_sftp.ResourceSchema(ctx)

	unknownPassword := tftypes.NewValue(tftypes.String, tftypes.UnknownValue)

	unknownBool := tftypes.NewValue(tftypes.Bool, tftypes.UnknownValue)
	unknownNum := tftypes.NewValue(tftypes.Number, tftypes.UnknownValue)

	// these cases must prove plan wiring: diagnostics, the planned
	// has_password, and the adoption warning. passwordActionFor's input space
	// lives in TestPasswordActionFor; do not re-enumerate it here
	cases := map[string]struct {
		config    map[string]tftypes.Value
		plan      map[string]tftypes.Value
		state     map[string]tftypes.Value
		nullPlan  bool
		nullState bool
		expectErr *regexp.Regexp
		// A non-nil value requires one matching warning.
		expectWarn *regexp.Regexp
		// nil leaves the framework value unchanged.
		hasPassword *bool
	}{
		"destroy is never blocked": {
			nullPlan: true,
			state: map[string]tftypes.Value{
				"name": str("keep"), "password_wo_version": num(1), "has_password": boolean(true),
			},
		},
		"create with both arguments": {
			nullState: true,
			config: map[string]tftypes.Value{
				"name": str("new"), "password_wo": unknownPassword, "password_wo_version": num(1),
			},
			plan:        map[string]tftypes.Value{"name": str("new"), "password_wo_version": num(1), "has_password": unknownBool},
			hasPassword: ptr(true),
		},
		"create with neither argument": {
			nullState:   true,
			config:      map[string]tftypes.Value{"name": str("new")},
			plan:        map[string]tftypes.Value{"name": str("new"), "has_password": unknownBool},
			hasPassword: ptr(false),
		},
		"create with a version and no password": {
			nullState: true,
			config:    map[string]tftypes.Value{"name": str("new"), "password_wo_version": num(1)},
			plan:      map[string]tftypes.Value{"name": str("new"), "password_wo_version": num(1), "has_password": unknownBool},
			expectErr: regexp.MustCompile(`Missing password_wo`),
		},
		"dropping only the version": {
			config: map[string]tftypes.Value{"name": str("keep"), "password_wo": unknownPassword},
			plan:   map[string]tftypes.Value{"name": str("keep"), "has_password": unknownBool},
			state: map[string]tftypes.Value{
				"name": str("keep"), "password_wo_version": num(1), "has_password": boolean(true),
			},
			expectErr: regexp.MustCompile(`Missing password_wo_version`),
		},
		"dropping both arguments": {
			config: map[string]tftypes.Value{"name": str("keep")},
			plan:   map[string]tftypes.Value{"name": str("keep"), "has_password": unknownBool},
			state: map[string]tftypes.Value{
				"name": str("keep"), "password_wo_version": num(1), "has_password": boolean(true),
			},
			hasPassword: ptr(false),
		},
		"rotation": {
			config: map[string]tftypes.Value{
				"name": str("keep"), "password_wo": unknownPassword, "password_wo_version": num(2),
			},
			plan: map[string]tftypes.Value{"name": str("keep"), "password_wo_version": num(2), "has_password": unknownBool},
			state: map[string]tftypes.Value{
				"name": str("keep"), "password_wo_version": num(1), "has_password": boolean(true),
			},
			hasPassword: ptr(true),
		},
		"unrelated update with the pair unchanged": {
			config: map[string]tftypes.Value{
				"name": str("keep"), "password_wo": unknownPassword, "password_wo_version": num(1), "is_http_disabled": boolean(true),
			},
			plan: map[string]tftypes.Value{
				"name": str("keep"), "password_wo_version": num(1), "is_http_disabled": boolean(true), "has_password": unknownBool,
			},
			state: map[string]tftypes.Value{
				"name": str("keep"), "password_wo_version": num(1), "is_http_disabled": boolean(false), "has_password": boolean(true),
			},
		},
		"adopted storage with a password": {
			config: map[string]tftypes.Value{
				"name": str("keep"), "password_wo": unknownPassword, "password_wo_version": num(1),
			},
			plan:        map[string]tftypes.Value{"name": str("keep"), "password_wo_version": num(1), "has_password": unknownBool},
			state:       map[string]tftypes.Value{"name": str("keep"), "has_password": boolean(true)},
			hasPassword: ptr(true),
			expectWarn:  regexp.MustCompile(`Replacing a password Terraform did not set`),
		},
		"password removed out of band": {
			config: map[string]tftypes.Value{
				"name": str("keep"), "password_wo": unknownPassword, "password_wo_version": num(1),
			},
			plan: map[string]tftypes.Value{"name": str("keep"), "password_wo_version": num(1), "has_password": boolean(false)},
			state: map[string]tftypes.Value{
				"name": str("keep"), "password_wo_version": num(1), "has_password": boolean(false),
			},
			hasPassword: ptr(true),
		},
		"adopted storage renamed and given a password in one plan": {
			config: map[string]tftypes.Value{
				"name": str("renamed"), "password_wo": unknownPassword, "password_wo_version": num(1),
			},
			plan:        map[string]tftypes.Value{"name": str("renamed"), "password_wo_version": num(1), "has_password": unknownBool},
			state:       map[string]tftypes.Value{"name": str("keep"), "has_password": boolean(true)},
			hasPassword: ptr(true),
		},
		"planned version unknown": {
			config: map[string]tftypes.Value{
				"name": str("keep"), "password_wo": unknownPassword, "password_wo_version": unknownNum,
			},
			plan: map[string]tftypes.Value{"name": str("keep"), "password_wo_version": unknownNum, "has_password": unknownBool},
			state: map[string]tftypes.Value{
				"name": str("keep"), "password_wo_version": num(1), "has_password": boolean(true),
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			configRaw := storageSftpValues(ctx, t, tc.config)
			planRaw := storageSftpValues(ctx, t, tc.plan)
			stateRaw := storageSftpValues(ctx, t, tc.state)
			if tc.nullPlan {
				configRaw = storageSftpNull(ctx, t)
				planRaw = storageSftpNull(ctx, t)
			}
			if tc.nullState {
				stateRaw = storageSftpNull(ctx, t)
			}

			modifier, ok := storage_sftp.NewResource().(resource.ResourceWithModifyPlan)
			if !ok {
				t.Fatal("the storage_sftp resource no longer implements ModifyPlan")
			}

			req := resource.ModifyPlanRequest{
				Config: tfsdk.Config{Raw: configRaw, Schema: schema},
				Plan:   tfsdk.Plan{Raw: planRaw, Schema: schema},
				State:  tfsdk.State{Raw: stateRaw, Schema: schema},
			}
			resp := &resource.ModifyPlanResponse{Plan: tfsdk.Plan{Raw: planRaw, Schema: schema}}
			modifier.ModifyPlan(ctx, req, resp)

			if tc.expectErr != nil {
				for _, d := range resp.Diagnostics.Errors() {
					if tc.expectErr.MatchString(d.Summary()) || tc.expectErr.MatchString(d.Detail()) {
						return
					}
				}
				t.Fatalf("expected an error matching %s, got: %v", tc.expectErr, resp.Diagnostics.Errors())
			}
			if resp.Diagnostics.HasError() {
				t.Fatalf("expected a clean plan, got: %v", resp.Diagnostics.Errors())
			}

			warnings := resp.Diagnostics.Warnings()
			switch {
			case tc.expectWarn == nil && len(warnings) != 0:
				t.Fatalf("expected no warning, got: %v", warnings)
			case tc.expectWarn != nil && len(warnings) != 1:
				t.Fatalf("expected exactly one warning matching %s, got: %v", tc.expectWarn, warnings)
			case tc.expectWarn != nil && !tc.expectWarn.MatchString(warnings[0].Summary()):
				t.Fatalf("expected a warning matching %s, got: %s", tc.expectWarn, warnings[0].Summary())
			}

			if tc.nullPlan {
				if !resp.Plan.Raw.IsNull() {
					t.Fatalf("a destroy plan was written to: %v", resp.Plan.Raw)
				}
				return
			}

			var got, given types.Bool
			if diags := resp.Plan.GetAttribute(ctx, path.Root("has_password"), &got); diags.HasError() {
				t.Fatalf("reading planned has_password: %v", diags)
			}
			if diags := req.Plan.GetAttribute(ctx, path.Root("has_password"), &given); diags.HasError() {
				t.Fatalf("reading given has_password: %v", diags)
			}
			want := given
			if tc.hasPassword != nil {
				want = types.BoolValue(*tc.hasPassword)
			}
			if !got.Equal(want) {
				t.Fatalf("planned has_password = %v, want %v", got, want)
			}
		})
	}
}
