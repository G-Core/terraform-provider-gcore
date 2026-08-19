package cloud_load_balancer

import (
	"context"
	"encoding/json"
	"fmt"
	"maps"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// TestStripAdoptedCreateOnlyFieldsAfterImport covers the plan the import-safe
// plan modifiers produce: state is null for the create-only attributes,
// the config has values, so the update body would otherwise PATCH fields the
// update endpoint does not accept (it only takes name, logging,
// preferred_connectivity and tags).
func TestStripAdoptedCreateOnlyFieldsAfterImport(t *testing.T) {
	t.Parallel()

	plan := CloudLoadBalancerModel{
		VipNetworkID: types.StringValue("net-a"),
		VipSubnetID:  types.StringValue("subnet-a"),
	}
	state := CloudLoadBalancerModel{}

	body, err := plan.MarshalJSONForUpdate(state)
	if err != nil {
		t.Fatalf("MarshalJSONForUpdate: %s", err)
	}
	for name := range createOnlyFields {
		if !jsonHasKey(t, body, name) {
			t.Fatalf("expected the unfiltered update body to carry %s, got %s", name, body)
		}
	}

	stripped, err := stripAdoptedCreateOnlyFields(body, state, true)
	if err != nil {
		t.Fatalf("stripAdoptedCreateOnlyFields: %s", err)
	}
	if string(stripped) != "{}" {
		t.Fatalf("expected an empty update body after import adoption, got %s", stripped)
	}
}

// TestStripAdoptedCreateOnlyFieldsOutsideAdoptionWindow is the regression test
// for the case that makes the adoption marker necessary: a load balancer that
// was created without these optional attributes also has null prior state, but
// setting one on day two is a real change. Nothing may be stripped there - the
// plan modifiers force replacement instead, and if a field ever does reach the
// body it must fail against the API rather than be silently swallowed.
func TestStripAdoptedCreateOnlyFieldsOutsideAdoptionWindow(t *testing.T) {
	t.Parallel()

	plan := CloudLoadBalancerModel{
		VipNetworkID: types.StringValue("net-a"),
		VipSubnetID:  types.StringValue("subnet-a"),
	}
	state := CloudLoadBalancerModel{}

	body, err := plan.MarshalJSONForUpdate(state)
	if err != nil {
		t.Fatalf("MarshalJSONForUpdate: %s", err)
	}

	stripped, err := stripAdoptedCreateOnlyFields(body, state, false)
	if err != nil {
		t.Fatalf("stripAdoptedCreateOnlyFields: %s", err)
	}
	for name := range createOnlyFields {
		if !jsonHasKey(t, stripped, name) {
			t.Errorf("expected %s to survive outside the adoption window, got %s", name, stripped)
		}
	}
}

// TestStripAdoptedCreateOnlyFieldsKeepsOtherChanges makes sure the filter only
// removes the adopted create-only fields and leaves a real update alone.
func TestStripAdoptedCreateOnlyFieldsKeepsOtherChanges(t *testing.T) {
	t.Parallel()

	plan := CloudLoadBalancerModel{
		Name:         types.StringValue("renamed"),
		VipNetworkID: types.StringValue("net-a"),
	}
	state := CloudLoadBalancerModel{Name: types.StringValue("original")}

	body, err := plan.MarshalJSONForUpdate(state)
	if err != nil {
		t.Fatalf("MarshalJSONForUpdate: %s", err)
	}

	stripped, err := stripAdoptedCreateOnlyFields(body, state, true)
	if err != nil {
		t.Fatalf("stripAdoptedCreateOnlyFields: %s", err)
	}
	if jsonHasKey(t, stripped, "vip_network_id") {
		t.Errorf("expected vip_network_id to be dropped, got %s", stripped)
	}
	if !jsonHasKey(t, stripped, "name") {
		t.Errorf("expected the name change to survive, got %s", stripped)
	}
}

// TestStripAdoptedCreateOnlyFieldsKeepsKnownPriorState pins the loud-failure
// path: with a known prior state these attributes force replacement, so they
// must never be silently dropped from an update body.
func TestStripAdoptedCreateOnlyFieldsKeepsKnownPriorState(t *testing.T) {
	t.Parallel()

	state := CloudLoadBalancerModel{
		VipNetworkID: types.StringValue("net-a"),
		VipSubnetID:  types.StringValue("subnet-a"),
	}
	body := []byte(`{"vip_network_id":"net-b","vip_subnet_id":"subnet-b"}`)

	stripped, err := stripAdoptedCreateOnlyFields(body, state, true)
	if err != nil {
		t.Fatalf("stripAdoptedCreateOnlyFields: %s", err)
	}
	for name := range createOnlyFields {
		if !jsonHasKey(t, stripped, name) {
			t.Errorf("expected %s to survive when the prior state is known, got %s", name, stripped)
		}
	}
}

func TestStripAdoptedCreateOnlyFieldsEmptyBody(t *testing.T) {
	t.Parallel()

	for _, body := range []string{"", "null", "{}"} {
		stripped, err := stripAdoptedCreateOnlyFields([]byte(body), CloudLoadBalancerModel{}, true)
		if err != nil {
			t.Fatalf("%q: stripAdoptedCreateOnlyFields: %s", body, err)
		}
		if string(stripped) != body {
			t.Errorf("%q: expected the body to be passed through, got %s", body, stripped)
		}
	}
}

// TestImportAdoptionMarkerRoundTrip covers the lifecycle of the private-state
// marker: absent by default, set by ImportState, cleared by the adopting Update.
func TestImportAdoptionMarkerRoundTrip(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	private := &fakePrivateState{}

	pending, diags := importAdoptionPending(ctx, private)
	if diags.HasError() {
		t.Fatalf("importAdoptionPending: %v", diags)
	}
	if pending {
		t.Fatal("expected no adoption pending before import")
	}

	if diags := markImportAdoptionPending(ctx, private); diags.HasError() {
		t.Fatalf("markImportAdoptionPending: %v", diags)
	}

	pending, diags = importAdoptionPending(ctx, private)
	if diags.HasError() {
		t.Fatalf("importAdoptionPending: %v", diags)
	}
	if !pending {
		t.Fatal("expected adoption pending after import")
	}

	if diags := clearImportAdoptionPending(ctx, private); diags.HasError() {
		t.Fatalf("clearImportAdoptionPending: %v", diags)
	}

	pending, diags = importAdoptionPending(ctx, private)
	if diags.HasError() {
		t.Fatalf("importAdoptionPending: %v", diags)
	}
	if pending {
		t.Fatal("expected the adoption window to close after the adopting update")
	}
}

// TestImportAdoptionMarkerNilPrivateState makes sure a nil private state - which
// the framework may hand us outside an import - is treated as "not adopting"
// rather than panicking or erroring.
func TestImportAdoptionMarkerNilPrivateState(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	pending, diags := importAdoptionPending(ctx, nil)
	if diags.HasError() {
		t.Fatalf("importAdoptionPending: %v", diags)
	}
	if pending {
		t.Fatal("expected no adoption pending without private state")
	}
	if diags := markImportAdoptionPending(ctx, nil); diags.HasError() {
		t.Fatalf("markImportAdoptionPending: %v", diags)
	}
	if diags := clearImportAdoptionPending(ctx, nil); diags.HasError() {
		t.Fatalf("clearImportAdoptionPending: %v", diags)
	}
}

// fakePrivateState mimics the framework's private state data, whose concrete
// type lives in an internal package. Setting an empty value deletes the key,
// matching the real SetKey behaviour.
type fakePrivateState struct {
	data map[string][]byte
}

func (f *fakePrivateState) GetKey(_ context.Context, key string) ([]byte, diag.Diagnostics) {
	return f.data[key], nil
}

func (f *fakePrivateState) SetKey(_ context.Context, key string, value []byte) diag.Diagnostics {
	if f.data == nil {
		f.data = map[string][]byte{}
	}
	if len(value) == 0 {
		delete(f.data, key)
		return nil
	}
	f.data[key] = value
	return nil
}

func jsonHasKey(t *testing.T, body []byte, key string) bool {
	t.Helper()

	fields := map[string]json.RawMessage{}
	if err := json.Unmarshal(body, &fields); err != nil {
		t.Fatalf("failed to parse %s: %s", body, err)
	}
	_, ok := fields[key]
	return ok
}

func adoptionTestLookup(t *testing.T, networks map[string]string) subnetNetworkLookup {
	t.Helper()
	return func(_ context.Context, subnetID string) (string, error) {
		networkID, ok := networks[subnetID]
		if !ok {
			return "", fmt.Errorf("unexpected subnet lookup %s", subnetID)
		}
		return networkID, nil
	}
}

// TestVerifyAdoptedCreateOnlyFields pins the property the whole design rests
// on: an adoption can only put values into state that the API corroborates, a
// contradicted value is a mismatch, and an evidence-free response is
// explicitly unverifiable rather than silently accepted.
func TestVerifyAdoptedCreateOnlyFields(t *testing.T) {
	t.Parallel()

	evidence := func(body string) adoptionEvidence {
		ev, err := parseAdoptionEvidence([]byte(body))
		if err != nil {
			t.Fatalf("parseAdoptionEvidence: %s", err)
		}
		return ev
	}
	fullEvidence := evidence(`{
		"vrrp_ips": [{"subnet_id": "subnet-real"}],
		"additional_vips": []
	}`)
	networks := map[string]string{"subnet-real": "net-real"}

	cases := map[string]struct {
		plan, state CloudLoadBalancerModel
		evidence    adoptionEvidence
		want        map[string]adoptionVerdict
	}{
		"genuine adoption of both": {
			plan: CloudLoadBalancerModel{
				VipNetworkID: types.StringValue("net-real"),
				VipSubnetID:  types.StringValue("subnet-real"),
			},
			evidence: fullEvidence,
			want:     map[string]adoptionVerdict{"vip_network_id": adoptionConfirmed, "vip_subnet_id": adoptionConfirmed},
		},
		"stale marker day-two wrong network": {
			plan:     CloudLoadBalancerModel{VipNetworkID: types.StringValue("net-B")},
			evidence: fullEvidence,
			want:     map[string]adoptionVerdict{"vip_network_id": adoptionMismatch},
		},
		"wrong subnet": {
			plan:     CloudLoadBalancerModel{VipSubnetID: types.StringValue("subnet-other")},
			evidence: fullEvidence,
			want:     map[string]adoptionVerdict{"vip_subnet_id": adoptionMismatch},
		},
		"no evidence at all degrades to unverifiable": {
			plan: CloudLoadBalancerModel{
				VipNetworkID: types.StringValue("net-real"),
				VipSubnetID:  types.StringValue("subnet-real"),
			},
			evidence: evidence(`{}`),
			want:     map[string]adoptionVerdict{"vip_network_id": adoptionUnverifiable, "vip_subnet_id": adoptionUnverifiable},
		},
		"nothing adopted yields nothing": {
			plan:     CloudLoadBalancerModel{},
			evidence: fullEvidence,
			want:     map[string]adoptionVerdict{},
		},
		"already in state means not an adoption": {
			plan:     CloudLoadBalancerModel{VipSubnetID: types.StringValue("subnet-other")},
			state:    CloudLoadBalancerModel{VipSubnetID: types.StringValue("subnet-real")},
			evidence: fullEvidence,
			want:     map[string]adoptionVerdict{},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			findings, err := verifyAdoptedCreateOnlyFields(context.Background(), tc.plan, tc.state, tc.evidence, adoptionTestLookup(t, networks))
			if err != nil {
				t.Fatalf("verifyAdoptedCreateOnlyFields: %s", err)
			}
			got := map[string]adoptionVerdict{}
			for _, finding := range findings {
				got[finding.attribute] = finding.verdict
			}
			if !maps.Equal(got, tc.want) {
				t.Fatalf("got %v, want %v", got, tc.want)
			}
		})
	}
}

