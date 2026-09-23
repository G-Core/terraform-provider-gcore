package cloud_gpu_baremetal_cluster

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/G-Core/gcore-go"
	"github.com/G-Core/gcore-go/cloud"
	"github.com/G-Core/gcore-go/option"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// serversSettingsWith extends serversSettingsWithCredentials
// (modify_plan_unknown_test.go) with overrides on servers_settings itself,
// which these tests need for user_data.
func serversSettingsWith(ctx context.Context, t *testing.T, settingsOverrides, credOverrides map[string]tftypes.Value) tftypes.Value {
	t.Helper()
	base := serversSettingsWithCredentials(ctx, t, credOverrides)

	var vals map[string]tftypes.Value
	if err := base.As(&vals); err != nil {
		t.Fatalf("decoding servers_settings fixture: %v", err)
	}
	for name, val := range settingsOverrides {
		vals[name] = val
	}

	return tftypes.NewValue(base.Type(), vals)
}

// serversSettingsWithNullCredentials builds a servers_settings whose
// credentials object is NULL, the shape a config that drops the block
// produces, as opposed to serversSettingsWith's known object of nulls.
func serversSettingsWithNullCredentials(ctx context.Context, t *testing.T) tftypes.Value {
	t.Helper()
	base := serversSettingsWith(ctx, t, nil, nil)

	var vals map[string]tftypes.Value
	if err := base.As(&vals); err != nil {
		t.Fatalf("decoding servers_settings fixture: %v", err)
	}
	vals["credentials"] = tftypes.NewValue(vals["credentials"].Type(), nil)

	return tftypes.NewValue(base.Type(), vals)
}

// unknownServersSettings returns a wholly unknown servers_settings value.
func unknownServersSettings(ctx context.Context, t *testing.T) tftypes.Value {
	t.Helper()
	sch := ResourceSchema(ctx)
	objType := sch.Type().TerraformType(ctx).(tftypes.Object)
	return tftypes.NewValue(objType.AttributeTypes["servers_settings"], tftypes.UnknownValue)
}

// unknownCredentials returns a servers_settings whose credentials object is
// wholly unknown.
func unknownCredentials(ctx context.Context, t *testing.T) tftypes.Value {
	t.Helper()
	base := serversSettingsWith(ctx, t, nil, nil)

	var vals map[string]tftypes.Value
	if err := base.As(&vals); err != nil {
		t.Fatalf("decoding servers_settings fixture: %v", err)
	}
	vals["credentials"] = tftypes.NewValue(vals["credentials"].Type(), tftypes.UnknownValue)

	return tftypes.NewValue(base.Type(), vals)
}

// pathToSSHKeyName is the full nested path the walkers deliberately do NOT use.
func pathToSSHKeyName() path.Path {
	return path.Root("servers_settings").AtName("credentials").AtName("ssh_key_name")
}

func str(s string) tftypes.Value { return tftypes.NewValue(tftypes.String, s) }
func unknownStr() tftypes.Value  { return tftypes.NewValue(tftypes.String, tftypes.UnknownValue) }
func num(n int64) tftypes.Value  { return tftypes.NewValue(tftypes.Number, n) }
func decide(ctx context.Context, t *testing.T, planRaw, stateRaw tftypes.Value) rebuildDecision {
	t.Helper()
	sch := ResourceSchema(ctx)
	plan := tfsdk.Plan{Raw: planRaw, Schema: sch}
	state := tfsdk.State{Raw: stateRaw, Schema: sch}
	return rebuildsOnApply(rebuildTriggersOf(ctx, plan), rebuildTriggersOf(ctx, state))
}

// TestRebuildsOnApply pins the decision that both ModifyPlan and Update read.
// The three triggers are exactly the fields the API hashes into a cluster's
// spec; username and password_wo_version are deliberately not among them.
func TestRebuildsOnApply(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	testCases := []struct {
		name     string
		plan     map[string]tftypes.Value
		state    map[string]tftypes.Value
		expected rebuildDecision
	}{
		{
			name:     "nothing changes",
			plan:     map[string]tftypes.Value{"image_id": str("img-a")},
			state:    map[string]tftypes.Value{"image_id": str("img-a")},
			expected: rebuildNo,
		},
		{
			name:     "image_id changes",
			plan:     map[string]tftypes.Value{"image_id": str("img-b")},
			state:    map[string]tftypes.Value{"image_id": str("img-a")},
			expected: rebuildYes,
		},
		{
			name:     "ssh_key_name changes",
			plan:     map[string]tftypes.Value{"servers_settings": serversSettingsWith(ctx, t, nil, map[string]tftypes.Value{"ssh_key_name": str("k2")})},
			state:    map[string]tftypes.Value{"servers_settings": serversSettingsWith(ctx, t, nil, map[string]tftypes.Value{"ssh_key_name": str("k1")})},
			expected: rebuildYes,
		},
		{
			name:     "user_data changes",
			plan:     map[string]tftypes.Value{"servers_settings": serversSettingsWith(ctx, t, map[string]tftypes.Value{"user_data": str("Y")}, nil)},
			state:    map[string]tftypes.Value{"servers_settings": serversSettingsWith(ctx, t, map[string]tftypes.Value{"user_data": str("X")}, nil)},
			expected: rebuildYes,
		},
		{
			// The decision keys on ssh_key_name, not on credentialsHaveChanged,
			// so a username change raises the blocking error alone and no
			// re-image warning prints beside it.
			name:     "username only",
			plan:     map[string]tftypes.Value{"servers_settings": serversSettingsWith(ctx, t, nil, map[string]tftypes.Value{"username": str("u2")})},
			state:    map[string]tftypes.Value{"servers_settings": serversSettingsWith(ctx, t, nil, map[string]tftypes.Value{"username": str("u1")})},
			expected: rebuildNo,
		},
		{
			name:     "password_wo_version only",
			plan:     map[string]tftypes.Value{"servers_settings": serversSettingsWith(ctx, t, nil, map[string]tftypes.Value{"password_wo_version": num(2)})},
			state:    map[string]tftypes.Value{"servers_settings": serversSettingsWith(ctx, t, nil, map[string]tftypes.Value{"password_wo_version": num(1)})},
			expected: rebuildNo,
		},
		{
			name:     "servers_settings wholly unknown",
			plan:     map[string]tftypes.Value{"servers_settings": unknownServersSettings(ctx, t)},
			state:    map[string]tftypes.Value{"servers_settings": serversSettingsWith(ctx, t, nil, map[string]tftypes.Value{"ssh_key_name": str("k1")})},
			expected: rebuildUnknown,
		},
		{
			name:     "credentials object unknown",
			plan:     map[string]tftypes.Value{"image_id": str("img-a"), "servers_settings": unknownCredentials(ctx, t)},
			state:    map[string]tftypes.Value{"image_id": str("img-a"), "servers_settings": serversSettingsWith(ctx, t, nil, map[string]tftypes.Value{"ssh_key_name": str("k1")})},
			expected: rebuildUnknown,
		},
		{
			name:     "ssh_key_name leaf unknown",
			plan:     map[string]tftypes.Value{"servers_settings": serversSettingsWith(ctx, t, nil, map[string]tftypes.Value{"ssh_key_name": unknownStr()})},
			state:    map[string]tftypes.Value{"servers_settings": serversSettingsWith(ctx, t, nil, map[string]tftypes.Value{"ssh_key_name": str("k1")})},
			expected: rebuildUnknown,
		},
		{
			name:     "image_id unknown",
			plan:     map[string]tftypes.Value{"image_id": unknownStr()},
			state:    map[string]tftypes.Value{"image_id": str("img-a")},
			expected: rebuildUnknown,
		},
		{
			// Unknown outranks Yes: Update may not act while any single
			// trigger is unresolved.
			name: "image_id unknown and ssh_key_name changed",
			plan: map[string]tftypes.Value{
				"image_id":         unknownStr(),
				"servers_settings": serversSettingsWith(ctx, t, nil, map[string]tftypes.Value{"ssh_key_name": str("k2")}),
			},
			state: map[string]tftypes.Value{
				"image_id":         str("img-a"),
				"servers_settings": serversSettingsWith(ctx, t, nil, map[string]tftypes.Value{"ssh_key_name": str("k1")}),
			},
			expected: rebuildUnknown,
		},
		{
			name:     "first key on a cluster that had none",
			plan:     map[string]tftypes.Value{"servers_settings": serversSettingsWith(ctx, t, nil, map[string]tftypes.Value{"ssh_key_name": str("k1")})},
			state:    map[string]tftypes.Value{"servers_settings": serversSettingsWith(ctx, t, nil, nil)},
			expected: rebuildYes,
		},
		{
			// A null planned value never re-images; sshKeyRemoved handles
			// the removal separately.
			name:     "ssh_key_name unset in plan over a known key",
			plan:     map[string]tftypes.Value{"servers_settings": serversSettingsWith(ctx, t, nil, nil)},
			state:    map[string]tftypes.Value{"servers_settings": serversSettingsWith(ctx, t, nil, map[string]tftypes.Value{"ssh_key_name": str("k1")})},
			expected: rebuildNo,
		},
		{
			name:     "user_data null over null",
			plan:     map[string]tftypes.Value{"servers_settings": serversSettingsWith(ctx, t, nil, nil)},
			state:    map[string]tftypes.Value{"servers_settings": serversSettingsWith(ctx, t, nil, nil)},
			expected: rebuildNo,
		},
		{
			name:     "user_data null over set",
			plan:     map[string]tftypes.Value{"servers_settings": serversSettingsWith(ctx, t, nil, nil)},
			state:    map[string]tftypes.Value{"servers_settings": serversSettingsWith(ctx, t, map[string]tftypes.Value{"user_data": str("X")}, nil)},
			expected: rebuildNo,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := decide(ctx, t, rawObjectWith(ctx, t, tc.plan), rawObjectWith(ctx, t, tc.state))
			if got != tc.expected {
				t.Fatalf("expected %s, got %s", tc.expected, got)
			}
		})
	}
}

