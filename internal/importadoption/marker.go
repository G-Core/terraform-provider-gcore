// Package importadoption is the import-adoption marker shared by resources whose
// create-only attributes the API never echoes (`no_refresh`). A resource declares
//
//	var importAdoption = importadoption.Marker{Key: "<resource>_import_adoption_pending"}
//
// wires four hooks (ImportState: Mark; Read: Advance; Update: BeginUpdate +
// defer settle; ModifyPlan: PlanHook.ModifyPlan) and passes importAdoption.Key
// to planmodifiers.*RequiresReplaceUnlessAdopting.
//
// Why a marker at all: "prior state is null" does NOT mean "just imported". It
// is equally true of a resource created with the attribute omitted, and a user
// may set such an attribute on day two. Adopting there would plan a clean
// in-place update, send nothing the API acts on, and record a value the
// infrastructure does not have - state that no refresh can ever correct. The
// marker keeps that case on the replacement path, the only one that applies the
// value.
//
// Lifecycle, advanced ONLY by Read - the one provider method whose private-state
// write Terraform persists on every apply, no-op applies included:
//
//	pending -> armed -> spent -> deleted
//
// ImportState writes pending; the Read Terraform runs right after the import
// persists armed; the refresh at the start of the next apply persists spent even
// when that apply is a no-op; the Read after that deletes it. Every reader tests
// presence, so "spent" still adopts - that is what lets the first apply after
// `terraform import` adopt. Do not collapse spent into deletion: deleting on
// armed breaks adoption on the `terraform import` path.
//
// Retirement proof on the do-nothing path: a resource imported with a config
// that adopts nothing and only ever no-op-applied retires its marker after two
// applies with zero user action. Plain `terraform plan` persists nothing, so
// any number of plans keep the window - and ModifyPlan therefore never writes
// private state (a plan whose only effect would be a marker write is a NoOp
// plan and Terraform drops its PlannedPrivate).
//
// Window: the first apply after import (one more for `import` blocks, whose
// plan runs in the import walk). A refresh-only apply between import and the
// adopting apply consumes it - loud (the next adopting apply plans a
// replacement), recoverable with -refresh=false or a re-import.
//
// Known limitations, listed for the reviewer rather than accepted:
//
//   - `-refresh=false` skips Read, so no-op applies under it never advance the
//     marker and the window stays open until the next refreshing apply or a
//     successful adopting Update. The plan-time warning still fires (it reads
//     the marker, not Read), so a late set is warned, not silent. No provider
//     hook can close this: private state persists only through Read or a
//     non-NoOp apply.
//   - An apply in which Update never runs (an upstream dependency fails, the
//     run is interrupted, the plugin crashes) still persists the refresh's
//     Advance without a settle. On the `terraform import` path that refresh
//     is the one that spends the marker, so the retry plans a replacement
//     (recover with -refresh=false on the retry, or re-import); on the
//     `import` block path one such apply is survived.
//
// Invariant every Update owes BeginUpdate: each non-error exit writes plan
// state. An early `return` without an error and without State.Set would clear
// the window without adopting. Review checklist: read every `return` in Update.
package importadoption

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

// Marker is one resource's import-adoption marker, identified by its private
// state key. The zero value is unusable; declare it with a Key.
type Marker struct{ Key string }

// Phase is the marker's lifecycle position.
type Phase string

const (
	PhaseNone    Phase = ""
	PhasePending Phase = "pending"
	PhaseArmed   Phase = "armed"
	PhaseSpent   Phase = "spent"
)

// Reader and Writer are the slices of the framework's private state used here;
// the concrete *privatestate.ProviderData lives in an internal framework
// package. Every method below is nil-safe: the framework hands a nil private
// state to some code paths (and to unit tests).
type Reader interface {
	GetKey(ctx context.Context, key string) ([]byte, diag.Diagnostics)
}

// Writer is the write half of the framework's private state.
type Writer interface {
	SetKey(ctx context.Context, key string, value []byte) diag.Diagnostics
}

// ReadWriter is both halves, as Read receives them.
type ReadWriter interface {
	Reader
	Writer
}

