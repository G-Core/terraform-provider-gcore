package sweep_test

import (
	"testing"

	"github.com/G-Core/terraform-provider-gcore/internal/acctest"
	"github.com/G-Core/terraform-provider-gcore/internal/sweep"
)

func TestToken(t *testing.T) {
	t.Parallel()

	// Both sides derive the token independently from the same run id, so the
	// mapping has to be stable and name-safe.
	const runID = "18273645509"
	got := sweep.Token(runID)

	if len(got) != sweep.TokenLen {
		t.Errorf("Token(%q) = %q, want %d characters", runID, got, sweep.TokenLen)
	}
	if again := sweep.Token(runID); again != got {
		t.Errorf("Token(%q) is not stable: %q then %q", runID, got, again)
	}
	for _, r := range got {
		if !((r >= 'a' && r <= 'z') || (r >= '0' && r <= '9')) {
			t.Errorf("Token(%q) = %q contains %q, which is not name-safe", runID, got, r)
		}
	}

	// Distinct inputs should mostly land on distinct tokens; a hash makes
	// collisions possible but they must not be routine.
	seen := map[string]int{}
	for i := 0; i < 500; i++ {
		seen[sweep.Token(string(rune('a'+i%26))+string(rune('0'+i%10))+string(rune(i)))]++
	}
	if len(seen) < 400 {
		t.Errorf("500 distinct run ids produced only %d tokens", len(seen))
	}
}

func TestTokenFromName(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name      string
		wantToken string
		wantOK    bool
	}{
		{"tf-test-h8ow-yofisa", "h8ow", true},
		{"tf-test-h8ow-yofisa-subnet", "h8ow", true},
		{"tf-test-h8ow-yofisa.com", "h8ow", true},
		{"tf_test_h8ow_yofisa", "h8ow", true},

		// Names from before tokens existed carry none, and must not be
		// mistaken for carrying one.
		{"tf-test-a1b2c3d4e5", "", false},
		{"tf-test-net-tf-test-a1b2c3d4e5", "", false},
		{"tf-test-bucket", "", false},

		// Not test resources at all.
		{"production", "", false},
		{"tf-testing-prod", "", false},
		{"", "", false},

		// Prefix present but nothing after it to parse.
		{"tf-test-", "", false},
		{"tf-test-abc", "", false},
	} {
		token, ok := sweep.TokenFromName(tc.name)
		if ok != tc.wantOK || token != tc.wantToken {
			t.Errorf("TokenFromName(%q) = (%q, %t), want (%q, %t)",
				tc.name, token, ok, tc.wantToken, tc.wantOK)
		}
	}
}

// TestGeneratedNamesCarryTheRunToken is the counterpart to
// TestRandomNameIsSweepable: it ties the generator to the token parser, so a
// change to either cannot leave resources the sweeper can match but not
// attribute to a run.
func TestGeneratedNamesCarryTheRunToken(t *testing.T) {
	t.Setenv(sweep.EnvRunID, "18273645509")
	want := sweep.Token("18273645509")

	for i := 0; i < 50; i++ {
		name := acctest.RandomName()
		token, ok := sweep.TokenFromName(name)
		if !ok {
			t.Fatalf("acctest.RandomName() produced %q, which carries no token", name)
		}
		if token != want {
			t.Fatalf("acctest.RandomName() produced %q with token %q, want %q", name, token, want)
		}
	}
}

func TestShouldSweepRequiresATokenMatch(t *testing.T) {
	const runID = "18273645509"
	name := "tf-test-" + sweep.Token(runID) + "-yofisa"

	t.Run("no sweep run id means no deletion", func(t *testing.T) {
		t.Setenv(sweep.EnvSweepRunID, "")
		if sweep.ShouldSweep("gcore_cloud_volume", name) {
			t.Error("swept without a run id; a sweep must name the finished run it reclaims")
		}
	})

	t.Run("matching token is swept", func(t *testing.T) {
		t.Setenv(sweep.EnvSweepRunID, runID)
		if !sweep.ShouldSweep("gcore_cloud_volume", name) {
			t.Error("did not sweep a resource belonging to the named run")
		}
	})

	t.Run("another run is left alone", func(t *testing.T) {
		t.Setenv(sweep.EnvSweepRunID, "99999999999")
		if sweep.ShouldSweep("gcore_cloud_volume", name) {
			t.Error("swept another run's resource")
		}
	})

	t.Run("untokenized names need an explicit opt-in", func(t *testing.T) {
		legacy := "tf-test-a1b2c3d4e5"
		t.Setenv(sweep.EnvSweepRunID, runID)
		if sweep.ShouldSweep("gcore_cloud_volume", legacy) {
			t.Error("swept a pre-token resource without the opt-in")
		}
		t.Setenv(sweep.EnvSweepUntokenized, "1")
		if !sweep.ShouldSweep("gcore_cloud_volume", legacy) {
			t.Error("opt-in did not reclaim a pre-token resource")
		}
	})

	t.Run("non-test names are never swept", func(t *testing.T) {
		t.Setenv(sweep.EnvSweepRunID, runID)
		t.Setenv(sweep.EnvSweepUntokenized, "1")
		if sweep.ShouldSweep("gcore_cloud_volume", "production") {
			t.Error("swept a non-test resource")
		}
	})
}

func TestShouldSweepTokenForResourcesWithoutNames(t *testing.T) {
	const runID = "18273645509"
	token := sweep.Token(runID)

	t.Run("no sweep run id means no deletion", func(t *testing.T) {
		t.Setenv(sweep.EnvSweepRunID, "")
		if sweep.ShouldSweepToken("gcore_cloud_floating_ip", "fip-1", token) {
			t.Error("swept without a run id")
		}
	})

	t.Run("matching token is swept", func(t *testing.T) {
		t.Setenv(sweep.EnvSweepRunID, runID)
		if !sweep.ShouldSweepToken("gcore_cloud_floating_ip", "fip-1", token) {
			t.Error("did not sweep a resource tagged with the named run")
		}
	})

	t.Run("another run is left alone", func(t *testing.T) {
		t.Setenv(sweep.EnvSweepRunID, "99999999999")
		if sweep.ShouldSweepToken("gcore_cloud_floating_ip", "fip-1", token) {
			t.Error("swept another run's resource")
		}
	})

	// An untagged floating IP is indistinguishable from a reserved production
	// address. Deleting every unattached one is what this replaces, so the
	// untokenized opt-in, which is bounded by the name prefix for everything
	// else, must not reach resources that have no name.
	t.Run("untagged is never swept, even with the opt-in", func(t *testing.T) {
		t.Setenv(sweep.EnvSweepRunID, runID)
		if sweep.ShouldSweepToken("gcore_cloud_floating_ip", "fip-2", "") {
			t.Error("swept an untagged address")
		}
		t.Setenv(sweep.EnvSweepUntokenized, "1")
		if sweep.ShouldSweepToken("gcore_cloud_floating_ip", "fip-2", "") {
			t.Error("untokenized opt-in swept an untagged address")
		}
	})
}