func TestRebuildDecisionFor(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		plan, state types.String
		expected    rebuildDecision
	}{
		{"plan null over known", types.StringNull(), types.StringValue("a"), rebuildNo},
		{"plan null over null", types.StringNull(), types.StringNull(), rebuildNo},
		{"plan null over unknown", types.StringNull(), types.StringUnknown(), rebuildNo},
		{"plan unknown", types.StringUnknown(), types.StringValue("a"), rebuildUnknown},
		{"state unknown", types.StringValue("a"), types.StringUnknown(), rebuildUnknown},
		{"both unknown", types.StringUnknown(), types.StringUnknown(), rebuildUnknown},
		{"equal", types.StringValue("a"), types.StringValue("a"), rebuildNo},
		{"changed", types.StringValue("b"), types.StringValue("a"), rebuildYes},
		{"set over null", types.StringValue("a"), types.StringNull(), rebuildYes},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if got := rebuildDecisionFor(tc.plan, tc.state); got != tc.expected {
				t.Fatalf("expected %s, got %s", tc.expected, got)
			}
		})
	}
}

func TestLeastResolved(t *testing.T) {
	t.Parallel()

	all := []rebuildDecision{rebuildNo, rebuildYes, rebuildUnknown}
	expected := map[rebuildDecision]map[rebuildDecision]rebuildDecision{
		rebuildNo:      {rebuildNo: rebuildNo, rebuildYes: rebuildYes, rebuildUnknown: rebuildUnknown},
		rebuildYes:     {rebuildNo: rebuildYes, rebuildYes: rebuildYes, rebuildUnknown: rebuildUnknown},
		rebuildUnknown: {rebuildNo: rebuildUnknown, rebuildYes: rebuildUnknown, rebuildUnknown: rebuildUnknown},
	}

	for _, a := range all {
		for _, b := range all {
			if got := leastResolved(a, b); got != expected[a][b] {
				t.Fatalf("leastResolved(%s, %s) = %s, want %s", a, b, got, expected[a][b])
			}
		}
	}
}

// TestStringAtObjectAtFailClosed: a value the walkers cannot read is Unknown,
// never Null; null would read as "no change".
func TestStringAtObjectAtFailClosed(t *testing.T) {
	t.Parallel()

	strType := map[string]attr.Type{"ssh_key_name": types.StringType}
	objType := map[string]attr.Type{"credentials": types.ObjectType{AttrTypes: strType}}

	t.Run("stringAt", func(t *testing.T) {
		t.Parallel()
		known, _ := types.ObjectValue(strType, map[string]attr.Value{"ssh_key_name": types.StringValue("k")})

		if got := stringAt(types.ObjectNull(strType), "ssh_key_name"); !got.IsNull() {
			t.Fatalf("null parent: expected null, got %v", got)
		}
		if got := stringAt(types.ObjectUnknown(strType), "ssh_key_name"); !got.IsUnknown() {
			t.Fatalf("unknown parent: expected unknown, got %v", got)
		}
		if got := stringAt(known, "not_there"); !got.IsUnknown() {
			t.Fatalf("missing attribute: expected unknown, got %v", got)
		}
		if got := stringAt(known, "ssh_key_name"); got.ValueString() != "k" {
			t.Fatalf("present attribute: expected k, got %v", got)
		}

		// An attribute whose value is not a types.String must fail closed too.
		notAString, _ := types.ObjectValue(objType, map[string]attr.Value{"credentials": types.ObjectNull(strType)})
		if got := stringAt(notAString, "credentials"); !got.IsUnknown() {
			t.Fatalf("wrong value type: expected unknown, got %v", got)
		}
	})

	t.Run("objectAt", func(t *testing.T) {
		t.Parallel()
		known, _ := types.ObjectValue(objType, map[string]attr.Value{"credentials": types.ObjectNull(strType)})

		if got := objectAt(types.ObjectNull(objType), "credentials"); !got.IsNull() {
			t.Fatalf("null parent: expected null, got %v", got)
		}
		if got := objectAt(types.ObjectUnknown(objType), "credentials"); !got.IsUnknown() {
			t.Fatalf("unknown parent: expected unknown, got %v", got)
		}
		if got := objectAt(known, "not_there"); !got.IsUnknown() {
			t.Fatalf("missing attribute: expected unknown, got %v", got)
		}

		notAnObject, _ := types.ObjectValue(strType, map[string]attr.Value{"ssh_key_name": types.StringValue("k")})
		if got := objectAt(notAnObject, "ssh_key_name"); !got.IsUnknown() {
			t.Fatalf("wrong value type: expected unknown, got %v", got)
		}
	})

	// The framework trap the walkers exist for: asking for the full nested path
	// through a wholly-unknown parent does not fail, it hands back a NULL leaf
	// (fwschemadata.Data.ValueAtPath swallows tftypes.ErrInvalidStep). Reading
	// that as "no change" is exactly the silent re-image this guards against.
	t.Run("nested GetAttribute through an unknown parent yields null", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		sch := ResourceSchema(ctx)
		planRaw := rawObjectWith(ctx, t, map[string]tftypes.Value{"servers_settings": unknownServersSettings(ctx, t)})
		plan := tfsdk.Plan{Raw: planRaw, Schema: sch}

		var leaf types.String
		diags := plan.GetAttribute(ctx, pathToSSHKeyName(), &leaf)
		if diags.HasError() {
			t.Fatalf("unexpected error reading the nested path: %v", diags)
		}
		if !leaf.IsNull() {
			t.Fatalf("framework behaviour changed: expected a null leaf, got %v - re-read rebuildTriggersOf's comment", leaf)
		}

		if got := rebuildTriggersOf(ctx, plan).sshKeyName; !got.IsUnknown() {
			t.Fatalf("rebuildTriggersOf must answer unknown where the direct read answers null, got %v", got)
		}
	})
}

