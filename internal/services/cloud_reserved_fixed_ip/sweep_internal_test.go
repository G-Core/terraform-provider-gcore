package cloud_reserved_fixed_ip

import (
	"testing"

	"github.com/G-Core/terraform-provider-gcore/internal/sweep"
)

// TestSweepReservedFixedIPsIsFailClosed pins the one property that makes this
// sweeper tolerable: it establishes no ownership, so it must delete nothing
// unless someone asked for it explicitly.
//
// The credentials below are deliberately fake. If the opt-in were ever honoured
// by accident the sweeper would reach the List call and fail against them, so
// a regression here surfaces as a test failure rather than as deletions in a
// real project.
func TestSweepReservedFixedIPsIsFailClosed(t *testing.T) {
	for _, tc := range []struct {
		name  string
		value string
		set   bool
	}{
		{name: "unset"},
		{name: "empty", value: "", set: true},
		{name: "zero", value: "0", set: true},
		{name: "true is not an opt-in", value: "true", set: true},
		{name: "yes is not an opt-in", value: "yes", set: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("GCORE_API_KEY", "not-a-real-key")
			t.Setenv("GCORE_CLOUD_PROJECT_ID", "1")
			t.Setenv("GCORE_CLOUD_REGION_ID", "1")
			if tc.set {
				t.Setenv(sweep.EnvSweepAllReservedFixedIPs, tc.value)
			}

			if err := sweepCloudReservedFixedIPs(""); err != nil {
				t.Fatalf("expected the sweeper to return without contacting the API, got: %s", err)
			}
		})
	}
}

// TestSweepReservedFixedIPsRequiresEnvironment keeps the credential check ahead
// of the opt-in, so setting the flag alone never starts a sweep against
// whatever account happens to be configured.
func TestSweepReservedFixedIPsRequiresEnvironment(t *testing.T) {
	t.Setenv("GCORE_API_KEY", "")
	t.Setenv("GCORE_CLOUD_PROJECT_ID", "")
	t.Setenv("GCORE_CLOUD_REGION_ID", "")
	t.Setenv(sweep.EnvSweepAllReservedFixedIPs, "1")

	if err := sweepCloudReservedFixedIPs(""); err == nil {
		t.Fatal("expected an error when the sweeper environment is incomplete")
	}
}
