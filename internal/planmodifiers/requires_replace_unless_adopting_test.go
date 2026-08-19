package planmodifiers

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

// TestRequiresReplaceUnlessAdopting walks the full transition matrix for a
// create-only attribute the API does not return on GET.
//
// The case that matters most is the last pair: a null prior value means
// "adopt this" only while an import adoption is pending. Without that
// distinction a load balancer created with the attribute omitted would silently
// accept a day-two change that never reaches the API.
func TestRequiresReplaceUnlessAdopting(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		stateRawNull   bool
		planRawNull    bool
		valueUnchanged bool
		priorValueNull bool
		adopting       bool
		want           bool
	}{
		{
			name:           "resource creation never replaces",
			stateRawNull:   true,
			priorValueNull: true,
		},
		{
			name:        "resource destruction never replaces",
			planRawNull: true,
		},
		{
			name:           "unchanged value does not replace",
			valueUnchanged: true,
		},
		{
			name:           "unset on both sides does not replace",
			valueUnchanged: true,
			priorValueNull: true,
		},
		{
			name: "changed known value replaces",
			want: true,
		},
		{
			name: "removal replaces",
			want: true,
		},
		{
			name: "unknown planned value replaces",
			want: true,
		},
		{
			name:           "null prior value adopts while an import is pending",
			priorValueNull: true,
			adopting:       true,
			want:           false,
		},
		{
			name:           "null prior value replaces when no import is pending",
			priorValueNull: true,
			adopting:       false,
			want:           true,
		},
		{
			name:           "adoption marker does not override a known prior value",
			priorValueNull: false,
			adopting:       true,
			want:           true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := requiresReplaceUnlessAdopting(
				tt.stateRawNull,
				tt.planRawNull,
				tt.valueUnchanged,
				tt.priorValueNull,
				tt.adopting,
			)
			if got != tt.want {
				t.Errorf("requiresReplaceUnlessAdopting = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAdoptionPending(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	const key = "adoption_key"

	pending, diags := adoptionPending(ctx, nil, key)
	if diags.HasError() {
		t.Fatalf("adoptionPending: %v", diags)
	}
	if pending {
		t.Error("expected no adoption pending without private state")
	}

	empty := fakePrivate{}
	if pending, _ = adoptionPending(ctx, empty, key); pending {
		t.Error("expected no adoption pending when the key is absent")
	}

	marked := fakePrivate{key: []byte(`{"pending":true}`)}
	if pending, _ = adoptionPending(ctx, marked, key); !pending {
		t.Error("expected adoption pending when the key is set")
	}

	zeroLength := fakePrivate{key: []byte{}}
	if pending, _ = adoptionPending(ctx, zeroLength, key); pending {
		t.Error("expected a zero-length value to read as absent")
	}
}

type fakePrivate map[string][]byte

func (f fakePrivate) GetKey(_ context.Context, key string) ([]byte, diag.Diagnostics) {
	return f[key], nil
}