// TestRebuildTriggerAttributeShapes fails if codegen renames a trigger or
// gives one a CustomType (the walkers would answer Unknown and every settings
// update would hard-error), or drops user_data's plan modifier (every update
// of a cluster without user_data, a rename included, would warn and then
// hard-error).
func TestRebuildTriggerAttributeShapes(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	sch := ResourceSchema(ctx)

	imageID, ok := sch.Attributes["image_id"].(schema.StringAttribute)
	if !ok {
		t.Fatalf("image_id is %T, expected schema.StringAttribute", sch.Attributes["image_id"])
	}
	if imageID.CustomType != nil {
		t.Fatalf("image_id gained a CustomType (%T); the walkers type-assert to types.String", imageID.CustomType)
	}

	settings, ok := sch.Attributes["servers_settings"].(schema.SingleNestedAttribute)
	if !ok {
		t.Fatalf("servers_settings is %T, expected schema.SingleNestedAttribute", sch.Attributes["servers_settings"])
	}
	if settings.CustomType != nil {
		t.Fatalf("servers_settings gained a CustomType (%T); the walkers type-assert to types.Object", settings.CustomType)
	}

	userData, ok := settings.Attributes["user_data"].(schema.StringAttribute)
	if !ok {
		t.Fatalf("servers_settings.user_data is %T, expected schema.StringAttribute", settings.Attributes["user_data"])
	}
	if userData.CustomType != nil {
		t.Fatalf("servers_settings.user_data gained a CustomType (%T)", userData.CustomType)
	}

	// The plan-vs-state comparison in rebuildDecisionFor rests on this: the
	// built-in UseStateForUnknown copies a null prior into the plan, so a
	// rename of a cluster with no user_data plans it as null, not unknown.
	// It must stay the built-in, not StringUseStateForUnknownInclNull: that
	// variant lacks the ConfigValue.IsUnknown() guard, so an interpolated
	// user_data would be replaced by the prior value at plan time and a real
	// change would never be applied.
	if !userData.Computed || !userData.Optional {
		t.Fatalf("servers_settings.user_data is no longer computed_optional (computed=%v optional=%v); re-check rebuildDecisionFor's comparison basis",
			userData.Computed, userData.Optional)
	}

	var hasUseStateForUnknown bool
	for _, modifier := range userData.PlanModifiers {
		if modifier.Description(ctx) == stringplanmodifier.UseStateForUnknown().Description(ctx) {
			hasUseStateForUnknown = true
		}
	}
	if !hasUseStateForUnknown {
		t.Fatal("servers_settings.user_data lost UseStateForUnknown; a rename would now plan user_data as unknown, warn, and hard-error at apply")
	}

	creds, ok := settings.Attributes["credentials"].(schema.SingleNestedAttribute)
	if !ok {
		t.Fatalf("servers_settings.credentials is %T, expected schema.SingleNestedAttribute", settings.Attributes["credentials"])
	}
	if creds.CustomType != nil {
		t.Fatalf("servers_settings.credentials gained a CustomType (%T)", creds.CustomType)
	}

	sshKey, ok := creds.Attributes["ssh_key_name"].(schema.StringAttribute)
	if !ok {
		t.Fatalf("credentials.ssh_key_name is %T, expected schema.StringAttribute", creds.Attributes["ssh_key_name"])
	}
	if sshKey.CustomType != nil {
		t.Fatalf("credentials.ssh_key_name gained a CustomType (%T)", sshKey.CustomType)
	}
}

// TestSSHKeyRemoved: a removal is a key in state and none in the plan;
// never-had-a-key and unresolved are not removals.
func TestSSHKeyRemoved(t *testing.T) {
	t.Parallel()

	key := func(v types.String) rebuildTriggers { return rebuildTriggers{sshKeyName: v} }

	testCases := []struct {
		name           string
		planned, prior types.String
		expected       bool
	}{
		{"key dropped", types.StringNull(), types.StringValue("k1"), true},
		{"key blanked", types.StringValue(""), types.StringValue("k1"), true},
		{"blank over blank", types.StringValue(""), types.StringValue(""), false},
		{"key kept", types.StringValue("k1"), types.StringValue("k1"), false},
		{"key changed", types.StringValue("k2"), types.StringValue("k1"), false},
		{"never had a key", types.StringNull(), types.StringNull(), false},
		{"first key", types.StringValue("k1"), types.StringNull(), false},
		{"prior unknown", types.StringNull(), types.StringUnknown(), false},
		{"plan unknown", types.StringUnknown(), types.StringValue("k1"), false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if got := sshKeyRemoved(key(tc.planned), key(tc.prior)); got != tc.expected {
				t.Fatalf("expected %v, got %v", tc.expected, got)
			}
		})
	}
}

// TestModifyPlanRebuildWarning covers what the user actually sees at plan time.
func TestModifyPlanRebuildWarning(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	testCases := []struct {
		name             string
		plan, state      map[string]tftypes.Value
		expectedWarnings int
		expectedError    bool
	}{
		{
			name:             "image_id change warns",
			plan:             map[string]tftypes.Value{"image_id": str("img-b")},
			state:            map[string]tftypes.Value{"image_id": str("img-a")},
			expectedWarnings: 1,
		},
		{
			name:             "ssh_key_name change warns",
			plan:             map[string]tftypes.Value{"servers_settings": serversSettingsWith(ctx, t, nil, map[string]tftypes.Value{"ssh_key_name": str("k2")})},
			state:            map[string]tftypes.Value{"servers_settings": serversSettingsWith(ctx, t, nil, map[string]tftypes.Value{"ssh_key_name": str("k1")})},
			expectedWarnings: 1,
		},
		{
			name:             "user_data change warns",
			plan:             map[string]tftypes.Value{"servers_settings": serversSettingsWith(ctx, t, map[string]tftypes.Value{"user_data": str("Y")}, nil)},
			state:            map[string]tftypes.Value{"servers_settings": serversSettingsWith(ctx, t, map[string]tftypes.Value{"user_data": str("X")}, nil)},
			expectedWarnings: 1,
		},
		{
			name:             "name-only change is silent",
			plan:             map[string]tftypes.Value{"name": str("new-name"), "image_id": str("img-a")},
			state:            map[string]tftypes.Value{"name": str("old-name"), "image_id": str("img-a")},
			expectedWarnings: 0,
		},
		{
			name:             "unknown servers_settings warns",
			plan:             map[string]tftypes.Value{"servers_settings": unknownServersSettings(ctx, t)},
			state:            map[string]tftypes.Value{"servers_settings": serversSettingsWith(ctx, t, nil, map[string]tftypes.Value{"ssh_key_name": str("k1")})},
			expectedWarnings: 1,
		},
		{
			// The blocking error stands alone: no re-image warning beside it.
			name:             "username change errors without warning",
			plan:             map[string]tftypes.Value{"servers_settings": serversSettingsWith(ctx, t, nil, map[string]tftypes.Value{"username": str("u2")})},
			state:            map[string]tftypes.Value{"servers_settings": serversSettingsWith(ctx, t, nil, map[string]tftypes.Value{"username": str("u1")})},
			expectedWarnings: 0,
			expectedError:    true,
		},
		{
			// Dropping the whole credentials block has no API path; it is an
			// error at plan time and the re-image warning stays quiet beside it.
			name: "credentials block removed errors without warning",
			plan: map[string]tftypes.Value{"servers_settings": serversSettingsWithNullCredentials(ctx, t)},
			state: map[string]tftypes.Value{"servers_settings": serversSettingsWith(ctx, t, nil, map[string]tftypes.Value{
				"ssh_key_name": str("k1"),
			})},
			expectedWarnings: 0,
			expectedError:    true,
		},
		{
			// An empty key that reached the plan (unknown at validate) is a
			// removal: the error fires and no re-image warning prints under it.
			name:             "ssh_key_name blanked errors without warning",
			plan:             map[string]tftypes.Value{"servers_settings": serversSettingsWith(ctx, t, nil, map[string]tftypes.Value{"ssh_key_name": str("")})},
			state:            map[string]tftypes.Value{"servers_settings": serversSettingsWith(ctx, t, nil, map[string]tftypes.Value{"ssh_key_name": str("k1")})},
			expectedWarnings: 0,
			expectedError:    true,
		},
		{
			// The credential error stands alone; no warning prints under it.
			name: "username and ssh_key_name together errors without warning",
			plan: map[string]tftypes.Value{"servers_settings": serversSettingsWith(ctx, t, nil, map[string]tftypes.Value{
				"username": str("u2"), "ssh_key_name": str("k2"),
			})},
			state: map[string]tftypes.Value{"servers_settings": serversSettingsWith(ctx, t, nil, map[string]tftypes.Value{
				"username": str("u1"), "ssh_key_name": str("k1"),
			})},
			expectedWarnings: 0,
			expectedError:    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			resp := runModifyPlan(ctx, t, rawObjectWith(ctx, t, tc.plan), rawObjectWith(ctx, t, tc.state))

			if got := resp.Diagnostics.WarningsCount(); got != tc.expectedWarnings {
				t.Fatalf("expected %d warning(s), got %d: %v", tc.expectedWarnings, got, resp.Diagnostics)
			}
			if got := resp.Diagnostics.HasError(); got != tc.expectedError {
				t.Fatalf("expected error=%v, got %v: %v", tc.expectedError, got, resp.Diagnostics)
			}
			if tc.expectedError && (strings.Contains(tc.name, "removed") || strings.Contains(tc.name, "blanked")) {
				if summary := resp.Diagnostics.Errors()[0].Summary(); summary != "Unsupported credential removal" {
					t.Fatalf("expected the removal error, got %q", summary)
				}
			}
			if tc.expectedWarnings > 0 {
				warning := resp.Diagnostics.Warnings()[0]
				if summary := warning.Summary(); summary != rebuildWarningSummary {
					t.Fatalf("unexpected warning summary: %q", summary)
				}
				// The detail is where the user actually learns what is at
				// stake; assert it so an empty or truncated body cannot pass.
				if detail := warning.Detail(); detail != rebuildWarningDetail {
					t.Fatalf("unexpected warning detail: %q", detail)
				}
			}
		})
	}
}

