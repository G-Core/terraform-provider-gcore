package cloud_gpu_baremetal_cluster_image

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestAdoptingURL(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		prior, planned types.String
		want           bool
	}{
		"null prior, configured":   {types.StringNull(), types.StringValue("http://x"), true},
		"null prior, unknown plan": {types.StringNull(), types.StringUnknown(), false},
		"null prior, null plan":    {types.StringNull(), types.StringNull(), false},
		"known prior, same value":  {types.StringValue("http://x"), types.StringValue("http://x"), false},
		"known prior, changed":     {types.StringValue("http://x"), types.StringValue("http://y"), false},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			if got := adoptingURL(tc.prior, tc.planned); got != tc.want {
				t.Errorf("adoptingURL(%v, %v) = %v, want %v", tc.prior, tc.planned, got, tc.want)
			}
		})
	}
}

// TestURLAdoptionWarning pins the promises the warning makes and that it never
// carries the (sensitive) value itself.
func TestURLAdoptionWarning(t *testing.T) {
	t.Parallel()

	w := urlAdoptionWarning()
	if w.Severity().String() != "Warning" {
		t.Fatalf("severity = %s, want Warning", w.Severity())
	}
	for _, want := range []string{"imported", "does not return", "Unless this plan replaces", "trusted as-is", "no future refresh", "-replace="} {
		if !strings.Contains(w.Detail(), want) {
			t.Errorf("warning detail is missing %q:\n%s", want, w.Detail())
		}
	}
	if strings.Contains(w.Detail(), "http") {
		t.Errorf("warning must not embed a url value:\n%s", w.Detail())
	}
}