// TestVerifyAdoptedCreateOnlyFieldsLookupFailure: a transient subnet lookup
// failure must abort verification (the caller fails the apply, marker intact)
// rather than being folded into a verdict.
func TestVerifyAdoptedCreateOnlyFieldsLookupFailure(t *testing.T) {
	t.Parallel()

	plan := CloudLoadBalancerModel{VipNetworkID: types.StringValue("net-real")}
	ev, err := parseAdoptionEvidence([]byte(`{"vrrp_ips": [{"subnet_id": "subnet-real"}]}`))
	if err != nil {
		t.Fatalf("parseAdoptionEvidence: %s", err)
	}
	_, err = verifyAdoptedCreateOnlyFields(context.Background(), plan, CloudLoadBalancerModel{}, ev,
		func(context.Context, string) (string, error) { return "", fmt.Errorf("boom") })
	if err == nil {
		t.Fatal("expected a lookup failure to surface as an error")
	}
}

// TestVerifyCoversEveryCreateOnlyField is the drift guard between the strip
// map and the verifier: every field the strip helper can swallow must produce
// a finding when adopted, so a new create-only attribute cannot be added to
// one without the other.
func TestVerifyCoversEveryCreateOnlyField(t *testing.T) {
	t.Parallel()

	plan := CloudLoadBalancerModel{
		VipNetworkID: types.StringValue("net-a"),
		VipSubnetID:  types.StringValue("subnet-a"),
	}
	findings, err := verifyAdoptedCreateOnlyFields(context.Background(), plan, CloudLoadBalancerModel{}, adoptionEvidence{}, adoptionTestLookup(t, nil))
	if err != nil {
		t.Fatalf("verifyAdoptedCreateOnlyFields: %s", err)
	}
	seen := map[string]bool{}
	for _, finding := range findings {
		seen[finding.attribute] = true
	}
	for name := range createOnlyFields {
		if !seen[name] {
			t.Errorf("createOnlyFields entry %s has no verification finding", name)
		}
	}
}
