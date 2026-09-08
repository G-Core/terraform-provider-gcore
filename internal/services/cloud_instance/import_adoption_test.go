package cloud_instance

import (
	"bytes"
	"context"
	"os"
	"slices"
	"strings"
	"testing"
	"unicode"

	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// fakePrivate stands in for the framework's private state, whose concrete type
// lives in an internal package and cannot be constructed from here.
type fakePrivate map[string][]byte

func (f fakePrivate) GetKey(_ context.Context, key string) ([]byte, diag.Diagnostics) {
	return f[key], nil
}

// SetKey mirrors the framework's semantics: a nil or empty value deletes the
// key rather than storing an empty entry. The marker lifecycle leans on that -
// retirement is a SetKey(key, nil) - so the fake must not diverge from it.
func (f fakePrivate) SetKey(_ context.Context, key string, value []byte) diag.Diagnostics {
	if len(value) == 0 {
		delete(f, key)
		return nil
	}
	f[key] = slices.Clone(value)

	return nil
}

func TestImportAdoptionPending(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	pending, diags := importAdoptionPending(ctx, nil)
	if diags.HasError() {
		t.Fatalf("importAdoptionPending: %v", diags)
	}
	if pending {
		t.Error("expected no adoption pending without private state")
	}

	if pending, _ = importAdoptionPending(ctx, fakePrivate{}); pending {
		t.Error("expected no adoption pending when the key is absent")
	}

	marked := fakePrivate{importAdoptionPrivateKey: importAdoptionValuePending}
	if pending, _ = importAdoptionPending(ctx, marked); !pending {
		t.Error("expected adoption pending when the marker is set")
	}

	// Clearing writes an empty value rather than removing the entry, so a
	// zero-length read must count as absent or the window would never close.
	cleared := fakePrivate{importAdoptionPrivateKey: {}}
	if pending, _ = importAdoptionPending(ctx, cleared); pending {
		t.Error("expected a zero-length value to read as absent")
	}
}

func config(values map[string]string) *map[string]customfield.MetaStringValue {
	m := map[string]customfield.MetaStringValue{}
	for k, v := range values {
		m[k] = customfield.NewMetaStringValue(v)
	}

	return &m
}

// TestAdoptedCreateOnlyAttributes pins which attributes the plan-time warning
// names. Only a transition from "absent in prior state" to "set in config"
// counts: that is the one shape where the value is recorded into state without
// ever reaching the API.
func TestAdoptedCreateOnlyAttributes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		plan  CloudInstanceModel
		state CloudInstanceModel
		want  []string
	}{
		{
			name: "nothing configured, nothing adopted",
		},
		{
			name:  "single attribute adopted over null state",
			plan:  CloudInstanceModel{Username: types.StringValue("admin")},
			state: CloudInstanceModel{Username: types.StringNull()},
			want:  []string{"username"},
		},
		{
			// ssh_key_name is the one adopted attribute the API can corroborate:
			// ImportState resolves the keypair UUID to its name, so this shape
			// only arises when that lookup failed. Adopting beats the
			// alternative, which is replacing a healthy instance.
			name:  "ssh_key_name adopted when the import lookup could not resolve it",
			plan:  CloudInstanceModel{SSHKeyName: types.StringValue("my-key")},
			state: CloudInstanceModel{SSHKeyName: types.StringNull()},
			want:  []string{"ssh_key_name"},
		},
		{
			// The normal import path: ImportState back-filled the name, so there
			// is nothing to adopt and an unchanged key plans no action at all.
			name:  "a back-filled ssh_key_name is not adopted",
			plan:  CloudInstanceModel{SSHKeyName: types.StringValue("my-key")},
			state: CloudInstanceModel{SSHKeyName: types.StringValue("my-key")},
			want:  nil,
		},
		{
			// A genuine key change must stay on the replacement path - the
			// keypair cannot be swapped on a running instance.
			name:  "changing a known ssh_key_name is a real change, not an adoption",
			plan:  CloudInstanceModel{SSHKeyName: types.StringValue("new-key")},
			state: CloudInstanceModel{SSHKeyName: types.StringValue("old-key")},
			want:  nil,
		},
		{
			name: "every create-only attribute at once, reported in a stable order",
			plan: CloudInstanceModel{
				AllowAppPorts: types.BoolValue(true),
				NameTemplate:  types.StringValue("{ip_octets}"),
				ServergroupID: types.StringValue("sg-1"),
				SSHKeyName:    types.StringValue("my-key"),
				UserData:      types.StringValue("Zm9v"),
				Username:      types.StringValue("admin"),
				Configuration: config(map[string]string{"k": "v"}),
			},
			want: []string{"allow_app_ports", "configuration", "name_template", "servergroup_id", "ssh_key_name", "user_data", "username"},
		},
		{
			name:  "a known prior value is a real change, not an adoption",
			plan:  CloudInstanceModel{Username: types.StringValue("root")},
			state: CloudInstanceModel{Username: types.StringValue("admin")},
			want:  nil,
		},
		{
			name:  "removal is not an adoption",
			plan:  CloudInstanceModel{Username: types.StringNull()},
			state: CloudInstanceModel{Username: types.StringValue("admin")},
			want:  nil,
		},
		{
			name:  "an empty configuration map counts as unset",
			plan:  CloudInstanceModel{Configuration: config(nil)},
			state: CloudInstanceModel{},
			want:  []string{"configuration"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := adoptedCreateOnlyAttributes(tt.plan, tt.state)
			if !slices.Equal(got, tt.want) {
				t.Errorf("adoptedCreateOnlyAttributes = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestImportAdoptionPhaseOf pins the parse rules the lifecycle depends on.
// The fallback matters as much as the happy path: a marker written by an
// earlier release ({"pending":true}), or bytes mangled in transit, must read
// as armed - one apply of window, then retirement - never as a phase that
// keeps the window open forever or one that retires it before it was used.
func TestImportAdoptionPhaseOf(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		value []byte
		want  string
	}{
		{name: "nil is absent", value: nil, want: ""},
		{name: "zero-length is absent", value: []byte{}, want: ""},
		{name: "pending", value: importAdoptionValuePending, want: importAdoptionPhasePending},
		{name: "armed", value: importAdoptionValueArmed, want: importAdoptionPhaseArmed},
		{name: "spent", value: importAdoptionValueSpent, want: importAdoptionPhaseSpent},
		{name: "legacy pre-phase marker reads as armed", value: []byte(`{"pending":true}`), want: importAdoptionPhaseArmed},
		{name: "garbage reads as armed", value: []byte(`not-json`), want: importAdoptionPhaseArmed},
		{name: "empty object reads as armed", value: []byte(`{}`), want: importAdoptionPhaseArmed},
		{name: "unknown phase reads as armed", value: []byte(`{"phase":"bogus"}`), want: importAdoptionPhaseArmed},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := importAdoptionPhaseOf(tt.value); got != tt.want {
				t.Errorf("importAdoptionPhaseOf(%q) = %q, want %q", tt.value, got, tt.want)
			}
		})
	}
}

// TestAdvanceImportAdoption walks every stored value one step and checks both
// what is written and - the load-bearing part - whether importAdoptionPending
// still sees a marker afterwards. Spent must remain PRESENT and pending:
// every reader tests only presence, and it is the spent-but-present state that
// lets the plan of the invocation that spends the marker still adopt. A
// refactor that "simplifies" armed straight to deletion passes every other
// test and silently breaks adoption on the `terraform import` path; the
// wantPending assertions below are what catch it.
func TestAdvanceImportAdoption(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		initial     []byte // nil = key absent
		wantValue   []byte // nil = key must be gone afterwards
		wantPending bool
	}{
		{
			name: "absent key stays absent",
		},
		{
			name:        "pending advances to armed",
			initial:     importAdoptionValuePending,
			wantValue:   importAdoptionValueArmed,
			wantPending: true,
		},
		{
			name:        "armed advances to spent, still pending",
			initial:     importAdoptionValueArmed,
			wantValue:   importAdoptionValueSpent,
			wantPending: true,
		},
		{
			name:    "spent retires the marker",
			initial: importAdoptionValueSpent,
		},
		{
			// The legacy value parses as armed, so its next stop is spent -
			// an already-open window keeps exactly one apply, then retires.
			name:        "legacy pre-phase marker advances to spent",
			initial:     []byte(`{"pending":true}`),
			wantValue:   importAdoptionValueSpent,
			wantPending: true,
		},
		{
			name:        "garbage advances to spent",
			initial:     []byte(`not-json`),
			wantValue:   importAdoptionValueSpent,
			wantPending: true,
		},
		{
			name:        "unknown phase advances to spent",
			initial:     []byte(`{"phase":"bogus"}`),
			wantValue:   importAdoptionValueSpent,
			wantPending: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			private := fakePrivate{}
			if tt.initial != nil {
				private[importAdoptionPrivateKey] = slices.Clone(tt.initial)
			}

			if diags := advanceImportAdoption(ctx, private); diags.HasError() {
				t.Fatalf("advanceImportAdoption: %v", diags)
			}

			got, present := private[importAdoptionPrivateKey]
			if tt.wantValue == nil {
				if present {
					t.Fatalf("expected the marker to be gone, got %s", got)
				}
			} else if !bytes.Equal(got, tt.wantValue) {
				t.Fatalf("marker = %s, want %s", got, tt.wantValue)
			}

			pending, diags := importAdoptionPending(ctx, private)
			if diags.HasError() {
				t.Fatalf("importAdoptionPending: %v", diags)
			}
			if pending != tt.wantPending {
				t.Errorf("importAdoptionPending = %t, want %t", pending, tt.wantPending)
			}
		})
	}

	// Nil private state is what the framework hands over when nothing was
	// ever stored; the advance must be a silent no-op, not a panic.
	if diags := advanceImportAdoption(context.Background(), nil); diags.HasError() {
		t.Errorf("advanceImportAdoption(nil): %v", diags)
	}
}

