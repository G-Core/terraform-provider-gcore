package sweep

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/G-Core/gcore-go"
)

// reset clears the package-level counters so each case starts clean.
func reset() {
	mu.Lock()
	defer mu.Unlock()
	stats = map[string]*counters{}
	apiSkips = 0
}

func TestFinalExitCode(t *testing.T) {
	t.Run("clean sweep exits zero", func(t *testing.T) {
		reset()
		RecordListed("gcore_cloud_volume", 40)
		RecordDeleted("gcore_cloud_volume", "tf-test-abcd-aaaaaa")
		recordSkip("gcore_cloud_volume")
		if got := FinalExitCode(); got != 0 {
			t.Errorf("FinalExitCode() = %d, want 0", got)
		}
	})

	// A sweeper that logged a failed delete and returned nil used to exit 0,
	// so a run where nothing could be deleted was indistinguishable from one
	// with nothing to delete.
	t.Run("a failed delete fails the run", func(t *testing.T) {
		reset()
		RecordListed("gcore_cloud_volume", 40)
		RecordFailure("gcore_cloud_volume", "tf-test-abcd-aaaaaa", errors.New("volume in use"))
		if got := FinalExitCode(); got != 1 {
			t.Errorf("FinalExitCode() = %d, want 1", got)
		}
	})

	// An account the sweep cannot see looks exactly like an account with
	// nothing left in it.
	t.Run("a refused list fails the run", func(t *testing.T) {
		reset()
		recordAPISkip(http.StatusForbidden)
		if got := FinalExitCode(); got != 1 {
			t.Errorf("FinalExitCode() = %d, want 1", got)
		}
	})
}

// A resource that passes the ownership check but is deliberately left alone
// must still be counted, or listed stops equalling matched plus skipped.
func TestRecordSkippedCountsDeliberateSkips(t *testing.T) {
	reset()
	RecordListed("gcore_cloud_floating_ip", 2)
	RecordSkipped("gcore_cloud_floating_ip", "fip-1", "still attached")
	RecordDeleted("gcore_cloud_floating_ip", "fip-2")

	c := count("gcore_cloud_floating_ip")
	if c.skipped != 1 || c.matched+c.skipped != c.listed {
		t.Errorf("listed=%d matched=%d skipped=%d, want skipped=1 and matched+skipped=listed",
			c.listed, c.matched, c.skipped)
	}
	if got := FinalExitCode(); got != 0 {
		t.Errorf("FinalExitCode() = %d, want 0: a deliberate skip is not a failure", got)
	}
}

func TestSkipSweepError(t *testing.T) {
	apiErr := func(status int) error {
		return &gcore.Error{StatusCode: status}
	}

	for _, tc := range []struct {
		name string
		err  error
		want bool
	}{
		{"nil", nil, false},
		{"403 is a skip", apiErr(http.StatusForbidden), true},
		{"401 is a skip", apiErr(http.StatusUnauthorized), true},
		{"404 is a skip", apiErr(http.StatusNotFound), true},
		{"500 is a real failure", apiErr(http.StatusInternalServerError), false},
		{"wrapped 403 is still a skip", fmt.Errorf("listing volumes: %w", apiErr(http.StatusForbidden)), true},

		// The regression this replaces: "403" was matched as a substring
		// against the whole error text, so an id containing it read as a
		// permission skip and the sweep reported success.
		{"an id containing 403 is not a skip", errors.New("resource sg-403abc not found in project"), false},
		{"prose mentioning forbidden is not a skip", errors.New("operation forbidden by policy"), false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			reset()
			if got := SkipSweepError(tc.err); got != tc.want {
				t.Errorf("SkipSweepError(%v) = %t, want %t", tc.err, got, tc.want)
			}
		})
	}
}
