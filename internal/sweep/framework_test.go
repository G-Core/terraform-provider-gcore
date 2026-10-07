package sweep_test

import (
	"testing"

	"github.com/G-Core/terraform-provider-gcore/internal/acctest"
	"github.com/G-Core/terraform-provider-gcore/internal/sweep"
)

func TestIsTestResource(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name string
		want bool
	}{
		// Names the repo's generators actually produce.
		{"tf-test-a1b2c3d4e5", true},
		{"tf-test-net-tf-test-a1b2c3d4e5", true},
		{"tf-test-a1b2c3d4e5.com", true},
		{"tf_test_something", true},

		// Bare random strings. A previous implementation matched any
		// 10-character alphanumeric name, which matched nothing the tests
		// create and every short production name that happened to fit.
		{"a1b2c3d4e5", false},
		{"production", false},
		{"monitoring", false},
		{"kubernetes", false},
		{"prodweb001", false},
		{"webserver1", false},

		// The separator is required, so a longer word starting with the
		// prefix is not a test resource.
		{"tf-testing-prod", false},
		{"tf-tests-live", false},
		{"tf_testing_prod", false},

		// Prefix, not substring.
		{"my-tf-test-thing", false},

		// Match is case-sensitive.
		{"TF-TEST-a1b2c3d4e5", false},

		{"", false},
	} {
		if got := sweep.IsTestResource(tc.name); got != tc.want {
			t.Errorf("IsTestResource(%q) = %t, want %t", tc.name, got, tc.want)
		}
	}
}

// TestRandomNameIsSweepable couples the name generator to the matcher, so the
// two cannot drift into a state where tests create resources no sweeper will
// ever reclaim.
func TestRandomNameIsSweepable(t *testing.T) {
	t.Parallel()

	for i := 0; i < 100; i++ {
		name := acctest.RandomName()
		if !sweep.IsTestResource(name) {
			t.Fatalf("acctest.RandomName() produced %q, which IsTestResource rejects", name)
		}
	}
}