// TestRearmImportAdoption covers Update's error path: whatever the marker held
// when the adopting apply failed - spent after the preceding refresh, or even
// nothing - the re-arm must leave a full apply of window behind, so the retry
// adopts instead of planning a destroy+recreate.
func TestRearmImportAdoption(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	for _, tt := range []struct {
		name    string
		initial []byte // nil = key absent
	}{
		{name: "absent key is armed"},
		{name: "spent marker is re-armed", initial: importAdoptionValueSpent},
		{name: "pending marker is re-armed", initial: importAdoptionValuePending},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			private := fakePrivate{}
			if tt.initial != nil {
				private[importAdoptionPrivateKey] = slices.Clone(tt.initial)
			}

			if diags := rearmImportAdoption(ctx, private); diags.HasError() {
				t.Fatalf("rearmImportAdoption: %v", diags)
			}
			if got := private[importAdoptionPrivateKey]; !bytes.Equal(got, importAdoptionValueArmed) {
				t.Fatalf("marker = %s, want %s", got, importAdoptionValueArmed)
			}

			pending, diags := importAdoptionPending(ctx, private)
			if diags.HasError() {
				t.Fatalf("importAdoptionPending: %v", diags)
			}
			if !pending {
				t.Error("expected the re-armed marker to read as pending")
			}
		})
	}

	if diags := rearmImportAdoption(ctx, nil); diags.HasError() {
		t.Errorf("rearmImportAdoption(nil): %v", diags)
	}
}

