package storage_sftp

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestPasswordActionFor(t *testing.T) {
	t.Parallel()

	configured := types.StringValue("Unit-Test-Pw-8Kd2xQ")
	ephemeral := types.StringUnknown() // ephemeral.random_password.x.result at plan time
	absent := types.StringNull()

	cases := map[string]struct {
		creating       bool
		configPassword types.String
		planVersion    types.Int64
		stateVersion   types.Int64
		stateHasPw     types.Bool
		wantAction     passwordAction
		wantDecided    bool
	}{
		"create with the pair": {
			creating: true, configPassword: configured, planVersion: types.Int64Value(1),
			wantAction: passwordSet, wantDecided: true,
		},
		"create with an ephemeral pair": {
			creating: true, configPassword: ephemeral, planVersion: types.Int64Value(1),
			wantAction: passwordSet, wantDecided: true,
		},
		"create without the pair": {
			creating: true, configPassword: absent,
			wantAction: passwordClear, wantDecided: true,
		},

		// also the out-of-band rotation case: it changes nothing Terraform sees
		"unrelated update, pair unchanged": {
			configPassword: configured, planVersion: types.Int64Value(1), stateVersion: types.Int64Value(1), stateHasPw: types.BoolValue(true),
			wantAction: passwordUntouched, wantDecided: true,
		},

		"rotation, version increased": {
			configPassword: configured, planVersion: types.Int64Value(2), stateVersion: types.Int64Value(1), stateHasPw: types.BoolValue(true),
			wantAction: passwordSet, wantDecided: true,
		},
		"rotation, version decreased": {
			configPassword: configured, planVersion: types.Int64Value(1), stateVersion: types.Int64Value(2), stateHasPw: types.BoolValue(true),
			wantAction: passwordSet, wantDecided: true,
		},

		"password removed out of band": {
			configPassword: configured, planVersion: types.Int64Value(1), stateVersion: types.Int64Value(1), stateHasPw: types.BoolValue(false),
			wantAction: passwordSet, wantDecided: true,
		},

		"pair removed": {
			configPassword: absent, planVersion: types.Int64Null(), stateVersion: types.Int64Value(1), stateHasPw: types.BoolValue(true),
			wantAction: passwordClear, wantDecided: true,
		},
		"pair removed, password already gone out of band": {
			configPassword: absent, planVersion: types.Int64Null(), stateVersion: types.Int64Value(1), stateHasPw: types.BoolValue(false),
			wantAction: passwordClear, wantDecided: true,
		},

		"adopted storage with a password takes the pair": {
			configPassword: configured, planVersion: types.Int64Value(1), stateVersion: types.Int64Null(), stateHasPw: types.BoolValue(true),
			wantAction: passwordSet, wantDecided: true,
		},
		"adopted storage without a password takes the pair": {
			configPassword: configured, planVersion: types.Int64Value(1), stateVersion: types.Int64Null(), stateHasPw: types.BoolValue(false),
			wantAction: passwordSet, wantDecided: true,
		},
		"adopted storage with unreadable has_password takes the pair": {
			configPassword: configured, planVersion: types.Int64Value(1), stateVersion: types.Int64Null(), stateHasPw: types.BoolNull(),
			wantAction: passwordSet, wantDecided: true,
		},
		"adopted storage with unknown has_password takes the pair": {
			configPassword: configured, planVersion: types.Int64Value(1), stateVersion: types.Int64Null(), stateHasPw: types.BoolUnknown(),
			wantAction: passwordSet, wantDecided: true,
		},

		"adopted storage with a password, no pair": {
			configPassword: absent, planVersion: types.Int64Null(), stateVersion: types.Int64Null(), stateHasPw: types.BoolValue(true),
			wantAction: passwordUntouched, wantDecided: true,
		},
		"adopted storage without a password, no pair": {
			configPassword: absent, planVersion: types.Int64Null(), stateVersion: types.Int64Null(), stateHasPw: types.BoolValue(false),
			wantAction: passwordUntouched, wantDecided: true,
		},

		"managed storage with unreadable has_password, versions equal": {
			configPassword: configured, planVersion: types.Int64Value(1), stateVersion: types.Int64Value(1), stateHasPw: types.BoolNull(),
			wantAction: passwordUntouched, wantDecided: true,
		},
		"managed storage with unknown has_password, versions equal": {
			configPassword: configured, planVersion: types.Int64Value(1), stateVersion: types.Int64Value(1), stateHasPw: types.BoolUnknown(),
			wantAction: passwordUntouched, wantDecided: true,
		},

		"password configured without a version": {
			configPassword: configured, planVersion: types.Int64Null(), stateVersion: types.Int64Value(1), stateHasPw: types.BoolValue(true),
			wantAction: passwordUntouched, wantDecided: true,
		},

		"planned version unknown": {
			configPassword: configured, planVersion: types.Int64Unknown(), stateVersion: types.Int64Value(1), stateHasPw: types.BoolValue(true),
			wantAction: passwordUntouched, wantDecided: false,
		},
		"planned version unknown, has_password unreadable": {
			configPassword: configured, planVersion: types.Int64Unknown(), stateVersion: types.Int64Value(1), stateHasPw: types.BoolNull(),
			wantAction: passwordUntouched, wantDecided: false,
		},
		"planned version unknown on an adopted storage with a password": {
			configPassword: configured, planVersion: types.Int64Unknown(), stateVersion: types.Int64Null(), stateHasPw: types.BoolValue(true),
			wantAction: passwordSet, wantDecided: true,
		},
		"planned version unknown on an adopted storage without a password": {
			configPassword: configured, planVersion: types.Int64Unknown(), stateVersion: types.Int64Null(), stateHasPw: types.BoolValue(false),
			wantAction: passwordSet, wantDecided: true,
		},
		"planned version unknown on an adopted storage with unreadable has_password": {
			configPassword: configured, planVersion: types.Int64Unknown(), stateVersion: types.Int64Null(), stateHasPw: types.BoolNull(),
			wantAction: passwordSet, wantDecided: true,
		},
		"planned version unknown after out-of-band removal": {
			configPassword: configured, planVersion: types.Int64Unknown(), stateVersion: types.Int64Value(1), stateHasPw: types.BoolValue(false),
			wantAction: passwordSet, wantDecided: true,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			action, decided := passwordActionFor(tc.creating, tc.configPassword, tc.planVersion, tc.stateVersion, tc.stateHasPw)
			if action != tc.wantAction || decided != tc.wantDecided {
				t.Fatalf("got (%s, decided=%t), want (%s, decided=%t)", action, decided, tc.wantAction, tc.wantDecided)
			}
		})
	}
}
