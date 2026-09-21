package cloud_gpu_baremetal_cluster_image

import (
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestCowFormatFromDiskFormat(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		diskFormat types.String
		want       types.Bool
	}{
		"raw is cow":          {types.StringValue("raw"), types.BoolValue(true)},
		"qcow2 is not cow":    {types.StringValue("qcow2"), types.BoolValue(false)},
		"other format":        {types.StringValue("vmdk"), types.BoolNull()},
		"null disk format":    {types.StringNull(), types.BoolNull()},
		"unknown disk format": {types.StringUnknown(), types.BoolNull()},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			if got := cowFormatFromDiskFormat(tc.diskFormat); !got.Equal(tc.want) {
				t.Errorf("cowFormatFromDiskFormat(%v) = %v, want %v", tc.diskFormat, got, tc.want)
			}
		})
	}
}

// TestCowFormatDerivationPinned guards the custom hunks in the generated
// resource: a reseal that drops them brings back the forced replacement of an
// imported image whose config sets cow_format = true, or a silent url
// adoption, and every acceptance test that does not set cow_format stays
// green.
func TestCowFormatDerivationPinned(t *testing.T) {
	t.Parallel()

	raw, err := os.ReadFile("resource.go")
	if err != nil {
		t.Fatalf("reading resource.go: %s", err)
	}
	source := string(raw)

	// Anchored at line start so a commented-out assignment does not match.
	assignment := regexp.MustCompile(`(?m)^\s*data\.CowFormat = cowFormatFromDiskFormat\(data\.DiskFormat\)$`)
	for _, method := range []string{"Create", "Update", "ImportState"} {
		if !assignment.MatchString(methodBody(t, source, method)) {
			t.Errorf("%s no longer sets data.CowFormat = cowFormatFromDiskFormat(data.DiskFormat)", method)
		}
	}

	refresh := regexp.MustCompile(`(?m)^\s*if cow := cowFormatFromDiskFormat\(data\.DiskFormat\); !cow\.IsNull\(\) \{\n\s*data\.CowFormat = cow$`)
	if !refresh.MatchString(methodBody(t, source, "Read")) {
		t.Error("Read no longer refreshes data.CowFormat from cowFormatFromDiskFormat(data.DiskFormat)")
	}

	warning := regexp.MustCompile(`(?m)^\s*if adoptingURL\(state\.URL, plan\.URL\) \{\n\s*resp\.Diagnostics\.Append\(urlAdoptionWarning\(\)\)$`)
	if !warning.MatchString(methodBody(t, source, "ModifyPlan")) {
		t.Error("ModifyPlan no longer appends urlAdoptionWarning() under adoptingURL(): the post-import adoption of url is silent again, or warns on every plan")
	}
}

// methodBody returns the source of the named resource method, up to its
// closing brace at column zero.
func methodBody(t *testing.T, source, method string) string {
	t.Helper()

	start := strings.Index(source, ") "+method+"(")
	if start < 0 {
		t.Fatalf("%s not found in resource.go", method)
	}
	body := source[start:]
	if end := strings.Index(body, "\n}\n"); end >= 0 {
		body = body[:end]
	}
	return body
}