// TestUpdateUnknownImageIDRejected: an unresolved trigger reaching Update
// raises rebuildAtUnknownValueError instead of rolling settings out. With no
// other change in the fixture, no API call is made (nil client).
func TestUpdateUnknownImageIDRejected(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	sch := ResourceSchema(ctx)
	planRaw := rawObjectWith(ctx, t, map[string]tftypes.Value{"image_id": unknownStr()})
	stateRaw := rawObjectWith(ctx, t, map[string]tftypes.Value{"image_id": str("img-a")})

	resp := &resource.UpdateResponse{State: tfsdk.State{Raw: stateRaw, Schema: sch}}
	(&CloudGPUBaremetalClusterResource{}).Update(ctx, resource.UpdateRequest{
		Plan:  tfsdk.Plan{Raw: planRaw, Schema: sch},
		State: tfsdk.State{Raw: stateRaw, Schema: sch},
	}, resp)

	if !resp.Diagnostics.HasError() {
		t.Fatal("expected Update to refuse an unresolved image_id")
	}
	if summary := resp.Diagnostics.Errors()[0].Summary(); summary != rebuildUnknownSummary {
		t.Fatalf("expected %q, got %q", rebuildUnknownSummary, summary)
	}
}

// TestUpdateNoRebuildSendsNothing proves rebuildNo never reaches the client:
// with a nil client, anything that called it would panic.
func TestUpdateNoRebuildSendsNothing(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	sch := ResourceSchema(ctx)
	raw := rawObjectWith(ctx, t, map[string]tftypes.Value{"image_id": str("img-a")})

	resp := &resource.UpdateResponse{State: tfsdk.State{Raw: raw, Schema: sch}}
	(&CloudGPUBaremetalClusterResource{}).Update(ctx, resource.UpdateRequest{
		Plan:  tfsdk.Plan{Raw: raw, Schema: sch},
		State: tfsdk.State{Raw: raw, Schema: sch},
	}, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("expected no diagnostics for an unchanged plan, got: %v", resp.Diagnostics)
	}
}

// TestRenameWithNullUserDataDoesNotRebuild: a rename with user_data null on
// both sides decides rebuildNo. TestUpdateRenameDoesNotAdoptAPIUserData covers
// the same case through a stubbed API.
func TestRenameWithNullUserDataDoesNotRebuild(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	plan := rawObjectWith(ctx, t, map[string]tftypes.Value{
		"name":             str("new-name"),
		"image_id":         str("img-a"),
		"servers_settings": serversSettingsWith(ctx, t, nil, map[string]tftypes.Value{"ssh_key_name": str("k1")}),
	})
	state := rawObjectWith(ctx, t, map[string]tftypes.Value{
		"name":             str("old-name"),
		"image_id":         str("img-a"),
		"servers_settings": serversSettingsWith(ctx, t, nil, map[string]tftypes.Value{"ssh_key_name": str("k1")}),
	})

	if got := decide(ctx, t, plan, state); got != rebuildNo {
		t.Fatalf("a rename must not re-image the cluster, got %s", got)
	}
}

// TestApplyTriggerPayloadUsesPlannedValues: every trigger field in the PATCH
// comes from the planned values; null and unknown are never sent.
func TestApplyTriggerPayloadUsesPlannedValues(t *testing.T) {
	t.Parallel()

	planned := rebuildTriggers{
		imageID:    types.StringValue("img-planned"),
		sshKeyName: types.StringValue("key-planned"),
		userData:   types.StringValue("ud-planned"),
	}

	params := cloud.GPUBaremetalClusterUpdateParams{}
	applyTriggerPayload(&params, planned, true)

	if got := params.ImageID.Value; got != "img-planned" {
		t.Errorf("image_id: expected the planned value, got %q", got)
	}
	if got := params.ServersSettings.Credentials.SSHKeyName.Value; got != "key-planned" {
		t.Errorf("ssh_key_name: expected the planned value, got %q", got)
	}
	if got := params.ServersSettings.UserData.Value; got != "ud-planned" {
		t.Errorf("user_data: expected the planned value, got %q", got)
	}

	// image_id rides on imageChanged, the other two do not.
	unchanged := cloud.GPUBaremetalClusterUpdateParams{}
	applyTriggerPayload(&unchanged, planned, false)
	if unchanged.ImageID.Valid() {
		t.Error("image_id must not be sent when it did not change")
	}
	if !unchanged.ServersSettings.Credentials.SSHKeyName.Valid() {
		t.Error("ssh_key_name must still be sent when only another trigger changed")
	}

	// Null and unknown values are never sent: an unknown would reach the API as
	// an empty string, and a null is not something the API can act on.
	for _, tc := range []struct {
		name    string
		planned rebuildTriggers
	}{
		{"null", rebuildTriggers{imageID: types.StringNull(), sshKeyName: types.StringNull(), userData: types.StringNull()}},
		{"unknown", rebuildTriggers{imageID: types.StringUnknown(), sshKeyName: types.StringUnknown(), userData: types.StringUnknown()}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			p := cloud.GPUBaremetalClusterUpdateParams{}
			applyTriggerPayload(&p, tc.planned, false)
			if p.ServersSettings.Credentials.SSHKeyName.Valid() {
				t.Errorf("%s ssh_key_name must not be sent", tc.name)
			}
			if p.ServersSettings.UserData.Valid() {
				t.Errorf("%s user_data must not be sent", tc.name)
			}
		})
	}
}

// newStubbedResource returns a resource whose SDK client talks to a local
// httptest server, plus the paths and request bodies that server received.
// Tests in this package can assign the client field directly, and
// option.WithBaseURL is the same option the provider uses.
func newStubbedResource(t *testing.T, handler http.HandlerFunc) (*CloudGPUBaremetalClusterResource, *[]string, *[]string) {
	t.Helper()

	var seen, bodies []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		seen = append(seen, r.Method+" "+r.URL.Path)
		if body, err := io.ReadAll(r.Body); err == nil {
			bodies = append(bodies, string(body))
		}
		handler(w, r)
	}))
	t.Cleanup(server.Close)

	client := gcore.NewClient(option.WithBaseURL(server.URL), option.WithAPIKey("test"))

	return &CloudGPUBaremetalClusterResource{client: &client}, &seen, &bodies
}

// clusterResponse is a minimal cluster GET/PATCH body. user_data is returned as
// the API returns it: present, required, and in PLAIN TEXT.
func clusterResponse(name, userData string) string {
	return fmt.Sprintf(`{
		"id": "cluster-1",
		"name": %q,
		"flavor": "bm3-ai-1xlarge-noib-a100-80-8",
		"image_id": "img-a",
		"servers_count": 1,
		"servers_settings": {"user_data": %q, "ssh_key_name": "k1", "interfaces": []},
		"servers_ids": [],
		"status": "ACTIVE",
		"tags": {}
	}`, name, userData)
}

