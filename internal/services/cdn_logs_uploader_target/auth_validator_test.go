package cdn_logs_uploader_target

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var authConfigAttrTypes = map[string]attr.Type{
	"token":             types.StringType,
	"header_name":       types.StringType,
	"account_key":       types.StringType,
	"access_key_id":     types.StringType,
	"secret_access_key": types.StringType,
}

var authObjectAttrTypes = map[string]attr.Type{
	"config": types.ObjectType{AttrTypes: authConfigAttrTypes},
	"type":   types.StringType,
}

// buildAuth constructs an auth object value for the given type and the set
// config fields (fields not present default to null).
func buildAuth(authType attr.Value, fields map[string]string) types.Object {
	cfg := map[string]attr.Value{}
	for name := range authConfigAttrTypes {
		if v, ok := fields[name]; ok {
			cfg[name] = types.StringValue(v)
		} else {
			cfg[name] = types.StringNull()
		}
	}
	return types.ObjectValueMust(authObjectAttrTypes, map[string]attr.Value{
		"config": types.ObjectValueMust(authConfigAttrTypes, cfg),
		"type":   authType,
	})
}

func TestAuthConfigValidator(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		auth    types.Object
		wantErr bool
	}{
		"token with header_name (valid HTTP)": {
			auth:    buildAuth(types.StringValue("token"), map[string]string{"token": "t", "header_name": "Authorization"}),
			wantErr: false,
		},
		"shared_key with header_name (Sofia's case)": {
			auth:    buildAuth(types.StringValue("shared_key"), map[string]string{"account_key": "k", "header_name": "X"}),
			wantErr: true,
		},
		"ak_sk with header_name": {
			auth:    buildAuth(types.StringValue("ak_sk"), map[string]string{"access_key_id": "a", "secret_access_key": "s", "header_name": "X"}),
			wantErr: true,
		},
		"sas_token with token (valid)": {
			auth:    buildAuth(types.StringValue("sas_token"), map[string]string{"token": "t"}),
			wantErr: false,
		},
		"sas_token with account_key (invalid)": {
			auth:    buildAuth(types.StringValue("sas_token"), map[string]string{"account_key": "k"}),
			wantErr: true,
		},
		"shared_key with account_key (valid)": {
			auth:    buildAuth(types.StringValue("shared_key"), map[string]string{"account_key": "k"}),
			wantErr: false,
		},
		"ak_sk with access/secret keys (valid)": {
			auth:    buildAuth(types.StringValue("ak_sk"), map[string]string{"access_key_id": "a", "secret_access_key": "s"}),
			wantErr: false,
		},
		"token case-insensitive (valid)": {
			auth:    buildAuth(types.StringValue("TOKEN"), map[string]string{"token": "t", "header_name": "X"}),
			wantErr: false,
		},
		"unknown type skips validation": {
			auth:    buildAuth(types.StringUnknown(), map[string]string{"header_name": "X"}),
			wantErr: false,
		},
		"null type skips validation": {
			auth:    buildAuth(types.StringNull(), map[string]string{"header_name": "X"}),
			wantErr: false,
		},
		"empty string field ignored": {
			auth:    buildAuth(types.StringValue("shared_key"), map[string]string{"account_key": "k", "header_name": ""}),
			wantErr: false,
		},
	}

	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			resp := &validator.ObjectResponse{}
			authConfigValidator{}.ValidateObject(context.Background(), validator.ObjectRequest{
				ConfigValue: tc.auth,
			}, resp)
			if got := resp.Diagnostics.HasError(); got != tc.wantErr {
				t.Fatalf("wantErr=%v, got error=%v (%v)", tc.wantErr, got, resp.Diagnostics)
			}
		})
	}
}

// null config skips validation (separate because buildAuth always builds a config).
func TestAuthConfigValidator_NullConfig(t *testing.T) {
	t.Parallel()
	auth := types.ObjectValueMust(authObjectAttrTypes, map[string]attr.Value{
		"config": basetypes.NewObjectNull(authConfigAttrTypes),
		"type":   types.StringValue("shared_key"),
	})
	resp := &validator.ObjectResponse{}
	authConfigValidator{}.ValidateObject(context.Background(), validator.ObjectRequest{
		ConfigValue: auth,
	}, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("expected no error for null config, got %v", resp.Diagnostics)
	}
}