func encode(p Phase) []byte {
	b, _ := json.Marshal(struct {
		Phase Phase `json:"phase"`
	}{p})
	return b
}

// PhaseOf parses a stored value. Absent -> PhaseNone. Any non-empty value that
// is not a known phase - including the legacy {"pending":true} written by the
// older per-resource copies of this pattern - is PhaseArmed: one apply of
// window, then retirement. That is the whole migration story for a resource
// moving from a mark/clear-only marker: no state upgrade needed.
func PhaseOf(v []byte) Phase {
	if len(v) == 0 {
		return PhaseNone
	}
	var m struct {
		Phase Phase `json:"phase"`
	}
	if err := json.Unmarshal(v, &m); err != nil {
		return PhaseArmed
	}
	switch m.Phase {
	case PhasePending, PhaseArmed, PhaseSpent:
		return m.Phase
	}
	return PhaseArmed
}

// Mark opens the window. Call from ImportState immediately before State.Set.
func (m Marker) Mark(ctx context.Context, w Writer) diag.Diagnostics {
	if w == nil {
		return nil
	}
	return w.SetKey(ctx, m.Key, encode(PhasePending))
}

// Advance moves the marker one phase: pending -> armed -> spent -> deleted.
// Call exactly once per successful Read, after any RemoveResource branch,
// immediately before State.Set. No-op when no marker exists.
func (m Marker) Advance(ctx context.Context, rw ReadWriter) diag.Diagnostics {
	if rw == nil {
		return nil
	}
	v, diags := rw.GetKey(ctx, m.Key)
	if diags.HasError() {
		return diags
	}
	switch PhaseOf(v) {
	case PhasePending:
		diags.Append(rw.SetKey(ctx, m.Key, encode(PhaseArmed))...)
	case PhaseArmed:
		diags.Append(rw.SetKey(ctx, m.Key, encode(PhaseSpent))...)
	case PhaseSpent:
		diags.Append(rw.SetKey(ctx, m.Key, nil)...) // nil deletes the key
	}
	return diags
}

// Rearm restores a full apply of window. Called by settle when an adopting
// Update returns with an error diagnostic, so an error raised inside Update
// (after BeginUpdate) does not let the next refresh retire the marker and turn
// the retry into a replacement. Failures outside Update - a panic, a framework
// or core check after Update returned, an apply in which Update never ran -
// are not covered; see the package doc.
func (m Marker) Rearm(ctx context.Context, w Writer) diag.Diagnostics {
	if w == nil {
		return nil
	}
	return w.SetKey(ctx, m.Key, encode(PhaseArmed))
}

// Pending reports whether the resource is inside the adoption window. It is
// presence, deliberately: spent still adopts. Byte-for-byte what
// planmodifiers.adoptionPending does for the typed plan modifiers.
func (m Marker) Pending(ctx context.Context, r Reader) (bool, diag.Diagnostics) {
	if r == nil {
		return false, nil
	}
	v, diags := r.GetKey(ctx, m.Key)
	return len(v) > 0, diags
}

// Clear closes the window. Setting a nil value deletes the key.
func (m Marker) Clear(ctx context.Context, w Writer) diag.Diagnostics {
	if w == nil {
		return nil
	}
	return w.SetKey(ctx, m.Key, nil)
}

// BeginUpdate is the entire Update-side contract. Call it as the FIRST
// statement of Update - before the plan and state decode, which can fail - and
// `defer settle()` on the next line: settle re-arms when Update returns with an
// error diagnostic (the framework returns prior state plus this private data,
// and core persists them when the state is non-null) and clears on a
// successful one. The guarantee covers exactly the errors raised inside
// Update after this call. See the package doc for the invariant the caller
// owes and for the failure paths that are not covered.
func (m Marker) BeginUpdate(ctx context.Context, r Reader, w Writer, diags *diag.Diagnostics) (adopting bool, settle func()) {
	adopting, d := m.Pending(ctx, r)
	diags.Append(d...)
	settle = func() {
		if !adopting {
			return
		}
		if diags.HasError() {
			diags.Append(m.Rearm(ctx, w)...)
			return
		}
		diags.Append(m.Clear(ctx, w)...)
	}
	return adopting, settle
}