// TestUpdateRenameDoesNotAdoptAPIUserData: a rename PATCHes the name only.
// The response carries the API's user_data, which UnmarshalComputed folds into
// the model; neither the settings payload nor state may adopt it, and
// apply_settings must never be called.
func TestUpdateRenameDoesNotAdoptAPIUserData(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	res, seen, _ := newStubbedResource(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, clusterResponse("new-name", "#!/bin/bash\necho from-the-api\n"))
	})

	sch := ResourceSchema(ctx)
	settings := serversSettingsWith(ctx, t, nil, map[string]tftypes.Value{"ssh_key_name": str("k1")})
	planRaw := rawObjectWith(ctx, t, map[string]tftypes.Value{
		"id": str("cluster-1"), "name": str("new-name"), "image_id": str("img-a"),
		"project_id": num(1), "region_id": num(2),
		"servers_count": num(1), "servers_settings": settings,
	})
	stateRaw := rawObjectWith(ctx, t, map[string]tftypes.Value{
		"id": str("cluster-1"), "name": str("old-name"), "image_id": str("img-a"),
		"project_id": num(1), "region_id": num(2),
		"servers_count": num(1), "servers_settings": settings,
	})

	resp := &resource.UpdateResponse{State: tfsdk.State{Raw: stateRaw, Schema: sch}}
	res.Update(ctx, resource.UpdateRequest{
		Plan:  tfsdk.Plan{Raw: planRaw, Schema: sch},
		State: tfsdk.State{Raw: stateRaw, Schema: sch},
	}, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("rename must succeed, got: %v", resp.Diagnostics)
	}

	for _, path := range *seen {
		if strings.Contains(path, "apply_settings") {
			t.Fatalf("a rename must never roll settings out; calls were %v", *seen)
		}
	}

	// And the API's user_data must not have been adopted into state.
	var got types.String
	resp.Diagnostics.Append(resp.State.GetAttribute(ctx, path.Root("servers_settings").AtName("user_data"), &got)...)
	if !got.IsNull() {
		t.Fatalf("state adopted the API's user_data (%v); the planned null must win", got)
	}
}

// TestUpdateUnknownAfterRenameLeavesStateAlone: the rename is still sent (an
// unreadable trigger must not block it), the settings rollout stops, and state
// is left at the prior value rather than the planned model.
func TestUpdateUnknownAfterRenameLeavesStateAlone(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	res, seen, _ := newStubbedResource(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, clusterResponse("new-name", "from-the-api"))
	})

	sch := ResourceSchema(ctx)
	planRaw := rawObjectWith(ctx, t, map[string]tftypes.Value{
		"id": str("cluster-1"), "name": str("new-name"), "image_id": str("img-a"), "servers_count": num(1),
		"project_id": num(1), "region_id": num(2),
		"servers_settings": serversSettingsWith(ctx, t,
			map[string]tftypes.Value{"user_data": unknownStr()},
			map[string]tftypes.Value{"ssh_key_name": str("k2")}),
	})
	stateRaw := rawObjectWith(ctx, t, map[string]tftypes.Value{
		"id": str("cluster-1"), "name": str("old-name"), "image_id": str("img-a"), "servers_count": num(1),
		"project_id": num(1), "region_id": num(2),
		"servers_settings": serversSettingsWith(ctx, t,
			map[string]tftypes.Value{"user_data": str("X")},
			map[string]tftypes.Value{"ssh_key_name": str("k1")}),
	})

	resp := &resource.UpdateResponse{State: tfsdk.State{Raw: stateRaw, Schema: sch}}
	res.Update(ctx, resource.UpdateRequest{
		Plan:  tfsdk.Plan{Raw: planRaw, Schema: sch},
		State: tfsdk.State{Raw: stateRaw, Schema: sch},
	}, resp)

	if !resp.Diagnostics.HasError() {
		t.Fatal("expected Update to refuse an unresolved trigger")
	}
	if summary := resp.Diagnostics.Errors()[0].Summary(); summary != rebuildUnknownSummary {
		t.Fatalf("expected %q, got %q", rebuildUnknownSummary, summary)
	}

	// The rename went out; the rollout did not.
	var renamed bool
	for _, p := range *seen {
		if strings.Contains(p, "apply_settings") {
			t.Fatalf("settings must not roll out on an unresolved trigger; calls were %v", *seen)
		}
		if p == "PATCH /cloud/v3/gpu/baremetal/1/2/clusters/cluster-1" {
			renamed = true
		}
	}
	if !renamed {
		t.Fatalf("the rename PATCH must still be sent; calls were %v", *seen)
	}

	// Nothing is published, so the planned ssh_key_name k2 never reaches state.
	if !resp.State.Raw.Equal(stateRaw) {
		t.Fatalf("state must be left at the prior value on this path, got %v", resp.State.Raw)
	}
}

// TestUpdateUnknownLeavesStateAlone is the no-other-change variant: the guard
// stops the rollout and leaves prior state in place.
func TestUpdateUnknownLeavesStateAlone(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	sch := ResourceSchema(ctx)
	planRaw := rawObjectWith(ctx, t, map[string]tftypes.Value{
		"image_id":         str("img-b"),
		"servers_settings": serversSettingsWith(ctx, t, map[string]tftypes.Value{"user_data": unknownStr()}, map[string]tftypes.Value{"ssh_key_name": str("k2")}),
	})
	stateRaw := rawObjectWith(ctx, t, map[string]tftypes.Value{
		"image_id":         str("img-a"),
		"servers_settings": serversSettingsWith(ctx, t, map[string]tftypes.Value{"user_data": str("X")}, map[string]tftypes.Value{"ssh_key_name": str("k1")}),
	})

	// Seeded with the prior state, exactly as the framework seeds it.
	resp := &resource.UpdateResponse{State: tfsdk.State{Raw: stateRaw, Schema: sch}}
	(&CloudGPUBaremetalClusterResource{}).Update(ctx, resource.UpdateRequest{
		Plan:  tfsdk.Plan{Raw: planRaw, Schema: sch},
		State: tfsdk.State{Raw: stateRaw, Schema: sch},
	}, resp)

	if !resp.Diagnostics.HasError() {
		t.Fatal("expected Update to refuse an unresolved trigger")
	}
	if summary := resp.Diagnostics.Errors()[0].Summary(); summary != rebuildUnknownSummary {
		t.Fatalf("expected %q, got %q", rebuildUnknownSummary, summary)
	}
	if !resp.State.Raw.Equal(stateRaw) {
		t.Fatalf("state must be left at the prior value on this path, got %v", resp.State.Raw)
	}
}

