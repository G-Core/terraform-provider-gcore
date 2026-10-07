package acctest

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"

	"github.com/G-Core/terraform-provider-gcore/internal/sweep"
)

// ProjectID returns the project ID from environment variable
// Used for constructing test configurations
func ProjectID() string {
	return os.Getenv("GCORE_CLOUD_PROJECT_ID")
}

// RegionID returns the region ID from environment variable
// Used for constructing test configurations
func RegionID() string {
	return os.Getenv("GCORE_CLOUD_REGION_ID")
}

// randomSuffixLen is the per-name random part. Six characters of the 35-symbol
// alphabet give a collision chance around 0.007% over a full run; five would
// give 0.24%, which is the wrong trade when the suite is trying to be less
// flaky.
const randomSuffixLen = 6

// RandomName returns a unique test resource name carrying this run's token, so
// a sweeper can tell which run created it: "tf-test-<token>-<random>".
//
// Append qualifiers as a suffix (rName + "-subnet") rather than prefixing them,
// which keeps the tf-test- prefix the sweeper matches on instead of repeating
// it further along the name.
func RandomName() string {
	return fmt.Sprintf("tf-test-%s-%s", runToken(), acctest.RandString(randomSuffixLen))
}

// RandomNameCompact returns a 16-character name for APIs that reject hyphens or
// cap length tightly, such as container-registry users.
//
// It is deliberately not matched by sweep.IsTestResource. Nothing named this way
// is swept directly: registry users are reclaimed with the registry that
// contains them, so these names only need to be unique, which is what stops
// concurrent runs colliding on a shared literal.
func RandomNameCompact() string {
	return fmt.Sprintf("tftest%s%s", runToken(), acctest.RandString(randomSuffixLen))
}

var (
	runTokenOnce  sync.Once
	runTokenValue string
)

// runToken identifies this run. scripts/acctest sets GCORE_ACCTEST_RUN_ID so
// every package in a run shares one token; without it each test binary mints
// its own, which still leaves every resource sweepable but splits one run
// across several tokens.
//
// The fallback logs the run id rather than the token: a sweep hashes
// GCORE_SWEEP_RUN_ID itself, so the token alone cannot reclaim anything.
func runToken() string {
	runTokenOnce.Do(func() {
		if v := sweep.RunToken(); v != "" {
			runTokenValue = v
			return
		}
		runID := fmt.Sprintf("local-%d-%s", os.Getpid(), acctest.RandString(8))
		runTokenValue = sweep.Token(runID)
		log.Printf("[WARN] %s is unset; this test binary is using run id %q. Reclaim its resources with: %s=%s ./scripts/sweep. Run through scripts/acctest so one run shares one id.",
			sweep.EnvRunID, runID, sweep.EnvSweepRunID, runID)
	})
	return runTokenValue
}