// TestImportAdoptionMarkerRoundTrip walks the marker through both of its whole
// lives: the Read-driven one (import -> three refreshes -> gone) and the
// Update-driven one (import -> successful adopting apply -> gone). The
// exact-bytes check on the pending value pins what ImportState actually
// persists - phase parsing is only as good as the writer agreeing with it.
func TestImportAdoptionMarkerRoundTrip(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	private := fakePrivate{}

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
	if got := private[importAdoptionPrivateKey]; !bytes.Equal(got, importAdoptionValuePending) {
		t.Fatalf("ImportState writes %s, want %s", got, importAdoptionValuePending)
	}

	// The Read-driven life: each refresh advances one phase, and the marker
	// must still read as pending right up until the third advance removes it.
	for i, want := range [][]byte{importAdoptionValueArmed, importAdoptionValueSpent, nil} {
		if diags := advanceImportAdoption(ctx, private); diags.HasError() {
			t.Fatalf("advanceImportAdoption #%d: %v", i+1, diags)
		}

		got, present := private[importAdoptionPrivateKey]
		if want == nil {
			if present {
				t.Fatalf("expected the marker to be retired after advance #%d, got %s", i+1, got)
			}
		} else if !bytes.Equal(got, want) {
			t.Fatalf("marker after advance #%d = %s, want %s", i+1, got, want)
		}

		pending, diags = importAdoptionPending(ctx, private)
		if diags.HasError() {
			t.Fatalf("importAdoptionPending after advance #%d: %v", i+1, diags)
		}
		if wantPending := want != nil; pending != wantPending {
			t.Errorf("importAdoptionPending after advance #%d = %t, want %t", i+1, pending, wantPending)
		}
	}

	// The Update-driven life: a successful adopting apply clears the marker
	// directly, whatever phase it was in.
	if diags := markImportAdoptionPending(ctx, private); diags.HasError() {
		t.Fatalf("markImportAdoptionPending: %v", diags)
	}
	if diags := clearImportAdoptionPending(ctx, private); diags.HasError() {
		t.Fatalf("clearImportAdoptionPending: %v", diags)
	}
	if _, ok := private[importAdoptionPrivateKey]; ok {
		t.Fatal("expected the cleared key to be removed from private state")
	}
	if pending, _ = importAdoptionPending(ctx, private); pending {
		t.Fatal("expected the adoption window to be closed after Update")
	}
}

