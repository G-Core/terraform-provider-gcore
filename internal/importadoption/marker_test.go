package importadoption_test

import (
	"context"
	"slices"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	"github.com/G-Core/terraform-provider-gcore/internal/importadoption"
)

// fakePrivate stands in for the framework's private state, whose concrete type
// lives in an internal package and cannot be constructed from here.
type fakePrivate map[string][]byte

func (f fakePrivate) GetKey(_ context.Context, key string) ([]byte, diag.Diagnostics) {
	return f[key], nil
}

// SetKey mirrors the framework's semantics: a nil or empty value deletes the
// key rather than storing an empty entry. Retirement is a SetKey(key, nil), so
// the fake must not diverge from that.
func (f fakePrivate) SetKey(_ context.Context, key string, value []byte) diag.Diagnostics {
	if len(value) == 0 {
		delete(f, key)
		return nil
	}
	f[key] = slices.Clone(value)
	return nil
}

var marker = importadoption.Marker{Key: "test_import_adoption_pending"}

func TestPhaseOf(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		value []byte
		want  importadoption.Phase
	}{
		{"absent", nil, importadoption.PhaseNone},
		{"empty", []byte{}, importadoption.PhaseNone},
		{"pending", []byte(`{"phase":"pending"}`), importadoption.PhasePending},
		{"armed", []byte(`{"phase":"armed"}`), importadoption.PhaseArmed},
		{"spent", []byte(`{"phase":"spent"}`), importadoption.PhaseSpent},
		// The legacy shape written by the older per-resource copies: one apply
		// of window, then retirement - no state upgrade needed.
		{"legacy pending:true", []byte(`{"pending":true}`), importadoption.PhaseArmed},
		{"garbage", []byte(`not json`), importadoption.PhaseArmed},
		{"unknown phase", []byte(`{"phase":"later"}`), importadoption.PhaseArmed},
		{"empty phase", []byte(`{"phase":""}`), importadoption.PhaseArmed},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if got := importadoption.PhaseOf(tc.value); got != tc.want {
				t.Errorf("PhaseOf(%s) = %q, want %q", tc.value, got, tc.want)
			}
		})
	}
}

// TestAdvanceSequence is the retirement proof on the do-nothing path: three
// Reads retire a marker with no Update ever running, and "spent" still counts
// as adopting so the first apply after `terraform import` can adopt.
func TestAdvanceSequence(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	p := fakePrivate{}

	if d := marker.Mark(ctx, p); d.HasError() {
		t.Fatalf("Mark: %v", d)
	}
	steps := []struct {
		after   importadoption.Phase
		pending bool
	}{
		{importadoption.PhaseArmed, true},
		{importadoption.PhaseSpent, true},
		{importadoption.PhaseNone, false},
		{importadoption.PhaseNone, false}, // no-op when absent
	}
	for i, s := range steps {
		if d := marker.Advance(ctx, p); d.HasError() {
			t.Fatalf("Advance #%d: %v", i+1, d)
		}
		if got := importadoption.PhaseOf(p[marker.Key]); got != s.after {
			t.Errorf("after Advance #%d: phase %q, want %q", i+1, got, s.after)
		}
		pending, _ := marker.Pending(ctx, p)
		if pending != s.pending {
			t.Errorf("after Advance #%d: Pending = %v, want %v", i+1, pending, s.pending)
		}
	}
	if _, ok := p[marker.Key]; ok {
		t.Error("retired marker must be deleted, not stored empty")
	}
}

func TestBeginUpdateSettle(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	tests := []struct {
		name         string
		initial      []byte
		fail         bool
		wantAdopting bool
		wantAfter    importadoption.Phase
	}{
		{"not adopting, ok", nil, false, false, importadoption.PhaseNone},
		{"not adopting, error", nil, true, false, importadoption.PhaseNone},
		{"adopting ok clears", []byte(`{"phase":"spent"}`), false, true, importadoption.PhaseNone},
		{"adopting error re-arms", []byte(`{"phase":"spent"}`), true, true, importadoption.PhaseArmed},
		{"adopting pending error re-arms", []byte(`{"phase":"pending"}`), true, true, importadoption.PhaseArmed},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			p := fakePrivate{}
			if tc.initial != nil {
				p[marker.Key] = tc.initial
			}
			var diags diag.Diagnostics
			adopting, settle := marker.BeginUpdate(ctx, p, p, &diags)
			if adopting != tc.wantAdopting {
				t.Fatalf("adopting = %v, want %v", adopting, tc.wantAdopting)
			}
			if tc.fail {
				diags.AddError("boom", "simulated apply failure")
			}
			settle()
			if got := importadoption.PhaseOf(p[marker.Key]); got != tc.wantAfter {
				t.Errorf("after settle: phase %q, want %q", got, tc.wantAfter)
			}
		})
	}
}

func TestNilPrivateStateIsSafe(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	if d := marker.Mark(ctx, nil); d.HasError() {
		t.Errorf("Mark(nil): %v", d)
	}
	if d := marker.Advance(ctx, nil); d.HasError() {
		t.Errorf("Advance(nil): %v", d)
	}
	if d := marker.Rearm(ctx, nil); d.HasError() {
		t.Errorf("Rearm(nil): %v", d)
	}
	if d := marker.Clear(ctx, nil); d.HasError() {
		t.Errorf("Clear(nil): %v", d)
	}
	if pending, d := marker.Pending(ctx, nil); pending || d.HasError() {
		t.Errorf("Pending(nil) = %v, %v", pending, d)
	}
	var diags diag.Diagnostics
	adopting, settle := marker.BeginUpdate(ctx, nil, nil, &diags)
	if adopting {
		t.Error("BeginUpdate(nil) reported adopting")
	}
	settle()
	if diags.HasError() {
		t.Errorf("BeginUpdate(nil) settle: %v", diags)
	}
}

func TestMarkRearmClear(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	p := fakePrivate{}

	_ = marker.Mark(ctx, p)
	if got := importadoption.PhaseOf(p[marker.Key]); got != importadoption.PhasePending {
		t.Errorf("after Mark: %q", got)
	}
	_ = marker.Rearm(ctx, p)
	if got := importadoption.PhaseOf(p[marker.Key]); got != importadoption.PhaseArmed {
		t.Errorf("after Rearm: %q", got)
	}
	_ = marker.Clear(ctx, p)
	if _, ok := p[marker.Key]; ok {
		t.Error("after Clear: key still present")
	}
	// A zero-length stored value must read as absent or the window never closes.
	p[marker.Key] = []byte{}
	if pending, _ := marker.Pending(ctx, p); pending {
		t.Error("zero-length value read as pending")
	}
}