// TestUpdateSettingsPayloadIsPlanSourced: with a rename and a key change in
// one apply, the settings PATCH that follows the rename must not echo the
// API's user_data back; the payload comes from the plan.
func TestUpdateSettingsPayloadIsPlanSourced(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	const apiUserData = "user-data-only-the-api-knows"

	res, seen, bodies := newStubbedResource(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if strings.Contains(r.URL.Path, "apply_settings") {
			_, _ = io.WriteString(w, `{"tasks":["task-1"]}`)

			return
		}
		if strings.Contains(r.URL.Path, "tasks") {
			_, _ = io.WriteString(w, `{"id":"task-1","state":"FINISHED"}`)

			return
		}
		_, _ = io.WriteString(w, clusterResponse("new-name", apiUserData))
	})

	sch := ResourceSchema(ctx)
	planRaw := rawObjectWith(ctx, t, map[string]tftypes.Value{
		"id": str("cluster-1"), "name": str("new-name"), "image_id": str("img-a"), "servers_count": num(1),
		"project_id": num(1), "region_id": num(2),
		"servers_settings": serversSettingsWith(ctx, t, nil, map[string]tftypes.Value{"ssh_key_name": str("k2")}),
	})
	stateRaw := rawObjectWith(ctx, t, map[string]tftypes.Value{
		"id": str("cluster-1"), "name": str("old-name"), "image_id": str("img-a"), "servers_count": num(1),
		"project_id": num(1), "region_id": num(2),
		"servers_settings": serversSettingsWith(ctx, t, nil, map[string]tftypes.Value{"ssh_key_name": str("k1")}),
	})

	resp := &resource.UpdateResponse{State: tfsdk.State{Raw: stateRaw, Schema: sch}}
	res.Update(ctx, resource.UpdateRequest{
		Plan:  tfsdk.Plan{Raw: planRaw, Schema: sch},
		State: tfsdk.State{Raw: stateRaw, Schema: sch},
	}, resp)

	// The key change is real, so the rollout must happen.
	var rolledOut bool
	for _, p := range *seen {
		if strings.Contains(p, "apply_settings") {
			rolledOut = true
		}
	}
	if !rolledOut {
		t.Fatalf("an ssh_key_name change must roll settings out; calls were %v", *seen)
	}

	// No request may carry the API's own user_data back to it; the planned key
	// must be in the settings PATCH; and the rollout must carry the API's
	// destructive-apply confirmation, without which the real API rejects it.
	var sentKey, confirmed bool
	for i, body := range *bodies {
		if strings.Contains(body, apiUserData) {
			t.Fatalf("request %d (%s) echoed the API's user_data back: %s", i, (*seen)[i], body)
		}
		if strings.Contains(body, `"ssh_key_name":"k2"`) {
			sentKey = true
		}
		if strings.Contains((*seen)[i], "apply_settings") && strings.Contains(body, `"max_disruption":"rebuild"`) {
			confirmed = true
		}
	}
	if !sentKey {
		t.Fatalf("the settings PATCH must carry the planned key; bodies were %v", *bodies)
	}
	if !confirmed {
		t.Fatalf("apply_settings must carry max_disruption=rebuild; bodies were %v", *bodies)
	}

	// And the apply succeeds and records the new key.
	if resp.Diagnostics.HasError() {
		t.Fatalf("expected a clean apply, got: %v", resp.Diagnostics)
	}
	var got types.String
	resp.Diagnostics.Append(resp.State.GetAttribute(ctx, pathToSSHKeyName(), &got)...)
	if got.ValueString() != "k2" {
		t.Fatalf("state must record the applied key k2, got %v", got)
	}
}

// TestUpdateSettingsOnlyRollsOut: a key change with no other change goes
// straight to the settings PATCH and the confirmed rollout, and persists.
func TestUpdateSettingsOnlyRollsOut(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	res, seen, bodies := newStubbedResource(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case strings.Contains(r.URL.Path, "apply_settings"):
			_, _ = io.WriteString(w, `{"tasks":["task-1"]}`)
		case strings.Contains(r.URL.Path, "tasks"):
			_, _ = io.WriteString(w, `{"id":"task-1","state":"FINISHED"}`)
		default:
			_, _ = io.WriteString(w, clusterResponse("same-name", "X"))
		}
	})

	sch := ResourceSchema(ctx)
	mk := func(key string) tftypes.Value {
		return rawObjectWith(ctx, t, map[string]tftypes.Value{
			"id": str("cluster-1"), "name": str("same-name"), "image_id": str("img-a"), "servers_count": num(1),
			"project_id": num(1), "region_id": num(2),
			"servers_settings": serversSettingsWith(ctx, t, map[string]tftypes.Value{"user_data": str("X")}, map[string]tftypes.Value{"ssh_key_name": str(key)}),
		})
	}
	planRaw, stateRaw := mk("k2"), mk("k1")

	resp := &resource.UpdateResponse{State: tfsdk.State{Raw: stateRaw, Schema: sch}}
	res.Update(ctx, resource.UpdateRequest{
		Plan:  tfsdk.Plan{Raw: planRaw, Schema: sch},
		State: tfsdk.State{Raw: stateRaw, Schema: sch},
	}, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("expected a clean apply, got: %v", resp.Diagnostics)
	}
	var patches, applies int
	for i, p := range *seen {
		switch {
		case strings.Contains(p, "apply_settings"):
			applies++
			if !strings.Contains((*bodies)[i], `"max_disruption":"rebuild"`) {
				t.Fatalf("apply_settings body lacks the confirmation: %s", (*bodies)[i])
			}
		case p == "PATCH /cloud/v3/gpu/baremetal/1/2/clusters/cluster-1":
			patches++
			if !strings.Contains((*bodies)[i], `"ssh_key_name":"k2"`) {
				t.Fatalf("settings PATCH lacks the planned key: %s", (*bodies)[i])
			}
		}
	}
	if patches != 1 || applies != 1 {
		t.Fatalf("expected exactly one settings PATCH and one apply_settings, got %d/%d: %v", patches, applies, *seen)
	}
	var got types.String
	resp.Diagnostics.Append(resp.State.GetAttribute(ctx, pathToSSHKeyName(), &got)...)
	if got.ValueString() != "k2" {
		t.Fatalf("state must record the applied key k2, got %v", got)
	}
}

// TestCreateRejectsResolvedCredentialConflict: a conflict the validator had
// to defer on an unknown is caught in Create once the value has resolved,
// before anything is sent (nil client: reaching the API would panic).
func TestCreateRejectsResolvedCredentialConflict(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	sch := ResourceSchema(ctx)
	raw := rawObjectWith(ctx, t, map[string]tftypes.Value{
		"name": str("c"), "flavor": str("f"), "image_id": str("img"), "servers_count": num(1),
		"servers_settings": serversSettingsWith(ctx, t, nil, map[string]tftypes.Value{
			"ssh_key_name": str("existing-key"), "password_wo_version": num(1),
		}),
	})

	resp := &resource.CreateResponse{State: tfsdk.State{Schema: sch}}
	(&CloudGPUBaremetalClusterResource{}).Create(ctx, resource.CreateRequest{
		Plan:   tfsdk.Plan{Raw: raw, Schema: sch},
		Config: tfsdk.Config{Raw: raw, Schema: sch},
	}, resp)

	if !resp.Diagnostics.HasError() {
		t.Fatal("expected Create to reject a resolved key + password_wo_version conflict")
	}
	if got := resp.Diagnostics.Errors()[0].Summary(); got != "Conflicting credentials configuration" {
		t.Fatalf("expected the conflict error, got %q", got)
	}
}