// TestImportAdoptionMarkerNilPrivateState covers the private state being absent
// altogether, which the framework permits.
func TestImportAdoptionMarkerNilPrivateState(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	if diags := markImportAdoptionPending(ctx, nil); diags.HasError() {
		t.Errorf("markImportAdoptionPending: %v", diags)
	}
	if diags := clearImportAdoptionPending(ctx, nil); diags.HasError() {
		t.Errorf("clearImportAdoptionPending: %v", diags)
	}
	if diags := advanceImportAdoption(ctx, nil); diags.HasError() {
		t.Errorf("advanceImportAdoption: %v", diags)
	}
	if diags := rearmImportAdoption(ctx, nil); diags.HasError() {
		t.Errorf("rearmImportAdoption: %v", diags)
	}

	pending, diags := importAdoptionPending(ctx, nil)
	if diags.HasError() {
		t.Errorf("importAdoptionPending: %v", diags)
	}
	if pending {
		t.Error("expected no adoption pending without private state")
	}
}

// methodBody returns the source of the named CloudInstanceResource method,
// located by its func header and bounded by the next top-level func. Header
// matching keeps the pin anchored to the real method rather than a comment or
// a similarly named helper.
func methodBody(t *testing.T, source, method string) string {
	t.Helper()

	header := "func (r *CloudInstanceResource) " + method + "("
	start := strings.Index(source, header)
	if start < 0 {
		t.Fatalf("could not find %q in resource.go", header)
	}
	body := source[start:]
	if end := strings.Index(body[len(header):], "\nfunc "); end >= 0 {
		body = body[:len(header)+end]
	}

	return body
}

// stripSpace removes all whitespace so the assertions below survive gofmt
// re-wrapping without loosening what they match.
func stripSpace(s string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, s)
}

// TestImportAdoptionLifecycleWiringPinned asserts the two lifecycle hooks are
// literally present in resource.go, read off disk.
//
// This is not a style check. resource.go is a generated file and these hooks
// are sealed custom code; deletions and small hunks in generated files are the
// known class of change that silently evaporates when the seal is rebuilt
// against a regen. Every behavioural test in this package calls the helpers
// directly, so none of them would notice the call sites vanishing - this test
// is the only thing between a bad reseal and shipping the regression.
func TestImportAdoptionLifecycleWiringPinned(t *testing.T) {
	t.Parallel()

	raw, err := os.ReadFile("resource.go")
	if err != nil {
		t.Fatalf("reading resource.go: %s", err)
	}
	source := string(raw)

	if read := stripSpace(methodBody(t, source, "Read")); !strings.Contains(read, "advanceImportAdoption(ctx,resp.Private)") {
		t.Error("Read no longer calls advanceImportAdoption(ctx, resp.Private): " +
			"a reseal dropped the marker-retirement hunk; the stale-marker defect returns silently - " +
			"the import-adoption marker outlives its import again and a day-two attribute set is " +
			"adopted into state instead of planning the replacement that would apply it")
	}

	if update := stripSpace(methodBody(t, source, "Update")); !strings.Contains(update, "rearmImportAdoption(ctx,resp.Private)") {
		t.Error("Update no longer calls rearmImportAdoption(ctx, resp.Private): " +
			"a reseal dropped the failed-apply re-arm hunk; a transient API error during the " +
			"adopting apply now consumes the adoption window and the retry silently plans a " +
			"destroy+recreate of the freshly imported instance")
	}
}

// TestImportAdoptionWarning pins the one promise in the warning that can be
// factually wrong. Adoption normally only writes state, so the warning says
// nothing reaches the API - but "servergroup_id" is placed for real by Update
// (it reads the group's members and calls AddToPlacementGroup when the instance
// is missing from it), so the absolute claim has to give way to the exception
// whenever that attribute is in the adopted set.
func TestImportAdoptionWarning(t *testing.T) {
	t.Parallel()

	const nothingSent = "Nothing is sent to the API for them"

	tests := []struct {
		name        string
		adopted     []string
		wantContain []string
		wantAbsent  []string
	}{
		{
			name:        "state-only attributes keep the absolute claim",
			adopted:     []string{"user_data", "username"},
			wantContain: []string{nothingSent, `"user_data", "username"`},
			wantAbsent:  []string{"placement group"},
		},
		{
			name:        "servergroup_id replaces the claim with the API exception",
			adopted:     []string{"servergroup_id"},
			wantContain: []string{"placement group", "actual placement", `"servergroup_id"`},
			wantAbsent:  []string{nothingSent},
		},
		{
			name:        "the exception survives being mixed with state-only attributes",
			adopted:     []string{"servergroup_id", "username"},
			wantContain: []string{"placement group", `"servergroup_id", "username"`},
			wantAbsent:  []string{nothingSent},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			detail := importAdoptionWarning(test.adopted).Detail()

			for _, want := range test.wantContain {
				if !strings.Contains(detail, want) {
					t.Errorf("warning detail is missing %q:\n%s", want, detail)
				}
			}

			for _, unwanted := range test.wantAbsent {
				if strings.Contains(detail, unwanted) {
					t.Errorf("warning detail should not contain %q:\n%s", unwanted, detail)
				}
			}
		})
	}
}
