package customvalidator

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var credentialsAttrTypes = map[string]attr.Type{
	"ssh_key_name":        types.StringType,
	"username":            types.StringType,
	"password_wo":         types.StringType,
	"password_wo_version": types.Int64Type,
}

func credentials(t *testing.T, sshKey, username, passwordWo attr.Value, version attr.Value) types.Object {
	t.Helper()
	obj, diags := types.ObjectValue(credentialsAttrTypes, map[string]attr.Value{
		"ssh_key_name":        sshKey,
		"username":            username,
		"password_wo":         passwordWo,
		"password_wo_version": version,
	})
	if diags.HasError() {
		t.Fatalf("building credentials fixture: %v", diags)
	}
	return obj
}

// TestCredentialsValidator: one of ssh_key_name or username+password_wo, never
// both, never neither, and never an empty key name (the API requires
// min_length=1 and would otherwise answer with a generic validation error).
func TestCredentialsValidator(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	null := types.StringNull()
	nullVersion := types.Int64Null()

	testCases := []struct {
		name          string
		value         types.Object
		expectedError string
	}{
		{"ssh key only", credentials(t, types.StringValue("k1"), null, null, nullVersion), ""},
		{"username and password", credentials(t, null, types.StringValue("u"), types.StringValue("p"), types.Int64Value(1)), ""},
		{"unknown ssh key is accepted at validate time", credentials(t, types.StringUnknown(), null, null, nullVersion), ""},
		{"null object is accepted", types.ObjectNull(credentialsAttrTypes), ""},

		// Unknown values defer only the checks that depend on them.
		{
			// Reviewer's case: key and username are known and already conflict;
			// an unknown password must not hide that.
			"known conflict with unknown password",
			credentials(t, types.StringValue("existing-key"), types.StringValue("ubuntu"), types.StringUnknown(), types.Int64Value(1)),
			"Conflicting credentials configuration",
		},
		// An unknown cannot supply a value that is null: these can never
		// become valid, so they are rejected now.
		{"null key and username, unknown password", credentials(t, null, null, types.StringUnknown(), nullVersion), "Invalid credentials configuration"},
		{"null key and password, unknown username", credentials(t, null, types.StringUnknown(), null, nullVersion), "Invalid credentials configuration"},
		// The version check runs on its own, whatever else is unknown.
		{"password without version, unknown key", credentials(t, types.StringUnknown(), null, types.StringValue("p"), nullVersion), "Missing password_wo_version"},
		{"password without version, unknown username", credentials(t, null, types.StringUnknown(), types.StringValue("p"), nullVersion), "Missing password_wo_version"},
		{"known empty key with unknown password", credentials(t, types.StringValue(""), null, types.StringUnknown(), nullVersion), "Empty ssh_key_name"},
		{"unknown key with known username and password", credentials(t, types.StringUnknown(), types.StringValue("u"), types.StringValue("p"), types.Int64Value(1)), ""},
		{"unknown username completes the pair", credentials(t, null, types.StringUnknown(), types.StringValue("p"), types.Int64Value(1)), ""},
		{"unknown password completes the pair", credentials(t, null, types.StringValue("u"), types.StringUnknown(), types.Int64Value(1)), ""},
		{"password with unknown version", credentials(t, null, types.StringValue("u"), types.StringValue("p"), types.Int64Unknown()), ""},
		{"key with unknown version is deferred", credentials(t, types.StringValue("k1"), null, null, types.Int64Unknown()), ""},
		{"key with known version conflicts", credentials(t, types.StringValue("k1"), null, null, types.Int64Value(1)), "Conflicting credentials configuration"},
		{"empty ssh key", credentials(t, types.StringValue(""), null, null, nullVersion), "Empty ssh_key_name"},
		{"neither method", credentials(t, null, null, null, nullVersion), "Invalid credentials configuration"},
		{"both methods", credentials(t, types.StringValue("k1"), types.StringValue("u"), types.StringValue("p"), types.Int64Value(1)), "Conflicting credentials configuration"},
		{"password without version", credentials(t, null, types.StringValue("u"), types.StringValue("p"), nullVersion), "Missing password_wo_version"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			resp := &validator.ObjectResponse{}
			CredentialsValidator{}.ValidateObject(ctx, validator.ObjectRequest{
				Path:        path.Root("servers_settings").AtName("credentials"),
				ConfigValue: tc.value,
			}, resp)

			if tc.expectedError == "" {
				if resp.Diagnostics.HasError() {
					t.Fatalf("expected no error, got: %v", resp.Diagnostics)
				}
				return
			}
			if !resp.Diagnostics.HasError() {
				t.Fatalf("expected error %q, got none", tc.expectedError)
			}
			if got := resp.Diagnostics.Errors()[0].Summary(); got != tc.expectedError {
				t.Fatalf("expected %q, got %q", tc.expectedError, got)
			}
		})
	}
}

// TestValidateResolvedCredentials: at apply time nothing is unknown, so every
// rule the validator deferred now decides. These are the configurations that
// pass validation with an unknown and must be caught once it resolves.
func TestValidateResolvedCredentials(t *testing.T) {
	t.Parallel()

	null := types.StringNull()
	p := path.Root("servers_settings").AtName("credentials")

	testCases := []struct {
		name          string
		value         types.Object
		expectedError string
	}{
		{"key with version resolved to 1 conflicts", credentials(t, types.StringValue("existing-key"), null, null, types.Int64Value(1)), "Conflicting credentials configuration"},
		{"key with password resolved conflicts", credentials(t, types.StringValue("k1"), null, types.StringValue("p"), types.Int64Null()), "Conflicting credentials configuration"},
		{"password with version resolved to null", credentials(t, null, types.StringValue("u"), types.StringValue("p"), types.Int64Null()), "Missing password_wo_version"},
		{"key resolved to empty", credentials(t, types.StringValue(""), null, null, types.Int64Null()), "Empty ssh_key_name"},
		{"key resolved to null with no password method", credentials(t, null, null, null, types.Int64Null()), "Invalid credentials configuration"},
		{"key resolved fine", credentials(t, types.StringValue("k1"), null, null, types.Int64Null()), ""},
		{"password method resolved fine", credentials(t, null, types.StringValue("u"), types.StringValue("p"), types.Int64Value(1)), ""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			diags := ValidateResolvedCredentials(p, tc.value)
			if tc.expectedError == "" {
				if diags.HasError() {
					t.Fatalf("expected no error, got: %v", diags)
				}
				return
			}
			if !diags.HasError() {
				t.Fatalf("expected error %q, got none", tc.expectedError)
			}
			if got := diags.Errors()[0].Summary(); got != tc.expectedError {
				t.Fatalf("expected %q, got %q", tc.expectedError, got)
			}
		})
	}
}