// TestCreateAcceptsPasswordCredentials: Terraform sends the write-only
// password_wo as null in the plan, even at apply. The re-check must read the
// config, where it is set, or it rejects every username/password cluster.
func TestCreateAcceptsPasswordCredentials(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	res, seen, _ := newStubbedResource(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"message": "stub"}`))
	})

	cluster := func(password tftypes.Value) tftypes.Value {
		return rawObjectWith(ctx, t, map[string]tftypes.Value{
			"project_id": num(1), "region_id": num(1),
			"name": str("c"), "flavor": str("f"), "image_id": str("img"), "servers_count": num(1),
			"servers_settings": serversSettingsWith(ctx, t, nil, map[string]tftypes.Value{
				"username": str("ubuntu"), "password_wo": password, "password_wo_version": num(1),
			}),
		})
	}
	sch := ResourceSchema(ctx)
	resp := &resource.CreateResponse{State: tfsdk.State{Schema: sch}}
	res.Create(ctx, resource.CreateRequest{
		Plan:   tfsdk.Plan{Raw: cluster(tftypes.NewValue(tftypes.String, nil)), Schema: sch},
		Config: tfsdk.Config{Raw: cluster(str("secret")), Schema: sch},
	}, resp)

	for _, d := range resp.Diagnostics.Errors() {
		if strings.Contains(d.Summary(), "credentials") || strings.Contains(d.Summary(), "password_wo") {
			t.Fatalf("valid password credentials rejected: %s: %s", d.Summary(), d.Detail())
		}
	}
	if len(*seen) == 0 {
		t.Fatalf("expected Create to reach the API, got: %v", resp.Diagnostics)
	}
}

// TestWarningsNeverClaimAPlanAction: ModifyPlan is handed the same request for
// `terraform plan -replace=ADDR` as for an ordinary update, so no diagnostic
// may state which action Terraform planned.
func TestWarningsNeverClaimAPlanAction(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		text     string
		banned   []string
		required []string
	}{
		{
			name:     "rebuildWarningSummary",
			text:     rebuildWarningSummary,
			banned:   []string{"This update", "in-place", "replaced", "Terraform plans", "keeps its id"},
			required: []string{"re-images servers"},
		},
		{
			name:   "rebuildWarningDetail",
			text:   rebuildWarningDetail,
			banned: []string{"in-place update", "will be replaced", "Terraform plans this", "rather than", "keeps its id", "keep their addresses"},
			required: []string{
				"image_id, servers_settings.credentials.ssh_key_name or servers_settings.user_data",
				"every server whose settings are out of date",
				"local disks is lost",
				"revert the attribute",
				"resource documentation",
			},
		},
		{
			name:     "rebuildUnknownDetail",
			text:     rebuildUnknownDetail,
			required: []string{"none of it was written to state", "the refresh at the start of the next plan"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			for _, phrase := range tc.banned {
				if strings.Contains(tc.text, phrase) {
					t.Errorf("must not contain %q", phrase)
				}
			}
			for _, phrase := range tc.required {
				if !strings.Contains(tc.text, phrase) {
					t.Errorf("must contain %q", phrase)
				}
			}
		})
	}
}

// configOf builds a tfsdk.Config from a raw object, for tests that need the
// config-driven replacement rules.
func configOf(ctx context.Context, raw tftypes.Value) tfsdk.Config {
	return tfsdk.Config{Raw: raw, Schema: ResourceSchema(ctx)}
}

func runModifyPlanWithConfig(ctx context.Context, t *testing.T, planRaw, stateRaw, configRaw tftypes.Value) *resource.ModifyPlanResponse {
	t.Helper()
	sch := ResourceSchema(ctx)
	resp := &resource.ModifyPlanResponse{}
	(&CloudGPUBaremetalClusterResource{}).ModifyPlan(ctx, resource.ModifyPlanRequest{
		Plan:   tfsdk.Plan{Raw: planRaw, Schema: sch},
		State:  tfsdk.State{Raw: stateRaw, Schema: sch},
		Config: tfsdk.Config{Raw: configRaw, Schema: sch},
	}, resp)
	return resp
}

// TestPlanReplacesCluster pins each RequiresReplace-family rule the predicate
// re-derives, and that a missing config disables only the config-driven ones.
func TestPlanReplacesCluster(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	sch := ResourceSchema(ctx)

	settings := func(ifaceType string) tftypes.Value {
		objType := sch.Type().TerraformType(ctx).(tftypes.Object)
		ifaceType_ := objType.AttributeTypes["servers_settings"].(tftypes.Object).AttributeTypes["interfaces"].(tftypes.List)
		elem := ifaceType_.ElementType.(tftypes.Object)
		vals := map[string]tftypes.Value{}
		for name, at := range elem.AttributeTypes {
			vals[name] = tftypes.NewValue(at, nil)
		}
		vals["type"] = str(ifaceType)
		return serversSettingsWith(ctx, t, map[string]tftypes.Value{
			"interfaces": tftypes.NewValue(ifaceType_, []tftypes.Value{tftypes.NewValue(elem, vals)}),
		}, map[string]tftypes.Value{"ssh_key_name": str("k1")})
	}

	base := map[string]tftypes.Value{
		"flavor": str("f1"), "project_id": num(1), "region_id": num(2), "servers_settings": settings("external"),
	}
	with := func(over map[string]tftypes.Value) map[string]tftypes.Value {
		m := map[string]tftypes.Value{}
		for k, v := range base {
			m[k] = v
		}
		for k, v := range over {
			m[k] = v
		}
		return m
	}

	testCases := []struct {
		name               string
		plan, state, cfg   map[string]tftypes.Value
		noConfig, expected bool
	}{
		{"no change", base, base, base, false, false},
		{"flavor changes (RequiresReplace)", with(map[string]tftypes.Value{"flavor": str("f2")}), base, with(map[string]tftypes.Value{"flavor": str("f2")}), false, true},
		{"flavor changes, no config", with(map[string]tftypes.Value{"flavor": str("f2")}), base, nil, true, true},
		{"flavor unknown (interpolated)", with(map[string]tftypes.Value{"flavor": unknownStr()}), base, with(map[string]tftypes.Value{"flavor": unknownStr()}), false, true},
		{"project_id changes while configured", with(map[string]tftypes.Value{"project_id": num(9)}), base, with(map[string]tftypes.Value{"project_id": num(9)}), false, true},
		{"project_id changes, not configured", with(map[string]tftypes.Value{"project_id": num(9)}), base, with(map[string]tftypes.Value{"project_id": tftypes.NewValue(tftypes.Number, nil)}), false, false},
		{"interface type changes (servers_settings, not ignored)", with(map[string]tftypes.Value{"servers_settings": settings("subnet")}), base, with(map[string]tftypes.Value{"servers_settings": settings("subnet")}), false, true},
		{"only ssh_key_name changes (ignored field)", with(map[string]tftypes.Value{"servers_settings": serversSettingsWith(ctx, t, nil, map[string]tftypes.Value{"ssh_key_name": str("k2")})}), base, with(map[string]tftypes.Value{"servers_settings": serversSettingsWith(ctx, t, nil, map[string]tftypes.Value{"ssh_key_name": str("k2")})}), false, false},
		{"servers_settings changes, no config", with(map[string]tftypes.Value{"servers_settings": settings("subnet")}), base, nil, true, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			plan := tfsdk.Plan{Raw: rawObjectWith(ctx, t, tc.plan), Schema: sch}
			state := tfsdk.State{Raw: rawObjectWith(ctx, t, tc.state), Schema: sch}
			var cfg tfsdk.Config
			if !tc.noConfig {
				cfg = configOf(ctx, rawObjectWith(ctx, t, tc.cfg))
			}
			if got := planReplacesCluster(ctx, plan, state, cfg); got != tc.expected {
				t.Fatalf("expected %v, got %v", tc.expected, got)
			}
		})
	}
}

// TestModifyPlanSkipsGuardsOnReplacement: a plan that already replaces the
// cluster runs Create, so dropping the key, changing the username, or changing
// a trigger must neither error nor warn.
func TestModifyPlanSkipsGuardsOnReplacement(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	withKey := serversSettingsWith(ctx, t, nil, map[string]tftypes.Value{"ssh_key_name": str("k1")})
	state := rawObjectWith(ctx, t, map[string]tftypes.Value{"flavor": str("f1"), "servers_settings": withKey})

	testCases := []struct {
		name string
		plan map[string]tftypes.Value
	}{
		{"flavor change and credentials removed", map[string]tftypes.Value{"flavor": str("f2"), "servers_settings": serversSettingsWithNullCredentials(ctx, t)}},
		{"flavor change and username change", map[string]tftypes.Value{"flavor": str("f2"), "servers_settings": serversSettingsWith(ctx, t, nil, map[string]tftypes.Value{"username": str("u1")})}},
		{"flavor change and ssh_key_name change", map[string]tftypes.Value{"flavor": str("f2"), "servers_settings": serversSettingsWith(ctx, t, nil, map[string]tftypes.Value{"ssh_key_name": str("k2")})}},
		{"flavor unknown and credentials removed", map[string]tftypes.Value{"flavor": unknownStr(), "servers_settings": serversSettingsWithNullCredentials(ctx, t)}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			planRaw := rawObjectWith(ctx, t, tc.plan)
			resp := runModifyPlanWithConfig(ctx, t, planRaw, state, planRaw)
			if resp.Diagnostics.HasError() {
				t.Fatalf("a replacement plan must not error, got: %v", resp.Diagnostics)
			}
			if n := resp.Diagnostics.WarningsCount(); n != 0 {
				t.Fatalf("a replacement plan must not warn, got %d: %v", n, resp.Diagnostics)
			}
		})
	}

	// And without a replacement the guard still fires.
	planRaw := rawObjectWith(ctx, t, map[string]tftypes.Value{"flavor": str("f1"), "servers_settings": serversSettingsWithNullCredentials(ctx, t)})
	resp := runModifyPlanWithConfig(ctx, t, planRaw, state, planRaw)
	if !resp.Diagnostics.HasError() {
		t.Fatal("removing the key without a replacement must still error")
	}
}

// TestReplacementRulesMatchSchema fails if codegen adds, drops or changes a
// RequiresReplace-family modifier anywhere in the schema, since
// planReplacesCluster re-derives them by hand. Modifier descriptions are the
// only thing a modifier exposes.
func TestReplacementRulesMatchSchema(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	sch := ResourceSchema(ctx)

	// Every attribute path that carries a replace-forcing modifier, found by
	// walking nested attributes too.
	found := map[string]bool{}
	var walk func(prefix string, attrs map[string]schema.Attribute)
	walk = func(prefix string, attrs map[string]schema.Attribute) {
		for name, a := range attrs {
			p := prefix + name
			for _, d := range modifierDescriptions(ctx, a) {
				// The built-in RequiresReplace family says "destroy and
				// recreate"; the repo's own says "Requires replacement".
				// Modifiers that merely mention replacement (e.g. "preserves
				// state unless resource is replaced") match neither.
				l := strings.ToLower(d)
				if strings.Contains(l, "destroy and recreate") || strings.HasPrefix(l, "requires replacement") {
					found[p] = true
				}
			}
			switch n := a.(type) {
			case schema.SingleNestedAttribute:
				walk(p+".", n.Attributes)
			case schema.ListNestedAttribute:
				walk(p+"[].", n.NestedObject.Attributes)
			case schema.SetNestedAttribute:
				walk(p+"[].", n.NestedObject.Attributes)
			}
		}
	}
	walk("", sch.Attributes)

	// Each path pinned to the exact rule planReplacesCluster implements for it.
	// Membership alone is not enough: swapping project_id from
	// RequiresReplaceIfConfigured to RequiresReplace would keep it in the set
	// while the predicate went on requiring a configured value.
	expected := map[string]string{
		"flavor":           "If the value of this attribute changes, Terraform will destroy and recreate the resource.",
		"project_id":       "If the value of this attribute is configured and changes, Terraform will destroy and recreate the resource.",
		"region_id":        "If the value of this attribute is configured and changes, Terraform will destroy and recreate the resource.",
		"servers_settings": "Requires replacement only when user-specified server settings change",
	}
	for p := range found {
		if _, ok := expected[p]; !ok {
			t.Errorf("%s carries a replacement modifier that planReplacesCluster does not implement", p)
		}
	}
	for p, want := range expected {
		if !found[p] {
			t.Errorf("%s lost its replacement modifier; planReplacesCluster still treats it as one", p)
			continue
		}
		var got []string
		for _, d := range modifierDescriptions(ctx, sch.Attributes[p]) {
			got = append(got, d)
		}
		var matched bool
		for _, d := range got {
			if d == want {
				matched = true
			}
		}
		if !matched {
			t.Errorf("%s: rule changed; planReplacesCluster implements %q, schema has %q", p, want, got)
		}
	}
}

// modifierDescriptions returns the descriptions of every plan modifier on a
// schema attribute, whatever its type.
func modifierDescriptions(ctx context.Context, a schema.Attribute) []string {
	var out []string
	switch n := a.(type) {
	case schema.StringAttribute:
		for _, m := range n.PlanModifiers {
			out = append(out, m.Description(ctx))
		}
	case schema.Int64Attribute:
		for _, m := range n.PlanModifiers {
			out = append(out, m.Description(ctx))
		}
	case schema.BoolAttribute:
		for _, m := range n.PlanModifiers {
			out = append(out, m.Description(ctx))
		}
	case schema.ListAttribute:
		for _, m := range n.PlanModifiers {
			out = append(out, m.Description(ctx))
		}
	case schema.MapAttribute:
		for _, m := range n.PlanModifiers {
			out = append(out, m.Description(ctx))
		}
	case schema.SingleNestedAttribute:
		for _, m := range n.PlanModifiers {
			out = append(out, m.Description(ctx))
		}
	case schema.ListNestedAttribute:
		for _, m := range n.PlanModifiers {
			out = append(out, m.Description(ctx))
		}
	case schema.SetNestedAttribute:
		for _, m := range n.PlanModifiers {
			out = append(out, m.Description(ctx))
		}
	}
	return out
}

// TestServersSettingsReplaceIgnoreMatchesModifier pins the ignore list
// planReplacesCluster duplicates from schema.go, behaviourally: the modifier
// on servers_settings and the predicate must agree on which nested fields
// force a replacement.
func TestServersSettingsReplaceIgnoreMatchesModifier(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	sch := ResourceSchema(ctx)
	settingsAttr := sch.Attributes["servers_settings"].(schema.SingleNestedAttribute)
	if len(settingsAttr.PlanModifiers) != 1 {
		t.Fatalf("expected exactly one plan modifier on servers_settings, got %d", len(settingsAttr.PlanModifiers))
	}
	modifier := settingsAttr.PlanModifiers[0]

	toObject := func(raw tftypes.Value) types.Object {
		v, err := settingsAttr.GetType().ValueFromTerraform(ctx, raw)
		if err != nil {
			t.Fatalf("converting fixture: %v", err)
		}
		return v.(types.Object)
	}
	ifaceOf := func(ifaceType string) tftypes.Value {
		objType := sch.Type().TerraformType(ctx).(tftypes.Object)
		listType := objType.AttributeTypes["servers_settings"].(tftypes.Object).AttributeTypes["interfaces"].(tftypes.List)
		elem := listType.ElementType.(tftypes.Object)
		vals := map[string]tftypes.Value{}
		for name, at := range elem.AttributeTypes {
			vals[name] = tftypes.NewValue(at, nil)
		}
		vals["type"] = str(ifaceType)
		return tftypes.NewValue(listType, []tftypes.Value{tftypes.NewValue(elem, vals)})
	}

	state := serversSettingsWith(ctx, t, map[string]tftypes.Value{"user_data": str("X"), "interfaces": ifaceOf("external")}, map[string]tftypes.Value{"ssh_key_name": str("k1")})
	testCases := []struct {
		name     string
		config   tftypes.Value
		replaces bool
	}{
		{"only user_data differs", serversSettingsWith(ctx, t, map[string]tftypes.Value{"user_data": str("Y"), "interfaces": ifaceOf("external")}, map[string]tftypes.Value{"ssh_key_name": str("k1")}), false},
		{"only ssh_key_name differs", serversSettingsWith(ctx, t, map[string]tftypes.Value{"user_data": str("X"), "interfaces": ifaceOf("external")}, map[string]tftypes.Value{"ssh_key_name": str("k2")}), false},
		{"interfaces differ", serversSettingsWith(ctx, t, map[string]tftypes.Value{"user_data": str("X"), "interfaces": ifaceOf("subnet")}, map[string]tftypes.Value{"ssh_key_name": str("k1")}), true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			resp := planmodifier.ObjectResponse{PlanValue: toObject(tc.config)}
			modifier.PlanModifyObject(ctx, planmodifier.ObjectRequest{
				ConfigValue: toObject(tc.config), PlanValue: toObject(tc.config), StateValue: toObject(state),
				State: tfsdk.State{Raw: rawObjectWith(ctx, t, map[string]tftypes.Value{"servers_settings": state}), Schema: sch},
			}, &resp)
			if resp.RequiresReplace != tc.replaces {
				t.Fatalf("modifier: expected RequiresReplace=%v, got %v", tc.replaces, resp.RequiresReplace)
			}

			plan := tfsdk.Plan{Raw: rawObjectWith(ctx, t, map[string]tftypes.Value{"flavor": str("f"), "servers_settings": tc.config}), Schema: sch}
			st := tfsdk.State{Raw: rawObjectWith(ctx, t, map[string]tftypes.Value{"flavor": str("f"), "servers_settings": state}), Schema: sch}
			cfg := configOf(ctx, rawObjectWith(ctx, t, map[string]tftypes.Value{"flavor": str("f"), "servers_settings": tc.config}))
			if got := planReplacesCluster(ctx, plan, st, cfg); got != tc.replaces {
				t.Fatalf("planReplacesCluster disagrees with the modifier: expected %v, got %v", tc.replaces, got)
			}
		})
	}
}
