package sweep

import (
	"errors"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/G-Core/gcore-go"
)

// The only name prefixes a sweeper will delete. Unexported consts rather than
// a slice, so the allowlist cannot be widened at run time; an empty prefix
// would make every name match.
//
// The trailing separator is load-bearing: "tf-test" alone also matches
// unrelated names such as "tf-testing-prod".
const (
	testPrefixHyphen     = "tf-test-"
	testPrefixUnderscore = "tf_test_"
)

// IsTestResource reports whether name was created by this repo's acceptance
// tests.
//
// It is deliberately the narrowest predicate that still covers every name
// generator in the repo, and it must never guess: API keys are client-scoped
// and fourteen sweepers enumerate the whole account with no project boundary,
// so for those this is the only thing standing between a sweep and real
// infrastructure.
func IsTestResource(name string) bool {
	return strings.HasPrefix(name, testPrefixHyphen) ||
		strings.HasPrefix(name, testPrefixUnderscore)
}

// ShouldSweep reports whether a resource may be deleted.
//
// A prefix match is not sufficient. Deletion also requires the resource to
// carry the run token named by EnvSweepRunID, so a sweep can only reclaim a run
// whose completion the caller has established — a scheduled job cannot know
// that, which is why no tokenised resource is swept when the variable is
// unset. The exception is EnvSweepUntokenized, which sweeps names carrying no
// token regardless; see its doc for why it cannot be scoped to a run.
//
// The opt-in is decided here, after the prefix check, and not in
// ShouldSweepToken: the prefix is what still bounds an untokenised sweep to
// test resources, and a resource with no name has no prefix to check.
func ShouldSweep(resourceType, name string) bool {
	if !IsTestResource(name) {
		log.Printf("[DEBUG] Skipping %s: %s (not a test resource)", resourceType, name)
		recordSkip(resourceType)
		return false
	}

	token, tokenized := TokenFromName(name)
	if !tokenized {
		if os.Getenv(EnvSweepUntokenized) == "1" {
			log.Printf("[INFO] Sweeping %s: %s (predates run tokens, %s=1)", resourceType, name, EnvSweepUntokenized)
			return true
		}
		log.Printf("[INFO] Skipping %s: %s (predates run tokens; set %s=1 to reclaim)", resourceType, name, EnvSweepUntokenized)
		recordSkip(resourceType)
		return false
	}

	return ShouldSweepToken(resourceType, name, token)
}

// SkipSweepError reports whether a List error means "this sweeper cannot see
// anything here" rather than "this sweeper failed".
//
// It matches on the HTTP status from the typed SDK error, never on substrings.
// The previous implementation matched a bare "403" anywhere in the error text,
// so an error mentioning a resource id containing "403" silently became a clean
// sweep -- as did a genuine credential failure.
//
// A skip is counted, and a run with any skip exits non-zero: an account the
// sweep could not see is indistinguishable from an account with nothing left in
// it, and treating the two the same is how cleanup appears to work while
// nothing is reclaimed.
func SkipSweepError(err error) bool {
	if err == nil {
		return false
	}

	var apiErr *gcore.Error
	if !errors.As(err, &apiErr) {
		// An error of unknown shape is a real failure, not a skip.
		return false
	}

	switch apiErr.StatusCode {
	case http.StatusUnauthorized, http.StatusForbidden, http.StatusNotFound:
		recordAPISkip(apiErr.StatusCode)
		return true
	}
	return false
}
